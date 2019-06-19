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

type user struct {
	GivenName  string
	FamilyName string
}

func main() {

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

		body := []user{
			user{
				GivenName:  "John",
				FamilyName: "Doe",
			},
			user{
				GivenName:  "Jane",
				FamilyName: "Doe",
			},
		}

		jsonResponse(w, body, 200)

	})

	log.Printf("GET /users listening on 9092")
	log.Fatal(http.ListenAndServe(":9092", nil))
}
