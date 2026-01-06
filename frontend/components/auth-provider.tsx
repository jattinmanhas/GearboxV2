"use client"

import { useEffect } from 'react'
import { AuthProvider as AuthContextProvider, useAuth } from '@/lib/contexts/auth-context'
import { setAuthTokenHandlers } from '@/lib/apiFunctions/http-client'

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

  useEffect(() => {
    // Set up token handlers for http-client
    setAuthTokenHandlers(
      () => accessToken,
      (token) => setAccessToken(token)
    )
  }, [accessToken, setAccessToken])

  return <>{children}</>
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  return (
    <AuthContextProvider>
      <AuthProviderInner>{children}</AuthProviderInner>
    </AuthContextProvider>
  )
}
