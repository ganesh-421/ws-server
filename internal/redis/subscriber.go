package redis

import (
    "context"
    "encoding/json"
    "log"
    "ws-server/internal/hub"

    "github.com/redis/go-redis/v9"
)

type Payload struct {
    Channel string          `json:"channel"`
    Event   string          `json:"event"`
    Data    json.RawMessage `json:"data"`
}

func StartSubscriber(h *hub.Hub) {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    sub := rdb.Subscribe(context.Background(), "ws_broadcast")
    ch := sub.Channel()

    log.Println("Redis subscriber started...")

    for msg := range ch {
        var payload Payload

        if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
            continue
        }

        h.BroadcastToChannel(payload.Channel, []byte(msg.Payload))
    }
}
