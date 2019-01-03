package httpclient

import (
	"errors"
	"github.com/liuhaoXD/go-httpclient/mimetype"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type H struct {
	Data string `json:"data"`
	Args struct {
		Param1 string `json:"param1"`
	}
	Headers struct {
		TestHader string `json:"testHader"`
	}
}

func TestGet(t *testing.T) {
	h := &H{}
	r, _ := New().
		Get("http://httpbin.org/get").
		Query("param1", "value").
		UnmarshalJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "value", h.Args.Param1)
}

func TestPut(t *testing.T) {
	r, _ := New().
		Put("http://httpbin.org/put").
		Do()
	assert.Equal(t, 200, r.StatusCode)
}

func TestDelete(t *testing.T) {
	r, _ := New().
		Delete("http://httpbin.org/delete").
		Do()
	assert.Equal(t, 200, r.StatusCode)
}

func TestHead(t *testing.T) {
	r, _ := New().
		Head("http://httpbin.org").
		Do()
	assert.Equal(t, 200, r.StatusCode)
}

func TestValidOk(t *testing.T) {
	_, err := New().Get("http://httpbin.org/get").Do()
	assert.Equal(t, nil, err)
}

func TestValidUrlFail(t *testing.T) {
	_, err := New().Do()
	e := errors.New("url is empty")
	assert.Equal(t, err, e)
}

func TestValidLoggerFail(t *testing.T) {
	_, err := New().Get("dummy").Logger(nil).Do()
	e := errors.New("logger is empty")
	assert.Equal(t, e, err)
}

func TestDummyUrl(t *testing.T) {
	_, err := New().Get("http://123").Do()
	assert.Error(t, err)
}
func TestTimeout(t *testing.T) {
	r, _ := New().
		Timeout(time.Duration(10 * time.Second)).
		Get("http://httpbin.org/get").
		Query("Param1", "value").
		Do()
	assert.Equal(t, 200, r.StatusCode)

}

func TestBasicAuth(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mimetype.ApplicationJson)
		auth := r.Header.Get("Authorization")

		if auth == "Basic dXNlcjpwYXNzd29yZA==" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	r, _ := New().
		BasicAuth("user", "password").
		Get(server.URL).
		Do()
	assert.Equal(t, 200, r.StatusCode)
}
