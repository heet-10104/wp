const chatBox = document.getElementById("chatBox");
const input = document.getElementById("chatInput");
const receiverSelect = document.getElementById("receiverSelect");
const sendButton = document.getElementById("sendButton");

function getUsername() {
    return new URLSearchParams(window.location.search).get("username") || "guest";
}

function sendMessage() {
    const ws = window.socket || window.websocket;

    if (!ws || ws.readyState !== WebSocket.OPEN) {
        console.log("WebSocket not ready");
        return;
    }

    const text = input.value.trim();
    if (!text) return;

    const receiver = receiverSelect.value;

    const msg = {
        sender: getUsername(),
        receiver: receiver,
        type: receiver === "*" ? "broadcast" : "personal",
        payload: {
            chatMessage: text
        }
    };
    console.log("SENDING:", msg);
    ws.send(JSON.stringify(msg));
    appendMessage(msg, true);
    input.value = "";
}

window.sendMessage = sendMessage;

if (sendButton) {
    sendButton.addEventListener("click", sendMessage);
}

export function appendMessage(msg, isSelf) {
    const div = document.createElement("div");
    div.className = "chat-message " + (isSelf ? "self" : "other");

    let label = msg.sender;
    if (msg.type === "personal") label += " (private)";

    div.innerHTML = `<strong>${label}:</strong> ${msg.payload.chatMessage}`;
appendMessage
    chatBox.appendChild(div);
    chatBox.scrollTop = chatBox.scrollHeight;
}

if (input) {
    input.addEventListener("keypress", function(e) {
        if (e.key === "Enter") {
            sendMessage();
        }
    });
}