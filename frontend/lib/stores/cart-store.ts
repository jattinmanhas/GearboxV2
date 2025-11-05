"use client"

import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { cartApi, productApi, couponApi } from '../api'
import { useEffect } from 'react'
import { showSuccess, showError, showLoading, updateLoading, NotificationMessages } from '../notifications'

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

export interface AppliedCoupon {
  id: number
  cart_id: number
  coupon_code: string
  discount_amount: number
  created_at: string
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
  applied_coupons?: AppliedCoupon[]
}

interface CartStore {
  cart: Cart | null
  items: CartItem[]
  appliedCoupons: AppliedCoupon[]
  availableCoupons: any[]
  isOpen: boolean
  isLoading: boolean
  error: string | null
  
  // Actions
  setCart: (cart: Cart | null) => void
  setItems: (items: CartItem[]) => void
  setAppliedCoupons: (coupons: AppliedCoupon[]) => void
  setAvailableCoupons: (coupons: any[]) => void
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
  recalculateCartTotals: () => void
  
  // Coupon functions
  applyCoupon: (couponCode: string) => Promise<void>
  removeCoupon: (couponCode: string) => Promise<void>
  loadAppliedCoupons: () => Promise<void>
  loadAvailableCoupons: () => Promise<void>
  recalculatePercentageDiscounts: () => Promise<void>
  loadCartSilently: () => Promise<void>
}

