package client

import (
    "github.com/gorilla/websocket"
)

type Client struct {
    ID       string
    Conn     *websocket.Conn
    Send     chan []byte
    Channels map[string]bool
}

func NewClient(conn *websocket.Conn, id string) *Client {
    return &Client{
        ID:       id,
        Conn:     conn,
        Send:     make(chan []byte, 256),
        Channels: make(map[string]bool),
    }
}
