"use client"

import React, { createContext, useContext, useState, useCallback } from 'react'

interface AuthContextType {
    accessToken: string | null
    setAccessToken: (token: string | null) => void
    clearAuth: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [accessToken, setAccessTokenState] = useState<string | null>(null)

    const setAccessToken = useCallback((token: string | null) => {
        setAccessTokenState(token)
    }, [])

    const clearAuth = useCallback(() => {
        setAccessTokenState(null)
    }, [])

    return (
        <AuthContext.Provider value={{ accessToken, setAccessToken, clearAuth }}>
            {children}
        </AuthContext.Provider>
    )
}

export function useAuth() {
    const context = useContext(AuthContext)
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}
