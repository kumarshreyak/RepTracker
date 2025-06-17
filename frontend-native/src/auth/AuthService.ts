import { GoogleSignin, statusCodes } from '@react-native-google-signin/google-signin';
import AsyncStorage from '@react-native-async-storage/async-storage';

// Configure Google Sign-in
GoogleSignin.configure({
  webClientId: process.env.EXPO_PUBLIC_GOOGLE_WEB_CLIENT_ID || '', // From Google Cloud Console
  iosClientId: process.env.EXPO_PUBLIC_GOOGLE_IOS_CLIENT_ID || '', // From Google Cloud Console
  offlineAccess: true,
  hostedDomain: '',
  forceCodeForRefreshToken: true,
});

interface User {
  id: string;
  email: string;
  name: string;
  picture?: string;
  sessionToken: string;
}

interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
}

class AuthService {
  private listeners: ((state: AuthState) => void)[] = [];
  private state: AuthState = {
    user: null,
    isLoading: false,
    isAuthenticated: false,
  };

  constructor() {
    this.loadStoredAuth();
  }

  // Load stored authentication state
  private async loadStoredAuth() {
    try {
      const storedUser = await AsyncStorage.getItem('user');
      if (storedUser) {
        const user = JSON.parse(storedUser);
        this.setState({
          user,
          isLoading: false,
          isAuthenticated: true,
        });
      }
    } catch (error) {
      console.error('Error loading stored auth:', error);
    }
  }

  // Subscribe to auth state changes
  subscribe(listener: (state: AuthState) => void) {
    this.listeners.push(listener);
    // Immediately call with current state
    listener(this.state);
    
    // Return unsubscribe function
    return () => {
      this.listeners = this.listeners.filter(l => l !== listener);
    };
  }

  // Update state and notify listeners
  private setState(newState: Partial<AuthState>) {
    this.state = { ...this.state, ...newState };
    this.listeners.forEach(listener => listener(this.state));
  }

  // Get current auth state
  getState(): AuthState {
    return this.state;
  }

  // Sign in with Google
  async signInWithGoogle(): Promise<boolean> {
    try {
      this.setState({ isLoading: true });

      // Check if Google Play Services are available
      await GoogleSignin.hasPlayServices({ showPlayServicesUpdateDialog: true });

      // Get user info from Google
      const userInfo = await GoogleSignin.signIn();
      
      if (userInfo.data?.user) {
        const googleUser = userInfo.data.user;
        
        // Create user in backend (or use local storage for now)
        // For demo purposes, we'll create a user object
        const user: User = {
          id: googleUser.id,
          email: googleUser.email,
          name: googleUser.name || '',
          picture: googleUser.photo || undefined,
          sessionToken: `session_${Date.now()}`, // In real app, get from backend
        };

        // Store user data
        await AsyncStorage.setItem('user', JSON.stringify(user));
        
        this.setState({
          user,
          isLoading: false,
          isAuthenticated: true,
        });

        return true;
      } else {
        this.setState({ isLoading: false });
        return false;
      }
    } catch (error: any) {
      console.error('Google sign-in error:', error);
      
      if (error.code === statusCodes.SIGN_IN_CANCELLED) {
        console.log('User cancelled the sign-in flow');
      } else if (error.code === statusCodes.IN_PROGRESS) {
        console.log('Sign-in is in progress already');
      } else if (error.code === statusCodes.PLAY_SERVICES_NOT_AVAILABLE) {
        console.log('Play services not available');
      } else {
        console.log('Some other error:', error);
      }
      
      this.setState({ isLoading: false });
      return false;
    }
  }

  // Sign out
  async signOut(): Promise<void> {
    try {
      await GoogleSignin.revokeAccess();
      await GoogleSignin.signOut();
      await AsyncStorage.removeItem('user');
      this.setState({
        user: null,
        isLoading: false,
        isAuthenticated: false,
      });
    } catch (error) {
      console.error('Sign out error:', error);
    }
  }

  // Get current user
  getCurrentUser(): User | null {
    return this.state.user;
  }

  // Check if user is authenticated
  isAuthenticated(): boolean {
    return this.state.isAuthenticated;
  }
}

export const authService = new AuthService();
export type { User, AuthState }; 