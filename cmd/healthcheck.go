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
	ctx.Write(js)

}

func main() {
	r := router.New()
	r.GET("/check", healthcheckHandler)

	var port = flag.String("port", "8001", "API server port")
	flag.Parse()

	addr := ":" + *port
	log.Printf("Server is listening on %s", addr)
	if err := fasthttp.ListenAndServe(addr, r.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe: %v", err)
	}
}
