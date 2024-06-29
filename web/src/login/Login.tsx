import React from 'react';

const Login = () => {
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
        // 既にログインしていれば、部屋選択ページなどへリダイレクトする
        window.location.href = '/room-selection';
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

  return (
    <button onClick={handleLogin}>ログイン</button>
  );
};

export default Login;
