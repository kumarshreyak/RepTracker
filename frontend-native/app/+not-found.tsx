import React from 'react';
import { View, StyleSheet, SafeAreaView } from 'react-native';
import { Link, Stack } from 'expo-router';
import { Typography, Button } from '../src/components';
import { getColor } from '../src/components/Colors';

export default function NotFoundScreen() {
  return (
    <SafeAreaView style={styles.container}>
      <Stack.Screen options={{ title: 'Oops! Not Found' }} />
      <View style={styles.content}>
              <Typography variant="heading-large" color="contentPrimary" style={styles.title}>
        Page Not Found
      </Typography>
      <Typography variant="paragraph-medium" color="contentSecondary" style={styles.description}>
        The page you're looking for doesn't exist.
      </Typography>
        <Link href="/(tabs)" asChild>
          <Button variant="primary" size="large" style={styles.button}>
            Go back to Home
          </Button>
        </Link>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundPrimary'),
  },
  content: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 16,
  },
  title: {
    marginBottom: 16,
    textAlign: 'center',
  },
  description: {
    marginBottom: 32,
    textAlign: 'center',
  },
  button: {
    minWidth: 200,
  },
}); 