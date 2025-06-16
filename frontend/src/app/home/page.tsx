"use client";

import { Typography, Button } from "@/components";
import Link from "next/link";
import { useState, useEffect } from "react";

interface Routine {
  id: string;
  user_id: string;
  name: string;
  description: string;
  exercises: any[];
  created_at: string;
  updated_at: string;
}

export default function HomePage() {
  const [routines, setRoutines] = useState<Routine[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // For now, using a placeholder user ID. In a real app, this would come from authentication
  const userId = "user_123";

  const fetchRoutines = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await fetch(`/api/routines?user_id=${userId}`);
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || "Failed to fetch routines");
      }
      
      if (data.success) {
        setRoutines(data.workouts || []);
      }
    } catch (err) {
      console.error("Error fetching routines:", err);
      setError(err instanceof Error ? err.message : "Failed to fetch routines");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRoutines();
  }, []);
  return (
    <main className="min-h-screen bg-light-gray-1">
      <div className="max-w-4xl mx-auto px-6 py-8">
        {/* Page Header */}
        <div className="mb-12">
          <Typography variant="heading-xxlarge" color="dark" className="mb-2">
            Home
          </Typography>
        </div>

        {/* Main Content */}
        <div className="space-y-6">
          {/* Quick Actions */}
          <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6">
            <div className="flex flex-col sm:flex-row gap-4">
              <Link href="/routines/create" className="flex-1">
                <Button 
                  variant="primary" 
                  size="large"
                  className="w-full justify-center"
                >
                  Create Routine
                </Button>
              </Link>
              <Button 
                variant="secondary" 
                size="large"
                className="flex-1 justify-center"
              >
                Browse exercises
              </Button>
            </div>
          </div>

          {/* Routines Section */}
          <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6">
            <div className="mb-6">
              <Typography variant="heading-large" color="dark" className="mb-2">
                Your Routines
              </Typography>
              <Typography variant="text-default" color="light">
                Manage and track your workout routines
              </Typography>
            </div>

            {/* Loading State */}
            {loading && (
              <div className="text-center py-8">
                <Typography variant="text-default" color="light">
                  Loading your routines...
                </Typography>
              </div>
            )}

            {/* Error State */}
            {error && (
              <div className="text-center py-8 space-y-4">
                <Typography variant="text-default" color="red">
                  {error}
                </Typography>
                <Button 
                  variant="secondary" 
                  size="default"
                  onClick={fetchRoutines}
                >
                  Try Again
                </Button>
              </div>
            )}

            {/* Empty State */}
            {!loading && !error && routines.length === 0 && (
              <div className="text-center py-8 space-y-4">
                <Typography variant="text-default" color="light">
                  No routines yet. Create your first routine to get started!
                </Typography>
                <Link href="/routines/create">
                  <Button variant="primary" size="default">
                    Create Your First Routine
                  </Button>
                </Link>
              </div>
            )}

            {/* Routines List */}
            {!loading && !error && routines.length > 0 && (
              <div className="space-y-4">
                {routines.map((routine) => (
                  <div 
                    key={routine.id}
                    className="border border-light-gray-3 rounded-lg p-4 hover:border-light-gray-2 transition-colors"
                  >
                    <div className="flex justify-between items-start mb-2">
                      <Typography variant="heading-small" color="dark">
                        {routine.name}
                      </Typography>
                      <Typography variant="text-small" color="light">
                        {new Date(routine.created_at).toLocaleDateString()}
                      </Typography>
                    </div>
                    {routine.description && (
                      <Typography variant="text-default" color="light" className="mb-3">
                        {routine.description}
                      </Typography>
                    )}
                    <div className="flex justify-between items-center">
                      <Typography variant="text-small" color="light">
                        {routine.exercises.length} exercise{routine.exercises.length !== 1 ? 's' : ''}
                      </Typography>
                      <div className="flex gap-2">
                        <Button variant="secondary" size="small">
                          View
                        </Button>
                        <Button variant="primary" size="small">
                          Start Workout
                        </Button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </main>
  );
} 