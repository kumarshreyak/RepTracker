import React, { useState } from 'react';
import {
  View,
  StyleSheet,
  Alert,
  Pressable,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Link, useRouter } from 'expo-router';
import { useSignUp } from '@clerk/clerk-expo';
import { Typography, Button, Input } from '../../src/components';
import { getColor } from '../../src/components/Colors';

export default function SignUpRoute() {
  const { isLoaded, signUp, setActive } = useSignUp();
  const router = useRouter();

  const [emailAddress, setEmailAddress] = useState('');
  const [password, setPassword] = useState('');
  const [pendingVerification, setPendingVerification] = useState(false);
  const [code, setCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const onSignUpPress = async () => {
    if (!isLoaded) return;

    if (!emailAddress || !password) {
      Alert.alert('Error', 'Please enter both email and password');
      return;
    }

    if (password.length < 8) {
      Alert.alert('Error', 'Password must be at least 8 characters long');
      return;
    }

    try {
      setIsLoading(true);

      await signUp.create({
        emailAddress,
        password,
      });

      await signUp.prepareEmailAddressVerification({ strategy: 'email_code' });

      setPendingVerification(true);
    } catch (err: any) {
      console.error('Sign up error:', JSON.stringify(err, null, 2));
      Alert.alert('Sign Up Failed', err.errors?.[0]?.message || 'Unable to create account. Please try again.');
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

      const signUpAttempt = await signUp.attemptEmailAddressVerification({
        code,
      });

      if (signUpAttempt.status === 'complete') {
        await setActive({ session: signUpAttempt.createdSessionId });
        router.replace('/');
      } else {
        console.error('Verification incomplete:', JSON.stringify(signUpAttempt, null, 2));
        Alert.alert('Error', 'Verification incomplete. Please try again.');
      }
    } catch (err: any) {
      console.error('Verification error:', JSON.stringify(err, null, 2));
      Alert.alert('Verification Failed', err.errors?.[0]?.message || 'Invalid verification code. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  if (pendingVerification) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.content}>
          <View style={styles.card}>
            <View style={styles.header}>
              <Typography variant="heading-large" color="contentPrimary" style={styles.title}>
                Verify Your Email
              </Typography>
              <Typography variant="label-small" color="contentSecondary" style={styles.subtitle}>
                We sent a verification code to {emailAddress}
              </Typography>
            </View>

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
                  maxLength={6}
                  editable={!isLoading}
                />
              </View>

              <Button
                variant="primary"
                size="large"
                style={styles.submitButton}
                onPress={onVerifyPress}
                disabled={isLoading || !isLoaded}
              >
                {isLoading ? 'Verifying...' : 'Verify Email'}
              </Button>
            </View>

            <View style={styles.footer}>
              <Typography variant="label-xsmall" color="contentTertiary" style={styles.termsText}>
                Didn't receive the code? Check your spam folder or try signing up again.
              </Typography>
            </View>
          </View>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.content}>
        <View style={styles.card}>
          <View style={styles.header}>
            <Typography variant="heading-large" color="contentPrimary" style={styles.title}>
              Create Account
            </Typography>
            <Typography variant="label-small" color="contentSecondary" style={styles.subtitle}>
              Start tracking your fitness journey
            </Typography>
          </View>

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
                placeholder="At least 8 characters"
                secureTextEntry
                autoComplete="password-new"
                editable={!isLoading}
              />
            </View>

            <Button
              variant="primary"
              size="large"
              style={styles.submitButton}
              onPress={onSignUpPress}
              disabled={isLoading || !isLoaded}
            >
              {isLoading ? 'Creating account...' : 'Sign Up'}
            </Button>
          </View>

          <View style={styles.footer}>
            <View style={styles.signInPrompt}>
              <Typography variant="label-small" color="contentSecondary">
                Already have an account?{' '}
              </Typography>
              <Link href="/(auth)/sign-in" asChild>
                <Pressable disabled={isLoading}>
                  <Typography variant="label-small" color="contentAccent">
                    Sign in
                  </Typography>
                </Pressable>
              </Link>
            </View>

            <Typography variant="label-xsmall" color="contentTertiary" style={styles.termsText}>
              By creating an account, you agree to our Terms of Service and Privacy Policy
            </Typography>
          </View>
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
  submitButton: {
    width: '100%',
    marginTop: 8,
  },
  footer: {
    alignItems: 'center',
  },
  signInPrompt: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
  },
  termsText: {
    textAlign: 'center',
    lineHeight: 16,
  },
});

