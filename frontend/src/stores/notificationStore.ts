import { create } from 'zustand';

interface Notification {
  id: string;
  message: string;
  type: 'success' | 'error' | 'warning' | 'info';
  timestamp: number;
}

interface NotificationStore {
  notifications: Notification[];
  addNotification: (message: string, type: Notification['type']) => void;
  removeNotification: (id: string) => void;
  clearNotifications: () => void;
}

export const useNotificationStore = create<NotificationStore>((set) => ({
  notifications: [],

  addNotification: (message: string, type: Notification['type']) => {
    const id = Date.now().toString();
    const notification: Notification = {
      id,
      message,
      type,
      timestamp: Date.now(),
    };

    set((state) => ({
      notifications: [...state.notifications, notification],
    }));

    // 5秒后自动移除
    setTimeout(() => {
      set((state) => ({
        notifications: state.notifications.filter((n) => n.id !== id),
      }));
    }, 5000);
  },

  removeNotification: (id: string) => {
    set((state) => ({
      notifications: state.notifications.filter((n) => n.id !== id),
    }));
  },

  clearNotifications: () => {
    set({ notifications: [] });
  },
}));