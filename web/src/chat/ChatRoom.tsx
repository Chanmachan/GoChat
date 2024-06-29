import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import ReconnectingWebSocket from 'reconnecting-websocket';

const ChatRoom = () => {
  const { roomNumber } = useParams<{ roomNumber?: string }>();
  const [message, setMessage] = useState('');
  const [chat, setChat] = useState<string[]>([]);
  const connRef = React.useRef<ReconnectingWebSocket>()

  useEffect(() => {
    if (!roomNumber) {
      console.error('Room number is undefined.');
      return;
    }

    const newConn = new ReconnectingWebSocket(
      `ws://localhost:9090/ws?room=${encodeURIComponent(roomNumber)}`
    );
    connRef.current = newConn;

    newConn.onopen = () => {
      console.log("Connection established to room " + roomNumber);
    };

    newConn.onmessage = (e) => {
      const data = JSON.parse(e.data);
      setChat(prev => [...prev, `${data.username}: ${data.message}`]);
    };

    newConn.onerror = (e) => {
      console.error('WebSocket error:', e);
    };

    return () => {
      newConn.close();
    };
  }, [roomNumber]);

  const sendMessage = () => {
    if (connRef.current && message) {
      connRef.current.send(JSON.stringify({ username: "username", message }));
      setMessage('');
    }
  };

  if (!roomNumber) {
    return <div>Invalid room number.</div>;
  }

  return (
    <div>
      <h1>WebSocket Chat in Room {roomNumber}</h1>
      <textarea value={chat.join('\n')} readOnly />
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
      />
      <button onClick={sendMessage}>Send</button>
    </div>
  );
};

export default ChatRoom;
