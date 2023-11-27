package nakama_client_go

import (
	"context"
	"fmt"
	"net/http"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama/v3/apigrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/joesonw/nakama-client-go/internal/pb"
)

type apiClient interface {
	pb.NakamaClientInterface
	Close() error
}

type grpcClient struct {
	apigrpc.NakamaClient
	cc *grpc.ClientConn
}

func NewGRPCClient(addr, serverKey string, secure bool, opts ...grpc.DialOption) (*Client, error) {
	newOpts := append([]grpc.DialOption{}, opts...)
	if !secure {
		newOpts = append(newOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(addr, newOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	client := apigrpc.NewNakamaClient(cc)
	return &Client{
		api: &grpcClient{
			NakamaClient: client,
			cc:           cc,
		},
		addr:      addr,
		serverKey: serverKey,
		secure:    secure,
	}, nil
}

func (c *grpcClient) Close() error {
	return c.cc.Close()
}

func (c *grpcClient) ListStorageObjects2(ctx context.Context, req *api.ListStorageObjectsRequest, opts ...grpc.CallOption) (*api.StorageObjectList, error) {
	return c.NakamaClient.ListStorageObjects(ctx, req, opts...)
}

func (c *grpcClient) RpcFunc2(ctx context.Context, req *api.Rpc, opts ...grpc.CallOption) (*api.Rpc, error) {
	return c.NakamaClient.RpcFunc(ctx, req, opts...)
}

func (c *grpcClient) WriteTournamentRecord2(ctx context.Context, req *api.WriteTournamentRecordRequest, opts ...grpc.CallOption) (*api.LeaderboardRecord, error) {
	return c.NakamaClient.WriteTournamentRecord(ctx, req, opts...)
}

type httpClient struct {
	*pb.NakamaClient
}

func NewHTTPClient(addr, serverKey string, secure bool, client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}
	url := addr
	if secure {
		url = "https://" + url
	} else {
		url = "http://" + url
	}
	return &Client{
		api: &httpClient{
			NakamaClient: &pb.NakamaClient{
				URL:    url,
				Client: client,
			},
		},
		addr:      addr,
		serverKey: serverKey,
		secure:    secure,
	}
}

func (c *httpClient) Close() error {
	return nil
}
