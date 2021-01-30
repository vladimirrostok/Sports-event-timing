package main

import (
	"log"
	"net/http"
	"sports/backend/srv/api/dashboard"
)

func main() {
	log.Printf("Server listening on http://localhost%s", ":8080")

	c := &dashboard.Dashboard{
		ConnHub: make(map[string]*dashboard.Connection),
		Results: make(chan *dashboard.Result),
		Join:    make(chan *dashboard.Connection),
		Leave:   make(chan *dashboard.Connection),
	}

	http.HandleFunc("/dashboard", c.Handler)

	go c.Run()

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
