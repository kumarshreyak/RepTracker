import React from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  Image,
} from 'react-native';
import { router } from 'expo-router';
import { Typography, Button } from '../../src/components';
import { getColor } from '../../src/components/Colors';
import { useAuth } from '../../src/hooks/useAuth';
import { authService } from '../../src/auth/AuthService';

export default function ProfileTab() {
  const { user } = useAuth();

  const handleSignOut = () => {
    router.replace('/sign-in');
  };

  if (!user) {
    router.replace('/');
    return null;
  }

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        {/* Header */}
        <View style={styles.header}>
          <Typography variant="heading-large" color="text-primary" style={styles.title}>
            Profile
          </Typography>
        
        </View>

        {/* User Profile Card */}
        <View style={styles.section}>
          <View style={styles.profileHeader}>
            {user.picture && (
              <Image source={{ uri: user.picture }} style={styles.avatar} />
            )}
            <View style={styles.userInfo}>
              <Typography variant="heading-medium" color="text-primary" style={styles.userName}>
                {user.name}
              </Typography>
              <Typography variant="paragraph-medium" color="text-secondary">
                {user.email}
              </Typography>
            </View>
          </View>
        </View>

        {/* Sign Out */}
        <View style={styles.signOutSection}>
          <Button
            variant="danger"
            size="large"
            style={styles.signOutButton}
            onPress={handleSignOut}
          >
            Sign Out
          </Button>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'),
  },
  scrollView: {
    flex: 1,
  },
  header: {
    paddingHorizontal: 16,
    paddingTop: 16,
    paddingBottom: 24,
  },
  title: {
    marginBottom: 8,
  },
  section: {
    paddingHorizontal: 16,
    marginBottom: 24,
  },
  profileHeader: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  avatar: {
    width: 80,
    height: 80,
    borderRadius: 40,
    marginRight: 16,
    backgroundColor: getColor('backgroundTertiary'),
  },
  userInfo: {
    flex: 1,
  },
  userName: {
    marginBottom: 4,
  },
  sectionTitle: {
    marginBottom: 16,
  },
  signOutSection: {
    paddingHorizontal: 16,
    marginBottom: 24,
  },
  signOutButton: {
    width: '100%',
  },
}); 