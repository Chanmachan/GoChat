import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

const ChatRoom: React.FC = () => {
  const { roomNumber } = useParams();
  const [message, setMessage] = useState('');
  const [chat, setChat] = useState<string[]>([]);

  useEffect(() => {
    // WebSocketの接続処理など
  }, [roomNumber]);

  const sendMessage = () => {
    // メッセージ送信処理
    setChat([...chat, message]);
    setMessage('');
  };

  return (
    <div>
      <h1>WebSocket Chat in Room {roomNumber}</h1>
      <textarea value={chat.join('\n')} readOnly />
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
      />
      <button onClick={sendMessage}>Send</button>
    </div>
  );
};

export default ChatRoom;
