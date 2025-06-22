import { useState, useEffect } from 'react';
import { authService, AuthState } from '../auth/AuthService';
import { userService, UserProfile } from '../auth/UserService';

export interface AuthProfileState {
  isLoading: boolean;
  isAuthenticated: boolean;
  needsOnboarding: boolean;
  userProfile: UserProfile | null;
}

export const useAuthAndProfile = () => {
  const [state, setState] = useState<AuthProfileState>({
    isLoading: true,
    isAuthenticated: false,
    needsOnboarding: false,
    userProfile: null,
  });

  useEffect(() => {
    let mounted = true;

    const checkAuthAndProfile = async (authState: AuthState) => {
      if (!mounted) return;

      // If auth is still loading, keep our loading state
      if (authState.isLoading) {
        setState(prev => ({ ...prev, isLoading: true }));
        return;
      }

      // If not authenticated, we're done
      if (!authState.isAuthenticated) {
        setState({
          isLoading: false,
          isAuthenticated: false,
          needsOnboarding: false,
          userProfile: null,
        });
        return;
      }

      // User is authenticated, now check profile completion
      try {
        setState(prev => ({ ...prev, isLoading: true }));
        
        const userProfile = await userService.getCurrentUserProfile();
        
        if (!mounted) return;

        if (!userProfile) {
          // Failed to fetch profile, treat as needs onboarding
          setState({
            isLoading: false,
            isAuthenticated: true,
            needsOnboarding: true,
            userProfile: null,
          });
          return;
        }

        const needsOnboarding = !userService.isProfileComplete(userProfile);
        
        setState({
          isLoading: false,
          isAuthenticated: true,
          needsOnboarding,
          userProfile,
        });
      } catch (error) {
        console.error('Error checking user profile:', error);
        if (!mounted) return;
        
        // On error, assume needs onboarding to be safe
        setState({
          isLoading: false,
          isAuthenticated: true,
          needsOnboarding: true,
          userProfile: null,
        });
      }
    };

    // Subscribe to auth state changes
    const unsubscribe = authService.subscribe((authState) => {
      checkAuthAndProfile(authState);
    });

    return () => {
      mounted = false;
      unsubscribe();
    };
  }, []);

  return state;
}; 