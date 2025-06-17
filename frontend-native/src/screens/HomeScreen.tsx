import React from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  Image,
} from 'react-native';
import { Typography, Button } from '../components';
import { getColor } from '../components/Colors';
import { authService, User } from '../auth/AuthService';

interface HomeScreenProps {
  user: User;
  onSignOut: () => void;
}

export const HomeScreen: React.FC<HomeScreenProps> = ({ user, onSignOut }) => {
  const handleSignOut = async () => {
    await authService.signOut();
    onSignOut();
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.content}>
        <View style={styles.header}>
          <View style={styles.userInfo}>
            {user.picture && (
              <Image source={{ uri: user.picture }} style={styles.avatar} />
            )}
            <Typography variant="heading-large" color="dark" style={styles.welcomeText}>
              Welcome, {user.name}!
            </Typography>
            <Typography variant="text-default" color="light" style={styles.emailText}>
              {user.email}
            </Typography>
          </View>
        </View>

        <View style={styles.mainContent}>
          <Typography variant="heading-default" color="dark" style={styles.sectionTitle}>
            Ready to start your workout journey?
          </Typography>
          <Typography variant="text-default" color="light" style={styles.description}>
            Track your exercises, monitor your progress, and achieve your fitness goals.
          </Typography>
          
          <View style={styles.actionButtons}>
            <Button
              variant="primary"
              size="large"
              style={styles.actionButton}
              onPress={() => {
                // TODO: Navigate to workout screen
                console.log('Start workout pressed');
              }}
            >
              Start Workout
            </Button>
            
            <Button
              variant="secondary"
              size="large"
              style={styles.actionButton}
              onPress={() => {
                // TODO: Navigate to progress screen
                console.log('View progress pressed');
              }}
            >
              View Progress
            </Button>
          </View>
        </View>

        <View style={styles.footer}>
          <Button
            variant="text"
            size="default"
            onPress={handleSignOut}
          >
            Sign Out
          </Button>
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
    paddingHorizontal: 16,
  },
  header: {
    paddingTop: 24,
    paddingBottom: 32,
    alignItems: 'center',
  },
  userInfo: {
    alignItems: 'center',
  },
  avatar: {
    width: 80,
    height: 80,
    borderRadius: 40,
    marginBottom: 16,
    backgroundColor: getColor('light-gray-3'),
  },
  welcomeText: {
    marginBottom: 8,
    textAlign: 'center',
  },
  emailText: {
    textAlign: 'center',
  },
  mainContent: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 24,
  },
  sectionTitle: {
    marginBottom: 12,
    textAlign: 'center',
  },
  description: {
    textAlign: 'center',
    marginBottom: 32,
    lineHeight: 22,
  },
  actionButtons: {
    width: '100%',
    gap: 16,
  },
  actionButton: {
    width: '100%',
  },
  footer: {
    paddingBottom: 24,
    alignItems: 'center',
  },
}); 