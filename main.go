package main

import (
	"net/http"
	_ "net/http/pprof"
	"ws/app"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	go http.ListenAndServe(":10000", nil)
	app.Start()

}
