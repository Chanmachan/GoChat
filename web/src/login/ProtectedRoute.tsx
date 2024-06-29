import React, {ReactNode} from 'react';
import { Navigate } from 'react-router-dom';
import { useUser } from '../contexts/UserContexts';

interface ProtectedRouteProps {
  children: ReactNode;
}

const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const { userInfo } = useUser(); // UserContextからユーザー情報を取得

  if (!userInfo) {
    // ユーザー情報がない場合、ログインページにリダイレクト
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;  // ユーザー情報がある場合、子コンポーネントを表示
};

export default ProtectedRoute;
