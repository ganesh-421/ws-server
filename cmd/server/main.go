package main

import (
    "log"
    "net/http"
    "ws-server/internal/hub"
    "ws-server/internal/redis"
    "ws-server/internal/websocket"
)

func main() {
    h := hub.NewHub()

    go redis.StartSubscriber(h)

    wsHandler := &websocket.WSHandler{Hub: h}

    http.HandleFunc("/ws", wsHandler.Handle)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
