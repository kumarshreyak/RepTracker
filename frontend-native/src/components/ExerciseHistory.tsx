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
import { apiGet } from '@/utils/api';

export interface ExerciseHistoryProps {
  exerciseId: string;
  exerciseName?: string;
}

export const ExerciseHistory: React.FC<ExerciseHistoryProps> = ({
  exerciseId,
  exerciseName,
}) => {
  const [historyData, setHistoryData] = useState<ExerciseHistoryResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [expandedCards, setExpandedCards] = useState<Set<number>>(new Set());

  useEffect(() => {
    fetchExerciseHistory();
  }, [exerciseId]);

  const fetchExerciseHistory = async () => {
    try {
      setLoading(true);
      setError(null);

      const data = await apiGet<ExerciseHistoryResponse>(
        `/api/exercises/${exerciseId}/history`
      );
      setHistoryData(data);
    } catch (err) {
      console.error('Error fetching exercise history:', err);
      setError(err instanceof Error ? err.message : 'Failed to load history');
    } finally {
      setLoading(false);
    }
  };

  const toggleCardExpansion = (index: number) => {
    const newExpandedCards = new Set(expandedCards);
    if (newExpandedCards.has(index)) {
      newExpandedCards.delete(index);
    } else {
      newExpandedCards.add(index);
    }
    setExpandedCards(newExpandedCards);
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
    const isExpanded = expandedCards.has(index);
    const hasAIAnalysis = entry.aiExercise && (entry.aiExercise.reasoning || entry.aiExercise.changesMade);

    return (
      <TouchableOpacity 
        key={index} 
        style={[styles.historyCard, isExpanded && styles.historyCardExpanded]}
        onPress={() => toggleCardExpansion(index)}
        activeOpacity={0.7}
      >
        {/* Card Header - Date and AI Indicator */}
        <View style={styles.cardHeader}>
          <View style={styles.cardHeaderLeft}>
            <Typography variant="label-medium" color="contentPrimary">
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
          <View style={styles.cardHeaderRight}>
            {hasAIAnalysis && (
              <MaterialIcons 
                name="psychology" 
                size={16} 
                color={getColor('contentTertiary')} 
              />
            )}
            {hasAIAnalysis && (
              <MaterialIcons 
                name={isExpanded ? "expand-less" : "expand-more"} 
                size={20} 
                color={getColor('contentTertiary')} 
              />
            )}
          </View>
        </View>

        {/* Performance Summary */}
        <View style={styles.performanceSection}>
          <View style={styles.performanceRow}>
            {bestSet && (
              <View style={styles.bestSetContainer}>
                <View style={styles.metricGroup}>
                  <Typography variant="heading-small" color="contentPrimary">
                    {bestSet.actualWeight}
                  </Typography>
                  <Typography variant="paragraph-xsmall" color="contentTertiary">
                    kg
                  </Typography>
                </View>
                <View style={styles.metricSeparator} />
                <View style={styles.metricGroup}>
                  <Typography variant="heading-small" color="contentPrimary">
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
              <Typography variant="paragraph-small" color="contentSecondary">
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
        </View>

        {/* Expanded AI Analysis Summary */}
        {isExpanded && hasAIAnalysis && (
          <View style={styles.aiSection}>
            <View style={styles.aiHeader}>
              <MaterialIcons 
                name={entry.aiExercise?.progressionApplied ? "trending-up" : "analytics"} 
                size={16} 
                color={getColor(entry.aiExercise?.progressionApplied ? 'contentPositive' : 'contentAccent')} 
              />
              <Typography 
                variant="label-small" 
                color={entry.aiExercise?.progressionApplied ? 'contentPositive' : 'contentAccent'}
              >
                {entry.aiExercise?.progressionApplied ? 'AI Updated' : 'AI Analysis'}
              </Typography>
            </View>
            
            {/* Show reasoning first */}
            {entry.aiExercise?.reasoning && (
              <View style={styles.aiTextContainer}>
                <Typography variant="paragraph-xsmall" color="contentTertiary">
                  Analysis
                </Typography>
                <Typography 
                  variant="paragraph-small" 
                  color="contentSecondary"
                  style={styles.aiText}
                >
                  {entry.aiExercise.reasoning}
                </Typography>
              </View>
            )}
            
            {/* Show changes made if any */}
            {entry.aiExercise?.changesMade && (
              <View style={styles.aiTextContainer}>
                <Typography variant="paragraph-xsmall" color="contentTertiary">
                  Changes Made
                </Typography>
                <Typography 
                  variant="paragraph-small" 
                  color={entry.aiExercise.progressionApplied ? 'contentPositive' : 'contentSecondary'}
                  style={styles.aiText}
                >
                  {entry.aiExercise.changesMade}
                </Typography>
              </View>
            )}
          </View>
        )}
      </TouchableOpacity>
    );
  };

  if (loading) {
    return (
      <View style={styles.container}>
        <View style={styles.header}>
          <Typography variant="heading-small" color="contentPrimary">
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
          <Typography variant="heading-small" color="contentPrimary">
            Recent History
          </Typography>
        </View>
        <View style={styles.emptyContainer}>
          <MaterialIcons 
            name="history" 
            size={32} 
            color={getColor('contentTertiary')} 
          />
          <Typography variant="paragraph-medium" color="contentTertiary" style={styles.emptyText}>
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
      <View style={styles.cardsContainer}>
        {historyData.history.map((entry, index) => renderHistoryCard(entry, index))}
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'),
  },
  
  // Header
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 16,
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderOpaque'),
  },

  // Loading and Empty States
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: 200,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    gap: 12,
    minHeight: 200,
  },
  emptyText: {
    textAlign: 'center',
  },

  // Cards Container
  cardsContainer: {
    padding: 16,
    gap: 12,
  },

  // History Cards
  historyCard: {
    backgroundColor: getColor('backgroundPrimary'),
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    borderRadius: 12,
    padding: 16,
    gap: 12,
  },
  historyCardExpanded: {
    borderColor: getColor('borderSelected'),
    shadowColor: getColor('contentPrimary'),
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },

  // Card Header
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  cardHeaderLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  cardHeaderRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  aiIndicator: {
    backgroundColor: getColor('backgroundLightAccent'),
    borderRadius: 8,
    padding: 3,
  },

  // Performance Section
  performanceSection: {
    gap: 12,
  },
  performanceRow: {
    gap: 12,
  },
  bestSetContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    gap: 16,
    paddingVertical: 8,
  },
  metricGroup: {
    alignItems: 'center',
    gap: 4,
  },
  metricSeparator: {
    width: 1,
    height: 30,
    backgroundColor: getColor('borderOpaque'),
  },
  
  // Sets Completion
  setsCompletion: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingTop: 8,
    borderTopWidth: 1,
    borderTopColor: getColor('borderOpaque'),
  },
  completionDots: {
    flexDirection: 'row',
    gap: 4,
  },
  completionDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
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
    borderRadius: 8,
    padding: 12,
    gap: 12,
    marginTop: 4,
  },
  aiHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  aiTextContainer: {
    gap: 4,
  },
  aiText: {
    lineHeight: 18,
  },
}); 