import React from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  Pressable,
} from 'react-native';
import { Typography } from '../src/components';
import { getColor } from '../src/components/Colors';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { AIProgressiveOverloadResponse } from '../src/types/ai-progressive';

export default function AIAnalysisDetailScreen() {
  const router = useRouter();
  const params = useLocalSearchParams();
  
  // Parse the analysis data from route params
  const analysis: AIProgressiveOverloadResponse = JSON.parse(params.analysis as string);

  const formatProcessingTime = (ms: number): string => {
    if (ms < 1000) return `${ms}ms`;
    return `${(ms / 1000).toFixed(1)}s`;
  };

  const formatTimeAgo = (dateString: string): string => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);
    
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    
    return date.toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric' 
    });
  };

  const renderExerciseCard = (exercise: AIProgressiveOverloadResponse['updatedExercises'][0], index: number) => (
    <View key={index} style={styles.exerciseCard}>
      {/* Progression indicator */}
      <View 
        style={[
          styles.progressionIndicator,
          { 
            backgroundColor: exercise.progressionApplied 
              ? getColor('positive') 
              : getColor('borderOpaque')
          }
        ]} 
      />
      
      <View style={styles.exerciseContent}>
        {/* Exercise header */}
        <View style={styles.exerciseHeader}>
          <Typography variant="label-small" color="contentPrimary">
            {exercise.exerciseName}
          </Typography>
          {exercise.progressionApplied && (
            <View style={styles.progressedBadge}>
              <Typography variant="label-xsmall" color="contentOnColor">
                Improved
              </Typography>
            </View>
          )}
        </View>
        
        {/* Changes made */}
        {exercise.progressionApplied && exercise.changesMade && (
          <View style={styles.changesSection}>
            <Typography variant="paragraph-small" color="contentSecondary">
              {exercise.changesMade}
            </Typography>
          </View>
        )}
        
        {/* Reasoning */}
        <View style={styles.reasoningSection}>
          <Typography variant="paragraph-xsmall" color="contentTertiary">
            {exercise.reasoning}
          </Typography>
        </View>
        
        {/* Exercise notes */}
        {exercise.notes && (
          <View style={styles.notesSection}>
            <Typography variant="label-xsmall" color="contentTertiary">
              Notes:
            </Typography>
            <Typography variant="paragraph-xsmall" color="contentSecondary">
              {exercise.notes}
            </Typography>
          </View>
        )}
        
        {/* Sets preview with notes */}
        {exercise.sets.length > 0 && (
          <View style={styles.setsSection}>
            <Typography variant="label-xsmall" color="contentTertiary" style={styles.setsHeader}>
              {exercise.sets.length} {exercise.sets.length === 1 ? 'set' : 'sets'}
            </Typography>
            {exercise.sets.some(set => set.notes) && (
              <View style={styles.setNotes}>
                {exercise.sets.map((set, setIndex) => 
                  set.notes ? (
                    <View key={setIndex} style={styles.setNote}>
                      <Typography variant="label-xsmall" color="contentTertiary">
                        Set {setIndex + 1}:
                      </Typography>
                      <Typography variant="paragraph-xsmall" color="contentSecondary">
                        {set.notes}
                      </Typography>
                    </View>
                  ) : null
                )}
              </View>
            )}
          </View>
        )}
      </View>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Pressable onPress={() => router.back()} style={styles.backButton}>
          <Typography variant="label-medium" color="contentPrimary">
            ←
          </Typography>
        </Pressable>
        <Typography variant="label-medium" color="contentPrimary">
          AI Analysis
        </Typography>
        <View style={styles.headerSpacer} />
      </View>

      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        {/* Status card */}
        <View style={styles.statusCard}>
          <View style={styles.statusHeader}>
            <View 
              style={[
                styles.statusIndicator,
                { 
                  backgroundColor: analysis.success 
                    ? getColor('backgroundPositive') 
                    : getColor('backgroundNegative')
                }
              ]}
            >
              <Typography variant="label-xsmall" color="contentOnColor">
                {analysis.success ? '✓' : '✗'}
              </Typography>
            </View>
            <Typography variant="label-medium" color="contentPrimary">
              {analysis.success ? 'Analysis Complete' : 'Analysis Failed'}
            </Typography>
          </View>
          
          {/* Meta info */}
          <View style={styles.metaInfo}>
            <View style={styles.metaItem}>
              <Typography variant="label-xsmall" color="contentTertiary">
                Processed
              </Typography>
              <Typography variant="paragraph-xsmall" color="contentSecondary">
                {formatTimeAgo(analysis.createdAt)}
              </Typography>
            </View>
            <View style={styles.metaItem}>
              <Typography variant="label-xsmall" color="contentTertiary">
                Duration
              </Typography>
              <Typography variant="paragraph-xsmall" color="contentSecondary">
                {formatProcessingTime(analysis.processingTimeMs)}
              </Typography>
            </View>
            <View style={styles.metaItem}>
              <Typography variant="label-xsmall" color="contentTertiary">
                Sessions Analyzed
              </Typography>
              <Typography variant="paragraph-xsmall" color="contentSecondary">
                {analysis.recentSessionsCount}
              </Typography>
            </View>
          </View>
        </View>

        {/* Analysis summary */}
        {analysis.analysisSummary && (
          <View style={styles.summaryCard}>
            <Typography variant="label-medium" color="contentPrimary" style={styles.sectionTitle}>
              Summary
            </Typography>
            <Typography variant="paragraph-medium" color="contentSecondary" style={styles.summaryText}>
              {analysis.analysisSummary}
            </Typography>
          </View>
        )}

        {/* Exercise updates */}
        <View style={styles.exercisesSection}>
          <View style={styles.exercisesHeader}>
            <Typography variant="label-medium" color="contentPrimary">
              Exercise Updates
            </Typography>
            <Typography variant="label-xsmall" color="contentTertiary">
              {analysis.updatedExercises.filter(e => e.progressionApplied).length} of {analysis.updatedExercises.length} improved
            </Typography>
          </View>
          
          <View style={styles.exercisesList}>
            {analysis.updatedExercises.map(renderExerciseCard)}
          </View>
        </View>

        {/* AI Model info */}
        <View style={styles.aiModelCard}>
          <Typography variant="label-xsmall" color="contentTertiary" style={styles.aiModelText}>
            Analyzed by {analysis.aiModel}
          </Typography>
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
  
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: getColor('backgroundPrimary'),
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderOpaque'),
  },
  
  backButton: {
    padding: 8,
    marginLeft: -8,
  },
  
  headerSpacer: {
    width: 24,
  },
  
  scrollView: {
    flex: 1,
  },
  
  statusCard: {
    backgroundColor: getColor('backgroundPrimary'),
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderOpaque'),
    paddingHorizontal: 16,
    paddingVertical: 20,
  },
  
  statusHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
  },
  
  statusIndicator: {
    width: 28,
    height: 28,
    borderRadius: 14,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  
  metaInfo: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  
  metaItem: {
    alignItems: 'center',
  },
  
  summaryCard: {
    backgroundColor: getColor('backgroundPrimary'),
    marginTop: 16,
    marginHorizontal: 16,
    padding: 16,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
  },
  
  sectionTitle: {
    marginBottom: 8,
  },
  
  summaryText: {
    lineHeight: 22,
  },
  
  exercisesSection: {
    marginTop: 16,
    marginHorizontal: 16,
  },
  
  exercisesHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  
  exercisesList: {
    gap: 12,
  },
  
  exerciseCard: {
    flexDirection: 'row',
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    overflow: 'hidden',
  },
  
  progressionIndicator: {
    width: 4,
  },
  
  exerciseContent: {
    flex: 1,
    padding: 16,
  },
  
  exerciseHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  
  progressedBadge: {
    backgroundColor: getColor('positive'),
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
  },
  
  changesSection: {
    marginBottom: 8,
  },
  
  reasoningSection: {
    marginBottom: 8,
  },
  
  notesSection: {
    marginBottom: 8,
  },
  
  setsSection: {
    marginTop: 4,
  },
  
  setsHeader: {
    marginBottom: 8,
  },
  
  setNotes: {
    gap: 6,
  },
  
  setNote: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    gap: 8,
  },
  
  aiModelCard: {
    marginTop: 24,
    marginHorizontal: 16,
    marginBottom: 32,
    padding: 12,
    backgroundColor: getColor('backgroundTertiary'),
    borderRadius: 6,
    alignItems: 'center',
  },
  
  aiModelText: {
    textAlign: 'center',
  },
});