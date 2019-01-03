package httpclient

import (
	"bytes"
	"encoding/json"
	"github.com/liuhaoXD/go-httpclient/mimetype"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"time"
)

// Builder is a object that help to build fluent style API.
type Builder struct {
	Url       string
	Method    string
	Path      string
	Headers   map[string]string
	Querys    map[string]string
	DebugMode bool

	logger    *log.Logger
	timeout   time.Duration
	basicAuth auth
	bodyByte  []byte
}

// New returns an new Builder object.
func New() *Builder {
	return &Builder{
		Headers: make(map[string]string),
		Querys:  make(map[string]string),
		logger:  log.New(os.Stdout, "", log.LstdFlags),
		timeout: time.Duration(20 * time.Second),
	}
}

// Get uses the http GET method with provided url and returns Builder.
func (b *Builder) Get(url string) *Builder {
	b.Url = url
	b.Method = http.MethodGet
	return b
}

// Post uses the http POST method with provided url and returns Builder.
func (b *Builder) Post(url string) *Builder {
	b.Url = url
	b.Method = http.MethodPost
	return b
}

// Put uses the http PUT method with provided url and returns Builder.
func (b *Builder) Put(url string) *Builder {
	b.Url = url
	b.Method = http.MethodPut
	return b
}

// Delete uses the http DELETE method with provided url and returns Builder.
func (b *Builder) Delete(url string) *Builder {
	b.Url = url
	b.Method = http.MethodDelete
	return b
}

// Head uses the http HEAD method with provided url and returns Builder.
func (b *Builder) Head(url string) *Builder {
	b.Url = url
	b.Method = http.MethodHead
	return b
}

// Header sets http header key with value and returns Builder.
func (b *Builder) Header(key, value string) *Builder {
	b.Headers[key] = value
	return b
}

// Query sets http QUERY parameter key with value and returns Builder.
func (b *Builder) Query(key, value string) *Builder {
	b.Querys[key] = value
	return b
}

// Do executes the http request client and returns http.Response and error.
func (b *Builder) Do() (*http.Response, error) {
	if err := b.valid(); err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: b.timeout}

	request := b.newRequest()

	if b.DebugMode {
		dump, _ := httputil.DumpRequest(request, true)
		b.logger.Println(string(dump))
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if b.DebugMode {
		dump, _ := httputil.DumpResponse(resp, true)
		b.logger.Println(string(dump))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return resp, resp.Body.Close()
}

// Logger sets the provided logger and returns Builder.
func (b *Builder) Logger(log *log.Logger) *Builder {
	b.logger = log
	return b
}

func (b *Builder) Timeout(timeout time.Duration) *Builder {
	b.timeout = timeout
	return b
}

func (b *Builder) BasicAuth(username, password string) *Builder {
	b.basicAuth = auth{username: username, password: password}
	return b
}

func (b *Builder) valid() error {
	if len(b.Url) == 0 {
		return errors.New("url is empty")
	}
	if b.logger == nil {
		return errors.New("logger is empty")
	}
	return nil
}

func (b *Builder) Debug(debug bool) *Builder {
	b.DebugMode = debug
	return b
}

type auth struct {
	username string
	password string
}

func (b *Builder) Body(v interface{}) *Builder {
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.String:
		b.bodyByte = []byte(rv.String())
	case reflect.Slice:
		sliceValue, _ := rv.Interface().([]byte)
		b.bodyByte = sliceValue
	case reflect.Map, reflect.Struct, reflect.Ptr:
		byteValue, _ := json.Marshal(v)
		b.bodyByte = byteValue
	}
	return b
}

func (b *Builder) newRequest() *http.Request {
	var reader io.Reader
	if len(b.bodyByte) > 0 {
		reader = bytes.NewBuffer(b.bodyByte)
	}

	req, _ := http.NewRequest(b.Method, b.Url, reader)

	//Set Default Content-Type Header
	if len(req.Header.Get("Content-Type")) == 0 {
		req.Header.Set("Content-Type", mimetype.ApplicationJson)
	}

	//Set Header
	for k, v := range b.Headers {
		req.Header.Set(k, v)
	}

	//Set Query
	q := req.URL.Query()
	for k, v := range b.Querys {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	if b.basicAuth != (auth{}) {
		req.SetBasicAuth(b.basicAuth.username, b.basicAuth.password)
	}

	return req
}
