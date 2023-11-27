package httputil

import (
	"context"
	"fmt"
	"net/http"
)

func CheckResponse(ctx context.Context, res *http.Response) error {
	if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("http error %d: %s", res.StatusCode, res.Status)
	}
	return nil
}
