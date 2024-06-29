import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const RoomSelection: React.FC = () => {
  const [roomNumber, setRoomNumber] = useState('');
  const navigate = useNavigate();

  const joinRoom = () => {
    if (!/^\d+$/.test(roomNumber)) {
      alert('有効なルーム番号を入力してください。');
      return;
    }
    navigate(`/chat-room/${roomNumber}`);
  };

  return (
    <div>
      <input
        type="text"
        value={roomNumber}
        onChange={(e) => setRoomNumber(e.target.value)}
        placeholder="Enter Room Number"
        autoFocus
      />
      <button onClick={joinRoom}>Join Room</button>
    </div>
  );
};

export default RoomSelection;
