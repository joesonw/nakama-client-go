package httputil

import (
	"context"
	"encoding/base64"

	"github.com/heroiclabs/nakama-common/api"
	"google.golang.org/grpc/metadata"
)

type (
	contextKeyBasicAuth   struct{}
	contextKeyBearerJWT   struct{}
	contextKeyHTTPKeyAuth struct{}

	contextBasicAuth struct {
		username string
		password string
	}
)

func WithBasicAuth(ctx context.Context, username, password string) context.Context {
	ctx = context.WithValue(ctx, contextKeyBasicAuth{}, &contextBasicAuth{username: username, password: password})
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
	return ctx
}

func WithBearerJWT(ctx context.Context, token string) context.Context {
	ctx = context.WithValue(ctx, contextKeyBearerJWT{}, token)
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	return ctx
}

func WithHTTPKeyAuth(ctx context.Context, key string) context.Context {
	ctx = context.WithValue(ctx, contextKeyHTTPKeyAuth{}, key)
	ctx = metadata.AppendToOutgoingContext(ctx, "q_http_key", key)
	return ctx
}

func WithSession(ctx context.Context, sess *api.Session) context.Context {
	ctx = WithBearerJWT(ctx, sess.Token)
	return ctx
}

func GetBasicAuth(ctx context.Context) (username, password string) {
	auth, ok := ctx.Value(contextKeyBasicAuth{}).(*contextBasicAuth)
	if !ok {
		return "", ""
	}
	return auth.username, auth.password
}

func GetBearerJWT(ctx context.Context) string {
	token, _ := ctx.Value(contextKeyBearerJWT{}).(string)
	return token
}

func GetHTTPKeyAuth(ctx context.Context) string {
	key, _ := ctx.Value(contextKeyHTTPKeyAuth{}).(string)
	return key
}
