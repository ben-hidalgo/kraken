// Package main implements a server for Greeter service.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func jsonResponse(w http.ResponseWriter, body interface{}, status int) {

	j, err := json.Marshal(body)

	if err != nil {
		log.Printf("JsonResponse() err=%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

type prox struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func newProxy(target string) *prox {
	url, _ := url.Parse(target)
	return &prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *prox) handle(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("X-Forwarded-Host", r.Host)
	r.Header.Add("X-Origin-Host", "127.0.0.1")
	p.proxy.ServeHTTP(w, r)
}

func main() {

	// proxy
	scenesProxy := newProxy("http://127.0.0.1:9091")
	usersProxy := newProxy("http://127.0.0.1:9092")

	http.HandleFunc("/scenes", scenesProxy.handle)
	http.HandleFunc("/users", usersProxy.handle)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		body := map[string]string{
			"rproxy": "listening",
		}
		jsonResponse(w, body, 200)
	})

	log.Printf("GET /       listening on 9090")
	log.Printf("GET /users  listening on 9090")
	log.Printf("GET /scenes listening on 9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
