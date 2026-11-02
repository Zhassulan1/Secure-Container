package main

import (
	"flag"
	"log"

	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func healthcheckHandler(ctx *fasthttp.RequestCtx) {
	js, err := json.Marshal(map[string]interface{}{"status": "active", "scope": "healthcheck"})
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	ctx.Write(js)

}

func newRouter() *router.Router {
	r := router.New()
	r.GET("/check", healthcheckHandler)
	return r
}

func serve(addr string) error {
	log.Printf("Server is listening on %s", addr)
	return fasthttp.ListenAndServe(addr, newRouter().Handler)
}

func main() {
	var port = flag.String("port", "8001", "API server port")
	flag.Parse()

	if err := serve(":" + *port); err != nil {
		log.Fatalf("Error in ListenAndServe: %v", err)
	}
}
