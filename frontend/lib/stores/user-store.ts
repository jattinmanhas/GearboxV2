import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { authApi } from '../api'

export interface User {
  id: number
  username: string
  email: string
  firstName: string
  middleName?: string
  lastName: string
  avatar?: string
  role?: string
  createdAt: string
  updatedAt: string
}

interface UserState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
  
  // Actions
  setUser: (user: User) => void
  clearUser: () => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  login: (userData: User) => void
  logout: () => void
}

export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      // Initial state
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      // Actions
      setUser: (user: User) => {
        set({
          user,
          isAuthenticated: true,
          error: null,
        })
      },

      clearUser: () => {
        set({
          user: null,
          isAuthenticated: false,
          error: null,
        })
      },

      setLoading: (loading: boolean) => {
        set({ isLoading: loading })
      },

      setError: (error: string | null) => {
        set({ error })
      },

      login: (userData: User) => {
        console.log("Zustand login called with:", userData)
        set({
          user: userData,
          isAuthenticated: true,
          isLoading: false,
          error: null,
        })
        console.log("Zustand state updated")
      },

      logout: async () => {
        try {
          // Call logout API to invalidate server-side session and clear cookies
          await authApi.logout()
          console.log("Logout API called successfully")
        } catch (error) {
          console.error("Logout API error:", error)
          // Continue with local logout even if API call fails
        } finally {
          // Clear user data from Zustand store
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
            error: null,
          })
          console.log("User data cleared from Zustand store")
        }
      },
    }),
    {
      name: 'user-storage', // unique name for localStorage key
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)

// Selectors for easier access to specific parts of the state
export const useUser = () => useUserStore((state) => state.user)
export const useIsAuthenticated = () => useUserStore((state) => state.isAuthenticated)
export const useIsLoading = () => useUserStore((state) => state.isLoading)
export const useUserError = () => useUserStore((state) => state.error)
