"use client"

import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { cartApi, productApi } from '../api'
import { useEffect } from 'react'

export interface CartItem {
  id: number
  cart_id: number
  product_id: number
  product_variant_id?: number
  quantity: number
  unit_price: number
  total_price: number
  created_at: string
  updated_at: string
  // Additional fields for display (will be fetched separately)
  name?: string
  sku?: string
  variant_name?: string
  variant_sku?: string
  image?: string
  maxQuantity?: number
}

export interface Cart {
  id: string
  session_id?: string
  user_id?: number
  items: CartItem[]
  subtotal: number
  tax_amount: number
  shipping_amount: number
  discount_amount: number
  total: number
  created_at: string
  updated_at: string
}

interface CartStore {
  cart: Cart | null
  items: CartItem[]
  isOpen: boolean
  isLoading: boolean
  error: string | null
  
  // Actions
  setCart: (cart: Cart | null) => void
  setItems: (items: CartItem[]) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  toggleCart: () => void
  setCartOpen: (open: boolean) => void
  
  // API Actions
  loadCart: () => Promise<void>
  addItem: (item: { product_id: number; product_variant_id?: number; quantity: number }) => Promise<void>
  removeItem: (itemId: number) => Promise<void>
  updateQuantity: (itemId: number, quantity: number) => Promise<void>
  clearCart: () => Promise<void>
  mergeCarts: (sourceCartId: string) => Promise<void>
  
  // Computed values
  getItemCount: () => number
  getTotalPrice: () => number
  getItemTotal: (id: number) => number
  
  // Helper functions
  fetchProductDetails: (items: CartItem[]) => Promise<CartItem[]>
  
  // Coupon functions
  applyCoupon: (couponCode: string) => Promise<void>
  removeCoupon: (couponCode: string) => Promise<void>
}

