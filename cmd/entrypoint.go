package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
)

var (
	address = flag.String("addr", "0.0.0.0:8000", "Host address where the app will bind to")
)

func main() {
	flag.Parse()

	// Websocket Manager initialization and Message Handlers
	wm := NewWebsocketManager(10)
	wm.AddMessageHandler("streamer", HandleStreamerMessage)
	wm.AddMessageHandler("offer", HandleOfferMessage)
	wm.AddMessageHandler("answer", HandleAnswerMessage)
	wm.AddMessageHandler("createOffer", HandleSendOfferMessage)
	wm.AddMessageHandler("createAnswer", HandleSendAnswerMessage)
	wm.AddMessageHandler("candidate", HandleCandidateMessage)

	// HTTP handlers
	http.Handle("GET /", http.FileServer(http.Dir("./pages/")))

	http.HandleFunc("GET /newconn", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("new connection", slog.String("peer", r.RemoteAddr))

		err := wm.AppendConnection(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("connection upgraded", slog.String("peer", r.RemoteAddr))
	})

	// err := http.ListenAndServe(*address, nil)
	err := http.ListenAndServeTLS(*address, "cert/cert.pem", "cert/key.pem", nil)

	if err != nil {
		log.Fatal(err)
	}
}
