import React from 'react';

const Login = () => {
  const handleLogin = async () => {
    fetch('http://localhost:9090/auth/', {
      credentials: 'include'
    })
      .then(response => response.json())
      .then(data => {
        console.log(data.url);
        window.location.href = data.url;  // URLからのリダイレクト
      })
      .catch(error => {
        console.error('ログイン失敗:', error);
      });
  };

  return (
    <button onClick={handleLogin}>ログイン</button>
  );
};

export default Login;