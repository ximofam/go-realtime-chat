const params = new URLSearchParams(window.location.search);
const username = params.get("username") || "Guest";

const roomId = params.get("room_id");

document.getElementById("roomIdText").innerText = roomId;

/* ===== WEBSOCKET ===== */

let connected = false;

function goHome() {
    window.location.href = "/";
}

const protocol =
    location.protocol === "https:" ? "wss" : "ws";

const ws = new WebSocket(
    `${protocol}://${location.host}/rooms/${roomId}?username=${username}`
);


const connectTimeout = setTimeout(() => {
    if (!connected) {
        ws.close();
        goHome();
    }
}, 5000);

ws.onopen = () => {
    connected = true;
    clearTimeout(connectTimeout);
};

ws.onerror = () => {
    goHome();
};

ws.onclose = () => {
    if (!connected) {
        goHome();
    }
};

const input = document.getElementById("msg");

input.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
        event.preventDefault();
        send();
    }
});

const chat = document.getElementById("chat");

const usersBox = document.getElementById("users");

// lưu user hiện tại trong room
const users = new Map();


ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);

    // ===== CHAT =====
    if (msg.type === "chat") {
        renderChat(msg.data);
        return;
    }

    // ===== USER JOIN =====
    if (msg.type === "connect_user") {
        users.set(msg.data.id, msg.data);
        renderUsers();
        return;
    }

    // ===== USER LEAVE =====
    if (msg.type === "disconnect_user") {
        users.delete(msg.data.id);
        renderUsers();
        return;
    }
};


function renderChat(chatMsg) {
    const row = document.createElement("div");
    row.classList.add("message-row");

    if (chatMsg.from === username) {
        row.classList.add("self");
    }

    const bubble = document.createElement("div");
    bubble.className = "message-bubble";

    const sender = document.createElement("div");
    sender.className = "sender";
    sender.textContent = chatMsg.from;

    const content = document.createElement("div");
    content.textContent = chatMsg.content;

    const time = document.createElement("div");
    time.className = "time";

    const date = new Date(chatMsg.send_at * 1000);
    time.textContent = date.toLocaleTimeString();

    bubble.appendChild(sender);
    bubble.appendChild(content);
    bubble.appendChild(time);

    row.appendChild(bubble);
    chat.appendChild(row);

    chat.scrollTop = chat.scrollHeight;
}

function renderUsers() {
    usersBox.innerHTML = "";

    users.forEach((u) => {
        const div = document.createElement("div");
        div.className = "user-item";
        div.textContent = u.username;
        usersBox.appendChild(div);
    });
}


function send() {
    const input = document.getElementById("msg");
    if (!input.value) return;

    ws.send(input.value);
    input.value = "";
}

async function loadUsers() {
    const res = await fetch(`/rooms/${roomId}/users`);
    const data = await res.json();

    data.forEach(u => users.set(u.id, u));
    renderUsers();
}

loadUsers();

function copyRoomID() {
    navigator.clipboard.writeText(roomId);
    alert("Room ID copied!");
}

function leaveRoom() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.close(1000, "user leave");
    }

    if (connected) goHome()
}