import React, { useState } from 'react';
import {
  View,
  StyleSheet,
  Alert,
  Pressable,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Link, useRouter } from 'expo-router';
import { useSignIn } from '@clerk/clerk-expo';
import { Typography, Button, Input } from '../../src/components';
import { getColor } from '../../src/components/Colors';

export default function SignInRoute() {
  const { signIn, setActive, isLoaded } = useSignIn();
  const router = useRouter();

  const [emailAddress, setEmailAddress] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [pendingVerification, setPendingVerification] = useState(false);
  const [code, setCode] = useState('');

  const onSignInPress = async () => {
    if (!isLoaded) return;

    if (!emailAddress || !password) {
      Alert.alert('Error', 'Please enter both email and password');
      return;
    }

    try {
      setIsLoading(true);
      
      const signInAttempt = await signIn.create({
        identifier: emailAddress,
        password,
      });

      if (signInAttempt.status === 'complete') {
        await setActive({ session: signInAttempt.createdSessionId });
        router.replace('/');
      } else if (signInAttempt.status === 'needs_second_factor') {
        // 2FA is required - prepare email code verification
        await signInAttempt.prepareSecondFactor({
          strategy: 'email_code',
        });
        setPendingVerification(true);
      } else {
        console.error('Sign in incomplete:', JSON.stringify(signInAttempt, null, 2));
        Alert.alert('Error', 'Sign in incomplete. Please try again.');
      }
    } catch (err: any) {
      console.error('Sign in error:', JSON.stringify(err, null, 2));
      Alert.alert('Sign In Failed', err.errors?.[0]?.message || 'Unable to sign in. Please check your credentials.');
    } finally {
      setIsLoading(false);
    }
  };

  const onVerifyPress = async () => {
    if (!isLoaded) return;

    if (!code) {
      Alert.alert('Error', 'Please enter the verification code');
      return;
    }

    try {
      setIsLoading(true);

      const signInAttempt = await signIn.attemptSecondFactor({
        strategy: 'email_code',
        code,
      });

      if (signInAttempt.status === 'complete') {
        await setActive({ session: signInAttempt.createdSessionId });
        router.replace('/');
      } else {
        console.error('Verification incomplete:', JSON.stringify(signInAttempt, null, 2));
        Alert.alert('Error', 'Verification incomplete. Please try again.');
      }
    } catch (err: any) {
      console.error('Verification error:', JSON.stringify(err, null, 2));
      Alert.alert('Verification Failed', err.errors?.[0]?.message || 'Invalid verification code. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.content}>
        <View style={styles.card}>
          <View style={styles.header}>
            <Typography variant="heading-large" color="contentPrimary" style={styles.title}>
              {pendingVerification ? 'Verify Email' : 'Welcome Back'}
            </Typography>
            <Typography variant="label-small" color="contentSecondary" style={styles.subtitle}>
              {pendingVerification
                ? 'Enter the verification code sent to your email'
                : 'Sign in to continue tracking your workouts'}
            </Typography>
          </View>

          {!pendingVerification ? (
            <View style={styles.form}>
              <View style={styles.inputGroup}>
                <Typography variant="label-small" color="contentPrimary" style={styles.label}>
                  Email
                </Typography>
                <Input
                  value={emailAddress}
                  onChangeText={setEmailAddress}
                  placeholder="Enter your email"
                  keyboardType="email-address"
                  autoCapitalize="none"
                  autoComplete="email"
                  editable={!isLoading}
                />
              </View>

              <View style={styles.inputGroup}>
                <Typography variant="label-small" color="contentPrimary" style={styles.label}>
                  Password
                </Typography>
                <Input
                  value={password}
                  onChangeText={setPassword}
                  placeholder="Enter your password"
                  secureTextEntry
                  autoComplete="password"
                  editable={!isLoading}
                />
              </View>

              <Button
                variant="primary"
                size="large"
                style={styles.signInButton}
                onPress={onSignInPress}
                disabled={isLoading || !isLoaded}
              >
                {isLoading ? 'Signing in...' : 'Sign In'}
              </Button>
            </View>
          ) : (
            <View style={styles.form}>
              <View style={styles.inputGroup}>
                <Typography variant="label-small" color="contentPrimary" style={styles.label}>
                  Verification Code
                </Typography>
                <Input
                  value={code}
                  onChangeText={setCode}
                  placeholder="Enter 6-digit code"
                  keyboardType="number-pad"
                  autoCapitalize="none"
                  editable={!isLoading}
                />
              </View>

              <Button
                variant="primary"
                size="large"
                style={styles.signInButton}
                onPress={onVerifyPress}
                disabled={isLoading || !isLoaded}
              >
                {isLoading ? 'Verifying...' : 'Verify'}
              </Button>

              <Button
                variant="text"
                size="default"
                style={styles.backButton}
                onPress={() => {
                  setPendingVerification(false);
                  setCode('');
                }}
                disabled={isLoading}
              >
                Back to Sign In
              </Button>
            </View>
          )}

          {!pendingVerification && (
            <View style={styles.footer}>
              <View style={styles.signUpPrompt}>
                <Typography variant="label-small" color="contentSecondary">
                  Don't have an account?{' '}
                </Typography>
                <Link href="/(auth)/sign-up" asChild>
                  <Pressable disabled={isLoading}>
                    <Typography variant="label-small" color="contentAccent">
                      Sign up
                    </Typography>
                  </Pressable>
                </Link>
              </View>

              <Typography variant="label-xsmall" color="contentTertiary" style={styles.termsText}>
                By signing in, you agree to our Terms of Service and Privacy Policy
              </Typography>
            </View>
          )}
        </View>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'),
  },
  content: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 16,
  },
  card: {
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
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
  },
  form: {
    marginBottom: 24,
  },
  inputGroup: {
    marginBottom: 16,
  },
  label: {
    marginBottom: 8,
  },
  signInButton: {
    width: '100%',
    marginTop: 8,
  },
  backButton: {
    width: '100%',
    marginTop: 8,
  },
  footer: {
    alignItems: 'center',
  },
  signUpPrompt: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
  },
  termsText: {
    textAlign: 'center',
    lineHeight: 16,
  },
}); 