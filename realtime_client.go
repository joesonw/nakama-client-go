package nakama_client_go

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/rtapi"
)

type RealtimeClient struct {
	socket *Socket
	chExit chan struct{}

	onExit              func(err error)
	onNotification      func(*api.Notification)
	onMatchData         func(*rtapi.MatchData)
	onMatchPresence     func(*rtapi.MatchPresenceEvent)
	onMatchmakerTicket  func(*rtapi.MatchmakerTicket)
	onMatchmakerMatched func(*rtapi.MatchmakerMatched)
	onStatusPresence    func(*rtapi.StatusPresenceEvent)
	onStreamPresence    func(*rtapi.StreamPresenceEvent)
	onStreamData        func(*rtapi.StreamData)
	onChannelMessage    func(*api.ChannelMessage)
	onChannelPresence   func(*rtapi.ChannelPresenceEvent)
	onPartyData         func(*rtapi.PartyData)
	onPartyPresence     func(*rtapi.PartyPresenceEvent)
	onPartyClose        func(*rtapi.PartyClose)
	onPartyJoinRequest  func(*rtapi.PartyJoinRequest)
	onPartyLeader       func(*rtapi.PartyLeader)
	onPartyMatchmaker   func(*rtapi.PartyMatchmakerTicket)
	onParty             func(*rtapi.Party)

	replyPending *sync.Map
}

func (s *Socket) Client() (*RealtimeClient, error) {
	if s.watching {
		return nil, errors.New("already watching")
	}
	rc := &RealtimeClient{
		socket:       s,
		chExit:       make(chan struct{}),
		replyPending: &sync.Map{},
	}
	return rc, nil
}

func (rc *RealtimeClient) Start() {
	if rc.socket.watching {
		return
	}
	go rc.watch()
}

//nolint:gocyclo
func (rc *RealtimeClient) watch() {
	rc.socket.watching = true
	defer func() {
		rc.socket.watching = false
	}()

	for {
		select {
		case <-rc.chExit:
			return
		default:
			ev, err := rc.socket.Read(context.Background())
			if err != nil {
				if rc.onExit != nil {
					rc.onExit(err)
				}
				return
			}
			switch m := ev.Message.(type) {
			case *rtapi.Envelope_Notifications:
				if rc.onNotification != nil {
					for _, notification := range m.Notifications.Notifications {
						rc.onNotification(notification)
					}
				}
			case *rtapi.Envelope_MatchData:
				if rc.onMatchData != nil {
					rc.onMatchData(m.MatchData)
				}

			case *rtapi.Envelope_MatchPresenceEvent:
				if rc.onMatchPresence != nil {
					rc.onMatchPresence(m.MatchPresenceEvent)
				}
			case *rtapi.Envelope_MatchmakerTicket:
				if rc.onMatchmakerTicket != nil {
					rc.onMatchmakerTicket(m.MatchmakerTicket)
				}
			case *rtapi.Envelope_MatchmakerMatched:
				if rc.onMatchmakerMatched != nil {
					rc.onMatchmakerMatched(m.MatchmakerMatched)
				}
			case *rtapi.Envelope_StatusPresenceEvent:
				if rc.onStatusPresence != nil {
					rc.onStatusPresence(m.StatusPresenceEvent)
				}
			case *rtapi.Envelope_StreamPresenceEvent:
				if rc.onStreamPresence != nil {
					rc.onStreamPresence(m.StreamPresenceEvent)
				}
			case *rtapi.Envelope_StreamData:
				if rc.onStreamData != nil {
					rc.onStreamData(m.StreamData)
				}
			case *rtapi.Envelope_ChannelMessage:
				if rc.onChannelMessage != nil {
					rc.onChannelMessage(m.ChannelMessage)
				}
			case *rtapi.Envelope_ChannelPresenceEvent:
				if rc.onChannelPresence != nil {
					rc.onChannelPresence(m.ChannelPresenceEvent)
				}
			case *rtapi.Envelope_PartyData:
				if rc.onPartyData != nil {
					rc.onPartyData(m.PartyData)
				}
			case *rtapi.Envelope_PartyPresenceEvent:
				if rc.onPartyPresence != nil {
					rc.onPartyPresence(m.PartyPresenceEvent)
				}
			case *rtapi.Envelope_PartyClose:
				if rc.onPartyClose != nil {
					rc.onPartyClose(m.PartyClose)
				}
			case *rtapi.Envelope_PartyJoinRequest:
				if rc.onPartyJoinRequest != nil {
					rc.onPartyJoinRequest(m.PartyJoinRequest)
				}
			case *rtapi.Envelope_PartyLeader:
				if rc.onPartyLeader != nil {
					rc.onPartyLeader(m.PartyLeader)
				}
			case *rtapi.Envelope_PartyMatchmakerTicket:
				if rc.onPartyMatchmaker != nil {
					rc.onPartyMatchmaker(m.PartyMatchmakerTicket)
				}
			case *rtapi.Envelope_Party:
				if rc.onParty != nil {
					rc.onParty(m.Party)
				}
			default:
				v, ok := rc.replyPending.LoadAndDelete(ev.Cid)
				if !ok {
					rc.Stop(fmt.Errorf("no reply handler for message: %v", ev))
				}
				ch := v.(chan *rtapi.Envelope)
				ch <- ev
			}
		}
	}
}

