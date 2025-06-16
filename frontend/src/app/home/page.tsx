"use client";

import { Typography, Button } from "@/components";
import Link from "next/link";
import { useState, useEffect } from "react";
import { useAuth } from "@/hooks/useAuth";

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
  const { backendUser, isAuthenticated, isLoading, logout } = useAuth();
  const [routines, setRoutines] = useState<Routine[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRoutines = async () => {
    if (!backendUser?.id) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      const response = await fetch(`/api/routines?user_id=${backendUser.id}`);
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
    if (isAuthenticated && backendUser?.id) {
      fetchRoutines();
    }
  }, [isAuthenticated, backendUser?.id]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-light-gray-1 flex items-center justify-center">
        <Typography variant="text-large" color="light">
          Loading...
        </Typography>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-light-gray-1 flex items-center justify-center">
        <Typography variant="text-large" color="light">
          Please sign in to continue
        </Typography>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-light-gray-1">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6 mb-6">
          <div className="flex justify-between items-center">
            <div>
              <Typography variant="heading-large" color="dark" className="mb-2">
                Welcome to GymLog
              </Typography>
              <Typography variant="text-default" color="light">
                Track your workouts and progress
              </Typography>
            </div>
            <Button
              variant="secondary"
              size="default"
              onClick={logout}
            >
              Sign Out
            </Button>
          </div>
        </div>

        {/* User Info */}
        <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6 mb-6">
          <Typography variant="heading-default" color="dark" className="mb-4">
            Your Profile
          </Typography>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <Typography variant="text-small" color="light" className="mb-1">
                Name
              </Typography>
              <Typography variant="text-default" color="dark">
                {backendUser?.name || 'Not provided'}
              </Typography>
            </div>
            
            <div>
              <Typography variant="text-small" color="light" className="mb-1">
                Email
              </Typography>
              <Typography variant="text-default" color="dark">
                {backendUser?.email || 'Not provided'}
              </Typography>
            </div>
            
            <div>
              <Typography variant="text-small" color="light" className="mb-1">
                User ID
              </Typography>
              <Typography variant="text-default" color="dark">
                {backendUser?.id || 'Not provided'}
              </Typography>
            </div>
            
            <div>
              <Typography variant="text-small" color="light" className="mb-1">
                Member Since
              </Typography>
              <Typography variant="text-default" color="dark">
                {backendUser?.createdAt 
                  ? new Date(backendUser.createdAt).toLocaleDateString()
                  : 'Not provided'
                }
              </Typography>
            </div>
          </div>

          {backendUser?.picture && (
            <div className="mt-4">
              <Typography variant="text-small" color="light" className="mb-2">
                Profile Picture
              </Typography>
              <img 
                src={backendUser.picture} 
                alt="Profile" 
                className="w-16 h-16 rounded-full border border-light-gray-3"
              />
            </div>
          )}
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6 mb-6">
          <Typography variant="heading-default" color="dark" className="mb-4">
            Quick Actions
          </Typography>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Link href="/routines/create">
              <Button variant="primary" size="large" className="w-full">
                Create Routine
              </Button>
            </Link>
            <Button variant="secondary" size="large" className="w-full">
              Start Workout
            </Button>
            <Button variant="secondary" size="large" className="w-full">
              Exercise Library
            </Button>
          </div>
        </div>

        {/* My Routines */}
        <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6 mb-6">
          <div className="flex justify-between items-center mb-4">
            <Typography variant="heading-default" color="dark">
              My Routines
            </Typography>
            <Link href="/routines/create">
              <Button variant="text" size="default">
                Create New
              </Button>
            </Link>
          </div>

          {loading ? (
            <div className="text-center py-8">
              <Typography variant="text-default" color="light">
                Loading routines...
              </Typography>
            </div>
          ) : error ? (
            <div className="text-center py-8">
              <Typography variant="text-default" color="red">
                {error}
              </Typography>
              <Button 
                variant="text" 
                size="small" 
                className="mt-2"
                onClick={fetchRoutines}
              >
                Try Again
              </Button>
            </div>
          ) : routines.length === 0 ? (
            <div className="text-center py-12">
              <Typography variant="text-default" color="light" className="mb-4">
                No routines created yet
              </Typography>
              <Typography variant="text-small" color="light" className="mb-6">
                Create your first workout routine to get started
              </Typography>
              <Link href="/routines/create">
                <Button variant="primary" size="default">
                  Create Your First Routine
                </Button>
              </Link>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {routines.map((routine) => (
                <div
                  key={routine.id}
                  className="p-4 border border-light-gray-3 rounded-lg hover:shadow-sm transition-shadow"
                >
                  <div className="mb-3">
                    <Typography variant="text-default" color="dark" className="font-medium mb-1">
                      {routine.name}
                    </Typography>
                    <Typography variant="text-small" color="light" className="mb-2">
                      {routine.description}
                    </Typography>
                    <Typography variant="text-small" color="light">
                      {routine.exercises?.length || 0} exercises • Created {new Date(routine.created_at).toLocaleDateString()}
                    </Typography>
                  </div>
                  
                  <div className="flex gap-2">
                    <Button variant="primary" size="small" className="flex-1">
                      Start Workout
                    </Button>
                    <Button variant="text" size="small">
                      Edit
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Status Info */}
        <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6">
          <Typography variant="heading-small" color="dark" className="mb-3">
            System Status
          </Typography>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <Typography variant="text-small" color="light" className="mb-1">
                Authentication Status
              </Typography>
              <Typography variant="text-default" color="green">
                ✅ Authenticated with Backend
              </Typography>
            </div>
            
            <div>
              <Typography variant="text-small" color="light" className="mb-1">
                Session Status
              </Typography>
              <Typography variant="text-default" color="green">
                ✅ Active Session
              </Typography>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
} 