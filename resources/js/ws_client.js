let websocket = null;

function initializeWebSocketListeners(ws) {
  ws.addEventListener("open", () => {
    console.log("CONNECTED");
  });

  ws.addEventListener("close", () => {
    console.log("DISCONNECTED");
  });

  ws.addEventListener("message", (e) => {
    console.log("RECEIVED:", e.data);
  });

  ws.addEventListener("error", (e) => {
    console.error("WebSocket ERROR:", e);
  });
}

console.log("OPENING");

// Extract username and room from URL parameters
const urlParams = new URLSearchParams(window.location.search);
const username = urlParams.get("username") || "guest";
const room = urlParams.get("room") || "general";

const wsUri = `ws://localhost:8080/ws?username=${encodeURIComponent(username)}&room=${encodeURIComponent(room)}`;
console.log(`Connecting to: ${wsUri}`);
websocket = new WebSocket(wsUri);
window.websocket = websocket;
window.socket = websocket;
window.chatUsername = username;
initializeWebSocketListeners(websocket);