import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import type {
  UserSession,
  UserStats,
  LoginUserRequest,
  CreateUserRequest,
} from '../types';
import {
  api,
  getStoredToken,
  removeStoredToken,
  parseJwt,
} from '../services/api';
import { useNotification } from './NotificationContext';

interface AuthContextType {
  user: UserSession | null;
  token: string | null;
  stats: UserStats | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isAuthModalOpen: boolean;
  authMode: 'login' | 'register';
  login: (data: LoginUserRequest) => Promise<void>;
  register: (data: CreateUserRequest) => Promise<void>;
  logout: () => void;
  deleteAccount: () => Promise<void>;
  refreshStats: () => Promise<void>;
  openAuthModal: (mode?: 'login' | 'register') => void;
  closeAuthModal: () => void;
  setAuthMode: (mode: 'login' | 'register') => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<UserSession | null>(null);
  const [token, setToken] = useState<string | null>(getStoredToken());
  const [stats, setStats] = useState<UserStats | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isAuthModalOpen, setIsAuthModalOpen] = useState<boolean>(false);
  const [authMode, setAuthMode] = useState<'login' | 'register'>('login');

  const { addToast } = useNotification();

  const refreshStats = useCallback(async () => {
    if (!getStoredToken()) return;
    try {
      const userStats = await api.getMyStats();
      setStats(userStats);
    } catch {
      // ignore
    }
  }, []);

  useEffect(() => {
    const rawToken = getStoredToken();
    if (rawToken) {
      const claims = parseJwt(rawToken);
      if (claims && claims.exp && claims.exp * 1000 > Date.now()) {
        const storedUser = localStorage.getItem('course_tracker_user_info');
        let sessionData: UserSession = {
          userId: claims.user_id || 0,
          token: rawToken,
          exp: claims.exp,
        };
        if (storedUser) {
          try {
            sessionData = { ...sessionData, ...JSON.parse(storedUser) };
          } catch {
            // ignore
          }
        }
        setUser(sessionData);
        setToken(rawToken);
        refreshStats();
      } else {
        removeStoredToken();
        setUser(null);
        setToken(null);
      }
    }
    setIsLoading(false);
  }, [refreshStats]);

  const login = async (data: LoginUserRequest) => {
    try {
      const session = await api.login(data);
      setUser(session);
      setToken(session.token);
      setIsAuthModalOpen(false);
      await refreshStats();
      addToast('Welcome back! ⚡', `Logged in as ${data.email} (+50 Daily Login XP awarded)`, 'reward', 50);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Login failed';
      throw new Error(msg);
    }
  };

  const register = async (data: CreateUserRequest) => {
    try {
      const res = await api.register(data);
      if (res.reactivated) {
        addToast(
          'Account Reactivated! 🎉',
          res.message || 'All your previous courses and progress have been restored.',
          'reward',
          50
        );
      } else {
        addToast('Account created! 🚀', res.message || 'Welcome to Course Flow.', 'success');
      }
      await login({ email: data.email, password: data.password });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Registration failed';
      throw new Error(msg);
    }
  };

  const logout = () => {
    removeStoredToken();
    setUser(null);
    setToken(null);
    setStats(null);
    addToast('Logged out', 'You have been signed out successfully.', 'info');
  };

  const deleteAccount = async () => {
    try {
      const res = await api.deleteAccount();
      setUser(null);
      setToken(null);
      setStats(null);
      addToast('Account Deleted', res.message || 'Your account has been deleted.', 'warning');
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Failed to delete account';
      throw new Error(msg);
    }
  };

  const openAuthModal = (mode: 'login' | 'register' = 'login') => {
    setAuthMode(mode);
    setIsAuthModalOpen(true);
  };

  const closeAuthModal = () => {
    setIsAuthModalOpen(false);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        stats,
        isAuthenticated: !!user && !!token,
        isLoading,
        isAuthModalOpen,
        authMode,
        login,
        register,
        logout,
        deleteAccount,
        refreshStats,
        openAuthModal,
        closeAuthModal,
        setAuthMode,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