export const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      cart: null,
      items: [],
      isOpen: false,
      isLoading: false,
      error: null,

      setCart: (cart) => set({ cart, items: cart?.items || [] }),
      setItems: (items) => set({ items }),
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),
      toggleCart: () => set({ isOpen: !get().isOpen }),
      setCartOpen: (open) => set({ isOpen: open }),

      loadCart: async () => {
        try {
          set({ isLoading: true, error: null })
          
          // First get the cart metadata
          const cartResponse = await cartApi.getOrCreateCart()
          const cart = cartResponse.data
          
          // Then get the cart items via summary
          const summaryResponse = await cartApi.getCartSummary(cart.id.toString())
          const summary = summaryResponse.data
          
          
          // Fetch product details for items
          const itemsWithDetails = await get().fetchProductDetails(summary.items || [])
          
          // Update cart with items from summary
          const cartWithItems = {
            ...cart,
            items: itemsWithDetails,
            subtotal: summary.subtotal || 0,
            tax_amount: summary.tax_amount || 0,
            shipping_amount: summary.shipping_amount || 0,
            discount_amount: summary.discount_amount || 0,
            total: summary.total_amount || 0
          }
          
          set({ cart: cartWithItems, items: itemsWithDetails, isLoading: false })
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to load cart',
            isLoading: false 
          })
        }
      },

      addItem: async (itemData) => {
        try {
          set({ isLoading: true, error: null })
          const state = get()
          
          if (!state.cart) {
            await state.loadCart()
          }
          
          const cart = get().cart
          if (!cart) {
            throw new Error('No cart available')
          }
          
          await cartApi.addItemToCart(cart.id, itemData)
          
          // Reload cart to get updated data
          await get().loadCart()
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to add item to cart',
            isLoading: false 
          })
        }
      },

      removeItem: async (itemId) => {
        try {
          set({ isLoading: true, error: null })
          await cartApi.deleteCartItem(itemId.toString())
          
          // Update local state
          set(state => ({
            items: state.items.filter(item => item.id !== itemId),
            cart: state.cart ? {
              ...state.cart,
              items: state.cart.items.filter(item => item.id !== itemId)
            } : null,
            isLoading: false
          }))
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to remove item from cart',
            isLoading: false 
          })
        }
      },

      updateQuantity: async (itemId, quantity) => {
        try {
          if (quantity <= 0) {
            await get().removeItem(itemId)
            return
          }

          set({ isLoading: true, error: null })
          await cartApi.updateCartItem(itemId.toString(), { quantity })
          
          // Reload cart to get updated totals from backend
          await get().loadCart()
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to update item quantity',
            isLoading: false 
          })
        }
      },

      clearCart: async () => {
        try {
          set({ isLoading: true, error: null })
          const state = get()
          
          if (state.cart) {
            await cartApi.clearCartItems(state.cart.id)
          }
          
          set({ cart: null, items: [], isLoading: false })
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to clear cart',
            isLoading: false 
          })
        }
      },

      mergeCarts: async (sourceCartId) => {
        try {
          set({ isLoading: true, error: null })
          const state = get()
          
          if (!state.cart) {
            throw new Error('No target cart available for merging')
          }
          
          await cartApi.mergeCarts(state.cart.id, sourceCartId)
          
          // Reload cart to get merged items
          await get().loadCart()
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to merge carts',
            isLoading: false 
          })
        }
      },

      getItemCount: () => {
        return get().items.reduce((total, item) => total + item.quantity, 0)
      },

      getTotalPrice: () => {
        return get().cart?.total || get().items.reduce((total, item) => total + item.total_price, 0)
      },

      getItemTotal: (id) => {
        const item = get().items.find(i => i.id === id)
        return item ? item.total_price : 0
      },

      fetchProductDetails: async (items) => {
        const itemsWithDetails = await Promise.all(
          items.map(async (item) => {
            try {
              // Fetch product details
              const product = await productApi.getProduct(item.product_id)
              
              // Fetch variant details if variant exists
              let variant = null
              if (item.product_variant_id) {
                try {
                  variant = await productApi.getProductVariant(item.product_variant_id)
                } catch (error) {
                  console.warn(`Failed to fetch variant ${item.product_variant_id}:`, error)
                }
              }
              
              
              return {
                ...item,
                name: product.name,
                sku: variant ? variant.sku : product.sku,
                variant_name: variant?.name,
                variant_sku: variant?.sku,
                image: undefined, // Product doesn't have image_url field
                maxQuantity: product.max_quantity || 999
              }
            } catch (error) {
              console.warn(`Failed to fetch product details for product ${item.product_id}:`, error)
              return {
                ...item,
                name: `Product ${item.product_id}`,
                sku: `SKU-${item.product_id}`,
                maxQuantity: 999
              }
            }
          })
        )
        
        return itemsWithDetails
      },

      applyCoupon: async (couponCode) => {
        try {
          set({ isLoading: true, error: null })
          const state = get()
          
          if (!state.cart) {
            throw new Error('No cart available')
          }
          
          await cartApi.applyCouponToCart(state.cart.id, couponCode)
          
          // Reload cart to get updated data with discount
          await get().loadCart()
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to apply coupon',
            isLoading: false 
          })
        }
      },

      removeCoupon: async (couponCode) => {
        try {
          set({ isLoading: true, error: null })
          const state = get()
          
          if (!state.cart) {
            throw new Error('No cart available')
          }
          
          await cartApi.removeCouponFromCart(state.cart.id, couponCode)
          
          // Reload cart to get updated data without discount
          await get().loadCart()
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to remove coupon',
            isLoading: false 
          })
        }
      }
    }),
    {
      name: 'cart-storage',
      partialize: (state) => ({ 
        cart: state.cart,
        items: state.items 
      }),
      skipHydration: true, // Prevent hydration mismatch
    }
  )
)

// Hook to handle hydration
export const useHydrateCartStore = () => {
  useEffect(() => {
    useCartStore.persist.rehydrate()
  }, [])
}
