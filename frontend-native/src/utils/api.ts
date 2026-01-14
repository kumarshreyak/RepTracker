/**
 * API utility for making authenticated requests with Clerk tokens
 */

const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';

export interface ApiRequestOptions extends RequestInit {
  token?: string;
}

// Global token getter function - set this from your app's auth context
let globalTokenGetter: (() => Promise<string | null>) | null = null;

/**
 * Set the global token getter function
 * Call this once in your app initialization with Clerk's getToken function
 */
export function setGlobalTokenGetter(getter: () => Promise<string | null>) {
  globalTokenGetter = getter;
}

/**
 * Make an authenticated API request
 * @param endpoint - API endpoint (e.g., '/api/workouts')
 * @param options - Request options including Clerk token
 * @returns Response data
 */
export async function apiRequest<T = any>(
  endpoint: string,
  options: ApiRequestOptions = {}
): Promise<T> {
  const { token, ...fetchOptions } = options;

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...fetchOptions.headers,
  };

  // Get token from options, or fall back to global token getter
  let authToken = token;
  if (!authToken && globalTokenGetter) {
    authToken = await globalTokenGetter();
  }

  // Add Authorization header if token is available
  if (authToken) {
    headers['Authorization'] = `Bearer ${authToken}`;
  }

  const url = `${API_BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...fetchOptions,
    headers,
  });

  // Handle non-JSON responses
  const contentType = response.headers.get('content-type');
  if (!contentType || !contentType.includes('application/json')) {
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    return {} as T;
  }

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || data.message || `HTTP ${response.status}: ${response.statusText}`);
  }

  return data;
}

/**
 * GET request helper
 */
export async function apiGet<T = any>(
  endpoint: string,
  token?: string,
  options: RequestInit = {}
): Promise<T> {
  return apiRequest<T>(endpoint, {
    ...options,
    method: 'GET',
    token,
  });
}

/**
 * POST request helper
 */
export async function apiPost<T = any>(
  endpoint: string,
  data: any,
  token?: string,
  options: RequestInit = {}
): Promise<T> {
  return apiRequest<T>(endpoint, {
    ...options,
    method: 'POST',
    body: JSON.stringify(data),
    token,
  });
}

/**
 * PUT request helper
 */
export async function apiPut<T = any>(
  endpoint: string,
  data: any,
  token?: string,
  options: RequestInit = {}
): Promise<T> {
  return apiRequest<T>(endpoint, {
    ...options,
    method: 'PUT',
    body: JSON.stringify(data),
    token,
  });
}

/**
 * DELETE request helper
 */
export async function apiDelete<T = any>(
  endpoint: string,
  token?: string,
  options: RequestInit = {}
): Promise<T> {
  return apiRequest<T>(endpoint, {
    ...options,
    method: 'DELETE',
    token,
  });
}

