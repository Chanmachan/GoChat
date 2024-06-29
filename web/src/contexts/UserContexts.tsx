// UserContext.js
import React, { createContext, useContext, useState, ReactNode } from 'react';

export type UserInfo = {
  id: string;
  email: string;
  verifiedEmail: boolean;
  name: string;
  givenName: string;
  familyName: string;
  picture: string;
};

interface UserContextType {
  userInfo: UserInfo | null;
  setUserInfo: (userInfo: UserInfo) => void;
}

// contextオブジェクトを生成
const UserContext = createContext<UserContextType | undefined>(undefined);

// 生成したオブジェクトのProvider定義
export const UserProvider = ({ children }: { children: ReactNode }) => {
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null);

  return (
    <UserContext.Provider value={{ userInfo, setUserInfo }}>
      {children}
    </UserContext.Provider>
  );
};

// contextを安全に取り出すためのカスタムフック
// contextがundefinedのとき(UserProviderが提供されていないとき)エラーをスローする
export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};
