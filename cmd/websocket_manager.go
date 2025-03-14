package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketManager struct {
	upgrader *websocket.Upgrader
	mu       sync.Mutex
	pool     []*Stream
	handlers map[string]MessageHandlerFunc
}

func NewWebsocketManager(capacity int) *WebsocketManager {
	return &WebsocketManager{
		pool: make([]*Stream, 0, capacity),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
		},
		handlers: make(map[string]MessageHandlerFunc),
	}
}

func (wm *WebsocketManager) AppendConnection(w http.ResponseWriter, r *http.Request) error {
	conn, err := wm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	stream := &Stream{
		Conn:     conn,
		Streamer: false,
	}

	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.pool = append(wm.pool, stream)

	go wm.handleConn(stream)

	return nil
}

func (wm *WebsocketManager) AddMessageHandler(_type string, h MessageHandlerFunc) error {
	if _, ok := wm.handlers[_type]; ok {
		return fmt.Errorf("already defined message [%s]", _type)
	}

	wm.handlers[_type] = h
	return nil
}

func (wm *WebsocketManager) handleConn(stream *Stream) {
	defer stream.Conn.Close()

	for {
		msgType, rawMessage, err := stream.Conn.ReadMessage()

		if err != nil {
			slog.Error("reading message from socket", slog.Any("error", err), slog.String("sender", stream.RemoteAddr().String()))
			break
		}

		if msgType != websocket.TextMessage {
			continue
		}

		message, err := wm.parseMessage(rawMessage)
		if err != nil {
			slog.Error("parsing message", slog.Any("error", err), slog.String("sender", stream.RemoteAddr().String()))
			continue
		}

		wm.handleMessage(stream, message)
	}

	// Removing connection from pool
	wm.deleteConnectionFromPool(stream)
}

func (wm *WebsocketManager) deleteConnectionFromPool(s *Stream) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	for index, stream := range wm.pool {
		if s == stream {
			wm.pool = append(wm.pool[:index], wm.pool[index+1:]...)
			break
		}
	}
}

func (wm *WebsocketManager) parseMessage(rawMessage []byte) (Message, error) {
	var msg Message

	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		return msg, err
	}

	return msg, nil
}

func (wm *WebsocketManager) handleMessage(sender *Stream, message Message) {
	handler, ok := wm.handlers[message.Type]

	if !ok {
		slog.Warn("not found message handler", slog.String("type", message.Type), slog.String("sender", sender.RemoteAddr().String()))
		return
	}

	handler(sender, wm.pool, message)
}
