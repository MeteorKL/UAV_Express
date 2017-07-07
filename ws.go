package main

import (
	"github.com/gorilla/websocket"
)

type User struct {
	name   string
	avatar string
	ws     *websocket.Conn
	msg    chan map[string]interface{}
	state  string
}
