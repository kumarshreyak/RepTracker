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

  const handleSignOut = async () => {
    await authService.signOut();
    router.replace('/');
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
          <Typography variant="heading-large" color="dark" style={styles.title}>
            Profile
          </Typography>
          <Typography variant="text-default" color="light">
            Manage your account and preferences
          </Typography>
        </View>

        {/* User Profile Card */}
        <View style={styles.card}>
          <View style={styles.profileHeader}>
            {user.picture && (
              <Image source={{ uri: user.picture }} style={styles.avatar} />
            )}
            <View style={styles.userInfo}>
              <Typography variant="heading-default" color="dark" style={styles.userName}>
                {user.name}
              </Typography>
              <Typography variant="text-default" color="light">
                {user.email}
              </Typography>
            </View>
          </View>
        </View>

        {/* Profile Details */}
        <View style={styles.card}>
          <Typography variant="heading-small" color="dark" style={styles.sectionTitle}>
            Account Information
          </Typography>
          
          <View style={styles.profileDetails}>
            <View style={styles.detailItem}>
              <Typography variant="text-small" color="light" style={styles.detailLabel}>
                Name
              </Typography>
              <Typography variant="text-default" color="dark">
                {user.name || 'Not provided'}
              </Typography>
            </View>
            
            <View style={styles.detailItem}>
              <Typography variant="text-small" color="light" style={styles.detailLabel}>
                Email
              </Typography>
              <Typography variant="text-default" color="dark">
                {user.email || 'Not provided'}
              </Typography>
            </View>
            
            <View style={styles.detailItem}>
              <Typography variant="text-small" color="light" style={styles.detailLabel}>
                User ID
              </Typography>
              <Typography variant="text-default" color="dark">
                {user.id || 'Not provided'}
              </Typography>
            </View>
          </View>
        </View>

        {/* Actions */}
        <View style={styles.card}>
          <Typography variant="heading-small" color="dark" style={styles.sectionTitle}>
            Actions
          </Typography>
          
          <View style={styles.actions}>
            <Button
              variant="secondary"
              size="large"
              style={styles.actionButton}
              onPress={() => {
                // TODO: Implement settings screen
                console.log('Navigate to settings');
              }}
            >
              Settings
            </Button>
            
            <Button
              variant="secondary"
              size="large"
              style={styles.actionButton}
              onPress={() => {
                // TODO: Implement help screen
                console.log('Navigate to help');
              }}
            >
              Help & Support
            </Button>
            
            <Button
              variant="danger"
              size="large"
              style={styles.actionButton}
              onPress={handleSignOut}
            >
              Sign Out
            </Button>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('light-gray-1'),
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
  card: {
    backgroundColor: 'white',
    marginHorizontal: 16,
    marginBottom: 16,
    borderRadius: 8,
    padding: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
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
    backgroundColor: getColor('light-gray-3'),
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
  profileDetails: {
    gap: 16,
  },
  detailItem: {
    marginBottom: 8,
  },
  detailLabel: {
    marginBottom: 4,
  },
  actions: {
    gap: 12,
  },
  actionButton: {
    width: '100%',
  },
}); 