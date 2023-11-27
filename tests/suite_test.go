package tests

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/efficientgo/e2e"
	"github.com/google/uuid"
	ginkgo "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nakama_client_go "github.com/joesonw/nakama-client-go"
)

var (
	nakamaHTTPEndpoint string
	nakamaGRPCEndpoint string
	nakamaServerKey    = "defaultkey"
)

func TestAll(t *testing.T) {
	e, err := e2e.New()
	testutil.Ok(t, err)
	t.Cleanup(e.Close)
	db := e.Runnable("cockroachdb").
		WithPorts(map[string]int{
			"grpc": 26257,
			"http": 8080,
		}).
		Init(e2e.StartOptions{
			Image:     "cockroachdb/cockroach:latest-v23.1",
			Command:   e2e.NewCommand("start-single-node", "--insecure", "--store=attrs=ssd,path=/var/lib/cockroach/"),
			Readiness: e2e.NewHTTPReadinessProbe("http", "/health?ready=1", 200, 200),
		})
	testutil.Ok(t, db.Start())
	testutil.Ok(t, db.WaitReady())

	dbAddr := db.InternalEndpoint("grpc")
	nakama := e.Runnable("nakama").
		WithPorts(map[string]int{
			"grpc": 7349,
			"http": 7350,
			"tcp":  7351,
		}).
		Init(e2e.StartOptions{
			Image:     "registry.heroiclabs.com/heroiclabs/nakama:3.19.0",
			User:      "",
			Command:   e2e.NewCommandWithoutEntrypoint("/bin/sh", "-ecx", fmt.Sprintf("/nakama/nakama migrate up --database.address root@%s && exec /nakama/nakama --name nakama --database.address root@%s --logger.level DEBUG --session.token_expiry_sec 7200", dbAddr, dbAddr)),
			Readiness: e2e.NewCmdReadinessProbe(e2e.NewCommand("/nakama/nakama", "healthcheck")),
		})
	testutil.Ok(t, e2e.StartAndWaitReady(nakama))

	nakamaHTTPEndpoint = nakama.Endpoint("http")
	nakamaGRPCEndpoint = nakama.Endpoint("grpc")

	RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Test Suite")
}

func newNakamaHTTPClient() *nakama_client_go.Client {
	return nakama_client_go.NewHTTPClient(nakamaHTTPEndpoint, nakamaServerKey, false, nil)
}

//nolint:unused
func createFacebookInstantGameAuthToken(id string) string {
	testSecret := "fb-instant-test-secret"

	payloadBytes, err := json.Marshal(map[string]any{
		"algorithm":       "HMAC-SHA256",
		"issued_at":       1594867628,
		"player_id":       id,
		"request_payload": "",
	})
	Expect(err).To(BeNil())

	encodedPayload := base64.URLEncoding.EncodeToString(payloadBytes)
	h := hmac.New(sha256.New, []byte(testSecret))
	_, err = h.Write(payloadBytes)
	Expect(err).To(BeNil())
	signature := h.Sum(nil)
	encodedSignature := base64.URLEncoding.EncodeToString(signature)
	return encodedSignature + "." + encodedPayload
}

func generateID() string {
	b, _ := uuid.New().MarshalBinary()
	return hex.EncodeToString(b)
}
