package hub

import (
    "sync"
    "ws-server/internal/client"
)

type Hub struct {
    Channels map[string]map[*client.Client]bool
    Mutex    sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        Channels: make(map[string]map[*client.Client]bool),
    }
}

func (h *Hub) Subscribe(c *client.Client, channel string) {
    h.Mutex.Lock()
    defer h.Mutex.Unlock()

    if h.Channels[channel] == nil {
        h.Channels[channel] = make(map[*client.Client]bool)
    }

    h.Channels[channel][c] = true
    c.Channels[channel] = true
}

func (h *Hub) Unsubscribe(c *client.Client, channel string) {
    h.Mutex.Lock()
    defer h.Mutex.Unlock()

    if clients, ok := h.Channels[channel]; ok {
        delete(clients, c)
        delete(c.Channels, channel)

        if len(clients) == 0 {
            delete(h.Channels, channel)
        }
    }
}

func (h *Hub) RemoveClient(c *client.Client) {
    h.Mutex.Lock()
    defer h.Mutex.Unlock()

    for ch := range c.Channels {
        delete(h.Channels[ch], c)
        if len(h.Channels[ch]) == 0 {
            delete(h.Channels, ch)
        }
    }
}

func (h *Hub) BroadcastToChannel(channel string, message []byte) {
    h.Mutex.RLock()
    clients := h.Channels[channel]
    h.Mutex.RUnlock()

    for c := range clients {
        select {
        case c.Send <- message:
        default:
            close(c.Send)
            h.RemoveClient(c)
        }
    }
}
