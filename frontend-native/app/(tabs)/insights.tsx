import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  RefreshControl,
  ActivityIndicator,
  TouchableOpacity,
} from 'react-native';
import { Typography } from '../../src/components';
import { getColor } from '../../src/components/Colors';
import { useAuth } from '../../src/hooks/useAuth';
import { useFocusEffect, useRouter } from 'expo-router';
import { 
  AIProgressiveOverloadResponse,
  ListAIProgressiveOverloadResponsesResponse
} from '../../src/types/ai-progressive';

export default function InsightsTab() {
  const { user } = useAuth();
  const router = useRouter();
  const [analyses, setAnalyses] = useState<AIProgressiveOverloadResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAIAnalyses = async (isRefreshing = false) => {
    if (!user?.id) return;

    try {
      if (!isRefreshing) setLoading(true);
      setError(null);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(
        `${API_BASE_URL}/api/users/${user.id}/ai-progressive-overload-responses?pageSize=20`
      );
      
      const data: ListAIProgressiveOverloadResponsesResponse = await response.json();
      
      if (!response.ok) {
        throw new Error("Failed to fetch AI analyses");
      }
      
      setAnalyses(data.responses || []);
    } catch (err) {
      console.error("Error fetching AI analyses:", err);
      setError(err instanceof Error ? err.message : "Failed to fetch analyses");
    } finally {
      setLoading(false);
      if (isRefreshing) setRefreshing(false);
    }
  };

  useEffect(() => {
    if (user?.id) {
      fetchAIAnalyses();
    }
  }, [user?.id]);

  // Refresh when screen comes into focus
  useFocusEffect(
    React.useCallback(() => {
      if (user?.id) {
        fetchAIAnalyses();
      }
    }, [user?.id])
  );

  const onRefresh = () => {
    setRefreshing(true);
    fetchAIAnalyses(true);
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

  const navigateToDetail = (analysis: AIProgressiveOverloadResponse) => {
    router.push({
      pathname: '/ai-analysis-detail',
      params: {
        analysis: JSON.stringify(analysis)
      }
    });
  };

  const renderAnalysisCard = (analysis: AIProgressiveOverloadResponse) => {
    const improvedCount = analysis.updatedExercises.filter(e => e.progressionApplied).length;
    const totalCount = analysis.updatedExercises.length;
    
    return (
      <TouchableOpacity
        key={analysis.id}
        style={styles.analysisCard}
        onPress={() => navigateToDetail(analysis)}
        activeOpacity={0.7}
      >
        {/* Status indicator */}
        <View 
          style={[
            styles.statusIndicator,
            { 
              backgroundColor: analysis.success 
                ? getColor('positive') 
                : getColor('negative')
            }
          ]} 
        />
        
        <View style={styles.cardContent}>
          {/* Header with status and time */}
          <View style={styles.cardHeader}>
            <View style={styles.statusSection}>
              <Typography variant="label-medium" color="contentPrimary">
                {analysis.success ? 'Analysis Complete' : 'Analysis Failed'}
              </Typography>
              <Typography variant="paragraph-xsmall" color="contentTertiary">
                {formatTimeAgo(analysis.createdAt)}
              </Typography>
            </View>
            
            {/* Progress indicator */}
            {analysis.success && (
              <View style={styles.progressBadge}>
                <Typography variant="label-xsmall" color="contentOnColor">
                  {improvedCount}/{totalCount}
                </Typography>
              </View>
            )}
          </View>
          
          {/* Summary or error message */}
          {analysis.success ? (
            <Typography 
              variant="paragraph-small" 
              color="contentSecondary"
              numberOfLines={2}
              style={styles.summaryText}
            >
              {analysis.analysisSummary || `Analyzed ${totalCount} exercises, improved ${improvedCount}`}
            </Typography>
          ) : (
            <Typography 
              variant="paragraph-small" 
              color="contentNegative"
              numberOfLines={2}
              style={styles.summaryText}
            >
              {analysis.message || 'Analysis failed to complete'}
            </Typography>
          )}
          
          {/* Exercise indicators */}
          {analysis.success && analysis.updatedExercises.length > 0 && (
            <View style={styles.exerciseIndicators}>
              {analysis.updatedExercises.slice(0, 3).map((exercise, index) => (
                <View 
                  key={index}
                  style={[
                    styles.exerciseChip,
                    {
                      backgroundColor: exercise.progressionApplied 
                        ? getColor('backgroundPositive') + '20'
                        : getColor('backgroundTertiary')
                    }
                  ]}
                >
                  <Typography variant="label-xsmall" color="contentSecondary">
                    {exercise.exerciseName.length > 12 
                      ? exercise.exerciseName.substring(0, 12) + '...'
                      : exercise.exerciseName}
                  </Typography>
                  {exercise.progressionApplied && (
                    <View style={styles.improvedDot} />
                  )}
                </View>
              ))}
              {analysis.updatedExercises.length > 3 && (
                <Typography variant="label-xsmall" color="contentTertiary">
                  +{analysis.updatedExercises.length - 3}
                </Typography>
              )}
            </View>
          )}
          
          {/* Meta info */}
          <View style={styles.metaInfo}>
            <Typography variant="paragraph-xsmall" color="contentTertiary">
              {analysis.recentSessionsCount} sessions analyzed
            </Typography>
            <Typography variant="paragraph-xsmall" color="contentTertiary">
              •
            </Typography>
            <Typography variant="paragraph-xsmall" color="contentTertiary">
              {analysis.aiModel}
            </Typography>
          </View>
        </View>
      </TouchableOpacity>
    );
  };

  if (loading && analyses.length === 0) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.appBar}>
          <Typography variant="heading-xsmall" color="contentPrimary">
            AI Insights
          </Typography>
        </View>
        <View style={styles.centerContent}>
          <ActivityIndicator size="large" color={getColor('accent')} />
          <Typography variant="paragraph-small" color="contentTertiary" style={styles.loadingText}>
            Loading analyses...
          </Typography>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      {/* App Bar */}
      <View style={styles.appBar}>
        <Typography variant="heading-xsmall" color="contentPrimary">
          AI Insights
        </Typography>
      </View>

      <ScrollView 
        style={styles.scrollView}
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            tintColor={getColor('accent')}
          />
        }
      >
        {error ? (
          <View style={styles.errorState}>
            <Typography variant="heading-xsmall" color="contentSecondary">
              Unable to load analyses
            </Typography>
            <Typography 
              variant="paragraph-small" 
              color="contentTertiary"
              style={styles.errorText}
            >
              {error}
            </Typography>
          </View>
        ) : analyses.length === 0 ? (
          <View style={styles.emptyState}>
            <Typography 
              variant="heading-xsmall" 
              color="contentSecondary" 
              style={styles.emptyTitle}
            >
              No AI analyses yet
            </Typography>
            <Typography 
              variant="paragraph-small" 
              color="contentTertiary"
              style={styles.emptySubtitle}
            >
              Complete workouts with AI progressive overload to see insights here
            </Typography>
          </View>
        ) : (
          <View style={styles.analysesList}>
            {analyses.map(renderAnalysisCard)}
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundPrimary'),
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
  
  scrollView: {
    flex: 1,
  },
  
  scrollContent: {
    flexGrow: 1,
    padding: 16,
  },
  
  analysesList: {
    gap: 12,
  },
  
  analysisCard: {
    flexDirection: 'row',
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    overflow: 'hidden',
  },
  
  statusIndicator: {
    width: 4,
  },
  
  cardContent: {
    flex: 1,
    padding: 16,
  },
  
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 8,
  },
  
  statusSection: {
    flex: 1,
  },
  
  progressBadge: {
    backgroundColor: getColor('positive'),
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
    marginLeft: 8,
  },
  
  summaryText: {
    marginBottom: 12,
    lineHeight: 20,
  },
  
  exerciseIndicators: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
    marginBottom: 8,
    flexWrap: 'wrap',
  },
  
  exerciseChip: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
    gap: 4,
  },
  
  improvedDot: {
    width: 6,
    height: 6,
    borderRadius: 3,
    backgroundColor: getColor('positive'),
  },
  
  metaInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  
  centerContent: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  
  loadingText: {
    marginTop: 12,
  },
  
  emptyState: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  
  emptyTitle: {
    marginBottom: 8,
    textAlign: 'center',
  },
  
  emptySubtitle: {
    textAlign: 'center',
    lineHeight: 20,
  },
  
  errorState: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  
  errorText: {
    marginTop: 8,
    textAlign: 'center',
  },
}); 