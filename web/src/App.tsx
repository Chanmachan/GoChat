import React from 'react';
import { ChakraProvider } from '@chakra-ui/react';
import './App.css';
import {BrowserRouter, Route, Router, Routes} from 'react-router-dom';
import Login from './login/Login';
import ChatRoom from './chat/ChatRoom';
import RoomSelection from './room/RoomSelection';
import LoginState from './login/LoginState';
import {UserProvider} from './contexts/UserContexts';
import ProtectedRoute from './login/ProtectedRoute';

const App: React.FC = () => {
  return (
    <ChakraProvider>
      <BrowserRouter>
        <UserProvider>
          <Routes>
            <Route path="/" element={<Login />} />
            <Route path="/room-selection" element={
              <ProtectedRoute>
                <RoomSelection />
              </ProtectedRoute>
            } />
            <Route path="/chat-room/:roomNumber" element={
              <ProtectedRoute>
                <ChatRoom />
              </ProtectedRoute>
            } />
            <Route path="/login-state" element={<LoginState />} />
          </Routes>
        </UserProvider>
      </BrowserRouter>
    </ChakraProvider>
  );
};

export default App;
