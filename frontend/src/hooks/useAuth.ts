"use client"

import { useSession, signOut } from "next-auth/react"
import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"

interface BackendUser {
  id: string
  email: string
  name: string
  googleId: string
  picture: string
  createdAt: string
  updatedAt: string
}

export function useAuth() {
  const { data: session, status } = useSession()
  const router = useRouter()
  const [backendUser, setBackendUser] = useState<BackendUser | null>(null)
  const [sessionToken, setSessionToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    if (status === "loading") return

    if (status === "unauthenticated") {
      setIsLoading(false)
      return
    }

    // Store session token if available from NextAuth session
    if (session?.user?.sessionToken) {
      localStorage.setItem("gymlog_session_token", session.user.sessionToken)
      setSessionToken(session.user.sessionToken)
      
      // Validate session with backend
      validateSession(session.user.sessionToken)
    } else {
      // Get session token from localStorage
      const token = localStorage.getItem("gymlog_session_token")
      if (token) {
        setSessionToken(token)
        validateSession(token)
      } else {
        setIsLoading(false)
      }
    }
  }, [status, session])

  const validateSession = async (token: string) => {
    try {
      const response = await fetch("/api/auth/validate", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (response.ok) {
        const user = await response.json()
        setBackendUser(user)
      } else if (response.status === 401) {
        // Session expired, logout
        handleSessionExpired()
      }
    } catch (error) {
      console.error("Error validating session:", error)
      handleSessionExpired()
    } finally {
      setIsLoading(false)
    }
  }

  const handleSessionExpired = () => {
    localStorage.removeItem("gymlog_session_token")
    setSessionToken(null)
    setBackendUser(null)
    signOut({ callbackUrl: "/auth/signin" })
  }

  const logout = async () => {
    const token = localStorage.getItem("gymlog_session_token")
    if (token) {
      try {
        await fetch("/api/auth/logout", {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        })
      } catch (error) {
        console.error("Error during logout:", error)
      }
    }

    localStorage.removeItem("gymlog_session_token")
    setSessionToken(null)
    setBackendUser(null)
    signOut({ callbackUrl: "/auth/signin" })
  }

  const isAuthenticated = status === "authenticated" && backendUser && sessionToken

  return {
    session,
    backendUser,
    sessionToken,
    isAuthenticated,
    isLoading: isLoading || status === "loading",
    logout,
    validateSession: () => sessionToken && validateSession(sessionToken),
  }
} 