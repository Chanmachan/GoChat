import React from 'react';
import logo from './logo.svg';
import './App.css';
import {BrowserRouter, Route, Router, Routes} from 'react-router-dom';
import Login from './login/Login';
import ChatRoom from './chat/ChatRoom';
import RoomSelection from './room/RoomSelection';
import LoginState from './login/LoginState';
import {UserProvider} from './contexts/UserContexts';

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <UserProvider>
        <Routes>
          <Route path="/" element={<Login />} />
          <Route path="/room-selection" element={<RoomSelection />} />
          <Route path="/chat-room/:roomNumber" element={<ChatRoom />} />
          <Route path="/login-state" element={<LoginState />} />
        </Routes>
      </UserProvider>
    </BrowserRouter>
  );
};

export default App;
