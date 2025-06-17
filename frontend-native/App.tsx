import React from 'react';
import { StatusBar } from 'expo-status-bar';
import { ActivityIndicator, View, StyleSheet } from 'react-native';
import { SignInScreen, HomeScreen } from './src/screens';
import { useAuth } from './src/hooks/useAuth';
import { getColor } from './src/components/Colors';

export default function App() {
  const { user, isLoading, isAuthenticated } = useAuth();

  // Show loading screen while checking authentication
  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={getColor('blue-bright')} />
        <StatusBar style="auto" />
      </View>
    );
  }

  // Show appropriate screen based on authentication state
  return (
    <>
      {isAuthenticated && user ? (
        <HomeScreen
          user={user}
          onSignOut={() => {
            // Auth state will be updated automatically by the auth service
          }}
        />
      ) : (
        <SignInScreen
          onSignInSuccess={() => {
            // Auth state will be updated automatically by the auth service
          }}
        />
      )}
      <StatusBar style="auto" />
    </>
  );
}

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    backgroundColor: getColor('light-gray-1'),
    alignItems: 'center',
    justifyContent: 'center',
  },
});
