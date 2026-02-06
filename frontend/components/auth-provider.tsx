"use client"

import { useEffect, useRef } from 'react'
import { AuthProvider as AuthContextProvider, useAuth } from '@/lib/contexts/auth-context'
import { setAuthTokenHandlers } from '@/lib/apiFunctions/http-client'
import { useUserStore, type User } from '@/lib/stores/user-store'

/**
 * Authentication Provider Component
 * 
 * This component wraps the app with AuthContext to provide access token management.
 * It sets up the token handlers for the http-client to enable automatic refresh.
 * 
 * Key Features:
 * - Provides access token storage in React state (memory only)
 * - Connects http-client with auth context for Authorization header
 * - Enables automatic token refresh on 401 responses
 */

function AuthProviderInner({ children }: { children: React.ReactNode }) {
  const { accessToken, setAccessToken } = useAuth()
  const { setUser, clearUser } = useUserStore()

  const tokenRef = useRef(accessToken)

  // Update ref synchronously during render so it's available immediately to children
  tokenRef.current = accessToken

  useEffect(() => {
    // Set up token handlers for http-client
    // We use a ref so the closure always has access to the current token
    // without needing to re-register the handler when the token changes
    setAuthTokenHandlers(
      () => tokenRef.current,
      (token) => setAccessToken(token)
    )
  }, [setAccessToken])

  useEffect(() => {
    // Automatically refresh access token on app initialization
    // This restores the session if a valid refresh token cookie exists
    const initializeAuth = async () => {
      // Skip if we already have an access token
      if (accessToken) return

      try {
        // Import authApi dynamically to avoid circular dependencies
        const { authApi } = await import('@/lib/apiFunctions')

        // Try to refresh token using HTTP-only cookie
        const response = await authApi.refreshToken()
        const newAccessToken = response?.data?.access_token
        const refreshedUser = response?.data?.user

        if (newAccessToken) {
          setAccessToken(newAccessToken)
          console.log('Access token restored from refresh token')
        }

        if (refreshedUser) {
          const existingUser = useUserStore.getState().user
          const mergedUser: User = {
            id: refreshedUser.id ?? existingUser?.id ?? 0,
            username: refreshedUser.username ?? existingUser?.username ?? '',
            email: refreshedUser.email ?? existingUser?.email ?? '',
            firstName: refreshedUser.firstName ?? existingUser?.firstName ?? '',
            middleName: refreshedUser.middleName ?? existingUser?.middleName ?? '',
            lastName: refreshedUser.lastName ?? existingUser?.lastName ?? '',
            avatar: refreshedUser.avatar ?? existingUser?.avatar ?? '',
            role: refreshedUser.role ?? existingUser?.role ?? 'user',
            createdAt: refreshedUser.createdAt ?? existingUser?.createdAt ?? new Date().toISOString(),
            updatedAt: refreshedUser.updatedAt ?? existingUser?.updatedAt ?? new Date().toISOString(),
          }
          setUser(mergedUser)
        }
      } catch (error) {
        // Refresh failed - user needs to log in again
        // This is normal if no refresh token exists or it's expired
        console.log('No valid session found, user needs to log in')
        clearUser()
      }
    }

    initializeAuth()
  }, []) // Run once on mount, empty dependency array

  return <>{children}</>
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  return (
    <AuthContextProvider>
      <AuthProviderInner>{children}</AuthProviderInner>
    </AuthContextProvider>
  )
}
