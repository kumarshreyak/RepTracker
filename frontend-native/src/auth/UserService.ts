import { apiGet, apiPut } from '../utils/api';

export interface UserProfile {
  id: string;
  clerkId: string;
  email: string;
  name: string;
  firstName: string;
  lastName: string;
  height: number; // in cm
  weight: number; // in kg
  age: number;
  goal: string;
  picture?: string;
  createdAt: string;
  updatedAt: string;
}

export interface UpdateUserData {
  email?: string;
  firstName?: string;
  lastName?: string;
  height: number;
  weight: number;
  age: number;
  goal: string;
  picture?: string;
}

class UserService {
  /**
   * Check if user profile is complete (has height, weight, age, and goal)
   */
  isProfileComplete(user: UserProfile): boolean {
    return !!(user.height && user.weight && user.age && user.goal);
  }

  /**
   * Get user profile for the authenticated user
   * Note: The backend will create a user record if it doesn't exist
   * @param clerkUserId - Clerk user ID (kept for backwards compatibility but not used in API call)
   * @param sessionToken - Clerk session token
   */
  async getUserProfile(clerkUserId: string, sessionToken: string): Promise<UserProfile | null> {
    try {
      // The backend extracts the user ID from the JWT token
      // No need to pass userId in the URL - it's more secure this way
      const userProfile = await apiGet<UserProfile>(
        `/api/users`,
        sessionToken
      );
      console.log('User profile fetched:', userProfile);
      return userProfile;
    } catch (error) {
      console.error('Error fetching user profile:', error);
      return null;
    }
  }

  /**
   * Create or update user profile with onboarding data (upsert)
   * @param clerkUserId - Clerk user ID (kept for backwards compatibility but not used in API call)
   * @param userData - Profile data to update
   * @param sessionToken - Clerk session token
   */
  async createOrUpdateUserProfile(
    clerkUserId: string,
    userData: UpdateUserData,
    sessionToken: string
  ): Promise<boolean> {
    try {
      // The backend extracts the user ID from the JWT token
      // No need to pass userId in the URL - it's more secure this way
      const updatedProfile = await apiPut<UserProfile>(
        `/api/users`,
        userData,
        sessionToken
      );
      console.log('User profile created/updated successfully:', updatedProfile);
      return true;
    } catch (error) {
      console.error('Error creating/updating user profile:', error);
      return false;
    }
  }

  /**
   * @deprecated Use createOrUpdateUserProfile instead
   * Update user profile with onboarding data
   * @param clerkUserId - Clerk user ID
   * @param userData - Profile data to update
   * @param sessionToken - Clerk session token
   */
  async updateUserProfile(
    clerkUserId: string,
    userData: UpdateUserData,
    sessionToken: string
  ): Promise<boolean> {
    return this.createOrUpdateUserProfile(clerkUserId, userData, sessionToken);
  }
}

export const userService = new UserService(); 