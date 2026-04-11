const video = document.getElementById("videoPlayer");
const videoList = document.querySelector(".video-list");

let isSelf = true;

function getUsername() {
    return new URLSearchParams(window.location.search).get("username") || "guest";
}

function appendCmdMessage(msg, isSelf) {
    const div = document.createElement("div");
    div.className = "chat-message control " + (isSelf ? "self" : "other");
    div.innerHTML = `<strong>${msg.sender}:</strong> <em>${msg.payload.command}</em>`;
    chatBox.appendChild(div);
    chatBox.scrollTop = chatBox.scrollHeight;
}

function sendControl(command, extra = {}) {
    const ws = window.websocket || window.socket;
    if (!ws || ws.readyState !== WebSocket.OPEN) return;

    const msg = {
        sender: getUsername(),
        receiver: "*",
        type: "control",
        payload: {
            command: command,
            ...extra
        }
    };
    appendCmdMessage(msg, isSelf);
    ws.send(JSON.stringify(msg));
    console.log("CONTROL SENT:", msg);
}

window.sendControl = sendControl;

export function handleControlMsg(msg) {
    isSelf = false;
    let command = msg.payload.command;
    console.log("Handling control command:", command, "isSelf:", isSelf);
    appendCmdMessage(msg, isSelf);
    switch (command) {
        case "play":
            if (video) video.play();
            break;

        case "pause":
            if (video) video.pause();
            break;
    }
    setTimeout(() => isSelf = true, 100);
}

if (video) {
    video.addEventListener("play", () => {
        console.log("VIDEO PLAYED");
        if (isSelf) sendControl("play");
    });

    video.addEventListener("pause", () => {
        console.log("VIDEO PAUSED");
        if (isSelf) sendControl("pause");
    });
}