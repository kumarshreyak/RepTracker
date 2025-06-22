import { authService } from './AuthService';

export interface UserProfile {
  id: string;
  email: string;
  name: string;
  firstName: string;
  lastName: string;
  height: number; // in cm
  weight: number; // in kg
  age: number;
  goal: string;
  googleId: string;
  picture?: string;
  createdAt: string;
  updatedAt: string;
}

export interface UpdateUserData {
  height: number;
  weight: number;
  age: number;
  goal: string;
}

class UserService {
  private apiBaseUrl: string;

  constructor() {
    this.apiBaseUrl = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
  }

  /**
   * Check if user profile is complete (has height, weight, age, and goal)
   */
  isProfileComplete(user: UserProfile): boolean {
    return !!(user.height && user.weight && user.age && user.goal);
  }

  /**
   * Fetch current user profile from backend
   */
  async getCurrentUserProfile(): Promise<UserProfile | null> {
    try {
      const currentUser = authService.getCurrentUser();
      if (!currentUser) {
        console.log('No authenticated user found');
        return null;
      }

      const response = await fetch(`${this.apiBaseUrl}/api/auth/validate`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${currentUser.sessionToken}`,
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        console.error('Failed to fetch user profile:', response.status, response.statusText);
        return null;
      }

      const userProfile: UserProfile = await response.json();
      console.log('User profile fetched:', userProfile);
      
      return userProfile;
    } catch (error) {
      console.error('Error fetching user profile:', error);
      return null;
    }
  }

  /**
   * Update user profile with onboarding data
   */
  async updateUserProfile(userId: string, userData: UpdateUserData): Promise<boolean> {
    try {
      const currentUser = authService.getCurrentUser();
      if (!currentUser) {
        throw new Error('User not authenticated');
      }

      const response = await fetch(`${this.apiBaseUrl}/api/users/${userId}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${currentUser.sessionToken}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Failed to update user profile:', response.status, response.statusText, errorText);
        return false;
      }

      const updatedProfile: UserProfile = await response.json();
      console.log('User profile updated successfully:', updatedProfile);
      
      return true;
    } catch (error) {
      console.error('Error updating user profile:', error);
      return false;
    }
  }
}

export const userService = new UserService(); 