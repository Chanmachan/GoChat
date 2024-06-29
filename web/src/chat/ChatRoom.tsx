import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import ReconnectingWebSocket from 'reconnecting-websocket';

const ChatRoom = () => {
  // パラメータからroom番号を取得ため
  const { roomNumber } = useParams<{ roomNumber?: string }>();
  // 入力するメッセージを保存するため
  const [message, setMessage] = useState('');
  // チャット履歴を配列として保存するため
  const [chat, setChat] = useState<string[]>([]);
  // 再レンダリングされたときにでもソケットを保存しておける
  const connRef = React.useRef<ReconnectingWebSocket>()

  // コンポーネントがマウントされた(roomNumberが変更されるたび)後に実行される
  useEffect(() => {
    if (!roomNumber) {
      console.error('Room number is undefined.');
      return;
    }

    const socket = new ReconnectingWebSocket(
      `ws://localhost:9090/ws?room=${encodeURIComponent(roomNumber)}`
    );
    // ソケットを保持しておく
    connRef.current = socket;

    socket.onopen = () => {
      console.log("Connection established to room " + roomNumber);
    };

    socket.onmessage = (e) => {
      const data = JSON.parse(e.data);
      setChat(prev => [...prev, `${data.username}: ${data.message}`]);
    };

    socket.onerror = (e) => {
      console.error('WebSocket error:', e);
    };

    return () => {
      socket.close();
    };
  }, [roomNumber]);

  const sendMessage = () => {
    if (connRef.current && message) {
      connRef.current.send(JSON.stringify({ username: "username", message }));
      // メッセージボックスをリセット
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
