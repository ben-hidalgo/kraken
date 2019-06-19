// Package main implements a server for Greeter service.
package main

import (
	"encoding/json"
	"log"
	"net/http"
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

type scene struct {
	Name string
}

func main() {

	http.HandleFunc("/scenes", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

		body := []scene{
			scene{
				Name: "John",
			},
			scene{
				Name: "Jane",
			},
		}

		jsonResponse(w, body, 200)

	})

	log.Printf("GET /scenes listening on 9091")
	log.Fatal(http.ListenAndServe(":9091", nil))
}
