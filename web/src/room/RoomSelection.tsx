import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Input, VStack, Container, Heading } from '@chakra-ui/react';

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
    <Container centerContent p={8}>
      <VStack spacing={4}>
        <Heading size="lg">Enter a Room</Heading>
        <Input
          placeholder="Enter Room Number"
          value={roomNumber}
          onChange={(e) => setRoomNumber(e.target.value)}
          size="lg"
          autoFocus
        />
        <Button colorScheme="teal" onClick={joinRoom} size="lg">
          Join Room
        </Button>
      </VStack>
    </Container>

  );
};

export default RoomSelection;
