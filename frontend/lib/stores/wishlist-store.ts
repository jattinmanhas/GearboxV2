"use client"

import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { wishlistApi } from '../api'

export interface WishlistItem {
  id: number
  product_id: number
  product_name: string
  product_sku: string
  price: number
  image?: string
  added_at: string
}

export interface Wishlist {
  id: number
  name: string
  description?: string
  is_public: boolean
  created_at: string
  updated_at: string
  items: WishlistItem[]
}

interface WishlistStore {
  wishlists: Wishlist[]
  currentWishlist: Wishlist | null
  isLoading: boolean
  error: string | null
  
  // Actions
  setWishlists: (wishlists: Wishlist[]) => void
  setCurrentWishlist: (wishlist: Wishlist | null) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  
  // API Actions
  loadWishlists: () => Promise<void>
  createWishlist: (data: { name: string; description?: string; is_public?: boolean }) => Promise<void>
  updateWishlist: (id: number, data: { name?: string; description?: string; is_public?: boolean }) => Promise<void>
  deleteWishlist: (id: number) => Promise<void>
  
  // Wishlist Items
  addItemToWishlist: (wishlistId: number, productId: number) => Promise<void>
  removeItemFromWishlist: (wishlistId: number, itemId: number) => Promise<void>
  moveItemToCart: (itemId: number) => Promise<void>
  
  // Utility
  isProductInWishlist: (productId: number, wishlistId?: number) => boolean
  getWishlistItemCount: (wishlistId?: number) => number
}

export const useWishlistStore = create<WishlistStore>()(
  persist(
    (set, get) => ({
      wishlists: [],
      currentWishlist: null,
      isLoading: false,
      error: null,

      setWishlists: (wishlists) => set({ wishlists }),
      setCurrentWishlist: (wishlist) => set({ currentWishlist: wishlist }),
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),

      loadWishlists: async () => {
        try {
          set({ isLoading: true, error: null })
          const response = await wishlistApi.getWishlists()
          // Handle different response structures
          const wishlistsData = response.data?.wishlists || response.wishlists || response.data || []
          set({ wishlists: wishlistsData })
        } catch (error) {
          console.error('Wishlist load error:', error)
          // If it's a 500 error, it might be because the backend expects user_id
          // For now, set empty array and show a helpful message
          if (error instanceof Error && error.message.includes('500')) {
            set({ 
              error: 'Wishlist service is not fully configured. Please try again later.',
              wishlists: []
            })
          } else {
            set({ error: error instanceof Error ? error.message : 'Failed to load wishlists' })
          }
        } finally {
          set({ isLoading: false })
        }
      },

      createWishlist: async (data) => {
        try {
          set({ isLoading: true, error: null })
          const response = await wishlistApi.createWishlist(data)
          const newWishlist = response.data
          set(state => ({ 
            wishlists: [...state.wishlists, newWishlist],
            isLoading: false 
          }))
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to create wishlist',
            isLoading: false 
          })
        }
      },

      updateWishlist: async (id, data) => {
        try {
          set({ isLoading: true, error: null })
          const response = await wishlistApi.updateWishlist(id.toString(), data)
          const updatedWishlist = response.data
          set(state => ({
            wishlists: state.wishlists.map(w => w.id === id ? updatedWishlist : w),
            currentWishlist: state.currentWishlist?.id === id ? updatedWishlist : state.currentWishlist,
            isLoading: false
          }))
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to update wishlist',
            isLoading: false 
          })
        }
      },

      deleteWishlist: async (id) => {
        try {
          set({ isLoading: true, error: null })
          await wishlistApi.deleteWishlist(id.toString())
          set(state => ({
            wishlists: state.wishlists.filter(w => w.id !== id),
            currentWishlist: state.currentWishlist?.id === id ? null : state.currentWishlist,
            isLoading: false
          }))
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to delete wishlist',
            isLoading: false 
          })
        }
      },

      addItemToWishlist: async (wishlistId, productId) => {
        try {
          set({ isLoading: true, error: null })
          await wishlistApi.addItemToWishlist(wishlistId.toString(), { product_id: productId })
          
          // Reload wishlists to get updated data
          await get().loadWishlists()
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to add item to wishlist',
            isLoading: false 
          })
        }
      },

      removeItemFromWishlist: async (wishlistId, itemId) => {
        try {
          set({ isLoading: true, error: null })
          await wishlistApi.deleteWishlistItem(itemId.toString())
          
          // Update local state
          set(state => ({
            wishlists: state.wishlists.map(w => 
              w.id === wishlistId 
                ? { ...w, items: w.items.filter(item => item.id !== itemId) }
                : w
            ),
            currentWishlist: state.currentWishlist?.id === wishlistId 
              ? { ...state.currentWishlist, items: state.currentWishlist.items.filter(item => item.id !== itemId) }
              : state.currentWishlist,
            isLoading: false
          }))
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to remove item from wishlist',
            isLoading: false 
          })
        }
      },

      moveItemToCart: async (itemId) => {
        try {
          set({ isLoading: true, error: null })
          await wishlistApi.moveItemToCart(itemId.toString())
          
          // Reload wishlists to get updated data
          await get().loadWishlists()
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to move item to cart',
            isLoading: false 
          })
        }
      },

      isProductInWishlist: (productId, wishlistId) => {
        const state = get()
        const targetWishlist = wishlistId 
          ? state.wishlists.find(w => w.id === wishlistId)
          : state.currentWishlist
        
        return targetWishlist?.items.some(item => item.product_id === productId) || false
      },

      getWishlistItemCount: (wishlistId) => {
        const state = get()
        const targetWishlist = wishlistId 
          ? state.wishlists.find(w => w.id === wishlistId)
          : state.currentWishlist
        
        return targetWishlist?.items.length || 0
      },
    }),
    {
      name: 'wishlist-storage',
      partialize: (state) => ({
        wishlists: state.wishlists,
        currentWishlist: state.currentWishlist,
      }),
    }
  )
)
