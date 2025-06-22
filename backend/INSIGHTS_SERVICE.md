# GymLog Insights Service Documentation

## Overview

The GymLog Insights Service uses Google's Gemini AI to analyze workout metrics and generate personalized, actionable insights. It processes data from the `workout_metrics` and `user_metrics` collections to provide intelligent feedback about training progress, volume management, recovery needs, muscle balance, and injury risk.

## Features

- **AI-Powered Analysis**: Uses Gemini 2.5 Flash to generate natural language insights
- **Multiple Insight Types**: Progress, volume, recovery, balance, and risk insights
- **Priority-Based**: Insights are prioritized (1-5) based on importance
- **Context-Aware**: Analyzes the last 7 days of workout data and user metrics
- **Actionable Recommendations**: Provides specific, one-sentence actionable insights

## Setup

### 1. Get a Gemini API Key

1. Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Sign in with your Google account
3. Click "Create API Key"
4. Copy the generated API key

### 2. Set Environment Variable

Add your Gemini API key to your environment:

```bash
export GEMINI_API_KEY="your_api_key_here"
```

Or add it to your `.env` file:

```
GEMINI_API_KEY=your_api_key_here
```

### 3. Build and Run

The insights service is automatically initialized when you start the backend server:

```bash
cd backend
go run cmd/server/main.go
```

## API Endpoints

### Generate Insights

```http
POST /api/users/{userId}/insights
```

Generates new insights based on recent workout data.

**Response:**
```json
{
  "insights": [
    {
      "id": "insight_id",
      "userId": "user_id",
      "type": "progress",
      "insight": "Your squat strength has increased by 8.5% this week, showing excellent progressive overload - keep pushing at this sustainable rate!",
      "basedOn": "Last 7 days",
      "priority": 4,
      "createdAt": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### Get Recent Insights

```http
GET /api/users/{userId}/insights?limit=5
```

Retrieves the most recent insights for a user, ordered by priority and recency.

**Query Parameters:**
- `limit` (optional): Number of insights to return (default: 5)

**Response:**
```json
{
  "insights": [
    {
      "id": "insight_id",
      "userId": "user_id",
      "type": "recovery",
      "insight": "Your training stress balance is -12.5, indicating accumulated fatigue - consider a deload week or extra rest day.",
      "basedOn": "Recovery metrics",
      "priority": 5,
      "createdAt": "2024-01-15T10:30:00Z"
    }
  ]
}
```

## Insight Types

### 1. Progress Insights (`progress`)
- Analyzes progressive overload index
- Tracks week-over-week progress rates
- Detects training plateaus
- Reviews 1RM estimates and RPE trends

Example: "Your bench press has plateaued for 3 weeks at 225lbs - try adding volume or changing rep ranges to break through."

### 2. Volume Insights (`volume`)
- Compares current volume to MEV/MAV/MRV landmarks
- Analyzes muscle group distribution
- Reviews hard sets and effective reps
- Provides volume recommendations

Example: "Your chest volume (18 sets) is approaching MRV (20 sets) - maintain current volume to avoid overtraining."

### 3. Recovery Insights (`recovery`)
- Evaluates training stress balance
- Calculates optimal training frequency
- Monitors form/freshness index
- Suggests recovery strategies

Example: "With 36 hours optimal recovery time between sessions, training every other day would maximize your adaptation."

### 4. Balance Insights (`balance`)
- Analyzes push/pull ratios
- Detects muscle imbalances
- Reviews antagonist muscle ratios
- Suggests corrective strategies

Example: "Your push:pull ratio of 1.8 indicates push dominance - add more rowing and pulling exercises to prevent shoulder issues."

### 5. Risk Insights (`risk`)
- Monitors injury risk score
- Detects training load spikes
- Identifies asymmetries
- Provides prevention strategies

Example: "Load spike detected: this week's volume is 165% of your 4-week average - reduce intensity to prevent overuse injuries."

## How It Works

1. **Data Collection**: The service fetches the last 7 days of workout metrics and current user metrics
2. **Context Preparation**: Relevant metrics are formatted into structured prompts for each insight type
3. **AI Analysis**: Gemini analyzes the data and generates a single-sentence insight
4. **Priority Assignment**: Insights are prioritized based on metric thresholds (e.g., high injury risk = priority 5)
5. **Storage**: Insights are saved to the `workout_insights` collection for future reference

## Priority Levels

- **5 (Critical)**: Immediate attention needed (injury risk, severe imbalances)
- **4 (High)**: Important for progress (plateaus, suboptimal training)
- **3 (Medium)**: General recommendations
- **2 (Low)**: Minor optimizations
- **1 (Info)**: Informational only

## Best Practices

1. **Generate insights after completing workouts** to get the most relevant feedback
2. **Review high-priority insights first** and take action on critical recommendations
3. **Track insight history** to see how recommendations change over time
4. **Use insights for program adjustments** rather than making drastic changes

## Troubleshooting

### "GEMINI_API_KEY environment variable not set"
- Ensure you've set the environment variable correctly
- Check that your `.env` file is being loaded

### "Failed to generate insights"
- Verify your API key is valid
- Check that you have sufficient workout data (at least 1 workout in the last 7 days)
- Ensure MongoDB collections (`workout_metrics`, `user_metrics`) contain data

### Rate Limiting
- Gemini API has rate limits (check current limits at [Google AI documentation](https://ai.google.dev/pricing))
- Consider caching insights rather than generating on every request

## Example Integration

```javascript
// Frontend example - Generate insights after workout completion
async function generateInsights(userId) {
  try {
    const response = await fetch(`/api/users/${userId}/insights`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    const data = await response.json();
    
    // Display high-priority insights prominently
    const criticalInsights = data.insights.filter(i => i.priority >= 4);
    criticalInsights.forEach(insight => {
      showNotification(insight.insight, insight.type);
    });
  } catch (error) {
    console.error('Failed to generate insights:', error);
  }
}

// Get recent insights on dashboard load
async function loadDashboardInsights(userId) {
  const response = await fetch(`/api/users/${userId}/insights?limit=3`);
  const data = await response.json();
  
  // Display insights in dashboard widget
  renderInsightsWidget(data.insights);
}
```

## Future Enhancements

- **Personalized Models**: Fine-tune insights based on user goals and training history
- **Comparative Analysis**: Compare performance to similar users
- **Predictive Insights**: Forecast future performance and potential issues
- **Real-time Alerts**: Push notifications for critical insights
- **Visual Insights**: Generate charts and graphs alongside text insights
- **Multi-language Support**: Generate insights in user's preferred language 