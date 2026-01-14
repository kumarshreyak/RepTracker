import React, { useEffect, useState } from 'react';
import { View, StyleSheet, ActivityIndicator } from 'react-native';
import { Redirect } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { useAuth, useUser } from '@clerk/clerk-expo';
import { getColor } from '../src/components/Colors';
import { userService } from '../src/auth/UserService';

export default function Index() {
  const { isSignedIn, isLoaded, getToken } = useAuth();
  const { user } = useUser();
  const [isCheckingProfile, setIsCheckingProfile] = useState(false);
  const [profileComplete, setProfileComplete] = useState<boolean | null>(null);

  // Check user profile when signed in
  useEffect(() => {
    async function checkUserProfile() {
      if (!isSignedIn || !user) {
        return;
      }

      setIsCheckingProfile(true);
      try {
        const token = await getToken();
        if (!token) {
          console.error('No token available');
          setProfileComplete(false);
          return;
        }

        const userProfile = await userService.getUserProfile(user.id, token);
        
        if (!userProfile) {
          // User profile not found, needs onboarding
          setProfileComplete(false);
          return;
        }

        // Check if profile is complete (has height, weight, age, goal)
        const isComplete = userService.isProfileComplete(userProfile);
        setProfileComplete(isComplete);
      } catch (error) {
        console.error('Error checking user profile:', error);
        // On error, assume profile is incomplete and send to onboarding
        setProfileComplete(false);
      } finally {
        setIsCheckingProfile(false);
      }
    }

    checkUserProfile();
  }, [isSignedIn, user, getToken]);

  // Show loading screen while checking authentication or profile
  if (!isLoaded || (isSignedIn && isCheckingProfile)) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={getColor('accent')} />
        <StatusBar style="auto" />
      </View>
    );
  }

  // Redirect based on authentication state
  if (!isSignedIn) {
    return <Redirect href="/(auth)/sign-in" />;
  }

  // Redirect based on profile completion
  if (profileComplete === false) {
    return <Redirect href="/onboarding" />;
  }

  if (profileComplete === true) {
    return <Redirect href="/(tabs)" />;
  }

  // Still checking profile, show loading
  return (
    <View style={styles.loadingContainer}>
      <ActivityIndicator size="large" color={getColor('accent')} />
      <StatusBar style="auto" />
    </View>
  );
}

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    backgroundColor: getColor('backgroundPrimary'),
    alignItems: 'center',
    justifyContent: 'center',
  },
}); 