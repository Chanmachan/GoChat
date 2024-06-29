var conn;

function showRoomSelection() {
    const app = document.getElementById('app');
    app.innerHTML = `
        <div id="roomSelection">
            <h1>Enter Chat Room</h1>
            <input type="text" id="roomNumber" placeholder="Enter Room Number" autofocus>
            <button onclick="joinRoom()">Join Room</button>
            <button onclick="fetchUserInfo()">Fetch User Info</button>
        </div>
    `;
}

function showChat() {
    const app = document.getElementById('app');
    app.innerHTML = `
        <div id="chatRoom">
            <h1>WebSocket Chat</h1>
            <textarea id="chat" rows="20" cols="50" readonly></textarea><br>
            <input type="text" id="message" size="50">
            <button onclick="sendMessage()">Send</button>
            <button onclick="leaveRoom()">Leave Room</button>
        </div>
    `;
}

function joinRoom() {
    var roomNumber = document.getElementById('roomNumber').value;
    if (!/^\d+$/.test(roomNumber)) {
        alert("Please enter a valid room number. Room number should only contain digits.");
        return; // 数字以外が入力された場合、ここで処理を中断しユーザーに警告
    }
    conn = new WebSocket('ws://localhost:9090/ws?room=' + encodeURIComponent(roomNumber));
    console.log('ws://localhost:9090/ws?room=' + encodeURIComponent(roomNumber))
    conn.onopen = function(e) {
        console.log("Connection established to room " + roomNumber);
    };
    conn.onmessage = function(e) {
        var chat = document.getElementById('chat');
        chat.value += e.data + "\n";
    };
    navigate("/chat");
}

function fetchUserInfo() {
    fetch('/api/userinfo', {
        method: 'Get',
        credentials: 'include',
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch user info');
            }
            console.log("hogehogo");
            return response.json();
        })
        .then(data => {
            console.log("User Info:", data);
            alert("User Info: " + JSON.stringify(data));
        })
        .catch(error => {
            console.error('Error fetching user info:', error);
            alert('Error fetching user info: ' + error.message);
        });
}

function sendMessage() {
    var message = document.getElementById('message').value;
    console.log("Sending message: " + message);
    conn.send(JSON.stringify({username: "username", message: message}));
    document.getElementById('message').value = '';
}

function leaveRoom() {
    conn.close();
    navigate('/room'); // チャットルーム選択に戻る
}