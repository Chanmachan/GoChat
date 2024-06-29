import React, { useState, useEffect } from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import ReconnectingWebSocket from 'reconnecting-websocket';
import { useUser } from '../contexts/UserContexts';

const ChatRoom = () => {
  // パラメータからroom番号を取得ため
  const { roomNumber } = useParams<{ roomNumber?: string }>();
  // 入力するメッセージを保存するため
  const [message, setMessage] = useState('');
  // チャット履歴を配列として保存するため
  const [chat, setChat] = useState<string[]>([]);
  // 再レンダリングされたときにでもソケットを保存しておける
  const connRef = React.useRef<ReconnectingWebSocket>()
  // ユーザー情報を取得
  const { userInfo } = useUser();
  const navigate = useNavigate();

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
      {userInfo && (
        <div>
          <img src={userInfo.picture} alt="User" style={{ width: 50, height: 50 }} />
          <h2>{userInfo.name}</h2>
        </div>
      )}
      <textarea value={chat.join('\n')} readOnly />
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
      />
      <button onClick={sendMessage}>Send</button>
      <button onClick={() => {
        navigate("/room-selection");
      }} >leave room</button>
    </div>
  );
};

export default ChatRoom;
