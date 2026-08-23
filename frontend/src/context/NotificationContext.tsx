import React, { createContext, useContext, useState, useCallback } from 'react';
import type { ToastItem, ToastType } from '../types';
import { fireConfetti } from '../utils/confetti';

interface NotificationContextType {
  toasts: ToastItem[];
  addToast: (title: string, message?: string, type?: ToastType, points?: number) => void;
  removeToast: (id: string) => void;
  triggerCelebration: (type?: 'reward' | 'course_complete' | 'level_up' | 'simple') => void;
}

const NotificationContext = createContext<NotificationContextType | undefined>(undefined);

export const NotificationProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [toasts, setToasts] = useState<ToastItem[]>([]);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const addToast = useCallback(
    (title: string, message?: string, type: ToastType = 'info', points?: number) => {
      const id = Math.random().toString(36).substring(2, 9);
      const newToast: ToastItem = {
        id,
        type,
        title,
        message,
        points,
        timestamp: Date.now(),
      };

      setToasts((prev) => [...prev, newToast]);

      setTimeout(() => {
        removeToast(id);
      }, type === 'reward' ? 5000 : 4000);
    },
    [removeToast]
  );

  const triggerCelebration = useCallback(
    (_type?: 'reward' | 'course_complete' | 'level_up' | 'simple') => {
      // Celebration animations disabled
    },
    []
  );

  return (
    <NotificationContext.Provider
      value={{
        toasts,
        addToast,
        removeToast,
        triggerCelebration,
      }}
    >
      {children}
    </NotificationContext.Provider>
  );
};

export function useNotification(): NotificationContextType {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error('useNotification must be used within a NotificationProvider');
  }
  return context;
}