func (rc *RealtimeClient) Stop(err error) {
	close(rc.chExit)
	if err != nil {
		rc.onExit(nil)
	}
}

func (rc *RealtimeClient) OnExit(f func(err error)) *RealtimeClient {
	rc.onExit = f
	return rc
}

func (rc *RealtimeClient) OnNotification(f func(*api.Notification)) *RealtimeClient {
	rc.onNotification = f
	return rc
}

func (rc *RealtimeClient) OnMatchData(f func(*rtapi.MatchData)) *RealtimeClient {
	rc.onMatchData = f
	return rc
}

func (rc *RealtimeClient) OnMatchPresence(f func(*rtapi.MatchPresenceEvent)) *RealtimeClient {
	rc.onMatchPresence = f
	return rc
}

func (rc *RealtimeClient) OnMatchmakerTicket(f func(*rtapi.MatchmakerTicket)) *RealtimeClient {
	rc.onMatchmakerTicket = f
	return rc
}

func (rc *RealtimeClient) OnMatchmakerMatched(f func(*rtapi.MatchmakerMatched)) *RealtimeClient {
	rc.onMatchmakerMatched = f
	return rc
}

func (rc *RealtimeClient) OnStatusPresence(f func(*rtapi.StatusPresenceEvent)) *RealtimeClient {
	rc.onStatusPresence = f
	return rc
}

func (rc *RealtimeClient) OnStreamPresence(f func(*rtapi.StreamPresenceEvent)) *RealtimeClient {
	rc.onStreamPresence = f
	return rc
}

func (rc *RealtimeClient) OnStreamData(f func(*rtapi.StreamData)) *RealtimeClient {
	rc.onStreamData = f
	return rc
}

func (rc *RealtimeClient) OnChannelMessage(f func(*api.ChannelMessage)) *RealtimeClient {
	rc.onChannelMessage = f
	return rc
}

func (rc *RealtimeClient) OnChannelPresence(f func(*rtapi.ChannelPresenceEvent)) *RealtimeClient {
	rc.onChannelPresence = f
	return rc
}

func (rc *RealtimeClient) OnPartyData(f func(*rtapi.PartyData)) *RealtimeClient {
	rc.onPartyData = f
	return rc
}

func (rc *RealtimeClient) OnPartyPresence(f func(*rtapi.PartyPresenceEvent)) *RealtimeClient {
	rc.onPartyPresence = f
	return rc
}

func (rc *RealtimeClient) OnPartyClose(f func(*rtapi.PartyClose)) *RealtimeClient {
	rc.onPartyClose = f
	return rc
}

func (rc *RealtimeClient) OnPartyJoinRequest(f func(*rtapi.PartyJoinRequest)) *RealtimeClient {
	rc.onPartyJoinRequest = f
	return rc
}

func (rc *RealtimeClient) OnPartyLeader(f func(*rtapi.PartyLeader)) *RealtimeClient {
	rc.onPartyLeader = f
	return rc
}

func (rc *RealtimeClient) OnPartyMatchmaker(f func(*rtapi.PartyMatchmakerTicket)) *RealtimeClient {
	rc.onPartyMatchmaker = f
	return rc
}

func (rc *RealtimeClient) OnParty(f func(*rtapi.Party)) *RealtimeClient {
	rc.onParty = f
	return rc
}

func (rc *RealtimeClient) sendForResponse(ctx context.Context, envelope *rtapi.Envelope) (*rtapi.Envelope, error) {
	cid := rc.socket.newCID()
	envelope.Cid = cid
	ch := make(chan *rtapi.Envelope, 1)
	defer close(ch)
	rc.replyPending.Store(cid, ch)
	if err := rc.socket.Write(ctx, envelope); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case reply := <-ch:
		if e, ok := reply.Message.(*rtapi.Envelope_Error); ok {
			return nil, &RealtimeError{err: e.Error}
		}
		return reply, nil
	}
}

