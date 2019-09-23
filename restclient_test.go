package httpclient

import (
	"github.com/liuhaoXD/go-httpclient/mimetype"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPost(t *testing.T) {
	h := &H{}
	b := string(`{"message":"ok"}`)
	r, _ := New().
		Post("http://httpbin.org/post").
		Body(b).
		UnmarshalJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, b, h.Data)

}

func TestHeader(t *testing.T) {
	h := &H{}
	r, _ := New().
		Get("http://httpbin.org/headers").
		Header("testHeader", "value").
		UnmarshalJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "value", h.Headers.TestHeader)
}

func TestDoJsonSuccess(t *testing.T) {
	type Json struct {
		Message string `json:"message"`
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mimetype.ApplicationJson)
		json := []byte(`{"message":"hi"}`)
		w.Write(json)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	json := &Json{}

	New().Get(server.URL).UnmarshalJson(json)
	assert.Equal(t, json.Message, "hi")
}

func TestDoJsonFail(t *testing.T) {
	_, err := New().Get("dummy").UnmarshalJson(nil)
	assert.Error(t, err)
}

func TestSendByteSliceBody(t *testing.T) {
	h := &H{}
	b := []byte(`{"message":"ok"}`) //send byte slice body
	r, _ := New().
		Post("http://httpbin.org/post").
		Body(b).
		UnmarshalJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, string(b), h.Data)

}

func TestSendStructBody(t *testing.T) {
	h := &H{}
	b := struct {
		Message string `json:"message"`
	}{"ok"}

	r, _ := New().
		Post("http://httpbin.org/post").
		Body(b).
		UnmarshalJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, `{"message":"ok"}`, h.Data)
}

func TestSendPtrStructBody(t *testing.T) {
	h := &H{}

	b := &struct {
		Message string `json:"message"`
	}{"ok"}

	r, _ := New().
		Post("http://httpbin.org/post").
		Body(b).
		UnmarshalJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, `{"message":"ok"}`, h.Data)
}

func TestSendMap(t *testing.T) {
	h := &H{}

	b := map[string]interface{}{
		"message": "ok",
	}

	r, _ := New().
		Post("http://httpbin.org/post").
		Body(b).
		UnmarshalJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, `{"message":"ok"}`, h.Data)
}

func TestDebugMode(t *testing.T) {
	h := &H{}
	b := struct {
		Message string `json:"message"`
	}{"ok"}

	r, _ := New().
		Debug(false).
		Header("h", "v").
		Post("http://httpbin.org/post").
		Body(b).
		UnmarshalJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, `{"message":"ok"}`, h.Data)

}

func TestLogger(t *testing.T) {
	h := &H{}
	r, _ := New().
		Logger(log.New(os.Stdout, "", log.LstdFlags)).
		Get("http://httpbin.org/get").
		QuerySet("Param1", "value").
		UnmarshalJson(h)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "value", h.Args.Param1)

}

/*
   apiUrl := "https://api.com"
    resource := "/user/"
    data := Url.Values{}
    data.Set("name", "foo")
    data.Set("surname", "bar")

    u, _ := Url.ParseRequestURI(apiUrl)
    u.Path = resource
    urlStr := u.String() // "https://api.com/user/"

    client := &http.Client{}
    r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
    r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    resp, _ := client.Do(r)
    fmt.Println(resp.Status)
*/
