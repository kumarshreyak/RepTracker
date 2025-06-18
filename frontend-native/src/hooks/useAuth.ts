import { useState, useEffect } from 'react';
import { authService, AuthState } from '../auth/AuthService';

export const useAuth = () => {
  const [authState, setAuthState] = useState<AuthState>(() => authService.getState());

  useEffect(() => {
    const unsubscribe = authService.subscribe(setAuthState);
    return unsubscribe;
  }, []);

  return authState;
}; 