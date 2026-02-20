function showCreate() {
    menu.classList.add("hidden");
    createBox.classList.remove("hidden");
}

function showJoin() {
    menu.classList.add("hidden");
    joinBox.classList.remove("hidden");
}


function goBack() {
    createBox.classList.add("hidden");
    joinBox.classList.add("hidden");
    menu.classList.remove("hidden");
}

async function createRoom() {
    const username =
        document.getElementById("createUsername").value.trim();

    if (!username) return alert("Enter your name");

    const res = await fetch("/rooms", { method: "POST" });
    const data = await res.json();

    accessRoom(data.room_id, username)
}

async function joinRoom() {
    const username =
        document.getElementById("joinUsername").value.trim();

    let value =
        document.getElementById("roomInput").value.trim();

    if (!username) return alert("Enter your name");
    if (!value) return alert("Enter room id");

    try {
        const res = await fetch(`/rooms/${value}/exists`);

        if (!res.ok) return alert("Invalid room id")

        const data = await res.json();

        if (!data.exists) return alert("Room not found")

        accessRoom(value, username);

    } catch (e) {
        alert("Server unavailable")
    }
}


function accessRoom(id, username) {
    window.location.href =
        `/room.html?room_id=${id}&username=${encodeURIComponent(username)}`;
}