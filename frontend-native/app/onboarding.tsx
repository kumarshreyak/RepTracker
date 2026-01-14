import React, { useState } from 'react';
import {
  View,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
  ActivityIndicator,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { router } from 'expo-router';
import { Typography, Button, Input } from '../src/components';
import { getColor } from '../src/components/Colors';
import { useAuth, useUser } from '@clerk/clerk-expo';
import { userService } from '../src/auth/UserService';

// Goal options with icons (matching backend constants)
const GOALS = [
  { id: 'gain_muscle', label: 'Gain Muscle', icon: '💪' },
  { id: 'lose_fat', label: 'Lose Fat', icon: '🔥' },
  { id: 'maintain', label: 'Maintain', icon: '⚖️' },
];

export default function OnboardingRoute() {
  const { getToken } = useAuth();
  const { user } = useUser();
  const [isLoading, setIsLoading] = useState(false);
  
  // Form state
  const [heightFeet, setHeightFeet] = useState('');
  const [heightInches, setHeightInches] = useState('');
  const [heightCm, setHeightCm] = useState('');
  const [weight, setWeight] = useState('');
  const [age, setAge] = useState('');
  const [goal, setGoal] = useState('');
  
  // Unit toggles
  const [useMetricHeight, setUseMetricHeight] = useState(false);
  const [useMetricWeight, setUseMetricWeight] = useState(false);

  const handleSubmit = async () => {
    // Validate inputs
    if (!age || !goal) {
      Alert.alert('Missing Information', 'Please fill in all fields');
      return;
    }

    if (useMetricHeight && !heightCm) {
      Alert.alert('Missing Information', 'Please enter your height');
      return;
    }

    if (!useMetricHeight && (!heightFeet || !heightInches)) {
      Alert.alert('Missing Information', 'Please enter your height');
      return;
    }

    if (!weight) {
      Alert.alert('Missing Information', 'Please enter your weight');
      return;
    }

    try {
      setIsLoading(true);
      
      // Convert height to cm for storage
      let heightInCm: number;
      if (useMetricHeight) {
        heightInCm = parseFloat(heightCm);
      } else {
        const feet = parseFloat(heightFeet) || 0;
        const inches = parseFloat(heightInches) || 0;
        heightInCm = (feet * 30.48) + (inches * 2.54);
      }

      // Convert weight to kg for storage
      const weightInKg = useMetricWeight 
        ? parseFloat(weight) 
        : parseFloat(weight) * 0.453592;

      // Update user profile via API
      if (!user) {
        throw new Error('User not authenticated');
      }

      const sessionToken = await getToken();
      if (!sessionToken) {
        throw new Error('No session token available');
      }
      
      const userData = {
        email: user.primaryEmailAddress?.emailAddress || '',
        firstName: user.firstName || '',
        lastName: user.lastName || '',
        height: heightInCm,
        weight: weightInKg,
        age: parseInt(age),
        goal: goal,
        picture: user.imageUrl || '',
      };

      console.log('Submitting user data:', userData);
      
      const success = await userService.createOrUpdateUserProfile(user.id, userData, sessionToken);
      
      if (!success) {
        throw new Error('Failed to create or update user profile');
      }
      
      // Navigate to main app
      router.replace('/(tabs)');
    } catch (error) {
      console.error('Onboarding error:', error);
      Alert.alert('Error', 'Failed to save your information. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView 
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
        keyboardShouldPersistTaps="handled"
      >
        {/* Progress indicator */}
        <View style={styles.progressContainer}>
          <View style={styles.progressDot} />
          <View style={[styles.progressDot, styles.progressDotActive]} />
          <View style={styles.progressDot} />
        </View>

        {/* Header */}
        <View style={styles.header}>
          <Typography variant="heading-large" color="contentPrimary">
            Let's personalize your experience
          </Typography>
          <Typography variant="paragraph-medium" color="contentSecondary" style={styles.subtitle}>
            This helps us provide better workout recommendations
          </Typography>
        </View>

        {/* Form Card */}
        <View style={styles.card}>
          {/* Height & Weight Row */}
          <View style={styles.row}>
            {/* Height Input */}
            <View style={styles.halfColumn}>
              <View style={styles.labelRow}>
                <Typography variant="label-small" color="contentPrimary">
                  Height
                </Typography>
                <TouchableOpacity 
                  onPress={() => setUseMetricHeight(!useMetricHeight)}
                  style={styles.unitToggle}
                >
                  <Typography variant="label-xsmall" color="accent">
                    {useMetricHeight ? 'cm' : 'ft/in'}
                  </Typography>
                </TouchableOpacity>
              </View>
              
              {useMetricHeight ? (
                <Input
                  value={heightCm}
                  onChangeText={setHeightCm}
                  placeholder="170 cm"
                  keyboardType="numeric"
                  maxLength={3}
                />
              ) : (
                <View style={styles.imperialHeight}>
                  <View style={styles.imperialInput}>
                    <Input
                      value={heightFeet}
                      onChangeText={setHeightFeet}
                      placeholder="5 ft"
                      keyboardType="numeric"
                      maxLength={1}
                    />
                  </View>
                  <View style={styles.imperialSpacer} />
                  <View style={styles.imperialInput}>
                    <Input
                      value={heightInches}
                      onChangeText={setHeightInches}
                      placeholder="10 in"
                      keyboardType="numeric"
                      maxLength={2}
                    />
                  </View>
                </View>
              )}
            </View>

            <View style={styles.columnSpacer} />

            {/* Weight Input */}
            <View style={styles.halfColumn}>
              <View style={styles.labelRow}>
                <Typography variant="label-small" color="contentPrimary">
                  Weight
                </Typography>
                <TouchableOpacity 
                  onPress={() => setUseMetricWeight(!useMetricWeight)}
                  style={styles.unitToggle}
                >
                  <Typography variant="label-xsmall" color="accent">
                    {useMetricWeight ? 'kg' : 'lbs'}
                  </Typography>
                </TouchableOpacity>
              </View>
              
              <Input
                value={weight}
                onChangeText={setWeight}
                placeholder={useMetricWeight ? "70 kg" : "155 lbs"}
                keyboardType="numeric"
                maxLength={3}
              />
            </View>
          </View>

          {/* Age Input */}
          <View style={styles.fieldContainer}>
            <Typography variant="label-small" color="contentPrimary" style={styles.label}>
              Age
            </Typography>
            <Input
              value={age}
              onChangeText={setAge}
              placeholder="25 years"
              keyboardType="numeric"
              maxLength={2}
            />
          </View>

          {/* Goal Selection */}
          <View style={styles.fieldContainer}>
            <Typography variant="label-small" color="contentPrimary" style={styles.label}>
              Fitness Goal
            </Typography>
            <View style={styles.goalGrid}>
              {GOALS.map((item) => (
                <TouchableOpacity
                  key={item.id}
                  style={[
                    styles.goalCard,
                    goal === item.id && styles.goalCardActive
                  ]}
                  onPress={() => setGoal(item.id)}
                  activeOpacity={0.7}
                >
                  <Typography variant="display-xsmall" style={styles.goalIcon}>
                    {item.icon}
                  </Typography>
                  <Typography 
                    variant="label-small" 
                    color={goal === item.id ? 'accent' : 'contentSecondary'}
                    style={styles.goalLabel}
                  >
                    {item.label}
                  </Typography>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          {/* Submit Button */}
          <Button
            variant="primary"
            size="large"
            onPress={handleSubmit}
            disabled={isLoading}
            style={styles.submitButton}
          >
            {isLoading ? 'Saving...' : 'Continue'}
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
  scrollContent: {
    flexGrow: 1,
    paddingHorizontal: 16,
    paddingBottom: 32,
  },
  progressContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: 24,
    gap: 8,
  },
  progressDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: getColor('borderTransparent'),
  },
  progressDotActive: {
    backgroundColor: getColor('accent'),
    width: 24,
  },
  header: {
    alignItems: 'center',
    marginBottom: 32,
    paddingHorizontal: 16,
  },
  subtitle: {
    marginTop: 8,
    textAlign: 'center',
  },
  card: {
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    padding: 24,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 1,
    },
    shadowOpacity: 0.05,
    shadowRadius: 2,
    elevation: 1,
  },
  row: {
    flexDirection: 'row',
    marginBottom: 24,
  },
  halfColumn: {
    flex: 1,
  },
  columnSpacer: {
    width: 12,
  },
  labelRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  label: {
    marginBottom: 8,
  },
  unitToggle: {
    padding: 4,
  },
  imperialHeight: {
    flexDirection: 'row',
  },
  imperialInput: {
    flex: 1,
  },
  imperialSpacer: {
    width: 8,
  },
  fieldContainer: {
    marginBottom: 24,
  },
  goalGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
    gap: 8,
  },
  goalCard: {
    width: '31%',
    backgroundColor: getColor('backgroundTertiary'),
    borderRadius: 8,
    padding: 16,
    alignItems: 'center',
    borderWidth: 2,
    borderColor: 'transparent',
    minHeight: 80,
  },
  goalCardActive: {
    borderColor: getColor('borderAccent'),
    backgroundColor: getColor('backgroundLightAccent'),
  },
  goalIcon: {
    marginBottom: 8,
  },
  goalLabel: {
    textAlign: 'center',
  },
  submitButton: {
    width: '100%',
    marginTop: 8,
  },
}); 