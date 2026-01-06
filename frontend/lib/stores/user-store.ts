import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { authApi } from '../apiFunctions/auth.api'
import { useCartStore } from './cart-store'

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
  updateProfile: (profileData: Partial<User>) => void
  debugAuthState: () => UserState
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

      login: async (userData: User) => {
        console.log("Zustand login called with:", userData)
        set({
          user: userData,
          isAuthenticated: true,
          isLoading: false,
          error: null,
        })
        console.log("Zustand state updated")

        // Handle cart merging when user logs in
        try {
          const cartStore = useCartStore.getState()

          // Get the current guest cart (if any)
          const guestCart = cartStore.cart

          if (guestCart && guestCart.session_id && !guestCart.user_id) {
            console.log("Guest cart found, attempting to merge with user cart")

            // Load/create user cart
            await cartStore.loadCart()

            // If we have a guest cart and a user cart, merge them
            if (cartStore.cart && cartStore.cart.id !== guestCart.id) {
              console.log("Merging guest cart with user cart")
              await cartStore.mergeCarts(guestCart.id)
            }
          } else {
            // No guest cart, just load the user's cart
            await cartStore.loadCart()
          }
        } catch (error) {
          console.error("Cart merging error:", error)
          // Don't fail login if cart merging fails
        }
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

      updateProfile: (profileData: Partial<User>) => {
        const currentUser = get().user
        if (currentUser) {
          set({
            user: { ...currentUser, ...profileData }
          })
        }
      },

      // Debug function to help troubleshoot auth issues
      debugAuthState: () => {
        const state = get()
        console.log('=== AUTH DEBUG INFO ===')
        console.log('isAuthenticated:', state.isAuthenticated)
        console.log('user:', state.user)
        console.log('isLoading:', state.isLoading)
        console.log('error:', state.error)
        console.log('localStorage user-storage:', typeof window !== 'undefined' ? localStorage.getItem('user-storage') : 'N/A')
        console.log('========================')
        return state
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
