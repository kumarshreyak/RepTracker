import React, { useEffect } from 'react';
import { View, StyleSheet, ActivityIndicator } from 'react-native';
import { router, Redirect } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { useAuthAndProfile } from '../src/hooks/useAuthAndProfile';
import { getColor } from '../src/components/Colors';

export default function Index() {
  const { isLoading, isAuthenticated, needsOnboarding } = useAuthAndProfile();

  // Show loading screen while checking authentication and profile
  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={getColor('accent')} />
        <StatusBar style="auto" />
      </View>
    );
  }

  // Redirect based on authentication and profile state
  if (!isAuthenticated) {
    return <Redirect href="/(auth)/sign-in" />;
  } else if (needsOnboarding) {
    return <Redirect href="/onboarding" />;
  } else {
    return <Redirect href="/(tabs)" />;
  }
}

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    backgroundColor: getColor('backgroundPrimary'),
    alignItems: 'center',
    justifyContent: 'center',
  },
}); 