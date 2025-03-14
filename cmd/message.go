package main

import (
	"log/slog"
)

type Message struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data,omitempty"`
}

type MessageHandlerFunc = func(*Stream, []*Stream, Message)

func HandleStreamerMessage(sender *Stream, pool []*Stream, message Message) {
	slog.Info("streamer message received", slog.String("sender", sender.Conn.RemoteAddr().String()))

	for _, stream := range pool {
		if stream.Streamer {
			return
		}
	}

	sender.Streamer = true
}

// HandleOfferMessage handles the message type that sends an WebRTC SDP offer, sending
// to the spec users the offer provided.
func HandleOfferMessage(sender *Stream, pool []*Stream, message Message) {
	slog.Info("offer message received", slog.String("sender", sender.Conn.RemoteAddr().String()))
	offer := message.Data

	for _, stream := range pool {
		if stream.Streamer != true {
			err := stream.WriteJSON(map[string]interface{}{
				"type": "description",
				"data": offer,
			})

			if err != nil {
				slog.Error("error sending json offer", slog.Any("error", err), slog.String("sender", stream.RemoteAddr().String()))
			}
		}
	}
}

// HandleAnswerMessage handles the message type that sends an WebRTC SDP answer from a peer,
// sending the answer to the pool Streamer.
func HandleAnswerMessage(sender *Stream, pool []*Stream, message Message) {
	slog.Info("answer message received", slog.String("sender", sender.Conn.RemoteAddr().String()))
	answer := message.Data

	for _, stream := range pool {
		if stream.Streamer {
			err := stream.WriteJSON(map[string]interface{}{
				"type": "description",
				"data": answer,
			})

			if err != nil {
				slog.Error("error sending json answer", slog.Any("error", err), slog.String("sender", stream.RemoteAddr().String()))
			}

			return
		}
	}
}

// HandleSendOfferMessage handles the message type that asks for an offer from the streamer.
// It proxies the request towards the connection streamer, waits for its offer response message.
func HandleSendOfferMessage(sender *Stream, pool []*Stream, message Message) {
	slog.Info("offer request message received", slog.String("sender", sender.Conn.RemoteAddr().String()))

	for _, stream := range pool {
		if stream.Streamer {
			err := stream.WriteJSON(map[string]interface{}{
				"type": "createOffer",
				"data": nil,
			})

			if err != nil {
				slog.Error("sending offer request", slog.Any("error", err), slog.String("sender", sender.Conn.RemoteAddr().String()))
				return
			}
		}
	}
}

// HandleSendAnswerMessage handles the message type that asks for an answer from the other connection
// rather than the Streamer.
func HandleSendAnswerMessage(sender *Stream, pool []*Stream, message Message) {
	slog.Info("answer request message received", slog.String("sender", sender.Conn.RemoteAddr().String()))

	for _, stream := range pool {
		if !(stream.Streamer) {
			stream.WriteJSON(map[string]any{
				"type": "createAnswer",
				"data": nil,
			})
		}
	}
}

// HandleCandidateMessage handles the message type that sends a ICE candidate sending the candidate
// to others peers.
func HandleCandidateMessage(sender *Stream, pool []*Stream, message Message) {
	candidate := message.Data

	for _, stream := range pool {
		if stream != sender {
			err := stream.WriteJSON(map[string]interface{}{
				"type": "candidate",
				"data": candidate,
			})

			if err != nil {
				slog.Error("sending candidate request", slog.Any("error", err), slog.String("sender", sender.Conn.RemoteAddr().String()))
			}
		}
	}
}
