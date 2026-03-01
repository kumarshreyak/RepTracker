import * as SecureStore from 'expo-secure-store';
import { Platform } from 'react-native';

// Token cache implementation for Clerk using expo-secure-store
export const tokenCache = {
  async getToken(key: string) {
    try {
      const item = await SecureStore.getItemAsync(key);
      if (item) {
        console.log(`${key} was used 🔐 \n`);
      } else {
        console.log('No values stored under key: ' + key);
      }
      return item;
    } catch (error) {
      console.error('SecureStore get item error: ', error);
      await SecureStore.deleteItemAsync(key);
      return null;
    }
  },
  async saveToken(key: string, value: string) {
    try {
      return SecureStore.setItemAsync(key, value);
    } catch (err) {
      console.error('SecureStore save item error: ', err);
      return;
    }
  },
};

/**
 * Clear all Clerk tokens from SecureStore
 * Use this when switching between Clerk instances (dev/production)
 * or when you need to force a fresh login
 */
export const clearClerkCache = async () => {
  try {
    // Clerk stores tokens with these keys
    const clerkKeys = [
      '__clerk_client_jwt',
      '__clerk_session_token',
      '__clerk_refresh_token',
      '__clerk_db_jwt',
    ];
    
    for (const key of clerkKeys) {
      try {
        await SecureStore.deleteItemAsync(key);
        console.log(`Cleared Clerk cache key: ${key}`);
      } catch (err) {
        console.log(`Key ${key} not found or already cleared`);
      }
    }
    
    console.log('✅ Clerk cache cleared successfully');
  } catch (error) {
    console.error('Error clearing Clerk cache:', error);
  }
};

