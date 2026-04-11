import { appendMessage } from "./chat_handler.js";
import { handleControlMsg } from "./cmd_handler.js";

const waitForWS = setInterval(() => {
    if (window.websocket) {
        clearInterval(waitForWS);

        window.websocket.addEventListener("message", (e) => {
            const msg = JSON.parse(e.data);
            console.log("RECEIVED:", msg);
            msgMux(msg);
            
        });
    }
}, 100);

function msgMux(msg) {
    switch (msg.type) {
        case "broadcast":
            appendMessage(msg, false);
            break;

        case "personal":
            appendMessage(msg, false);
            break;

        case "control":
            handleControlMsg(msg);

        default:
            console.warn("Unknown message type:", msg.type);
    }
}