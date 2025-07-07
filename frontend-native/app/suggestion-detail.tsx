import React, { useState } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { Typography, Button } from '../src/components';
import { getColor } from '../src/components/Colors';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { 
  WorkoutChange, 
  SuggestedWorkout, 
  StoredSuggestedWorkout 
} from '../src/types/suggestion';

export default function SuggestionDetailScreen() {
  const router = useRouter();
  const params = useLocalSearchParams();
  const [expandedChanges, setExpandedChanges] = useState<Set<number>>(new Set());
  const [processing, setProcessing] = useState(false);
  
  // Parse the individual suggested workout
  const suggestion: StoredSuggestedWorkout = JSON.parse(params.suggestion as string || '{}');
  
  if (!suggestion.id) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.errorContainer}>
          <Typography variant="paragraph-medium" color="contentSecondary">
            Suggestion not found
          </Typography>
        </View>
      </SafeAreaView>
    );
  }

  const toggleChangeExpanded = (index: number) => {
    const newExpanded = new Set(expandedChanges);
    if (newExpanded.has(index)) {
      newExpanded.delete(index);
    } else {
      newExpanded.add(index);
    }
    setExpandedChanges(newExpanded);
  };

  const getChangeTypeIcon = (type: string) => {
    switch (type.toLowerCase()) {
      case 'add':
        return 'add-circle-outline';
      case 'remove':
        return 'remove-circle-outline';
      case 'modify':
      case 'update':
        return 'create-outline';
      case 'reorder':
        return 'swap-vertical-outline';
      default:
        return 'ellipse-outline';
    }
  };

  const getChangeTypeColor = (type: string) => {
    switch (type.toLowerCase()) {
      case 'add':
        return getColor('positive');
      case 'remove':
        return getColor('negative');
      case 'modify':
      case 'update':
        return getColor('accent');
      default:
        return getColor('contentSecondary');
    }
  };

  const handleAccept = async () => {
    setProcessing(true);
    // Implement accept logic
    setTimeout(() => {
      setProcessing(false);
      router.back();
    }, 1000);
  };

  const handleReject = () => {
    // Implement reject logic
    router.back();
  };

  const formatDaysAnalyzed = (days: number) => {
    if (days === 1) return '1 day';
    if (days === 7) return '1 week';
    if (days === 14) return '2 weeks';
    if (days === 30) return '1 month';
    return `${days} days`;
  };

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity onPress={() => router.back()} style={styles.backButton}>
          <Ionicons name="arrow-back" size={24} color={getColor('contentPrimary')} />
        </TouchableOpacity>
        <Typography variant="label-small" color="contentPrimary" style={styles.headerTitle}>
          Workout Suggestion
        </Typography>
        <View style={styles.headerSpacer} />
      </View>

      <ScrollView 
        style={styles.scrollView}
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
      >
        {/* Workout Info Section */}
        <View style={styles.workoutSection}>
          <View style={styles.priorityBadge}>
            <View 
              style={[
                styles.priorityIndicator,
                { backgroundColor: suggestion.priority >= 4 ? getColor('accent') : getColor('borderOpaque') }
              ]}
            />
            <Typography variant="label-xsmall" color="contentSecondary">
              Priority {suggestion.priority}/5
            </Typography>
          </View>
          
          <Typography variant="label-medium" color="contentPrimary" style={styles.workoutName}>
            {suggestion.name}
          </Typography>
          
          {suggestion.description && (
            <Typography variant="paragraph-small" color="contentSecondary" style={styles.workoutDescription}>
              {suggestion.description}
            </Typography>
          )}
        </View>

        {/* Analysis Summary */}
        <View style={styles.summaryCard}>
          <View style={styles.summaryHeader}>
            <Ionicons name="analytics-outline" size={20} color={getColor('accent')} />
            <Typography variant="label-small" color="contentPrimary" style={styles.summaryTitle}>
              Analysis Summary
            </Typography>
          </View>
          
          <Typography variant="paragraph-small" color="contentSecondary" style={styles.summaryText}>
            {suggestion.overallReasoning}
          </Typography>
          
          <View style={styles.summaryMeta}>
            <View style={styles.metaItem}>
              <Ionicons name="calendar-outline" size={16} color={getColor('contentTertiary')} />
              <Typography variant="paragraph-xsmall" color="contentTertiary">
                Based on {formatDaysAnalyzed(suggestion.daysAnalyzed)}
              </Typography>
            </View>
            <View style={styles.metaItem}>
              <Ionicons name="fitness-outline" size={16} color={getColor('contentTertiary')} />
              <Typography variant="paragraph-xsmall" color="contentTertiary">
                {suggestion.changes.length} improvements
              </Typography>
            </View>
          </View>
        </View>

        {/* Changes Section */}
        <View style={styles.changesSection}>
          <Typography variant="label-medium" color="contentPrimary" style={styles.sectionTitle}>
            Suggested Changes
          </Typography>
          
          {suggestion.changes.map((change: WorkoutChange, index: number) => {
            const isExpanded = expandedChanges.has(index);
            const changeColor = getChangeTypeColor(change.type);
            
            return (
              <TouchableOpacity
                key={index}
                style={styles.changeCard}
                onPress={() => toggleChangeExpanded(index)}
                activeOpacity={0.7}
              >
                <View style={styles.changeHeader}>
                  <View style={styles.changeIconContainer}>
                    <Ionicons 
                      name={getChangeTypeIcon(change.type)} 
                      size={24} 
                      color={changeColor}
                    />
                  </View>
                  
                  <View style={styles.changeInfo}>
                    <Typography variant="label-medium" color="contentPrimary">
                      {change.exerciseName}
                    </Typography>
                    <View style={styles.changeTypeBadge}>
                      <Typography 
                        variant="label-xsmall" 
                        color="contentSecondary"
                        style={{ color: changeColor }}
                      >
                        {change.type}
                      </Typography>
                    </View>
                  </View>
                  
                  <Ionicons 
                    name={isExpanded ? "chevron-up" : "chevron-down"} 
                    size={20} 
                    color={getColor('contentTertiary')}
                  />
                </View>
                
                {isExpanded && (
                  <View style={styles.changeDetails}>
                    {change.oldValue && (
                      <View style={styles.changeComparison}>
                        <View style={styles.changeValue}>
                          <Typography variant="label-xsmall" color="contentTertiary">
                            Current
                          </Typography>
                          <Typography variant="paragraph-small" color="contentSecondary">
                            {change.oldValue}
                          </Typography>
                        </View>
                        
                        <Ionicons 
                          name="arrow-forward" 
                          size={16} 
                          color={getColor('contentTertiary')}
                          style={styles.changeArrow}
                        />
                        
                        <View style={styles.changeValue}>
                          <Typography variant="label-xsmall" color="contentTertiary">
                            Suggested
                          </Typography>
                          <Typography variant="paragraph-small" color="contentPrimary">
                            {change.newValue}
                          </Typography>
                        </View>
                      </View>
                    )}
                    
                    <View style={styles.changeReason}>
                      <Typography variant="paragraph-xsmall" color="contentSecondary">
                        {change.reason}
                      </Typography>
                    </View>
                  </View>
                )}
              </TouchableOpacity>
            );
          })}
        </View>
      </ScrollView>

      {/* Fixed Bottom Actions */}
      <View style={styles.bottomActions}>
        <Button
          variant="secondary"
          size="large"
          onPress={handleReject}
          style={styles.actionButton}
        >
          Dismiss
        </Button>
        <Button
          variant="primary"
          size="large"
          onPress={handleAccept}
          disabled={processing}
          style={styles.actionButton}
        >
          {processing ? (
            <ActivityIndicator size="small" color={getColor('contentOnColor')} />
          ) : (
            'Apply Changes'
          )}
        </Button>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'),
  },
  
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: getColor('backgroundPrimary'),
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderOpaque'),
  },
  
  backButton: {
    padding: 4,
  },
  
  headerTitle: {
    flex: 1,
    textAlign: 'center',
  },
  
  headerSpacer: {
    width: 32,
  },
  
  scrollView: {
    flex: 1,
  },
  
  scrollContent: {
    paddingBottom: 100,
  },
  
  workoutSection: {
    backgroundColor: getColor('backgroundPrimary'),
    padding: 16,
    marginBottom: 2,
  },
  
  priorityBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 8,
  },
  
  priorityIndicator: {
    width: 12,
    height: 12,
    borderRadius: 6,
    marginRight: 6,
  },
  
  workoutName: {
    marginBottom: 4,
  },
  
  workoutDescription: {
    lineHeight: 22,
  },
  
  summaryCard: {
    backgroundColor: getColor('backgroundPrimary'),
    padding: 16,
    marginBottom: 16,
  },
  
  summaryHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 12,
  },
  
  summaryTitle: {
    marginLeft: 8,
  },
  
  summaryText: {
    lineHeight: 20,
    marginBottom: 12,
  },
  
  summaryMeta: {
    flexDirection: 'row',
    gap: 16,
  },
  
  metaItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  
  changesSection: {
    paddingHorizontal: 16,
  },
  
  sectionTitle: {
    marginBottom: 12,
  },
  
  changeCard: {
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    marginBottom: 8,
    overflow: 'hidden',
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
  },
  
  changeHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 16,
  },
  
  changeIconContainer: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: getColor('backgroundTertiary'),
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  
  changeInfo: {
    flex: 1,
  },
  
  changeTypeBadge: {
    marginTop: 4,
  },
  
  changeDetails: {
    paddingHorizontal: 16,
    paddingBottom: 16,
    borderTopWidth: 1,
    borderTopColor: getColor('borderOpaque'),
  },
  
  changeComparison: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: 12,
    marginBottom: 12,
  },
  
  changeValue: {
    flex: 1,
  },
  
  changeArrow: {
    marginHorizontal: 12,
  },
  
  changeReason: {
    backgroundColor: getColor('backgroundTertiary'),
    padding: 12,
    borderRadius: 6,
  },
  
  bottomActions: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    flexDirection: 'row',
    padding: 16,
    gap: 12,
    backgroundColor: getColor('backgroundPrimary'),
    borderTopWidth: 1,
    borderTopColor: getColor('borderOpaque'),
    paddingBottom: 32,
  },
  
  actionButton: {
    flex: 1,
  },
  
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
}); 