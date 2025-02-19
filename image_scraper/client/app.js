const ws = new WebSocket("ws://localhost:8080/ws/client")

ws.onopen = function (){
    document.getElementById("ws-status").textContent = "WebSocket connected";
    document.getElementById("ws-status").classList.add("connected");
}

ws.onclose = function () {
    document.getElementById("ws-status").textContent = "WebSocket closed";
};

ws.onerror = function (error) {
    console.error("error WebSocket:", error);
};

ws.onmessage = function (event){
    let img = document.createElement("img");
    img.src = event.data;
    img.alt = event.data;
    let container = document.getElementById("image-container");
    container.appendChild(img);
}

function sendCommand(command) {
    const url = document.getElementById("url-input").value;
    const width = document.getElementById("width-input").value;
    const height = document.getElementById("height-input").value;

    if (!url) {
        alert("enter url before send command");
        return;
    }

    const message = {
        command: command,
        url: url,
        width: parseInt(width) || 0,
        height: parseInt(height) || 0
    }

    ws.send(JSON.stringify(message));
}