import { useAuth as useClerkAuth, useUser } from '@clerk/clerk-expo';

export interface AuthState {
  user: {
    id: string;
    email: string;
    name: string;
    picture?: string;
  } | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  getToken: () => Promise<string | null>;
}

export const useAuth = (): AuthState => {
  const { isSignedIn, isLoaded, getToken } = useClerkAuth();
  const { user } = useUser();

  return {
    user: user ? {
      id: user.id,
      email: user.primaryEmailAddress?.emailAddress || '',
      name: user.fullName || user.firstName || '',
      picture: user.imageUrl,
    } : null,
    isLoading: !isLoaded,
    isAuthenticated: isSignedIn || false,
    getToken: async () => {
      try {
        return await getToken();
      } catch (error) {
        console.error('Error getting Clerk token:', error);
        return null;
      }
    },
  };
}; 