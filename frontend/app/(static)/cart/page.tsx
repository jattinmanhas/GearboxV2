"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"
import { 
  ShoppingCart, 
  Plus, 
  Minus, 
  Trash2, 
  X,
  Package,
  CreditCard,
  ArrowRight,
  Loader2,
  ArrowLeft,
  Tag
} from "lucide-react"
import { useCartStore, useHydrateCartStore } from "@/lib/stores/cart-store"
import { formatCurrency } from "@/lib/currency"
import Link from "next/link"
import Image from "next/image"

export default function CartPage() {
  const [isMounted, setIsMounted] = useState(false)
  
  // Hydrate the store on client side
  useHydrateCartStore()
  
  const { 
    cart, 
    items, 
    appliedCoupons,
    availableCoupons,
    isLoading, 
    error, 
    loadCart, 
    updateQuantity, 
    removeItem, 
    clearCart,
    getItemCount,
    getTotalPrice,
    applyCoupon,
    removeCoupon,
    loadAvailableCoupons
  } = useCartStore()
  
  const [couponCode, setCouponCode] = useState("")
  const [isApplyingCoupon, setIsApplyingCoupon] = useState(false)
  const [couponMessage, setCouponMessage] = useState("")
  const [showAvailableCoupons, setShowAvailableCoupons] = useState(false)

  const itemCount = getItemCount()
  const totalPrice = getTotalPrice()

  // Prevent hydration mismatch by only rendering on client
  useEffect(() => {
    setIsMounted(true)
  }, [])

  useEffect(() => {
    if (!cart && isMounted) {
      loadCart()
    }
  }, [cart, loadCart, isMounted])

  useEffect(() => {
    if (isMounted && availableCoupons.length === 0) {
      loadAvailableCoupons().catch((error) => {
        console.error("Failed to load available coupons:", error)
      })
    }
  }, [isMounted, availableCoupons.length, loadAvailableCoupons])

  const handleQuantityChange = async (itemId: number, newQuantity: number) => {
    await updateQuantity(itemId, newQuantity)
  }

  const handleRemoveItem = async (itemId: number) => {
    await removeItem(itemId)
  }

  const handleClearCart = async () => {
    if (confirm('Are you sure you want to clear your cart?')) {
      await clearCart()
    }
  }

  const handleApplyCoupon = async () => {
    if (!couponCode.trim() || !cart) return
    
    setIsApplyingCoupon(true)
    setCouponMessage("")
    try {
      await applyCoupon(couponCode.trim())
      setCouponCode("") // Clear the input on success
      setCouponMessage("Coupon applied successfully!")
      setTimeout(() => setCouponMessage(""), 3000) // Clear after 3 seconds
    } catch (error) {
      console.error("Failed to apply coupon:", error)
      setCouponMessage(error instanceof Error ? error.message : "Failed to apply coupon")
    } finally {
      setIsApplyingCoupon(false)
    }
  }

  const handleRemoveCoupon = async (couponCode: string) => {
    setCouponMessage("") // Clear message when removing coupon
    await removeCoupon(couponCode)
  }

  const handleSelectCoupon = async (coupon: any) => {
    if (!coupon || !coupon.code) {
      setCouponMessage("Invalid coupon selected")
      return
    }
    
    setCouponCode(coupon.code)
    setShowAvailableCoupons(false)
    // Auto-apply after selection
    setTimeout(async () => {
      if (!cart) {
        setCouponMessage("Cart not available")
        return
      }
      
      setIsApplyingCoupon(true)
      setCouponMessage("")
      
      try {
        await applyCoupon(coupon.code)
        setCouponCode("") // Clear the input on success
        setCouponMessage("Coupon applied successfully!")
        setTimeout(() => setCouponMessage(""), 3000)
      } catch (error) {
        console.error("Failed to apply coupon:", error)
        const errorMessage = error instanceof Error 
          ? error.message 
          : "Failed to apply coupon"
        setCouponMessage(errorMessage)
        // Don't clear the coupon code on error so user can try again
      } finally {
        setIsApplyingCoupon(false)
      }
    }, 100)
  }

  if (!isMounted) {
    return null
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-7xl">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-4 mb-4">
          <Button variant="ghost" size="sm" asChild>
            <Link href="/shop">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Continue Shopping
            </Link>
          </Button>
        </div>
        <h1 className="text-3xl font-bold flex items-center gap-3">
          <ShoppingCart className="h-8 w-8" />
          Shopping Cart
          {itemCount > 0 && (
            <Badge variant="secondary" className="text-lg px-3 py-1">
              {itemCount} {itemCount === 1 ? 'item' : 'items'}
            </Badge>
          )}
        </h1>
      </div>

      {isLoading && !cart ? (
        <div className="flex items-center justify-center py-12">
          <div className="flex items-center gap-2 text-muted-foreground">
            <Loader2 className="h-6 w-6 animate-spin" />
            <span>Loading cart...</span>
          </div>
        </div>
      ) : error ? (
        <Card>
          <CardContent className="py-12">
            <div className="text-center">
              <div className="text-destructive mb-4">Failed to load cart</div>
              <Button variant="outline" onClick={loadCart}>
                Try Again
              </Button>
            </div>
          </CardContent>
        </Card>
      ) : itemCount === 0 ? (
        <Card>
          <CardContent className="py-12">
            <div className="flex flex-col items-center justify-center text-center">
              <Package className="h-16 w-16 text-muted-foreground mb-4" />
              <h2 className="text-2xl font-semibold mb-2">Your cart is empty</h2>
              <p className="text-muted-foreground mb-6">
                Add some items to get started
              </p>
              <Button asChild size="lg">
                <Link href="/shop">
                  Browse Products
                </Link>
              </Button>
            </div>
          </CardContent>
        </Card>
      ) : (
        <div className="grid lg:grid-cols-3 gap-8">
          {/* Cart Items - Left Column */}
          <div className="lg:col-span-2 space-y-4">
            {/* Cart Items List */}
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle>Cart Items</CardTitle>
                  <Button 
                    variant="ghost" 
                    size="sm"
                    onClick={handleClearCart}
                    disabled={isLoading}
                    className="text-destructive hover:text-destructive"
                  >
                    <Trash2 className="h-4 w-4 mr-2" />
                    Clear Cart
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {items.map((item) => (
                    <div key={item.id} className="flex gap-4 p-4 border rounded-lg hover:bg-muted/50 transition-colors">
                      {/* Product Image */}
                      <div className="w-24 h-24 bg-muted rounded-md flex items-center justify-center flex-shrink-0 overflow-hidden">
                        {item.image ? (
                          <Image 
                            src={item.image} 
                            alt={item.name || 'Product'} 
                            width={96}
                            height={96}
                            className="w-full h-full object-cover"
                          />
                        ) : (
                          <Package className="h-8 w-8 text-muted-foreground" />
                        )}
                      </div>
                      
                      {/* Product Info */}
                      <div className="flex-1 min-w-0">
                        <h3 className="font-semibold text-lg mb-1">
                          {item.name || `Product ${item.product_id}`}
                        </h3>
                        <div className="flex flex-col gap-1 text-sm text-muted-foreground mb-3">
                          <span>SKU: {item.sku || `SKU-${item.product_id}`}</span>
                          {item.variant_name && (
                            <span className="text-blue-600 dark:text-blue-400">
                              Variant: {item.variant_name}
                            </span>
                          )}
                          <span className="text-base font-medium text-foreground">
                            {formatCurrency(item.unit_price)} each
                          </span>
                        </div>
                        
                        {/* Quantity Controls */}
                        <div className="flex items-center gap-3">
                          <div className="flex items-center border rounded-md">
                            <Button
                              variant="ghost"
                              size="sm"
                              className="h-9 w-9 p-0"
                              onClick={() => handleQuantityChange(item.id, item.quantity - 1)}
                              disabled={isLoading}
                            >
                              <Minus className="h-4 w-4" />
                            </Button>
                            <span className="px-4 py-2 text-base min-w-[3rem] text-center font-medium">
                              {item.quantity}
                            </span>
                            <Button
                              variant="ghost"
                              size="sm"
                              className="h-9 w-9 p-0"
                              onClick={() => handleQuantityChange(item.id, item.quantity + 1)}
                              disabled={isLoading || (item.maxQuantity ? item.quantity >= item.maxQuantity : false)}
                            >
                              <Plus className="h-4 w-4" />
                            </Button>
                          </div>
                          
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleRemoveItem(item.id)}
                            disabled={isLoading}
                            className="text-destructive hover:text-destructive"
                          >
                            <Trash2 className="h-4 w-4 mr-2" />
                            Remove
                          </Button>
                        </div>
                      </div>
                      
                      {/* Price */}
                      <div className="text-right">
                        <div className="text-2xl font-bold">
                          {formatCurrency(item.total_price)}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Order Summary - Right Column */}
          <div className="space-y-4">
            {/* Coupon Section */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Tag className="h-5 w-5" />
                  Coupons & Discounts
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Applied Coupons */}
                {appliedCoupons.length > 0 && (
                  <div className="space-y-2">
                    <div className="text-sm font-medium text-muted-foreground">Applied Coupon</div>
                    {appliedCoupons.map((coupon) => (
                      <div key={coupon.id} className="flex items-center justify-between bg-green-50 dark:bg-green-900/20 px-4 py-3 rounded-lg border border-green-200 dark:border-green-800">
                        <div>
                          <div className="font-semibold text-green-700 dark:text-green-300">
                            {coupon.coupon_code}
                          </div>
                          <div className="text-sm text-green-600 dark:text-green-400">
                            Discount: {formatCurrency(coupon.discount_amount)}
                          </div>
                        </div>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleRemoveCoupon(coupon.coupon_code)}
                          className="text-red-500 hover:text-red-700"
                        >
                          <X className="h-4 w-4" />
                        </Button>
                      </div>
                    ))}
                  </div>
                )}

                {/* Manual Coupon Entry */}
                <div className="space-y-2">
                  <div className="text-sm font-medium text-muted-foreground">Enter Coupon Code</div>
                  <div className="flex gap-2">
                    <Input
                      placeholder="Enter code"
                      value={couponCode}
                      onChange={(e) => setCouponCode(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && handleApplyCoupon()}
                      className="flex-1"
                    />
                    <Button
                      onClick={handleApplyCoupon}
                      disabled={!couponCode.trim() || isApplyingCoupon}
                    >
                      {isApplyingCoupon ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        "Apply"
                      )}
                    </Button>
                  </div>
                  {couponMessage && (
                    <div className={`text-sm ${couponMessage.includes('successfully') ? 'text-green-600' : 'text-red-600'}`}>
                      {couponMessage}
                    </div>
                  )}
                </div>

                {/* Available Coupons */}
                <div className="space-y-2">
                  <Button
                    variant="outline"
                    className="w-full"
                    onClick={() => setShowAvailableCoupons(!showAvailableCoupons)}
                  >
                    {showAvailableCoupons ? "Hide Available Coupons" : "View Available Coupons"}
                  </Button>
                  
                  {showAvailableCoupons && (
                    <div className="space-y-2 max-h-64 overflow-y-auto">
                      {appliedCoupons.length > 0 && (
                        <div className="text-xs text-muted-foreground px-3 py-2 bg-blue-50 dark:bg-blue-950/20 rounded">
                          Only one coupon can be applied per cart. Remove the current coupon to apply a different one.
                        </div>
                      )}
                      {availableCoupons.length === 0 ? (
                        <div className="text-sm text-muted-foreground text-center py-4">
                          No coupons available
                        </div>
                      ) : (
                        availableCoupons
                          .filter((coupon) => coupon && coupon.code && coupon.id)
                          .map((coupon) => (
                            <div key={coupon.id} className="flex items-center justify-between px-3 py-3 border rounded hover:bg-muted/50">
                              <div className="flex-1">
                                <div className="font-semibold">{coupon.code}</div>
                                <div className="text-sm text-muted-foreground">
                                  {(() => {
                                    try {
                                      const discountValue = coupon.value || 0;
                                      const discountType = coupon.type || 'fixed_amount';
                                      
                                      if (isNaN(discountValue) || discountValue === null || discountValue === undefined || discountValue === 0) {
                                        return 'Discount available';
                                      }
                                      
                                      return discountType === 'percentage' 
                                        ? `${discountValue}% off`
                                        : discountType === 'free_shipping'
                                        ? 'Free shipping'
                                        : `${formatCurrency(discountValue)} off`;
                                    } catch (err) {
                                      console.error("Error formatting coupon discount:", err)
                                      return 'Discount available'
                                    }
                                  })()}
                                </div>
                              </div>
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handleSelectCoupon(coupon)}
                                disabled={appliedCoupons.length > 0 || isApplyingCoupon}
                              >
                                {appliedCoupons.some(ac => ac.coupon_code === coupon.code) ? "Applied" : 
                                 appliedCoupons.length > 0 ? "Remove existing first" : "Apply"}
                              </Button>
                            </div>
                          ))
                      )}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Order Summary */}
            <Card>
              <CardHeader>
                <CardTitle>Order Summary</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span>Subtotal</span>
                    <span className="font-medium">{formatCurrency(cart?.subtotal || totalPrice)}</span>
                  </div>
                  {cart && cart.tax_amount > 0 && (
                    <div className="flex justify-between">
                      <span>Tax</span>
                      <span className="font-medium">{formatCurrency(cart.tax_amount)}</span>
                    </div>
                  )}
                  {cart && cart.shipping_amount > 0 && (
                    <div className="flex justify-between">
                      <span>Shipping</span>
                      <span className="font-medium">{formatCurrency(cart.shipping_amount)}</span>
                    </div>
                  )}
                  {cart && cart.discount_amount > 0 && (
                    <div className="flex justify-between text-green-600">
                      <span>Discount</span>
                      <span className="font-medium">-{formatCurrency(cart.discount_amount)}</span>
                    </div>
                  )}
                </div>
                
                <Separator />
                
                <div className="flex justify-between text-lg font-bold">
                  <span>Total</span>
                  <span>{formatCurrency(cart?.total || totalPrice)}</span>
                </div>

                <Button 
                  className="w-full" 
                  size="lg"
                  asChild
                >
                  <Link href="/checkout">
                    <CreditCard className="h-5 w-5 mr-2" />
                    Proceed to Checkout
                    <ArrowRight className="h-5 w-5 ml-2" />
                  </Link>
                </Button>
                
                <Button 
                  variant="outline" 
                  className="w-full"
                  asChild
                >
                  <Link href="/shop">
                    Continue Shopping
                  </Link>
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      )}
    </div>
  )
}

