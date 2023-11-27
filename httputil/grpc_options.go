package httputil

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AttachGRPCToRequest(ctx context.Context, req *http.Request, opts ...grpc.CallOption) *http.Request {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		for k, vv := range md {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
	}

	return req
}

func RetrieveGRPCFromResponse(ctx context.Context, res *http.Response, opts ...grpc.CallOption) *http.Response {
	for _, opt := range opts {
		switch o := opt.(type) { //nolint:gocritic
		case grpc.HeaderCallOption:
			md := metadata.MD{}
			for k, vv := range res.Header {
				md[k] = append(md[k], vv...)
			}
			*o.HeaderAddr = md
		}
	}

	return res
}
