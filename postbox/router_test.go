package postbox

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestRequest struct {
	Foo string `json:"foo"`
}

func (t TestRequest) GetName() string {
	return "TestRequest"
}

func (t TestRequest) GetVersion() string {
	return "1.0.0"
}

type TestResponse struct {
	RetFoo string `json:"ret_foo"`
}

func (t TestResponse) GetName() string {
	return "TestResponse"
}

func (t TestResponse) GetVersion() string {
	return "1.0.0"
}

func testHandler(ctx context.Context, req *Envelope) *Envelope {
	r, err := GetPayload[TestRequest](req)
	if err != nil {
		panic(err) // TODO
	}

	resp, err := NewEnvelope[TestResponse](TestResponse{RetFoo: r.Foo})
	if err != nil {
		panic(err) // TODO
	}
	return resp
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	router.Register("TestRequest", "1.x", testHandler)
	env, err := readEnvelope([]byte(`{"n": "TestRequest", "v": "1.2.3", "p": {"msg": "hello world"}}`))

	require.NoError(t, err)
	require.NotNil(t, env)
}
