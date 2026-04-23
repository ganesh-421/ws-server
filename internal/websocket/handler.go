package websocket

import (
    "encoding/json"
    "log"
    "net/http"
    "ws-server/internal/client"
    "ws-server/internal/hub"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
    Hub *hub.Hub
}

func (h *WSHandler) Handle(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }

    // Example: user ID from query (?user_id=42)
    userID := r.URL.Query().Get("user_id")
    c := client.NewClient(conn, userID)

    go h.writePump(c)
    h.readPump(c)
}

type WSMessage struct {
    Action  string `json:"action"`
    Channel string `json:"channel"`
}

func (h *WSHandler) readPump(c *client.Client) {
    defer func() {
        h.Hub.RemoveClient(c)
        c.Conn.Close()
    }()

    for {
        _, msg, err := c.Conn.ReadMessage()
        if err != nil {
            break
        }

        var m WSMessage
        if err := json.Unmarshal(msg, &m); err != nil {
            continue
        }

        switch m.Action {
        case "subscribe":
            if authorize(c, m.Channel) {
                h.Hub.Subscribe(c, m.Channel)
            }
        case "unsubscribe":
            h.Hub.Unsubscribe(c, m.Channel)
        }
    }
}

func (h *WSHandler) writePump(c *client.Client) {
    defer c.Conn.Close()

    for msg := range c.Send {
        err := c.Conn.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            return
        }
    }
}

func authorize(c *client.Client, channel string) bool {
    // Only allow user to subscribe to their own private channel
    if len(channel) > 5 && channel[:5] == "user:" {
        return channel == "user:"+c.ID
    }
    return true
}
