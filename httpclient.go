package httpclient

import (
	"bytes"
	"encoding/json"
	"github.com/liuhaoXD/go-httpclient/mimetype"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	Queries   url.Values
	debugMode bool

	logger    *log.Logger
	timeout   time.Duration
	basicAuth auth
	bodyByte  []byte
}

// New returns an new Builder object.
func New() *Builder {
	return &Builder{
		Headers: make(map[string]string),
		Queries: make(url.Values),
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

// QuerySet sets http QUERY parameter key with value returns Builder.
func (b *Builder) QuerySet(key string, value string) *Builder {
	b.Queries.Set(key, value)
	return b
}

// QueryAdd adds http QUERY parameter key with value and returns Builder.
func (b *Builder) QueryAdd(key, value string) *Builder {
	b.Queries.Add(key, value)
	return b
}

func (b *Builder) ContentType(mimeType string) *Builder {
	b.Header("Content-Type", mimeType)
	return b
}

// Do executes the http request client and returns http.Response and error.
func (b *Builder) Do() (*http.Response, error) {
	if err := b.valid(); err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: b.timeout}

	request, err := b.newRequest()
	if err != nil {
		return nil, err
	}

	if b.debugMode {
		dump, _ := httputil.DumpRequest(request, true)
		b.logger.Println(string(dump))
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if b.debugMode {
		dump, _ := httputil.DumpResponse(resp, true)
		b.logger.Println(string(dump))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

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
		return ErrUrlIsEmpty
	}
	if b.logger == nil {
		return ErrLoggerIsEmpty
	}
	return nil
}

func (b *Builder) Debug(debug bool) *Builder {
	b.debugMode = debug
	return b
}

type auth struct {
	username string
	password string
}

// Body convert v to []byte and using ad http request body
func (b *Builder) Body(v interface{}) *Builder {
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.String:
		b.bodyByte = []byte(rv.String())
	case reflect.Slice:
		sliceValue, _ := rv.Interface().([]byte)
		b.bodyByte = sliceValue
	case reflect.Map, reflect.Struct, reflect.Ptr:
		byteValue, err := json.Marshal(v)
		if err != nil {
			log.Println("failed to marshal body ", err)
		}
		b.bodyByte = byteValue
	}
	return b
}

func (b *Builder) UrlEncodedBody(v map[string]string) *Builder {
	var buffer bytes.Buffer
	if len(v) > 0 {
		i := 0
		for key, val := range v {
			buffer.WriteString(key)
			buffer.WriteByte('=')
			buffer.WriteString(val)
			if i < len(v)-1 {
				buffer.WriteByte('&')
				i++
			}
		}
		b.bodyByte = buffer.Bytes()
	}
	return b
}

func (b *Builder) StringBody(v string) *Builder {
	b.bodyByte = []byte(v)
	return b
}

func (b *Builder) JsonBody(v interface{}) *Builder {
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Map, reflect.Struct, reflect.Ptr:
		byteValue, err := json.Marshal(v)
		if err != nil {
			log.Println("failed to marshal body ", err)
		}
		b.bodyByte = byteValue
	}
	return b
}

func (b *Builder) newRequest() (*http.Request, error) {
	var reader io.Reader
	if len(b.bodyByte) > 0 {
		reader = bytes.NewBuffer(b.bodyByte)
	}

	req, err := http.NewRequest(b.Method, b.Url, reader)
	if err != nil {
		return nil, err
	}

	//Set Default Content-Type Header
	if len(req.Header.Get("Content-Type")) == 0 {
		req.Header.Set("Content-Type", mimetype.ApplicationJson)
	}

	//Set Header
	for k, v := range b.Headers {
		req.Header.Set(k, v)
	}

	req.URL.RawQuery = b.Queries.Encode()

	if b.basicAuth != (auth{}) {
		req.SetBasicAuth(b.basicAuth.username, b.basicAuth.password)
	}

	return req, nil
}
