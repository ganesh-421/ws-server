//
//  subscribe.js
//  
//
//  Created by Ganesh Adhikari on 23/04/2026.
//

function connectWebSocket(userId) {
  const ws = new WebSocket(`ws://localhost:8080/ws?user_id=${userId}`);

  ws.onopen = () => {
    console.log("Connected to WS");

    // Subscribe to private channel
    ws.send(JSON.stringify({
      action: "subscribe",
      channel: `user:${userId}`
    }));
  };

  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data);
      console.log("Received:", msg);
    } catch (e) {
      console.log("Raw message:", event.data);
    }
  };

  ws.onclose = () => {
    console.log("Disconnected");
  };

  ws.onerror = (err) => {
    console.error("WS Error:", err);
  };

  return ws;
}

connectWebsocket(1);
