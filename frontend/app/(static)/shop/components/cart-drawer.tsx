"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"
import { 
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"
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
  Eye,
  ChevronDown,
  ChevronUp
} from "lucide-react"
import { useCartStore, useHydrateCartStore } from "@/lib/stores/cart-store"
import { formatCurrency } from "@/lib/currency"
import Link from "next/link"

interface CartDrawerProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CartDrawer({ open, onOpenChange }: CartDrawerProps) {
  const [isMounted, setIsMounted] = useState(false)
  const [showAllItems, setShowAllItems] = useState(false)
  const [maxVisibleItems] = useState(3) // Show max 3 items initially
  
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
    if (open && !cart) {
      loadCart()
    }
  }, [open, cart, loadCart])

  useEffect(() => {
    if (open && availableCoupons.length === 0) {
      loadAvailableCoupons()
    }
  }, [open, availableCoupons.length, loadAvailableCoupons])

  const handleQuantityChange = async (itemId: number, newQuantity: number) => {
    await updateQuantity(itemId, newQuantity)
  }

  const handleRemoveItem = async (itemId: number) => {
    await removeItem(itemId)
  }

  const handleClearCart = async () => {
    await clearCart()
  }

  const handleApplyCoupon = async () => {
    if (!couponCode.trim() || !cart) return
    
    setIsApplyingCoupon(true)
    setCouponMessage("")
    try {
      await applyCoupon(couponCode.trim())
      setCouponCode("") // Clear the input on success
      setCouponMessage("Coupon applied successfully!")
    } catch (error) {
      console.error("Failed to apply coupon:", error)
      setCouponMessage(error instanceof Error ? error.message : "Failed to apply coupon")
    } finally {
      setIsApplyingCoupon(false)
    }
  }

  const handleRemoveCoupon = async (couponCode: string) => {
    await removeCoupon(couponCode)
  }

  const handleSelectCoupon = async (coupon: any) => {
    setCouponCode(coupon.code)
    await handleApplyCoupon()
  }

  if (!isMounted) {
    return null
  }

  const visibleItems = showAllItems ? items : items.slice(0, maxVisibleItems)
  const hasMoreItems = items.length > maxVisibleItems

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-lg h-full flex flex-col">
        <SheetHeader className="flex-shrink-0">
          <SheetTitle className="flex items-center gap-2">
            <ShoppingCart className="h-5 w-5" />
            Shopping Cart
            {itemCount > 0 && (
              <Badge variant="secondary" className="ml-2">
                {itemCount}
              </Badge>
            )}
          </SheetTitle>
          <SheetDescription>
            {itemCount > 0 
              ? `${itemCount} item${itemCount === 1 ? '' : 's'} in your cart`
              : "Your cart is empty"
            }
          </SheetDescription>
        </SheetHeader>

        <div className="flex-1 flex flex-col min-h-0">
          {isLoading ? (
            <div className="flex-1 flex items-center justify-center">
              <div className="flex items-center gap-2 text-muted-foreground">
                <Loader2 className="h-4 w-4 animate-spin" />
                Loading cart...
              </div>
            </div>
          ) : error ? (
            <div className="flex-1 flex items-center justify-center">
              <div className="text-center">
                <div className="text-destructive mb-2">Failed to load cart</div>
                <Button variant="outline" size="sm" onClick={loadCart}>
                  Try Again
                </Button>
              </div>
            </div>
          ) : itemCount === 0 ? (
            <div className="flex-1 flex flex-col items-center justify-center text-center">
              <Package className="h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="text-lg font-semibold mb-2">Your cart is empty</h3>
              <p className="text-muted-foreground mb-4">
                Add some items to get started
              </p>
              <Button asChild onClick={() => onOpenChange(false)}>
                <Link href="/shop">
                  Continue Shopping
                </Link>
              </Button>
            </div>
          ) : (
            <>
              {/* Cart Items - Compact List */}
              <div className="flex-1 py-4 min-h-0">
                <div className="space-y-2">
                  {visibleItems.map((item) => (
                    <div key={item.id} className="flex items-center gap-3 p-3 border rounded-lg hover:bg-muted/50 transition-colors group">
                      {/* Product Image - Smaller */}
                      <div className="w-12 h-12 bg-muted rounded-md flex items-center justify-center flex-shrink-0">
                        {item.image ? (
                          <img 
                            src={item.image} 
                            alt={item.name || 'Product'} 
                            className="w-full h-full object-cover rounded-md"
                          />
                        ) : (
                          <Package className="h-4 w-4 text-muted-foreground" />
                        )}
                      </div>
                      
                      {/* Product Info - More Compact */}
                      <div className="flex-1 min-w-0">
                        <h4 className="font-medium text-sm line-clamp-1 mb-1">
                          {item.name || `Product ${item.product_id}`}
                        </h4>
                        <div className="flex items-center gap-2 text-xs text-muted-foreground">
                          <span>SKU: {item.sku || `SKU-${item.product_id}`}</span>
                          {item.variant_name && (
                            <>
                              <span>•</span>
                              <span className="text-blue-600">{item.variant_name}</span>
                            </>
                          )}
                        </div>
                      </div>
                      
                      {/* Quantity and Price - Horizontal Layout */}
                      <div className="flex items-center gap-3">
                        {/* Quantity Controls - Compact */}
                        <div className="flex items-center border rounded-md">
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-7 w-7 p-0"
                            onClick={() => handleQuantityChange(item.id, item.quantity - 1)}
                            disabled={isLoading}
                          >
                            <Minus className="h-3 w-3" />
                          </Button>
                          <span className="px-2 py-1 text-xs min-w-[1.5rem] text-center">
                            {item.quantity}
                          </span>
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-7 w-7 p-0"
                            onClick={() => handleQuantityChange(item.id, item.quantity + 1)}
                            disabled={isLoading || (item.maxQuantity ? item.quantity >= item.maxQuantity : false)}
                          >
                            <Plus className="h-3 w-3" />
                          </Button>
                        </div>
                        
                        {/* Price */}
                        <div className="text-right min-w-[4rem]">
                          <div className="font-medium text-sm">
                            {formatCurrency(item.total_price)}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            {formatCurrency(item.unit_price)} ea
                          </div>
                        </div>
                        
                        {/* Remove Button */}
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-7 w-7 p-0 text-destructive hover:text-destructive opacity-0 group-hover:opacity-100 transition-opacity"
                          onClick={() => handleRemoveItem(item.id)}
                          disabled={isLoading}
                        >
                          <Trash2 className="h-3 w-3" />
                        </Button>
                      </div>
                    </div>
                  ))}

                  {/* Show More/Less Button */}
                  {hasMoreItems && (
                    <div className="pt-2">
                      <Button
                        variant="outline"
                        size="sm"
                        className="w-full"
                        onClick={() => setShowAllItems(!showAllItems)}
                      >
                        {showAllItems ? (
                          <>
                            <ChevronUp className="h-4 w-4 mr-2" />
                            Show Less
                          </>
                        ) : (
                          <>
                            <ChevronDown className="h-4 w-4 mr-2" />
                            Show {items.length - maxVisibleItems} More Items
                          </>
                        )}
                      </Button>
                    </div>
                  )}
                </div>
              </div>

              {/* Coupon Section - Compact */}
              <div className="py-3 border-t">
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <h4 className="font-medium text-sm">Coupons</h4>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setShowAvailableCoupons(!showAvailableCoupons)}
                    >
                      {showAvailableCoupons ? "Hide" : "Browse"}
                    </Button>
                  </div>

                  {/* Applied Coupons - Compact */}
                  {appliedCoupons.length > 0 && (
                    <div className="space-y-1">
                      {appliedCoupons.map((coupon) => (
                        <div key={coupon.id} className="flex items-center justify-between bg-green-50 dark:bg-green-900/20 px-2 py-1 rounded text-xs">
                          <div className="flex items-center gap-2">
                            <span className="font-medium text-green-700 dark:text-green-300">
                              {coupon.coupon_code}
                            </span>
                            <span className="text-green-600 dark:text-green-400">
                              -{formatCurrency(coupon.discount_amount)}
                            </span>
                          </div>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleRemoveCoupon(coupon.coupon_code)}
                            className="h-5 w-5 p-0 text-red-500 hover:text-red-700"
                          >
                            <X className="h-3 w-3" />
                          </Button>
                        </div>
                      ))}
                    </div>
                  )}

                  {/* Available Coupons - Compact */}
                  {showAvailableCoupons && (
                    <div className="space-y-1">
                      {appliedCoupons.length > 0 && (
                        <div className="text-xs text-muted-foreground px-2 py-1 bg-blue-50 dark:bg-blue-950/20 rounded">
                          Only one coupon can be applied per cart. Remove the current coupon to apply a different one.
                        </div>
                      )}
                      <div className="max-h-20 overflow-y-auto space-y-1">
                        {availableCoupons.map((coupon) => (
                          <div key={coupon.id} className="flex items-center justify-between px-2 py-1 border rounded text-xs hover:bg-muted/50">
                            <div className="flex-1">
                              <div className="font-medium">{coupon.code}</div>
                              <div className="text-muted-foreground">
                                {(() => {
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
                                })()}
                              </div>
                            </div>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleSelectCoupon(coupon)}
                              disabled={appliedCoupons.length > 0}
                              className="h-6 px-2 text-xs"
                            >
                              {appliedCoupons.some(ac => ac.coupon_code === coupon.code) ? "Applied" : 
                               appliedCoupons.length > 0 ? "Remove existing first" : "Apply"}
                            </Button>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Manual Coupon Entry - Compact */}
                  <div className="flex gap-2">
                    <Input
                      placeholder="Enter coupon code"
                      value={couponCode}
                      onChange={(e) => setCouponCode(e.target.value)}
                      className="flex-1 h-8 text-sm"
                    />
                    <Button
                      onClick={handleApplyCoupon}
                      disabled={!couponCode.trim() || isApplyingCoupon}
                      size="sm"
                      className="h-8 px-3"
                    >
                      {isApplyingCoupon ? (
                        <Loader2 className="h-3 w-3 animate-spin" />
                      ) : (
                        "Apply"
                      )}
                    </Button>
                  </div>
                  {couponMessage && (
                    <div className={`text-xs ${couponMessage.includes('successfully') ? 'text-green-600' : 'text-red-600'}`}>
                      {couponMessage}
                    </div>
                  )}
                </div>
              </div>

              {/* Cart Summary - Fixed at Bottom */}
              <div className="py-3 border-t space-y-2 bg-background">
                <div className="space-y-1">
                  <div className="flex justify-between text-sm">
                    <span>Subtotal</span>
                    <span>{formatCurrency(cart?.subtotal || totalPrice)}</span>
                  </div>
                  {cart && cart.tax_amount > 0 && (
                    <div className="flex justify-between text-sm">
                      <span>Tax</span>
                      <span>{formatCurrency(cart.tax_amount)}</span>
                    </div>
                  )}
                  {cart && cart.shipping_amount > 0 && (
                    <div className="flex justify-between text-sm">
                      <span>Shipping</span>
                      <span>{formatCurrency(cart.shipping_amount)}</span>
                    </div>
                  )}
                  {cart && cart.discount_amount > 0 && (
                    <div className="flex justify-between text-sm text-green-600">
                      <span>Discount</span>
                      <span>-{formatCurrency(cart.discount_amount)}</span>
                    </div>
                  )}
                  <Separator />
                  <div className="flex justify-between font-semibold">
                    <span>Total</span>
                    <span>{formatCurrency(cart?.total || totalPrice)}</span>
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="space-y-2">
                  <Button 
                    className="w-full" 
                    size="sm"
                    asChild
                  >
                    <Link href="/checkout">
                      <CreditCard className="h-4 w-4 mr-2" />
                      Proceed to Checkout
                      <ArrowRight className="h-4 w-4 ml-2" />
                    </Link>
                  </Button>
                  
                  <div className="flex gap-2">
                    <Button 
                      variant="outline" 
                      className="flex-1"
                      size="sm"
                      asChild
                    >
                      <Link href="/shop" onClick={() => onOpenChange(false)}>
                        Continue Shopping
                      </Link>
                    </Button>
                    
                    <Button 
                      variant="outline" 
                      size="sm"
                      onClick={handleClearCart}
                      disabled={isLoading}
                      className="text-destructive hover:text-destructive"
                    >
                      <Trash2 className="h-4 w-4 mr-1" />
                      Clear
                    </Button>
                  </div>
                </div>
              </div>
            </>
          )}
        </div>
      </SheetContent>
    </Sheet>
  )
}