import { create } from 'zustand';

interface ToastState {
  open: boolean;
  message: string;
}

interface UIState {
  toast: ToastState;
  openToast: (message: string) => void;
  closeToast: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  toast: { open: false, message: '' },
  openToast: (message) => set({ toast: { open: true, message } }),
  closeToast: () => set({ toast: { open: false, message: '' } }),
}));
