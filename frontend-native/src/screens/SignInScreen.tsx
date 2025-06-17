import React, { useState } from 'react';
import {
  View,
  StyleSheet,
  ActivityIndicator,
  Alert,
  SafeAreaView,
} from 'react-native';
import { Typography, Button } from '../components';
import { getColor } from '../components/Colors';
import { authService } from '../auth/AuthService';

// Google Icon SVG component
const GoogleIcon = () => (
  <View style={styles.googleIcon}>
    <View style={[styles.iconPath, { backgroundColor: getColor('blue') }]} />
    <View style={[styles.iconPath, { backgroundColor: getColor('red') }]} />
    <View style={[styles.iconPath, { backgroundColor: getColor('yellow') }]} />
    <View style={[styles.iconPath, { backgroundColor: getColor('green') }]} />
  </View>
);

interface SignInScreenProps {
  onSignInSuccess: () => void;
}

export const SignInScreen: React.FC<SignInScreenProps> = ({ onSignInSuccess }) => {
  const [isLoading, setIsLoading] = useState(false);

  const handleGoogleSignIn = async () => {
    try {
      setIsLoading(true);
      const success = await authService.signInWithGoogle();
      
      if (success) {
        onSignInSuccess();
      } else {
        Alert.alert('Sign In Failed', 'Unable to sign in with Google. Please try again.');
      }
    } catch (error) {
      console.error('Sign in error:', error);
      Alert.alert('Error', 'An unexpected error occurred. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.content}>
        <View style={styles.card}>
          <View style={styles.header}>
            <Typography variant="heading-large" color="dark" style={styles.title}>
              Welcome to GymLog
            </Typography>
            <Typography variant="text-default" color="light" style={styles.subtitle}>
              Track your workouts and progress
            </Typography>
          </View>

          <View style={styles.buttonContainer}>
            <Button
              variant="primary"
              size="large"
              style={styles.signInButton}
              onPress={handleGoogleSignIn}
              disabled={isLoading}
              icon={isLoading ? <ActivityIndicator size="small" color="white" /> : <GoogleIcon />}
            >
              {isLoading ? 'Signing in...' : 'Continue with Google'}
            </Button>
          </View>

          <View style={styles.footer}>
            <Typography variant="text-small" color="light" style={styles.termsText}>
              By signing in, you agree to our Terms of Service and Privacy Policy
            </Typography>
          </View>
        </View>
      </View>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('light-gray-1'),
  },
  content: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 16,
  },
  card: {
    backgroundColor: getColor('white'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('light-gray-3'),
    padding: 32,
    width: '100%',
    maxWidth: 384,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 1,
    },
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
  },
  header: {
    alignItems: 'center',
    marginBottom: 32,
  },
  title: {
    marginBottom: 8,
    textAlign: 'center',
  },
  subtitle: {
    textAlign: 'center',
    marginBottom: 24,
  },
  buttonContainer: {
    marginBottom: 24,
  },
  signInButton: {
    width: '100%',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 12,
  },
  googleIcon: {
    width: 20,
    height: 20,
    flexDirection: 'row',
    flexWrap: 'wrap',
  },
  iconPath: {
    width: 8,
    height: 8,
    borderRadius: 2,
    margin: 1,
  },
  footer: {
    alignItems: 'center',
  },
  termsText: {
    textAlign: 'center',
    lineHeight: 18,
  },
}); 