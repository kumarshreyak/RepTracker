import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  RefreshControl,
  ActivityIndicator,
  Text,
  TouchableOpacity,
} from 'react-native';
import { Typography, Button } from '../../src/components';
import { getColor } from '../../src/components/Colors';
import { useAuth } from '../../src/hooks/useAuth';
import { useFocusEffect, useRouter } from 'expo-router';
import { 
  WorkoutInsight, 
  StoredSuggestedWorkout
} from '../../src/types/suggestion';

// Renaming for backward compatibility with existing code
type Insight = WorkoutInsight;

export default function InsightsTab() {
  const { user } = useAuth();
  const router = useRouter();
  const [insights, setInsights] = useState<Insight[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [generating, setGenerating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [storedSuggestions, setStoredSuggestions] = useState<StoredSuggestedWorkout[]>([]);
  const [generatingSuggestions, setGeneratingSuggestions] = useState(false);
  const [loadingSuggestions, setLoadingSuggestions] = useState(false);
  const [processingWorkoutId, setProcessingWorkoutId] = useState<string | null>(null);

  const fetchInsights = async (isRefreshing = false) => {
    if (!user?.id) return;

    try {
      if (!isRefreshing) setLoading(true);
      setError(null);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(
        `${API_BASE_URL}/api/users/${user.id}/insights?limit=10`
      );
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || "Failed to fetch insights");
      }
      
      setInsights(data.insights || []);
    } catch (err) {
      console.error("Error fetching insights:", err);
      setError(err instanceof Error ? err.message : "Failed to fetch insights");
    } finally {
      setLoading(false);
      if (isRefreshing) setRefreshing(false);
    }
  };

  const fetchStoredSuggestions = async () => {
    if (!user?.id) return;

    try {
      setLoadingSuggestions(true);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(
        `${API_BASE_URL}/api/users/${user.id}/suggestions/stored?pageSize=5`
      );
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || "Failed to fetch suggestions");
      }
      
      setStoredSuggestions(data.suggestedWorkouts || []);
    } catch (err) {
      console.error("Error fetching suggestions:", err);
    } finally {
      setLoadingSuggestions(false);
    }
  };

  useEffect(() => {
    if (user?.id) {
      fetchInsights();
      fetchStoredSuggestions();
    }
  }, [user?.id]);

  // Refresh insights when screen comes into focus
  useFocusEffect(
    React.useCallback(() => {
      if (user?.id) {
        fetchInsights();
        fetchStoredSuggestions();
      }
    }, [user?.id])
  );

  const onRefresh = () => {
    setRefreshing(true);
    Promise.all([
      fetchInsights(true),
      fetchStoredSuggestions()
    ]).finally(() => setRefreshing(false));
  };

  const generateInsights = async () => {
    if (!user?.id || generating) {
      console.log("generateInsights early return:", { userId: user?.id, generating });
      return;
    }

    console.log("Starting insight generation for user:", user.id);
    
    try {
      setGenerating(true);
      setError(null);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const url = `${API_BASE_URL}/api/users/${user.id}/insights`;
      
      console.log("Making POST request to:", url);
      
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      
      console.log("Response status:", response.status);
      console.log("Response headers:", response.headers);
      
      let data;
      try {
        data = await response.json();
        console.log("Response data:", data);
      } catch (parseError) {
        console.error("Failed to parse response as JSON:", parseError);
        throw new Error("Invalid response from server");
      }
      
      if (!response.ok) {
        console.error("API error:", { status: response.status, data });
        throw new Error(data.error || `Failed to generate insights (${response.status})`);
      }
      
      console.log("Insights generated successfully, fetching updated list...");
      
      // After generating, fetch the updated insights list
      await fetchInsights();
      
      console.log("Insight generation completed successfully");
      
    } catch (err) {
      console.error("Error generating insights:", err);
      
      // More detailed error messaging
      let errorMessage = "Failed to generate insights";
      if (err instanceof Error) {
        if (err.message.includes("fetch")) {
          errorMessage = "Network error: Could not connect to server";
        } else if (err.message.includes("Invalid response")) {
          errorMessage = "Server error: Invalid response received";
        } else {
          errorMessage = err.message;
        }
      }
      
      setError(errorMessage);
    } finally {
      console.log("Setting generating to false");
      setGenerating(false);
    }
  };

  const generateWorkoutSuggestions = async () => {
    if (!user?.id || generatingSuggestions) return;

    try {
      setGeneratingSuggestions(true);
      setError(null);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(
        `${API_BASE_URL}/api/users/${user.id}/suggestions/workouts`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            daysToAnalyze: 14,
            maxSuggestions: 3
          })
        }
      );
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || "Failed to generate suggestions");
      }
      
      // Refresh the stored suggestions list
      await fetchStoredSuggestions();
      
    } catch (err) {
      console.error("Error generating suggestions:", err);
      setError(err instanceof Error ? err.message : "Failed to generate suggestions");
    } finally {
      setGeneratingSuggestions(false);
    }
  };

  const acceptSuggestion = async (suggestionId: string) => {
    if (!user?.id || processingWorkoutId) return;

    try {
      setProcessingWorkoutId(suggestionId);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      
      const response = await fetch(
        `${API_BASE_URL}/api/users/${user.id}/suggestions/${suggestionId}/confirm`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            accept: true
          })
        }
      );
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || "Failed to accept suggestion");
      }
      
      // Refresh suggestions after accepting
      await fetchStoredSuggestions();
      
    } catch (err) {
      console.error("Error accepting suggestion:", err);
      setError("Failed to apply suggested changes");
    } finally {
      setProcessingWorkoutId(null);
    }
  };

  const rejectSuggestion = async (suggestionId: string) => {
    if (!user?.id || processingWorkoutId) return;

    try {
      setProcessingWorkoutId(suggestionId);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      
      const response = await fetch(
        `${API_BASE_URL}/api/users/${user.id}/suggestions/${suggestionId}/confirm`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            accept: false
          })
        }
      );
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || "Failed to reject suggestion");
      }
      
      // Refresh suggestions after rejecting
      await fetchStoredSuggestions();
      
    } catch (err) {
      console.error("Error rejecting suggestion:", err);
      setError("Failed to reject suggestion");
    } finally {
      setProcessingWorkoutId(null);
    }
  };

  const navigateToSuggestionDetail = (suggestion: StoredSuggestedWorkout) => {
    router.push({
      pathname: '/suggestion-detail',
      params: {
        suggestion: JSON.stringify(suggestion)
      }
    });
  };

  const getTypeColor = (type: string): string => {
    switch (type.toLowerCase()) {
      case 'progress':
        return getColor('positive');
      case 'recommendation':
        return getColor('accent');
      case 'warning':
        return getColor('warning');
      case 'achievement':
        return getColor('backgroundAccent');
      default:
        return getColor('contentTertiary');
    }
  };

  const getPriorityColor = (priority: number): string => {
    if (priority >= 4) return getColor('accent');
    if (priority >= 3) return getColor('contentSecondary');
    return getColor('contentTertiary');
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

  const formatType = (type: string): string => {
    return type.charAt(0).toUpperCase() + type.slice(1);
  };

  const renderInsightCard = (insight: Insight) => (
    <View key={insight.id} style={styles.insightCard}>
      {/* Priority indicator bar */}
      <View 
        style={[
          styles.priorityBar, 
          { backgroundColor: getPriorityColor(insight.priority) }
        ]} 
      />
      
      <View style={styles.cardContent}>
        {/* Type badge and time */}
        <View style={styles.metaRow}>
          <View 
            style={[
              styles.typeBadge, 
              { backgroundColor: getTypeColor(insight.type) + '20' }
            ]}
          >
            <Typography 
              variant="label-xsmall" 
              color="contentPrimary"
              style={{ color: getTypeColor(insight.type) }}
            >
              {formatType(insight.type)}
            </Typography>
          </View>
          
          <Typography variant="paragraph-xsmall" color="contentTertiary">
            {formatTimeAgo(insight.createdAt)}
          </Typography>
        </View>
        
        {/* Main insight text */}
        <Typography 
          variant="paragraph-medium" 
          color="contentPrimary" 
          style={styles.insightText}
        >
          {insight.insight}
        </Typography>
        
        {/* Based on info */}
        <Typography variant="paragraph-xsmall" color="contentSecondary">
          Based on {insight.basedOn.toLowerCase()}
        </Typography>
      </View>
    </View>
  );

  const renderSuggestionCard = (suggestion: StoredSuggestedWorkout) => {
    // Skip if not pending
    if (suggestion.status && suggestion.status !== 'pending') return null;

    const changeCount = suggestion.changes.length;
    const changeTypes = [...new Set(suggestion.changes.map(c => c.type))];
    
    return (
      <View key={suggestion.id} style={styles.suggestedWorkoutItem}>
        {/* Priority indicator */}
        <View 
          style={[
            styles.workoutPriorityIndicator,
            { backgroundColor: suggestion.priority >= 4 ? getColor('accent') : getColor('borderOpaque') }
          ]} 
        />
        
        <View style={styles.suggestedWorkoutContent}>
          {/* Touchable content area */}
          <TouchableOpacity 
            onPress={() => navigateToSuggestionDetail(suggestion)}
            activeOpacity={0.7}
            style={styles.suggestionContentArea}
          >
            {/* Header with workout name and change count */}
            <View style={styles.workoutHeader}>
              <View style={{ flex: 1 }}>
                <Typography variant="label-medium" color="contentPrimary">
                  {suggestion.name}
                </Typography>
              </View>
              <View style={styles.changeCountBadge}>
                <Typography variant="label-xsmall" color="contentOnColor">
                  {changeCount} {changeCount === 1 ? 'change' : 'changes'}
                </Typography>
              </View>
            </View>
            
            {/* Change type indicators */}
            <View style={styles.changeTypesRow}>
              {changeTypes.slice(0, 3).map((type, typeIndex) => (
                <View key={typeIndex} style={styles.changeTypeBadge}>
                  <Typography variant="label-xsmall" color="contentSecondary">
                    {type}
                  </Typography>
                </View>
              ))}
              {changeTypes.length > 3 && (
                <Typography variant="label-xsmall" color="contentTertiary">
                  +{changeTypes.length - 3}
                </Typography>
              )}
            </View>
            
            {/* Brief reasoning */}
            <View style={styles.workoutReasoning}>
              <Text 
                style={{
                  fontFamily: 'Uber Move Text',
                  fontSize: 12,
                  lineHeight: 20,
                  fontWeight: '400',
                  color: getColor('contentSecondary'),
                }}
                numberOfLines={2}
                ellipsizeMode="tail"
              >
                {suggestion.overallReasoning}
              </Text>
            </View>
          </TouchableOpacity>
          
          {/* Action buttons */}
          <View style={styles.workoutActions}>
            <Button
              variant="secondary"
              size="small"
              onPress={() => rejectSuggestion(suggestion.id)}
              style={styles.suggestionButton}
            >
              Dismiss
            </Button>
            <Button
              variant="primary"
              size="small"
              onPress={() => acceptSuggestion(suggestion.id)}
              disabled={processingWorkoutId === suggestion.id}
              style={styles.suggestionButton}
            >
              {processingWorkoutId === suggestion.id ? 'Applying...' : 'Apply'}
            </Button>
          </View>
        </View>
      </View>
    );
  };

  if (loading && insights.length === 0) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.header}>
          <Typography variant="heading-small" color="contentPrimary">
            Insights
          </Typography>
        </View>
        <View style={styles.centerContent}>
          <ActivityIndicator size="large" color={getColor('accent')} />
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Typography variant="heading-small" color="contentPrimary">
          Insights
        </Typography>
        <Button
          variant="primary"
          size="small"
          onPress={generateInsights}
          disabled={generating}
        >
          {generating ? 'Generating...' : 'Generate'}
        </Button>
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
        {/* Workout Suggestions Section */}
        {(storedSuggestions.length > 0 || !loadingSuggestions) && (
          <View style={styles.suggestionsSection}>
            <View style={styles.sectionHeader}>
              <Typography variant="label-medium" color="contentPrimary">
                Workout Improvements
              </Typography>
              <Button
                variant="text"
                size="small"
                onPress={generateWorkoutSuggestions}
                disabled={generatingSuggestions}
              >
                {generatingSuggestions ? 'Analyzing...' : 'New Analysis'}
              </Button>
            </View>
            
            {loadingSuggestions ? (
              <View style={styles.loadingContainer}>
                <ActivityIndicator size="small" color={getColor('accent')} />
              </View>
            ) : storedSuggestions.length === 0 ? (
              <View style={styles.emptySuggestions}>
                <Typography variant="paragraph-small" color="contentTertiary">
                  No suggestions yet. Complete more workouts for personalized recommendations.
                </Typography>
              </View>
            ) : (
              <View style={styles.suggestionsList}>
                {storedSuggestions.map(renderSuggestionCard)}
              </View>
            )}
          </View>
        )}
        
        {/* Divider */}
        {storedSuggestions.length > 0 && insights.length > 0 && (
          <View style={styles.sectionDivider} />
        )}

        {/* Insights Section */}
        <View style={styles.insightsSection}>
          <Typography variant="label-medium" color="contentPrimary" style={styles.sectionTitle}>
            Recent Activity
          </Typography>
          
          {error ? (
            <View style={styles.centerContent}>
              <Typography variant="paragraph-small" color="contentNegative">
                {error}
              </Typography>
            </View>
          ) : insights.length === 0 ? (
            <View style={styles.emptyState}>
              <Typography 
                variant="heading-xsmall" 
                color="contentSecondary" 
                style={styles.emptyTitle}
              >
                No insights yet
              </Typography>
              <Typography 
                variant="paragraph-small" 
                color="contentTertiary"
                style={{ textAlign: 'center' }}
              >
                Complete more workouts to get personalized insights
              </Typography>
            </View>
          ) : (
            <View style={styles.insightsList}>
              {insights.map(renderInsightCard)}
            </View>
          )}
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
  
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 16,
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderOpaque'),
  },
  
  scrollView: {
    flex: 1,
  },
  
  scrollContent: {
    flexGrow: 1,
  },
  
  // Suggestions Section Styles
  suggestionsSection: {
    paddingTop: 16,
  },
  
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    marginBottom: 12,
  },
  
  suggestionsList: {
    paddingHorizontal: 16,
  },
  
  suggestedWorkoutItem: {
    flexDirection: 'row',
    backgroundColor: getColor('backgroundSecondary'),
    borderRadius: 6,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    overflow: 'hidden',
    marginBottom: 16,
  },
  
  workoutPriorityIndicator: {
    width: 4,
  },
  
  suggestedWorkoutContent: {
    flex: 1,
    padding: 12,
  },
  
  suggestionContentArea: {
    flex: 1,
  },
  
  workoutHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 8,
  },
  
  workoutReasoning: {
    marginBottom: 12,
    lineHeight: 20,
  },
  
  workoutActions: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: 8,
  },
  
  changeCountBadge: {
    backgroundColor: getColor('accent'),
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
    marginLeft: 8,
  },
  
  changeTypesRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 8,
    gap: 6,
  },
  
  changeTypeBadge: {
    backgroundColor: getColor('backgroundTertiary'),
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 4,
  },
  
  suggestionButton: {
    minWidth: 80,
  },
  
  emptySuggestions: {
    padding: 24,
    alignItems: 'center',
  },
  
  loadingContainer: {
    padding: 32,
    alignItems: 'center',
  },
  
  sectionDivider: {
    height: 1,
    backgroundColor: getColor('borderOpaque'),
    marginHorizontal: 16,
    marginVertical: 24,
  },
  
  // Insights Section Styles
  insightsSection: {
    flex: 1,
  },
  
  sectionTitle: {
    paddingHorizontal: 16,
    marginBottom: 12,
  },
  
  insightsList: {
    padding: 16,
  },
  
  insightCard: {
    flexDirection: 'row',
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    overflow: 'hidden',
  },
  
  priorityBar: {
    width: 4,
  },
  
  cardContent: {
    flex: 1,
    padding: 16,
  },
  
  metaRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  
  typeBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
  },
  
  insightText: {
    marginBottom: 8,
    lineHeight: 22,
  },
  
  centerContent: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  
  emptyState: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  
  emptyTitle: {
    marginBottom: 8,
  },
}); 