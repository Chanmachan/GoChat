var conn;

function joinRoom() {
    var roomNumber = document.getElementById('roomNumber').value;
    document.getElementById('roomSelection').style.display = 'none';
    document.getElementById('chatRoom').style.display = 'block';
    conn = new WebSocket('ws://localhost:9090/ws?room=' + encodeURIComponent(roomNumber));
    console.log('ws://localhost:9090/ws?room=' + encodeURIComponent(roomNumber))
    conn.onopen = function(e) {
        console.log("Connection established to room " + roomNumber);
    };
    conn.onmessage = function(e) {
        var chat = document.getElementById('chat');
        chat.value += e.data + "\n";
    };
}

function sendMessage() {
    var message = document.getElementById('message').value;
    console.log("Sending message: " + message);
    conn.send(JSON.stringify({username: "username", message: message}));
    document.getElementById('message').value = '';
}

function leaveRoom() {
    conn.close();
    document.getElementById('chat').value = '';
    document.getElementById('roomSelection').style.display = 'block';
    document.getElementById('chatRoom').style.display = 'none';
}