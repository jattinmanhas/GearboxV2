"use client"

import { useEffect } from 'react'
import { useUserStore } from '@/lib/stores/user-store'

/**
 * Authentication Provider Component
 * 
 * This component handles the initialization of authentication state when the app loads.
 * It validates persisted user data and clears it if the tokens are expired.
 * 
 * Key Features:
 * - Validates persisted authentication state on app startup
 * - Automatically clears user state if tokens are expired
 * - Prevents showing user data when authentication is invalid
 * - Runs only once when the component mounts
 */
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const initializeAuth = useUserStore((state) => state.initializeAuth)

  useEffect(() => {
    // Initialize authentication validation when the component mounts
    initializeAuth()
  }, [initializeAuth])

  return <>{children}</>
}
