package main

import (
	"github.com/gorilla/websocket"
)

var pool = []*Stream{}

type Stream struct {
	*websocket.Conn
	Streamer bool
}
