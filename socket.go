package nakama_client_go

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/heroiclabs/nakama-common/rtapi"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	socketJSONMarshaler = protojson.MarshalOptions{
		UseProtoNames: true,
	}
	socketJSONUnmarshaler = protojson.UnmarshalOptions{
		AllowPartial: true,
	}
)

type Socket struct {
	session *Session
	conn    net.Conn
	reader  *wsutil.Reader

	watching  bool
	idCounter *atomic.Int64
}

func (s *Session) NewSocket(ctx context.Context, secure, create bool) (*Socket, error) {
	scheme := "ws"
	if secure {
		scheme = "wss"
	}
	token, err := s.getUpToDatedToken(ctx)
	if err != nil {
		return nil, err
	}

	conn, _, _, err := ws.Dial(ctx, fmt.Sprintf("%s://%s/ws?lang=en&status=%v&token=%s", scheme, s.client.addr, create, token.Token))
	if err != nil {
		return nil, err
	}

	sock := &Socket{
		session:   s,
		conn:      conn,
		reader:    wsutil.NewClientSideReader(conn),
		idCounter: &atomic.Int64{},
	}

	return sock, nil
}

func (s *Socket) Close() error {
	err := s.conn.Close()
	s.conn = nil
	return err
}

func (s *Socket) KeepAlive(interval time.Duration) (func(), error) {
	if s.conn == nil {
		return nil, fmt.Errorf("socket closed")
	}

	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			if s.conn == nil {
				return
			}

			if err := s.Write(context.Background(), &rtapi.Envelope{
				Message: &rtapi.Envelope_Ping{
					Ping: &rtapi.Ping{},
				},
			}); err != nil {
				return
			}
		}
	}()

	return func() {
		ticker.Stop()
	}, nil
}

func (s *Socket) Read(ctx context.Context) (*rtapi.Envelope, error) {
	for {
		frame, err := s.read(ctx)
		if err != nil {
			return nil, err
		}
		if _, ok := frame.Message.(*rtapi.Envelope_Pong); ok {
			continue
		}
		return frame, nil
	}
}

func (s *Socket) read(ctx context.Context) (*rtapi.Envelope, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("socket closed")
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if deadline, ok := ctx.Deadline(); ok {
				if err := s.conn.SetReadDeadline(deadline); err != nil {
					return nil, err
				}
			}
			header, err := s.reader.NextFrame()
			if err != nil {
				return nil, err
			}

			if !header.OpCode.IsData() {
				continue
			}

			buf := make([]byte, header.Length)
			_, err = io.ReadFull(s.reader, buf)
			if err != nil {
				return nil, err
			}

			ev := &rtapi.Envelope{}
			if err := socketJSONUnmarshaler.Unmarshal(buf, ev); err != nil {
				return nil, err
			}
			return ev, nil
		}
	}
}

func (s *Socket) newCID() string {
	return strconv.FormatInt(s.idCounter.Add(1), 10)
}

func (s *Socket) Write(ctx context.Context, message *rtapi.Envelope) error {
	if s.conn == nil {
		return fmt.Errorf("socket closed")
	}

	if message.Cid == "" {
		message.Cid = s.newCID()
	}
	buf, err := socketJSONMarshaler.Marshal(message)
	if err != nil {
		return err
	}

	if deadline, ok := ctx.Deadline(); ok {
		if err := s.conn.SetWriteDeadline(deadline); err != nil {
			return err
		}
	}

	return wsutil.WriteClientMessage(s.conn, ws.OpText, buf)
}
