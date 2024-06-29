import React, { useState, useEffect } from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import ReconnectingWebSocket from 'reconnecting-websocket';
import { useUser } from '../contexts/UserContexts';
import { Box, Button, VStack, Flex, Text, Input, useToast, Avatar, Container } from '@chakra-ui/react';

interface Message {
  username: string;
  message: string;
  picture?: string;
  timestamp: string;
}

const ChatRoom = () => {
  // パラメータからroom番号を取得ため
  const { roomNumber } = useParams<{ roomNumber?: string }>();
  // 入力するメッセージを保存するため
  const [message, setMessage] = useState('');
  // チャット履歴を配列として保存するため
  const [chat, setChat] = useState<Message[]>([]);
  // 再レンダリングされたときにでもソケットを保存しておける
  const connRef = React.useRef<ReconnectingWebSocket>()
  // ユーザー情報を取得
  const { userInfo } = useUser();
  const navigate = useNavigate();
  const toast = useToast();

  // コンポーネントがマウントされた(roomNumberが変更されるたび)後に実行される
  useEffect(() => {
    if (!roomNumber) {
      toast({
        title: 'Error',
        description: 'Room number is undefined.',
        status: 'error',
        duration: 9000,
        isClosable: true,
      });
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
      setChat(prev => [...prev, {
        username: data.username,
        message: data.message,
        picture: data.picture, // 画像があれば追加、なければ undefined
        timestamp: new Date().toISOString(),
      }]);
    };

    socket.onerror = (e) => {
      console.error('WebSocket error:', e);
    };

    return () => {
      socket.close();
    };
  }, [roomNumber, toast]);

  const sendMessage = () => {
    if (connRef.current && message && userInfo) {
      const payload = {
        username: userInfo.name,
        message,
        picture: userInfo.picture,
        timestamp: new Date().toISOString(),
      };
      connRef.current.send(JSON.stringify(payload));
      // メッセージボックスをリセット
      setMessage('');
    }
  };

  const leaveRoom = () => {
    if (window.confirm("Are you sure you want to leave the room?")) {
      navigate("/room-selection");
    }
  };

  if (!roomNumber) {
    return <div>Invalid room number.</div>;
  }

  return (
    <Container maxW="container.md" centerContent p={4}>
      <VStack spacing={4} align="stretch" w="100%" position="relative">
        <Box position="absolute" right="0">
          <Button colorScheme="red" size="sm" onClick={leaveRoom}>
            Leave Room
          </Button>
        </Box>
        <Text fontSize="2xl" fontWeight="bold">GoChat in Room {roomNumber}</Text>
        <VStack spacing={4} align="stretch" overflowY="auto" h="lg" p={4} borderWidth="1px" borderRadius="lg">
          {chat.map((msg, index) => (
            <Flex key={index} align="center">
              <Avatar size="sm" src={msg.picture} name={msg.username} mr={2} />
              <Box p={2} bg="blue.100" borderRadius="lg">
                <Text fontWeight="bold">{msg.username}</Text>
                <Text>{msg.message}</Text>
              </Box>
              <Text p={1} fontSize="xs" color="gray.500">
                {new Date(msg.timestamp).toLocaleDateString()} {/* 年月日 */}
                <br/>
                {new Date(msg.timestamp).toLocaleTimeString()} {/* 時間 */}
              </Text>
            </Flex>
          ))}
        </VStack>
        <Flex mt={2}>
          <Input placeholder="Type a message..." value={message} onChange={e => setMessage(e.target.value)} onKeyPress={e => e.key === 'Enter' && sendMessage()} />
          <Button ml={2} colorScheme="blue" onClick={sendMessage}>Send</Button>
        </Flex>
      </VStack>
    </Container>
  );
};

export default ChatRoom;
