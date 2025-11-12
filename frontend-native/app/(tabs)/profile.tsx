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
      {/* App Bar */}
      <View style={styles.appBar}>
        <Typography variant="heading-xsmall" color="contentPrimary">
          Profile
        </Typography>
      </View>

      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        {/* User Profile Card */}
        <View style={styles.section}>
          <View style={styles.profileHeader}>
            {user.picture && (
              <Image source={{ uri: user.picture }} style={styles.avatar} />
            )}
            <View style={styles.userInfo}>
              <Typography variant="label-medium" color="contentPrimary" style={styles.userName}>
                {user.name}
              </Typography>
              <Typography variant="label-small" color="contentSecondary">
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
    backgroundColor: getColor('backgroundPrimary'),
  },
  scrollView: {
    flex: 1,
  },
  appBar: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: getColor('backgroundPrimary'),
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderOpaque'),
  },
  section: {
    paddingHorizontal: 16,
    paddingTop: 24,
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