export const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      cart: null,
      items: [],
      appliedCoupons: [],
      availableCoupons: [],
      isOpen: false,
      isLoading: false,
      error: null,

      setCart: (cart) => set({ cart, items: cart?.items || [] }),
      setItems: (items) => set({ items }),
      setAppliedCoupons: (coupons) => set({ appliedCoupons: coupons }),
      setAvailableCoupons: (coupons) => set({ availableCoupons: coupons }),
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
          
          // Get applied coupons
          const couponsResponse = await cartApi.getCartCoupons(cart.id.toString())
          const appliedCoupons = couponsResponse.data || []
          
          
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
            total: summary.total_amount || 0,
            applied_coupons: appliedCoupons
          }
          
          set({ 
            cart: cartWithItems, 
            items: itemsWithDetails, 
            appliedCoupons: appliedCoupons,
            isLoading: false 
          })
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to load cart',
            isLoading: false 
          })
        }
      },

      addItem: async (itemData) => {
        const state = get()
        
        // Ensure cart exists
        if (!state.cart) {
          await state.loadCart()
        }
        
        const cart = get().cart
        if (!cart) {
          showError(NotificationMessages.cart.addError)
          return
        }
        
        // Try to fetch product price for optimistic update
        let estimatedPrice = 0
        try {
          const product = await productApi.getProduct(itemData.product_id)
          if (itemData.product_variant_id) {
            try {
              const variant = await productApi.getProductVariant(itemData.product_variant_id)
              estimatedPrice = variant.price || product.price
            } catch {
              estimatedPrice = product.price
            }
          } else {
            estimatedPrice = product.price
          }
        } catch (error) {
          console.warn('Failed to fetch price for optimistic update:', error)
        }
        
        // Create optimistic item for immediate UI update
        const optimisticItem: CartItem = {
          id: Date.now(), // Temporary ID
          cart_id: typeof cart.id === 'string' ? parseInt(cart.id) : cart.id,
          product_id: itemData.product_id,
          product_variant_id: itemData.product_variant_id,
          quantity: itemData.quantity,
          unit_price: estimatedPrice,
          total_price: estimatedPrice * itemData.quantity,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }
        
        // Optimistically update the UI immediately
        set(state => ({
          items: [...state.items, optimisticItem],
          cart: state.cart ? {
            ...state.cart,
            items: [...state.cart.items, optimisticItem]
          } : null
        }))
        
        // Recalculate totals immediately
        get().recalculateCartTotals()
        
        // Show success notification immediately
        showSuccess(NotificationMessages.cart.itemAdded)
        
        // Make API call in background
        try {
          await cartApi.addItemToCart(cart.id, itemData)
          
          // Reload cart silently to get accurate data
          await get().loadCartSilently()
          
          // If there are applied coupons, we need to recalculate percentage discounts
          const currentState = get()
          if (currentState.appliedCoupons.length > 0) {
            await get().recalculatePercentageDiscounts()
          }
          
          set({ isLoading: false })
        } catch (error) {
          // Revert optimistic update on error
          set(state => ({
            items: state.items.filter(item => item.id !== optimisticItem.id),
            cart: state.cart ? {
              ...state.cart,
              items: state.cart.items.filter(item => item.id !== optimisticItem.id)
            } : null,
            error: error instanceof Error ? error.message : 'Failed to add item to cart',
            isLoading: false
          }))
          
          // Show error notification
          showError(NotificationMessages.cart.addError)
        }
      },

      removeItem: async (itemId) => {
        const state = get()
        
        // Store the item for potential rollback
        const itemToRemove = state.items.find(item => item.id === itemId)
        if (!itemToRemove) return
        
        // Optimistically remove item from UI
        set(state => ({
          items: state.items.filter(item => item.id !== itemId),
          cart: state.cart ? {
            ...state.cart,
            items: state.cart.items.filter(item => item.id !== itemId)
          } : null
        }))
        
        // Recalculate totals immediately
        get().recalculateCartTotals()
        
        // Show success notification immediately
        showSuccess(NotificationMessages.cart.itemRemoved)
        
        // Make API call in background
        try {
          await cartApi.deleteCartItem(itemId.toString())
          
          // Reload cart silently to get accurate totals
          await get().loadCartSilently()
          
          set({ isLoading: false })
        } catch (error) {
          // Revert optimistic update on error
          set(state => ({
            items: [...state.items, itemToRemove],
            cart: state.cart ? {
              ...state.cart,
              items: [...state.cart.items, itemToRemove]
            } : null,
            error: error instanceof Error ? error.message : 'Failed to remove item from cart',
            isLoading: false
          }))
          
          // Show error notification
          showError(NotificationMessages.cart.removeError)
        }
      },

      updateQuantity: async (itemId, quantity) => {
        if (quantity <= 0) {
          await get().removeItem(itemId)
          return
        }
        
        const state = get()
        
        // Store the old item for potential rollback
        const oldItem = state.items.find(item => item.id === itemId)
        if (!oldItem) return
        
        // Optimistically update quantity in UI
        set(state => ({
          items: state.items.map(item => 
            item.id === itemId 
              ? { ...item, quantity, total_price: item.unit_price * quantity }
              : item
          ),
          cart: state.cart ? {
            ...state.cart,
            items: state.cart.items.map(item => 
              item.id === itemId 
                ? { ...item, quantity, total_price: item.unit_price * quantity }
                : item
            )
          } : null
        }))
        
        // Recalculate totals immediately
        get().recalculateCartTotals()
        
        // Make API call in background
        try {
          await cartApi.updateCartItem(itemId.toString(), { quantity })
          
          // Reload cart silently to get accurate totals
          await get().loadCartSilently()
          
          // If there are applied coupons, we need to recalculate percentage discounts
          const currentState = get()
          if (currentState.appliedCoupons.length > 0) {
            await get().recalculatePercentageDiscounts()
          }
          
          set({ isLoading: false })
        } catch (error) {
          // Revert optimistic update on error
          set(state => ({
            items: state.items.map(item => item.id === itemId ? oldItem : item),
            cart: state.cart ? {
              ...state.cart,
              items: state.cart.items.map(item => item.id === itemId ? oldItem : item)
            } : null,
            error: error instanceof Error ? error.message : 'Failed to update item quantity',
            isLoading: false
          }))
          
          showError('Failed to update quantity')
        }
      },

      clearCart: async () => {
        const state = get()
        
        // Store cart for potential rollback
        const oldCart = state.cart
        const oldItems = state.items
        const oldAppliedCoupons = state.appliedCoupons
        
        // Optimistically clear cart in UI
        set({ items: [], appliedCoupons: [] })
        
        // Recalculate totals immediately (will set to 0)
        get().recalculateCartTotals()
        
        // Show success notification immediately
        showSuccess(NotificationMessages.cart.cartCleared)
        
        // Make API call in background
        try {
          if (state.cart) {
            await cartApi.clearCartItems(state.cart.id)
          }
          
          // Reload cart silently to ensure sync
          await get().loadCartSilently()
          
          set({ isLoading: false })
        } catch (error) {
          // Revert optimistic update on error
          set({
            cart: oldCart,
            items: oldItems,
            appliedCoupons: oldAppliedCoupons,
            error: error instanceof Error ? error.message : 'Failed to clear cart',
            isLoading: false
          })
          
          showError('Failed to clear cart')
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
                image: product.images && product.images.length > 0 ? product.images[0].url : undefined,
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

      recalculateCartTotals: () => {
        const state = get()
        if (!state.cart || !state.items) return

        // Calculate subtotal from items
        const subtotal = state.items.reduce((total, item) => total + item.total_price, 0)

        // Use discount from cart (backend calculated) - don't recalculate locally
        // The backend handles maximum discount limits correctly
        const discount = state.cart.discount_amount || 0

        // Tax and shipping remain the same or recalculate if you have logic
        const taxAmount = state.cart.tax_amount || 0
        const shippingAmount = state.cart.shipping_amount || 0

        // Calculate total
        const total = subtotal - discount + taxAmount + shippingAmount

        // Update cart with new totals (keep discount from backend)
        set(state => ({
          cart: state.cart ? {
            ...state.cart,
            subtotal,
            // Keep discount_amount from backend - don't override
            discount_amount: discount,
            total
          } : null
        }))
      },

      applyCoupon: async (couponCode) => {
        const state = get()
        
        if (!state.cart) {
          showError('No cart available')
          return
        }
        
        // Store old coupons for potential rollback
        const oldAppliedCoupons = state.appliedCoupons
        
        // Optimistically show coupon as applied (we don't know the discount yet)
        const optimisticCoupon: AppliedCoupon = {
          id: Date.now(),
          cart_id: typeof state.cart.id === 'string' ? parseInt(state.cart.id) : state.cart.id,
          coupon_code: couponCode,
          discount_amount: 0, // Will be updated from API
          created_at: new Date().toISOString()
        }
        
        // If there's already a coupon, replace it
        set({ appliedCoupons: [optimisticCoupon] })
        
        // Recalculate totals immediately (will use coupon info from availableCoupons)
        get().recalculateCartTotals()
        
        // Show loading notification
        const loadingToast = showLoading('Applying coupon...')
        
        // Make API call in background
        try {
          // If there's already a coupon applied, remove it first
          if (oldAppliedCoupons.length > 0) {
            const existingCoupon = oldAppliedCoupons[0]
            try {
              await cartApi.removeCouponFromCart(state.cart.id, existingCoupon.coupon_code)
            } catch (removeError) {
              console.warn('Failed to remove existing coupon:', removeError)
            }
          }
          
          await cartApi.applyCouponToCart(state.cart.id, couponCode)
          
          // Reload cart silently to get accurate discount
          await get().loadCartSilently()
          
          updateLoading(loadingToast, 'Coupon applied successfully!', 'success')
          set({ isLoading: false })
        } catch (error) {
          // Revert optimistic update on error
          let errorMessage = 'Failed to apply coupon'
          
          if (error instanceof Error) {
            // Error message is already parsed by handleResponse, use it directly
            errorMessage = error.message
          }
          
          set({ 
            appliedCoupons: oldAppliedCoupons,
            error: errorMessage,
            isLoading: false 
          })
          
          updateLoading(loadingToast, 'Failed to apply coupon', 'error', {
            description: errorMessage
          })
        }
      },

      removeCoupon: async (couponCode) => {
        const state = get()
        
        if (!state.cart) {
          showError('No cart available')
          return
        }
        
        // Store old coupons for potential rollback
        const oldAppliedCoupons = state.appliedCoupons
        
        // Optimistically remove coupon from UI
        set({ 
          appliedCoupons: state.appliedCoupons.filter(c => c.coupon_code !== couponCode)
        })
        
        // Recalculate totals immediately (will remove discount)
        get().recalculateCartTotals()
        
        // Show success notification immediately
        showSuccess('Coupon removed')
        
        // Make API call in background
        try {
          await cartApi.removeCouponFromCart(state.cart.id, couponCode)
          
          // Reload cart silently to get updated totals
          await get().loadCartSilently()
          
          set({ isLoading: false })
        } catch (error) {
          // Revert optimistic update on error
          set({ 
            appliedCoupons: oldAppliedCoupons,
            error: error instanceof Error ? error.message : 'Failed to remove coupon',
            isLoading: false 
          })
          
          showError('Failed to remove coupon')
        }
      },

      loadAppliedCoupons: async () => {
        try {
          const state = get()
          if (!state.cart) {
            return
          }
          
          const response = await cartApi.getCartCoupons(state.cart.id)
          set({ appliedCoupons: response.data || [] })
        } catch (error) {
          console.error('Failed to load applied coupons:', error)
        }
      },

      loadAvailableCoupons: async () => {
        try {
          const response = await couponApi.getCoupons({
            page: 1,
            limit: 20,
            status: 'active'
          })
          
          // Filter out expired coupons
          const now = new Date()
          const validCoupons = (response.data?.coupons || []).filter((coupon: any) => {
            // If coupon has no expiration date, it's valid
            if (!coupon.expires_at) {
              return true
            }
            
            // Check if coupon has expired
            const expiresAt = new Date(coupon.expires_at)
            return expiresAt > now
          })
          
          set({ availableCoupons: validCoupons })
        } catch (error) {
          console.error('Failed to load available coupons:', error)
          // Set empty array on error to prevent crashes
          set({ availableCoupons: [] })
        }
      },

      recalculatePercentageDiscounts: async () => {
        try {
          const state = get()
          if (!state.cart || state.appliedCoupons.length === 0) {
            return
          }

          // Get the current applied coupon codes
          const appliedCouponCodes = state.appliedCoupons.map(coupon => coupon.coupon_code)
          
          if (appliedCouponCodes.length === 0) {
            return
          }

          // For percentage-based coupons, we need to reapply them to get correct discounts
          // This is necessary because the backend doesn't automatically recalculate percentage discounts
          for (const couponCode of appliedCouponCodes) {
            try {
              // Remove the coupon first
              await cartApi.removeCouponFromCart(state.cart.id, couponCode)
              
              // Reapply the coupon (this will recalculate the discount with new cart total)
              await cartApi.applyCouponToCart(state.cart.id, couponCode)
            } catch (error) {
              console.warn(`Failed to recalculate coupon ${couponCode}:`, error)
            }
          }

          // Use silent reload to avoid double loading states
          await get().loadCartSilently()
          
        } catch (error) {
          console.error('Failed to recalculate percentage discounts:', error)
          // Fallback: silent reload
          await get().loadCartSilently()
        }
      },

      loadCartSilently: async () => {
        try {
          // First get the cart metadata
          const cartResponse = await cartApi.getOrCreateCart()
          const cart = cartResponse.data
          
          // Then get the cart items via summary
          const summaryResponse = await cartApi.getCartSummary(cart.id.toString())
          const summary = summaryResponse.data
          
          // Get applied coupons
          const couponsResponse = await cartApi.getCartCoupons(cart.id.toString())
          const appliedCoupons = couponsResponse.data || []
          
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
            total: summary.total_amount || 0,
            applied_coupons: appliedCoupons
          }
          
          set({ 
            cart: cartWithItems, 
            items: itemsWithDetails, 
            appliedCoupons: appliedCoupons
            // Note: Not setting isLoading to false here to avoid UI flicker
          })
        } catch (error) {
          console.error('Failed to load cart silently:', error)
        }
      }
    }),
    {
      name: 'cart-storage',
      partialize: (state) => ({ 
        cart: state.cart,
        items: state.items,
        appliedCoupons: state.appliedCoupons
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
