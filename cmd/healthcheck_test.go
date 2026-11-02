package main

import (
	"net"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestCheck(t *testing.T) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()
	go fasthttp.Serve(ln, newRouter().Handler)
	client := &fasthttp.Client{Dial: func(string) (net.Conn, error) { return ln.Dial() }}

	cases := []struct {
		method, path string
		status       int
		body         string
	}{
		{"GET", "/check", 200, `{"scope":"healthcheck","status":"active"}`},
		{"POST", "/check", 405, ""},
		{"GET", "/nope", 404, ""},
	}
	for _, c := range cases {
		req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
		req.Header.SetMethod(c.method)
		req.SetRequestURI("http://localhost" + c.path)
		if err := client.Do(req, resp); err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode() != c.status {
			t.Errorf("%s %s: status %d, want %d", c.method, c.path, resp.StatusCode(), c.status)
		}
		if c.body != "" && string(resp.Body()) != c.body {
			t.Errorf("%s %s: body %s", c.method, c.path, resp.Body())
		}
		if c.status == 200 && string(resp.Header.ContentType()) != "application/json" {
			t.Errorf("content-type %s", resp.Header.ContentType())
		}
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}
}