func (rc *RealtimeClient) AcceptPartyMember(ctx context.Context, partyID string, presence *rtapi.UserPresence) error {
	return rc.socket.Write(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyAccept{
			PartyAccept: &rtapi.PartyAccept{
				PartyId:  partyID,
				Presence: presence,
			},
		},
	})
}

func (rc *RealtimeClient) AddMatchmaker(ctx context.Context, minCount, maxCount int, query string, stringProperties map[string]string, numericProperties map[string]float64) (*rtapi.MatchmakerTicket, error) {
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchmakerAdd{
			MatchmakerAdd: &rtapi.MatchmakerAdd{
				MinCount:          int32(minCount),
				MaxCount:          int32(maxCount),
				Query:             query,
				StringProperties:  stringProperties,
				NumericProperties: numericProperties,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetMatchmakerTicket(), nil
}

func (rc *RealtimeClient) AddMatchmakerParty(ctx context.Context, partyID string, minCount, maxCount int, query string, stringProperties map[string]string, numericProperties map[string]float64) (*rtapi.PartyMatchmakerTicket, error) {
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyMatchmakerAdd{
			PartyMatchmakerAdd: &rtapi.PartyMatchmakerAdd{
				PartyId:           partyID,
				MinCount:          int32(minCount),
				MaxCount:          int32(maxCount),
				Query:             query,
				StringProperties:  stringProperties,
				NumericProperties: numericProperties,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetPartyMatchmakerTicket(), nil
}

func (rc *RealtimeClient) CloseParty(ctx context.Context, partyID string) error {
	_, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyClose{
			PartyClose: &rtapi.PartyClose{
				PartyId: partyID,
			},
		},
	})
	return err
}

func (rc *RealtimeClient) CreateMatch(ctx context.Context, name string) (*rtapi.Match, error) {
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchCreate{
			MatchCreate: &rtapi.MatchCreate{
				Name: name,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetMatch(), nil
}

func (rc *RealtimeClient) CreateParty(ctx context.Context, open bool, maxSize int32) (*rtapi.Party, error) {
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyCreate{
			PartyCreate: &rtapi.PartyCreate{
				Open:    open,
				MaxSize: maxSize,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetParty(), nil
}

func (rc *RealtimeClient) FollowUsers(ctx context.Context, userIds []string) (*rtapi.Status, error) {
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_StatusFollow{
			StatusFollow: &rtapi.StatusFollow{
				UserIds: userIds,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetStatus(), nil
}

type JoinChatOption struct {
	Persistence *bool
	Hidden      *bool
}

func (rc *RealtimeClient) JoinChat(ctx context.Context, target string, typ rtapi.ChannelJoin_Type, opts ...JoinChatOption) (*rtapi.Channel, error) {
	opt := JoinChatOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelJoin{
			ChannelJoin: &rtapi.ChannelJoin{
				Target:      target,
				Type:        int32(typ),
				Persistence: boolValue(opt.Persistence),
				Hidden:      boolValue(opt.Hidden),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetChannel(), nil
}

func (rc *RealtimeClient) JoinMatch(ctx context.Context, matchID, token string, metadata map[string]string) (*rtapi.Match, error) {
	m := &rtapi.MatchJoin{
		Metadata: metadata,
	}
	if matchID != "" {
		m.Id = &rtapi.MatchJoin_MatchId{
			MatchId: matchID,
		}
	} else {
		m.Id = &rtapi.MatchJoin_Token{
			Token: token,
		}
	}
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchJoin{
			MatchJoin: m,
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetMatch(), nil
}

func (rc *RealtimeClient) JoinParty(ctx context.Context, partyID string) error {
	_, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyJoin{
			PartyJoin: &rtapi.PartyJoin{
				PartyId: partyID,
			},
		},
	})
	return err
}

func (rc *RealtimeClient) LeaveChat(ctx context.Context, channelID string) error {
	_, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelLeave{
			ChannelLeave: &rtapi.ChannelLeave{
				ChannelId: channelID,
			},
		},
	})
	return err
}

func (rc *RealtimeClient) LeaveMatch(ctx context.Context, matchID string) error {
	_, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchLeave{
			MatchLeave: &rtapi.MatchLeave{
				MatchId: matchID,
			},
		},
	})
	return err
}

func (rc *RealtimeClient) LeaveParty(ctx context.Context, partyID string) error {
	_, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyLeave{
			PartyLeave: &rtapi.PartyLeave{
				PartyId: partyID,
			},
		},
	})
	return err
}

func (rc *RealtimeClient) ListPartyJoinRequests(ctx context.Context, partyID string) (*rtapi.PartyJoinRequest, error) {
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyJoinRequestList{
			PartyJoinRequestList: &rtapi.PartyJoinRequestList{
				PartyId: partyID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetPartyJoinRequest(), nil
}

func (rc *RealtimeClient) PromotePartyMember(ctx context.Context, partyID string, partyMember *rtapi.UserPresence) (*rtapi.PartyLeader, error) {
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyPromote{
			PartyPromote: &rtapi.PartyPromote{
				PartyId:  partyID,
				Presence: partyMember,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetPartyLeader(), nil
}

func (rc *RealtimeClient) RemoveChatMessage(ctx context.Context, channelID, messageID string) (*rtapi.ChannelMessageAck, error) {
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelMessageRemove{
			ChannelMessageRemove: &rtapi.ChannelMessageRemove{
				ChannelId: channelID,
				MessageId: messageID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetChannelMessageAck(), nil
}

func (rc *RealtimeClient) RemoveMatchmaker(ctx context.Context, ticket string) error {
	_, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchmakerRemove{
			MatchmakerRemove: &rtapi.MatchmakerRemove{
				Ticket: ticket,
			},
		},
	})
	return err
}

func (rc *RealtimeClient) RemoveMatchmakerParty(ctx context.Context, partyID, ticket string) error {
	_, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyMatchmakerRemove{
			PartyMatchmakerRemove: &rtapi.PartyMatchmakerRemove{
				PartyId: partyID,
				Ticket:  ticket,
			},
		},
	})
	return err
}

func (rc *RealtimeClient) RemovePartyMember(ctx context.Context, partyID string, presence *rtapi.UserPresence) error {
	_, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyRemove{
			PartyRemove: &rtapi.PartyRemove{
				PartyId:  partyID,
				Presence: presence,
			},
		},
	})
	return err
}

func (rc *RealtimeClient) RPC(ctx context.Context, id, payload, httpKey string) (*api.Rpc, error) {
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_Rpc{
			Rpc: &api.Rpc{
				Id:      id,
				Payload: payload,
				HttpKey: httpKey,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetRpc(), nil
}

func (rc *RealtimeClient) SendMatchState(ctx context.Context, matchID string, opCode int64, data []byte, presences []*rtapi.UserPresence, reliable bool) error {
	return rc.socket.Write(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchDataSend{
			MatchDataSend: &rtapi.MatchDataSend{
				MatchId:   matchID,
				OpCode:    opCode,
				Data:      data,
				Presences: presences,
				Reliable:  reliable,
			},
		},
	})
}
func (rc *RealtimeClient) SendPartyData(ctx context.Context, partyID string, opCode int64, data []byte) error {
	return rc.socket.Write(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyDataSend{
			PartyDataSend: &rtapi.PartyDataSend{
				PartyId: partyID,
				OpCode:  opCode,
				Data:    data,
			},
		},
	})
}

func (rc *RealtimeClient) UpdateStatus(ctx context.Context, status *string) error {
	_, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_StatusUpdate{
			StatusUpdate: &rtapi.StatusUpdate{
				Status: stringValue(status),
			},
		},
	})
	return err
}

func (rc *RealtimeClient) WriteChatMessage(ctx context.Context, channelID, content string) (*rtapi.ChannelMessageAck, error) {
	ev, err := rc.sendForResponse(ctx, &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelMessageSend{
			ChannelMessageSend: &rtapi.ChannelMessageSend{
				ChannelId: channelID,
				Content:   content,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return ev.GetChannelMessageAck(), nil
}

type RealtimeError struct {
	err *rtapi.Error
}

func (e *RealtimeError) Code() int32 {
	return e.err.Code
}

func (e *RealtimeError) Message() string {
	return e.err.Message
}

func (e *RealtimeError) Context() map[string]string {
	return e.err.Context
}

func (e *RealtimeError) Error() string {
	return fmt.Sprintf("realtime error(code=%d, message=%s, context=%v)", e.err.Code, e.err.Message, e.err.Context)
}

func IsRealtimeError(err error) bool {
	_, ok := err.(*RealtimeError)
	return ok
}

func AsRealtimeError(err error) *RealtimeError {
	if e, ok := err.(*RealtimeError); ok {
		return e
	}
	return nil
}
