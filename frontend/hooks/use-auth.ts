"use client"

import { useEffect } from 'react'
import { useUserStore } from '@/lib/stores/user-store'

/**
 * Custom hook for authentication management
 * 
 * This hook provides a centralized way to handle authentication state
 * and automatically validates tokens when needed.
 * 
 * Features:
 * - Automatic token validation on mount
 * - Handles authentication state changes
 * - Provides loading states for auth operations
 */
export function useAuth() {
  const {
    user,
    isAuthenticated,
    isLoading,
    error,
    login,
    logout,
    clearUser,
    validateToken,
    checkAuthStatus,
    initializeAuth
  } = useUserStore()

  // Initialize authentication on mount
  useEffect(() => {
    initializeAuth()
  }, [initializeAuth])

  return {
    user,
    isAuthenticated,
    isLoading,
    error,
    login,
    logout,
    clearUser,
    validateToken,
    checkAuthStatus,
    initializeAuth
  }
}
