package nakama_client_go

import (
	"context"
	"time"

	"github.com/heroiclabs/nakama-common/api"

	"github.com/joesonw/nakama-client-go/httputil"
)

type Client struct {
	api       apiClient
	serverKey string
	addr      string
	secure    bool
}

func (c *Client) Close() error {
	return c.api.Close()
}

type AuthenticateOption struct {
	Create   *bool
	Vars     map[string]string
	Username string
}

func parseAuthenticateOptions(opts ...AuthenticateOption) AuthenticateOption {
	opt := AuthenticateOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (c *Client) AuthenticateApple(ctx context.Context, token string, opts ...AuthenticateOption) (*Session, error) {
	opt := parseAuthenticateOptions(opts...)

	res, err := c.api.AuthenticateApple(httputil.WithBasicAuth(ctx, c.serverKey, ""), &api.AuthenticateAppleRequest{
		Account: &api.AccountApple{
			Token: token,
			Vars:  opt.Vars,
		},
		Create:   boolValue(opt.Create),
		Username: opt.Username,
	})
	if err != nil {
		return nil, err
	}

	return newSession(res, c)
}

func (c *Client) AuthenticateCustom(ctx context.Context, id string, opts ...AuthenticateOption) (*Session, error) {
	opt := parseAuthenticateOptions(opts...)

	res, err := c.api.AuthenticateCustom(httputil.WithBasicAuth(ctx, c.serverKey, ""), &api.AuthenticateCustomRequest{
		Account: &api.AccountCustom{
			Id:   id,
			Vars: opt.Vars,
		},
		Create:   boolValue(opt.Create),
		Username: opt.Username,
	})
	if err != nil {
		return nil, err
	}

	return newSession(res, c)
}

func (c *Client) AuthenticateDevice(ctx context.Context, id string, opts ...AuthenticateOption) (*Session, error) {
	opt := parseAuthenticateOptions(opts...)

	res, err := c.api.AuthenticateDevice(httputil.WithBasicAuth(ctx, c.serverKey, ""), &api.AuthenticateDeviceRequest{
		Account: &api.AccountDevice{
			Id:   id,
			Vars: opt.Vars,
		},
		Create:   boolValue(opt.Create),
		Username: opt.Username,
	})
	if err != nil {
		return nil, err
	}

	return newSession(res, c)
}

func (c *Client) AuthenticateEmail(ctx context.Context, email, password string, opts ...AuthenticateOption) (*Session, error) {
	opt := parseAuthenticateOptions(opts...)

	res, err := c.api.AuthenticateEmail(httputil.WithBasicAuth(ctx, c.serverKey, ""), &api.AuthenticateEmailRequest{
		Account: &api.AccountEmail{
			Email:    email,
			Password: password,
			Vars:     opt.Vars,
		},
		Create:   boolValue(opt.Create),
		Username: opt.Username,
	})
	if err != nil {
		return nil, err
	}

	return newSession(res, c)
}

func (c *Client) AuthenticateFacebook(ctx context.Context, token string, opts ...AuthenticateOption) (*Session, error) {
	opt := parseAuthenticateOptions(opts...)

	res, err := c.api.AuthenticateFacebook(httputil.WithBasicAuth(ctx, c.serverKey, ""), &api.AuthenticateFacebookRequest{
		Account: &api.AccountFacebook{
			Token: token,
			Vars:  opt.Vars,
		},
		Create:   boolValue(opt.Create),
		Username: opt.Username,
	})
	if err != nil {
		return nil, err
	}

	return newSession(res, c)
}

func (c *Client) AuthenticateFacebookInstantGame(ctx context.Context, signedPlayerInfo string, opts ...AuthenticateOption) (*Session, error) {
	opt := parseAuthenticateOptions(opts...)

	res, err := c.api.AuthenticateFacebookInstantGame(httputil.WithBasicAuth(ctx, c.serverKey, ""), &api.AuthenticateFacebookInstantGameRequest{
		Account: &api.AccountFacebookInstantGame{
			SignedPlayerInfo: signedPlayerInfo,
			Vars:             opt.Vars,
		},
		Create:   boolValue(opt.Create),
		Username: opt.Username,
	})
	if err != nil {
		return nil, err
	}

	return newSession(res, c)
}

func (c *Client) AuthenticateGameCenter(ctx context.Context, bundleId, playerId, publicKeyUrl, salt, signature string, timestamp time.Time, opts ...AuthenticateOption) (*Session, error) {
	opt := parseAuthenticateOptions(opts...)

	res, err := c.api.AuthenticateGameCenter(httputil.WithBasicAuth(ctx, c.serverKey, ""), &api.AuthenticateGameCenterRequest{
		Account: &api.AccountGameCenter{
			BundleId:         bundleId,
			PlayerId:         playerId,
			PublicKeyUrl:     publicKeyUrl,
			Salt:             salt,
			Signature:        signature,
			TimestampSeconds: timestamp.Unix(),
			Vars:             opt.Vars,
		},
		Create:   boolValue(opt.Create),
		Username: opt.Username,
	})
	if err != nil {
		return nil, err
	}

	return newSession(res, c)
}

func (c *Client) AuthenticateGoogle(ctx context.Context, token string, opts ...AuthenticateOption) (*Session, error) {
	opt := parseAuthenticateOptions(opts...)

	res, err := c.api.AuthenticateGoogle(httputil.WithBasicAuth(ctx, c.serverKey, ""), &api.AuthenticateGoogleRequest{
		Account: &api.AccountGoogle{
			Token: token,
			Vars:  opt.Vars,
		},
		Create:   boolValue(opt.Create),
		Username: opt.Username,
	})
	if err != nil {
		return nil, err
	}

	return newSession(res, c)
}

func (c *Client) AuthenticateSteam(ctx context.Context, token string, opts ...AuthenticateOption) (*Session, error) {
	opt := parseAuthenticateOptions(opts...)

	res, err := c.api.AuthenticateSteam(httputil.WithBasicAuth(ctx, c.serverKey, ""), &api.AuthenticateSteamRequest{
		Account: &api.AccountSteam{
			Token: token,
			Vars:  opt.Vars,
		},
		Create:   boolValue(opt.Create),
		Username: opt.Username,
	})
	if err != nil {
		return nil, err
	}

	return newSession(res, c)
}

func (c *Client) RPC(ctx context.Context, httpKey, id string, payload []byte) (*api.Rpc, error) {
	return c.api.RpcFunc2(httputil.WithHTTPKeyAuth(ctx, httpKey), &api.Rpc{
		Id:      id,
		Payload: string(payload),
	})
}
