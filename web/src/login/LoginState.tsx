// SuccessPage.js
import {useLocation, useNavigate} from 'react-router-dom';
import {useEffect} from 'react';

const LoginState = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const queryParams = new URLSearchParams(location.search);
  const status = queryParams.get('status');

  useEffect(() => {
    if (status === 'success') {
      setTimeout(() => {
        navigate('/room-selection');
      }, 2000);
    } else if (status === 'fail') {
      setTimeout(() => {
        navigate('/login');
      }, 2000);
    }
  }, [status, navigate]);

  return (
    <div>
      {status === 'success' ? (
        <h1>Login Successful!</h1>
      ) : (
        <h1>Login Failed</h1>
      )}
    </div>
  );
};

export default LoginState;
