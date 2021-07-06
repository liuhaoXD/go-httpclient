# go-httpclient

A more simple and elegant http client in Go

(fork from [esrest](github.com/easonlin404/esrest))

## Features

- [x] Support basic HTTP __GET__/__POST__/__PUT__/__DELETE__/__HEAD__  in a fluent style
- [x] Only use `Body` fluent function to send payload(JSON/string/slice/pointer/map)
- [x] Basic Authentication
- [x] Request timeout
- [x] Debug with customized Logger
- [x] Receive unmarshal JSON
- [ ] Multipart request
- [ ] [Context](https://golang.org/pkg/context/)
- [ ] application/x-www-form-urlencoded
- [ ] todo

## Installation

```sh
$ go get -u github.com/liuhaoxd/go-httpclient
```

```go
import ghc "github.com/liuhaoxd/go-httpclient"
```

## Usage

__GET__/__POST__/__PUT__/__DELETE__

```go

res, err := ghc.New().Get("http://httpbin.org/get").Do()

```

Add header (Default ContentType is "application/json")

``` go
res, err := ghc.New().
		    Get("http://httpbin.org/get").
		    Header("MyHader", "headvalue").
		    Do()
```

Sending _JSON_ payload use `Body` chain method same as other:

``` go
//JSON struct
res, err := ghc.New().
		    Post("http://httpbin.org/post").
		    Body(struct {
                 		Message string `json:"message"`
                 	}{"ok"}).
		    Do()
//pointer to JSON struct
res, err := ghc.New().
		    Post("http://httpbin.org/post").
		    Body(&struct {
                 		Message string `json:"message"`
                 	}{"ok"}).
		    Do()		    
//slice
res, err := ghc.New().
		    Post("http://httpbin.org/post").
		    Body([]byte(`{"message":"ok"}`)).
		    Do()
		    
//string
res, err := ghc.New().
		    Post("http://httpbin.org/post").
		    Body(string(`{"message":"ok"}`)).
		    Do()
		    
//map
m := map[string]interface{}{
		"message": "ok",
	}
	
res, err := ghc.New().
		    Post("http://httpbin.org/post").
		    Body(m).
		    Do()
```

Add Query parameter:

``` go
res, err := ghc.New().
		    Get("http://httpbin.org/get").
		    Query("Param1", "value").
		    Do()
```

Receive unmarshal JSON:

``` go
json := &struct {
		Message string `json:"message"`
	}{}
res, err := ghc.New().
		    Post("http://httpbin.org/post").
		    DoJson(json)
```

Basic Authentication:

``` go
res, err := ghc.New().
		    BasicAuth("user", "password").
		    Get("http://httpbin.org/get").
		    Do()
```

Debug:

Print http request and response debug payload at stdout, and you also can use your logger by using `Logger` chain

``` go
mylogger:=log.New(os.Stdout, "", log.LstdFlags)

res, err := ghc.New().
		    Debug(true).
		    Logger(mylogger).  //optional
		    Get("http://httpbin.org/get").
		    Do()
```