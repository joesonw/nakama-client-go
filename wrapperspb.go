package nakama_client_go

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Bool(b bool) *bool {
	return &b
}

func False() *bool {
	return Bool(false)
}

func True() *bool {
	return Bool(true)
}

func String(s string) *string {
	return &s
}

func boolValue(b *bool) *wrapperspb.BoolValue {
	if b == nil {
		return nil
	}
	return wrapperspb.Bool(*b)
}

func stringValue(s *string) *wrapperspb.StringValue {
	if s == nil {
		return nil
	}
	return wrapperspb.String(*s)
}
