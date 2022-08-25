package httphelper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
)

type RequestVerifier = func(r *http.Request) bool

type Ctrl struct {
	*httptest.Server
	FullURL          string
	statusCode       int
	response         any
	requestVerifiers []RequestVerifier
}

func (c *Ctrl) mockHandler(w http.ResponseWriter, r *http.Request) {
	for _, v := range c.requestVerifiers {
		v(r)
	}

	var resp []byte

	rt := reflect.TypeOf(c.response)
	if rt.Kind() == reflect.String {
		resp = []byte(c.response.(string))
	} else if rt.Kind() == reflect.Struct || rt.Kind() == reflect.Ptr {
		resp, _ = json.Marshal(c.response)
	}

	w.WriteHeader(c.statusCode)
	_, err := w.Write(resp)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func HttpMock(pattern string, statusCode int, response any, verifiers ...RequestVerifier) *Ctrl {
	c := &Ctrl{statusCode: statusCode, response: response, requestVerifiers: verifiers}

	handler := http.NewServeMux()
	handler.HandleFunc(pattern, c.mockHandler)

	srv := httptest.NewServer(handler)

	c.Server = srv
	c.FullURL = srv.URL + pattern

	return c
}
