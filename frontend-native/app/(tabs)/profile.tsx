import React from 'react';
import {
  View,
  StyleSheet,
  ScrollView,
  Image,
  Alert,
  Pressable,
  Linking,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { router } from 'expo-router';
import { useClerk } from '@clerk/clerk-expo';
import { Typography, Button } from '../../src/components';
import { getColor } from '../../src/components/Colors';
import { useAuth } from '../../src/hooks/useAuth';

export default function ProfileTab() {
  const { user } = useAuth();
  const { signOut } = useClerk();

  const handleSignOut = async () => {
    try {
      await signOut();
      router.replace('/(auth)/sign-in');
    } catch (err) {
      console.error('Sign out error:', JSON.stringify(err, null, 2));
      Alert.alert('Error', 'Failed to sign out. Please try again.');
    }
  };

  const handleDeleteAccount = async () => {
    const deleteFormUrl = 'https://docs.google.com/forms/d/e/1FAIpQLSdS-SsF6rY91h53evUPfd6RAjOn1dbpGJVmM498HQolcRl26g/viewform?usp=sharing&ouid=108074063199428464488';
    
    try {
      const supported = await Linking.canOpenURL(deleteFormUrl);
      if (supported) {
        await Linking.openURL(deleteFormUrl);
      } else {
        Alert.alert('Error', 'Unable to open the link.');
      }
    } catch (err) {
      console.error('Error opening delete form:', err);
      Alert.alert('Error', 'Failed to open the deletion form.');
    }
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

        {/* Delete Account */}
        <View style={styles.deleteAccountSection}>
          <Pressable onPress={handleDeleteAccount}>
            <Typography variant="label-xsmall" color="contentTertiary" style={styles.deleteAccountText}>
              Delete my account
            </Typography>
          </Pressable>
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
  deleteAccountSection: {
    paddingHorizontal: 16,
    paddingBottom: 32,
    alignItems: 'center',
  },
  deleteAccountText: {
    textDecorationLine: 'underline',
  },
}); 