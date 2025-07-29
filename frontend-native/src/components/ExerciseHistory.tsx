import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { Typography } from './Typography';
import { getColor } from './Colors';
import { MaterialIcons } from '@expo/vector-icons';
import { ExerciseHistoryResponse, ExerciseHistoryEntry } from '@/types/exercise';

export interface ExerciseHistoryProps {
  exerciseId: string;
  userId: string;
  exerciseName?: string;
}

export const ExerciseHistory: React.FC<ExerciseHistoryProps> = ({
  exerciseId,
  userId,
  exerciseName,
}) => {
  const [historyData, setHistoryData] = useState<ExerciseHistoryResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchExerciseHistory();
  }, [exerciseId, userId]);

  const fetchExerciseHistory = async () => {
    try {
      setLoading(true);
      setError(null);

      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(`${API_BASE_URL}/api/exercises/${exerciseId}/history?user_id=${userId}`);
      
      if (!response.ok) {
        throw new Error('Failed to fetch exercise history');
      }

      const data: ExerciseHistoryResponse = await response.json();
      setHistoryData(data);
    } catch (err) {
      console.error('Error fetching exercise history:', err);
      setError(err instanceof Error ? err.message : 'Failed to load history');
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (date: Date | string): string => {
    const now = new Date();
    const targetDate = typeof date === 'string' ? new Date(date) : date;
    
    // Check if date is valid
    if (isNaN(targetDate.getTime())) {
      return 'Unknown';
    }
    
    const diffDays = Math.floor((now.getTime() - targetDate.getTime()) / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays}d ago`;
    if (diffDays < 30) return `${Math.floor(diffDays / 7)}w ago`;
    return targetDate.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  };

  const getBestSet = (entry: ExerciseHistoryEntry) => {
    if (!entry.sessionExercise.sets.length) return null;
    
    return entry.sessionExercise.sets.reduce((best, set) => {
      const currentVolume = set.actualWeight * set.actualReps;
      const bestVolume = best.actualWeight * best.actualReps;
      return currentVolume > bestVolume ? set : best;
    });
  };

  const renderHistoryCard = (entry: ExerciseHistoryEntry, index: number) => {
    const bestSet = getBestSet(entry);
    const totalSets = entry.sessionExercise.sets.length;
    const completedSets = entry.sessionExercise.sets.filter(set => set.completed).length;
    const hasAIChanges = entry.aiExercise?.progressionApplied || false;

    return (
      <View key={index} style={styles.historyCard}>
        {/* Card Header - Date and AI Indicator */}
        <View style={styles.cardHeader}>
          <Typography variant="label-small" color="contentTertiary">
            {formatDate(entry.sessionInfo.finishedAt)}
          </Typography>
          {hasAIChanges && (
            <View style={styles.aiIndicator}>
              <MaterialIcons 
                name="auto-awesome" 
                size={12} 
                color={getColor('contentAccent')} 
              />
            </View>
          )}
        </View>

        {/* Performance Summary */}
        <View style={styles.performanceSection}>
          {bestSet && (
            <View style={styles.bestSetRow}>
              <View style={styles.metricGroup}>
                <Typography variant="label-medium" color="contentPrimary">
                  {bestSet.actualWeight}
                </Typography>
                <Typography variant="paragraph-xsmall" color="contentTertiary">
                  kg
                </Typography>
              </View>
              <View style={styles.metricSeparator} />
              <View style={styles.metricGroup}>
                <Typography variant="label-medium" color="contentPrimary">
                  {bestSet.actualReps}
                </Typography>
                <Typography variant="paragraph-xsmall" color="contentTertiary">
                  reps
                </Typography>
              </View>
            </View>
          )}
          
          {/* Sets Completion */}
          <View style={styles.setsCompletion}>
            <Typography variant="paragraph-xsmall" color="contentSecondary">
              {completedSets}/{totalSets} sets
            </Typography>
            <View style={styles.completionDots}>
              {Array.from({ length: totalSets }).map((_, i) => (
                <View 
                  key={i}
                  style={[
                    styles.completionDot,
                    i < completedSets ? styles.completionDotCompleted : styles.completionDotPending
                  ]}
                />
              ))}
            </View>
          </View>
        </View>

        {/* AI Analysis Summary */}
        {entry.aiExercise && (
          <View style={styles.aiSection}>
            <View style={styles.aiHeader}>
              <MaterialIcons 
                name={entry.aiExercise.progressionApplied ? "trending-up" : "analytics"} 
                size={14} 
                color={getColor(entry.aiExercise.progressionApplied ? 'contentPositive' : 'contentAccent')} 
              />
              <Typography 
                variant="paragraph-xsmall" 
                color={entry.aiExercise.progressionApplied ? 'contentPositive' : 'contentAccent'}
              >
                {entry.aiExercise.progressionApplied ? 'AI Updated' : 'AI Analysis'}
              </Typography>
            </View>
            
            {/* Show reasoning first */}
            {entry.aiExercise.reasoning && (
              <Typography 
                variant="paragraph-xsmall" 
                color="contentSecondary"
                style={styles.aiText}
                numberOfLines={2}
              >
                {entry.aiExercise.reasoning}
              </Typography>
            )}
            
            {/* Show changes made if any */}
            {entry.aiExercise.changesMade && (
              <Typography 
                variant="paragraph-xsmall" 
                color={entry.aiExercise.progressionApplied ? 'contentPositive' : 'contentSecondary'}
                style={styles.aiText}
                numberOfLines={2}
              >
                {entry.aiExercise.changesMade}
              </Typography>
            )}
          </View>
        )}
      </View>
    );
  };

  if (loading) {
    return (
      <View style={styles.container}>
        <View style={styles.header}>
          <Typography variant="label-small" color="contentPrimary">
            Recent History
          </Typography>
        </View>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="small" color={getColor('contentTertiary')} />
        </View>
      </View>
    );
  }

  if (error || !historyData?.history.length) {
    return (
      <View style={styles.container}>
        <View style={styles.header}>
          <Typography variant="label-small" color="contentPrimary">
            Recent History
          </Typography>
        </View>
        <View style={styles.emptyContainer}>
          <MaterialIcons 
            name="history" 
            size={32} 
            color={getColor('contentTertiary')} 
          />
          <Typography variant="paragraph-small" color="contentTertiary" style={styles.emptyText}>
            {error ? 'Failed to load' : 'No recent sessions'}
          </Typography>
        </View>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      {/* Section Header */}
      <View style={styles.header}>
        <Typography variant="label-small" color="contentPrimary">
          Recent History
        </Typography>
        <Typography variant="paragraph-xsmall" color="contentTertiary">
          Last {historyData.history.length} sessions
        </Typography>
      </View>

      {/* History Cards */}
      <ScrollView 
        horizontal 
        showsHorizontalScrollIndicator={false}
        contentContainerStyle={styles.scrollContent}
        style={styles.scrollView}
      >
        {historyData.history.map((entry, index) => renderHistoryCard(entry, index))}
      </ScrollView>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    paddingVertical: 16,
    backgroundColor: getColor('backgroundSecondary'),
  },
  
  // Header
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    marginBottom: 12,
  },

  // Loading and Empty States
  loadingContainer: {
    height: 120,
    justifyContent: 'center',
    alignItems: 'center',
  },
  emptyContainer: {
    height: 120,
    justifyContent: 'center',
    alignItems: 'center',
    gap: 8,
  },
  emptyText: {
    textAlign: 'center',
  },

  // Scroll View
  scrollView: {
    height: 140,
  },
  scrollContent: {
    paddingHorizontal: 16,
    gap: 12,
  },

  // History Cards
  historyCard: {
    width: 160,
    backgroundColor: getColor('backgroundPrimary'),
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    borderRadius: 8,
    padding: 12,
    gap: 10,
  },

  // Card Header
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  aiIndicator: {
    backgroundColor: getColor('backgroundLightAccent'),
    borderRadius: 8,
    padding: 2,
  },

  // Performance Section
  performanceSection: {
    gap: 8,
  },
  bestSetRow: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    gap: 8,
  },
  metricGroup: {
    alignItems: 'center',
    gap: 2,
  },
  metricSeparator: {
    width: 1,
    height: 20,
    backgroundColor: getColor('borderOpaque'),
  },
  
  // Sets Completion
  setsCompletion: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  completionDots: {
    flexDirection: 'row',
    gap: 3,
  },
  completionDot: {
    width: 6,
    height: 6,
    borderRadius: 3,
  },
  completionDotCompleted: {
    backgroundColor: getColor('positive'),
  },
  completionDotPending: {
    backgroundColor: getColor('borderOpaque'),
  },

  // AI Section
  aiSection: {
    backgroundColor: getColor('backgroundLightPositive'),
    borderRadius: 6,
    padding: 8,
    gap: 4,
  },
  aiHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  aiText: {
    lineHeight: 14,
  },
}); 