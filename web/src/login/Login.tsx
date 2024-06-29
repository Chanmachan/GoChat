import { useUser } from '../contexts/UserContexts';
import { useNavigate } from 'react-router-dom';

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

  return (
    <button onClick={handleLogin}>ログイン</button>
  );
};

export default Login;
