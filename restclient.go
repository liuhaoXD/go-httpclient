package httpclient

import (
	"bytes"
	"encoding/json"
	"github.com/liuhaoXD/go-httpclient/mimetype"
	"io"
	"io/ioutil"
	"net/http"
)

// UnmarshalJson executes the http request client and returns http.Response and error.
func (b *Builder) UnmarshalJson(v interface{}) (*http.Response, error) {
	resp, err := b.Do()
	if err != nil {
		return resp, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err = json.Unmarshal(body, v)
	return resp, err
}

func (b *Builder) newJsonRequest() *http.Request {
	var reader io.Reader
	if len(b.bodyByte) > 0 {
		reader = bytes.NewBuffer(b.bodyByte)
	}

	req, _ := http.NewRequest(b.Method, b.Url, reader)

	//Set Default Content-Type Header
	req.Header.Set("Content-Type", mimetype.ApplicationJson)

	//Set Header
	for k, v := range b.Headers {
		req.Header.Set(k, v)
	}

	//Set Query
	req.URL.RawQuery = b.Queries.Encode()

	if b.basicAuth != (auth{}) {
		req.SetBasicAuth(b.basicAuth.username, b.basicAuth.password)
	}

	return req
}
