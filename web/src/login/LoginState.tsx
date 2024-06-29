import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Box, CircularProgress, Text, Container, VStack } from '@chakra-ui/react';
import {useUser} from '../contexts/UserContexts';

const LoginState = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const queryParams = new URLSearchParams(location.search);
  const status = queryParams.get('status');
  const [countdown, setCountdown] = useState(2);
  const { setUserInfo } = useUser();

  useEffect(() => {
    const timer = setInterval(() => {
      setCountdown(prev => prev - 1);
    }, 1000);

    const timeout = setTimeout(() => {
      if (status === 'success') {
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
            console.log('Unauthorized: ', error);
          });
        navigate('/room-selection');
      } else {
        navigate('/login');
      }
    }, 2000);

    return () => {
      clearTimeout(timeout);
      clearInterval(timer);
    };
  }, [status, navigate]);

  return (
    <Container centerContent p={4}>
      <VStack spacing={4} align="center" justify="center" minH="100vh">
        <CircularProgress isIndeterminate color={status === 'success' ? 'green.300' : 'red.400'} thickness="12px" />
        <Box p={4} borderWidth="1px" borderRadius="lg" bg={status === 'success' ? 'green.100' : 'red.100'}>
          <Text fontSize="2xl" fontWeight="bold">{status === 'success' ? 'Authorization Successful!' : 'Authorization Failed'}</Text>
          <Text fontSize="md">You will be redirected in {countdown} seconds...</Text>
        </Box>
      </VStack>
    </Container>
  );
};

export default LoginState;
