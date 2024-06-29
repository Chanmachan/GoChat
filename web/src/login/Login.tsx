import { useUser } from '../contexts/UserContexts';
import { useNavigate } from 'react-router-dom';
import {Box, Button, VStack, Image, Text, Container, Flex, useBreakpointValue} from '@chakra-ui/react';
import MyAppIcon from '../static/zenigame.png'

const Login = () => {
  // 分割代入でsetUserInfoを取り出して使えるようにする
  const { setUserInfo } = useUser();
  const navigate = useNavigate();

  const handleLogin = async () => {
    // 先にユーザー情報をチェック
    fetch('http://localhost:9090/api/login', {
      credentials: 'include'
    })
      .then(response => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error('Unauthorized');
        }
      })
      .then(data => {
        console.log('Already logged in, redirecting:', data);
        // ユーザー情報をcontextに保存
        setUserInfo(data);
        // 既にログインしていれば、部屋選択ページなどへリダイレクトする
        navigate('/room-selection');
      })
      .catch(error => {
        console.log('Need to authenticate:', error);
        // 認証されていない場合、OAuthを開始
        authenticate();
      });
  };

  const authenticate = () => {
    fetch('http://localhost:9090/auth/', {
      credentials: 'include'
    })
      .then(response => response.json())
      .then(data => {
        window.location.href = data.url;  // URLからのリダイレクト
      })
      .catch(error => {
        console.error('ログイン失敗:', error);
      });
  };

  const buttonSize = useBreakpointValue({ base: 'md', md: 'lg' });

  return (
    <Container centerContent p={8}>
      <Flex direction="column" align="center" justify="center" minH="100vh">
        <VStack spacing={8}>
          <Image src={MyAppIcon} boxSize="150px" alt="GoChat Logo" mb={4} />
          <Text fontSize="3xl" fontWeight="bold">Welcome to GoChat!</Text>
          <Text fontSize="2xl" fontWeight="bold">This is a chat application using WebSockets.</Text>
          <Box p={5} shadow="md" borderWidth="1px" borderRadius="md" w="100%">
            <Text fontSize="lg" fontWeight="semibold" mb={2}>&lt;How to Use&gt;</Text>
            <Text>1. Log in with your Google account.</Text>
            <Text>2. Enter a room number to open the chat screen.</Text>
            <Text>3. Let's send messages!</Text>
            <Text fontSize="lg" fontWeight="semibold" mb={2} mt={4}>&lt;Features&gt;</Text>
            <Text>・OAuth with Google authentication</Text>
            <Text>・Backend: Go</Text>
            <Text>・Frontend: React</Text>
          </Box>
          <Button colorScheme="teal" size={buttonSize} onClick={handleLogin}>
            Login
          </Button>
        </VStack>
      </Flex>
    </Container>
  );
};

export default Login;
