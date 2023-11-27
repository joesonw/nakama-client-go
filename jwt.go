package nakama_client_go

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func unmarshalJWTBody(token string, v interface{}) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 { //nolint:gomnd
		return fmt.Errorf("invalid jwt token")
	}
	body, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
