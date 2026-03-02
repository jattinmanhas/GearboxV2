"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"
import { Badge } from "@/components/ui/badge"
import {
  ShoppingCart,
  MapPin,
  CreditCard,
  Check,
  Loader2,
  ArrowLeft,
  Plus,
  AlertCircle
} from "lucide-react"
import { useCartStore, useHydrateCartStore } from "@/lib/stores/cart-store"
import { profileApi, orderApi } from "@/lib/apiFunctions"
import { formatCurrency } from "@/lib/currency"
import { Address, OrderAddressRequest } from "@/lib/types"
import { showSuccess, showError, showLoading, updateLoading } from "@/lib/notifications"
import Link from "next/link"
import Image from "next/image"

export default function CheckoutPage() {
  const router = useRouter()
  const [isMounted, setIsMounted] = useState(false)

  // Hydrate the store on client side
  useHydrateCartStore()

  const {
    cart,
    items,
    appliedCoupons,
    isLoading,
    loadCart,
    getTotalPrice
  } = useCartStore()

  const [addresses, setAddresses] = useState<Address[]>([])
  const [paymentMethods, setPaymentMethods] = useState<any[]>([])
  const [paymentGateways, setPaymentGateways] = useState<any[]>([])
  const [userEmail, setUserEmail] = useState<string>("")
  const [loadingAddresses, setLoadingAddresses] = useState(true)

  const [selectedShippingAddress, setSelectedShippingAddress] = useState<number | null>(null)
  const [selectedBillingAddress, setSelectedBillingAddress] = useState<number | null>(null)
  const [selectedPaymentMethod, setSelectedPaymentMethod] = useState<string>("card")
  const [selectedGateway, setSelectedGateway] = useState<string>("stripe")
  const [sameAsShipping, setSameAsShipping] = useState(true)
  const [isProcessing, setIsProcessing] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const totalPrice = getTotalPrice()

  // Prevent hydration mismatch
  useEffect(() => {
    setIsMounted(true)
  }, [])

  // Load cart and addresses
  useEffect(() => {
    if (!isMounted) return

    const loadData = async () => {
      try {
        // Load cart
        if (!cart) {
          await loadCart()
        }

        // Load user profile for email
        const profileData = await profileApi.getProfile()
        setUserEmail(profileData.email)

        // Load addresses
        try {
          const addressesData = await profileApi.getAddresses()
          // getAddresses() always returns an array
          setAddresses(addressesData)

          // Set defaults
          if (addressesData.length > 0) {
            const defaultAddress = addressesData.find((a: Address) => a.is_default) || addressesData[0]
            setSelectedShippingAddress(defaultAddress.id)
            if (sameAsShipping) {
              setSelectedBillingAddress(defaultAddress.id)
            }
          }
        } catch (err) {
          console.error("Failed to load addresses:", err)
          setAddresses([]) // Set to empty array on error
        } finally {
          setLoadingAddresses(false)
        }

        // Initialize payment defaults (static now)
        setSelectedPaymentMethod("card")
        setSelectedGateway("stripe")
      } catch (err) {
        console.error("Failed to load checkout data:", err)
        setError(err instanceof Error ? err.message : "Failed to load checkout data")
      }
    }

    loadData()
  }, [isMounted, cart, loadCart])

  // Convert address to order address format
  const addressToOrderAddress = (address: Address, email: string): OrderAddressRequest => {
    return {
      first_name: address.first_name,
      last_name: address.last_name,
      company: address.company,
      address1: address.address_line_1,
      address2: address.address_line_2,
      city: address.city,
      state: address.state,
      country: address.country,
      postal_code: address.postal_code,
      phone: address.phone_number || "",
      email: email,
    }
  }

  const handleCheckout = async () => {
    if (!cart) {
      showError("Cart not found")
      return
    }

    if (!selectedShippingAddress) {
      showError("Please select a shipping address")
      return
    }

    if (!selectedPaymentMethod) {
      showError("Please select a payment method")
      return
    }

    if (!selectedGateway) {
      showError("Please select a payment gateway")
      return
    }

    setIsProcessing(true)
    setError(null)

    const loadingToast = showLoading("Processing your order...")

    try {
      // Get selected addresses
      const shippingAddress = addresses.find(a => a.id === selectedShippingAddress)
      const billingAddress = sameAsShipping
        ? shippingAddress
        : addresses.find(a => a.id === selectedBillingAddress)

      if (!shippingAddress || !billingAddress) {
        throw new Error("Address not found")
      }

      if (!userEmail) {
        throw new Error("User email not found")
      }

      // Prepare order data
      const orderData = {
        shipping_address: addressToOrderAddress(shippingAddress, userEmail),
        billing_address: sameAsShipping ? undefined : addressToOrderAddress(billingAddress, userEmail),
        user_shipping_address_id: selectedShippingAddress,
        user_billing_address_id: sameAsShipping ? selectedShippingAddress : (selectedBillingAddress || null),
        payment_method: selectedPaymentMethod,
        currency: "INR", // Get from cart or user preference
        apply_coupons: appliedCoupons.map(c => c.coupon_code),
      }

      // Create order from cart
      const orderResponse = await orderApi.createOrderFromCart(
        typeof cart.id === 'string' ? parseInt(cart.id) : cart.id,
        orderData
      )
      // Extract order from response (backend wraps in data field)
      const order = orderResponse?.data || orderResponse

      // Create payment via order service to keep order/payment state in sync
      const paymentResponse = await orderApi.createOrderPayment(order.id, {
        payment_method: selectedPaymentMethod,
        gateway_id: selectedGateway,
        metadata: {
          order_number: order.order_number,
          cart_id: cart.id,
        },
      })
      const createdPayment = paymentResponse?.payment || paymentResponse?.data?.payment || paymentResponse

      // Process payment via order service so payment status updates on the order
      const processResponse = await orderApi.processOrderPayment(order.id, {
        payment_data: {},
        return_url: `${window.location.origin}/checkout/success`,
        cancel_url: `${window.location.origin}/checkout`,
      })

      // Check if payment gateway requires redirect
      const processedPayment = processResponse?.payment || processResponse?.data?.payment || processResponse
      if (!createdPayment?.id || !processedPayment?.id) {
        throw new Error("Payment processing failed")
      }
      if (processedPayment?.gateway_response?.redirect_url) {
        // Redirect to payment gateway
        window.location.href = processedPayment.gateway_response.redirect_url
        return
      }

      updateLoading(loadingToast, "Order placed successfully!", "success")

      // Redirect to success page
      setTimeout(() => {
        router.push(`/checkout/success?order_id=${order.id}`)
      }, 1500)
    } catch (err) {
      console.error("Checkout error:", err)
      const errorMessage = err instanceof Error ? err.message : "Failed to process checkout"
      setError(errorMessage)
      updateLoading(loadingToast, "Checkout failed", "error", {
        description: errorMessage
      })
    } finally {
      setIsProcessing(false)
    }
  }

  if (!isMounted) {
    return null
  }

  if (!cart || items.length === 0) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-8">
          <Card>
            <CardContent className="py-12">
              <div className="text-center">
                <ShoppingCart className="h-16 w-16 mx-auto mb-4 text-muted-foreground" />
                <h2 className="text-2xl font-semibold mb-2">Your cart is empty</h2>
                <p className="text-muted-foreground mb-6">
                  Add some items to your cart before checking out
                </p>
                <Button asChild>
                  <Link href="/shop">Continue Shopping</Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    )
  }

  const shippingAddress = addresses.find(a => a.id === selectedShippingAddress)
  const billingAddress = sameAsShipping
    ? shippingAddress
    : addresses.find(a => a.id === selectedBillingAddress)

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto px-4 py-8 max-w-7xl">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-4 mb-4">
            <Button variant="ghost" size="sm" asChild>
              <Link href="/cart">
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Cart
              </Link>
            </Button>
          </div>
          <h1 className="text-3xl font-bold flex items-center gap-3">
            <CreditCard className="h-8 w-8" />
            Checkout
          </h1>
        </div>

        {error && (
          <Card className="mb-6 border-destructive">
            <CardContent className="pt-6">
              <div className="flex items-center gap-2 text-destructive">
                <AlertCircle className="h-5 w-5" />
                <span>{error}</span>
              </div>
            </CardContent>
          </Card>
        )}

        <div className="grid lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Shipping Address */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MapPin className="h-5 w-5" />
                  Shipping Address
                </CardTitle>
                <CardDescription>
                  Select or add a shipping address
                </CardDescription>
              </CardHeader>
              <CardContent>
                {loadingAddresses ? (
                  <div className="text-center py-8">
                    <Loader2 className="h-6 w-6 animate-spin mx-auto mb-4 text-muted-foreground" />
                    <p className="text-muted-foreground">Loading addresses...</p>
                  </div>
                ) : !addresses || addresses.length === 0 ? (
                  <div className="text-center py-8">
                    <p className="text-muted-foreground mb-4">No addresses found</p>
                    <Button asChild variant="outline">
                      <Link href="/profile">Add Address</Link>
                    </Button>
                  </div>
                ) : (
                  <RadioGroup
                    value={selectedShippingAddress?.toString() || ""}
                    onValueChange={(value) => setSelectedShippingAddress(parseInt(value))}
                  >
                    <div className="space-y-3">
                      {addresses.map((address) => (
                        <div key={address.id} className="flex items-start space-x-3">
                          <RadioGroupItem value={address.id.toString()} id={`shipping-${address.id}`} />
                          <Label
                            htmlFor={`shipping-${address.id}`}
                            className={`flex-1 cursor-pointer p-4 border rounded-lg hover:bg-muted/50 transition-colors ${selectedShippingAddress === address.id ? 'border-primary bg-primary/5' : ''
                              }`}
                          >
                            <div className="flex items-start justify-between w-full">
                              <div className="flex-1">
                                <div className="font-medium flex items-center gap-2">
                                  {address.first_name} {address.last_name}
                                  {address.is_default && (
                                    <Badge variant="secondary" className="text-xs">
                                      <Check className="h-3 w-3 mr-1" />
                                      Default
                                    </Badge>
                                  )}
                                </div>
                                <div className="text-sm text-muted-foreground mt-1">
                                  {address.address_line_1}
                                  {address.address_line_2 && `, ${address.address_line_2}`}
                                  <br />
                                  {address.city}, {address.state} {address.postal_code}
                                  <br />
                                  {address.country}
                                </div>
                                {address.phone_number && (
                                  <div className="text-sm text-muted-foreground mt-1">
                                    {address.phone_number}
                                  </div>
                                )}
                              </div>
                            </div>
                          </Label>
                        </div>
                      ))}
                    </div>
                  </RadioGroup>
                )}
                <Button variant="outline" className="mt-4" asChild>
                  <Link href="/profile">
                    <Plus className="h-4 w-4 mr-2" />
                    Add New Address
                  </Link>
                </Button>
              </CardContent>
            </Card>

            {/* Billing Address */}
            <Card>
              <CardHeader>
                <CardTitle>Billing Address</CardTitle>
                <CardDescription>
                  Select or add a billing address
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex items-center space-x-2 mb-4">
                  <input
                    type="checkbox"
                    id="same-as-shipping"
                    checked={sameAsShipping}
                    onChange={(e) => {
                      setSameAsShipping(e.target.checked)
                      if (e.target.checked) {
                        setSelectedBillingAddress(selectedShippingAddress)
                      }
                    }}
                    className="h-4 w-4"
                  />
                  <Label htmlFor="same-as-shipping" className="cursor-pointer">
                    Same as shipping address
                  </Label>
                </div>

                {!sameAsShipping && (
                  <>
                    {loadingAddresses ? (
                      <div className="text-center py-8">
                        <Loader2 className="h-6 w-6 animate-spin mx-auto mb-4 text-muted-foreground" />
                        <p className="text-muted-foreground">Loading addresses...</p>
                      </div>
                    ) : !addresses || addresses.length === 0 ? (
                      <div className="text-center py-8">
                        <p className="text-muted-foreground mb-4">No addresses found</p>
                        <Button asChild variant="outline">
                          <Link href="/profile">Add Address</Link>
                        </Button>
                      </div>
                    ) : (
                      <RadioGroup
                        value={selectedBillingAddress?.toString() || ""}
                        onValueChange={(value) => setSelectedBillingAddress(parseInt(value))}
                      >
                        <div className="space-y-3">
                          {addresses.map((address) => (
                            <div key={address.id} className="flex items-start space-x-3">
                              <RadioGroupItem value={address.id.toString()} id={`billing-${address.id}`} />
                              <Label
                                htmlFor={`billing-${address.id}`}
                                className={`flex-1 cursor-pointer p-4 border rounded-lg hover:bg-muted/50 transition-colors ${selectedBillingAddress === address.id ? 'border-primary bg-primary/5' : ''
                                  }`}
                              >
                                <div className="flex items-start justify-between w-full">
                                  <div className="flex-1">
                                    <div className="font-medium flex items-center gap-2">
                                      {address.first_name} {address.last_name}
                                      {address.is_default && (
                                        <Badge variant="secondary" className="text-xs">
                                          <Check className="h-3 w-3 mr-1" />
                                          Default
                                        </Badge>
                                      )}
                                    </div>
                                    <div className="text-sm text-muted-foreground mt-1">
                                      {address.address_line_1}
                                      {address.address_line_2 && `, ${address.address_line_2}`}
                                      <br />
                                      {address.city}, {address.state} {address.postal_code}
                                      <br />
                                      {address.country}
                                    </div>
                                    {address.phone_number && (
                                      <div className="text-sm text-muted-foreground mt-1">
                                        {address.phone_number}
                                      </div>
                                    )}
                                  </div>
                                </div>
                              </Label>
                            </div>
                          ))}
                        </div>
                      </RadioGroup>
                    )}
                    <Button variant="outline" className="mt-4" asChild>
                      <Link href="/profile">
                        <Plus className="h-4 w-4 mr-2" />
                        Add New Address
                      </Link>
                    </Button>
                  </>
                )}
              </CardContent>
            </Card>

            {/* Payment & Gateway Selection (Simplified) */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <CreditCard className="h-5 w-5" />
                  Payment Details
                </CardTitle>
                <CardDescription>
                  Secure payment via Stripe
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="p-4 border rounded-lg bg-primary/5 border-primary flex items-center justify-between">
                  <div>
                    <div className="font-medium">Credit / Debit Card</div>
                    <div className="text-sm text-muted-foreground mt-1">
                      Processed securely by Stripe
                    </div>
                  </div>
                  <Badge>Selected</Badge>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Order Summary */}
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Order Summary</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Cart Items */}
                <div className="space-y-3">
                  {items.map((item) => (
                    <div key={item.id} className="flex gap-3">
                      <div className="w-16 h-16 bg-muted rounded-md flex items-center justify-center flex-shrink-0 overflow-hidden">
                        {item.image ? (
                          <Image
                            src={item.image}
                            alt={item.name || 'Product'}
                            width={64}
                            height={64}
                            className="w-full h-full object-cover"
                          />
                        ) : (
                          <ShoppingCart className="h-6 w-6 text-muted-foreground" />
                        )}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="font-medium">{item.name || `Product ${item.product_id}`}</div>
                        <div className="text-sm text-muted-foreground">
                          {formatCurrency(item.unit_price)} × {item.quantity}
                        </div>
                      </div>
                      <div className="font-medium">
                        {formatCurrency(item.total_price)}
                      </div>
                    </div>
                  ))}
                </div>

                <Separator />

                {/* Totals */}
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span>Subtotal</span>
                    <span className="font-medium">{formatCurrency(cart.subtotal || totalPrice)}</span>
                  </div>
                  {cart.discount_amount > 0 && (
                    <div className="flex justify-between text-green-600">
                      <span>Discount</span>
                      <span className="font-medium">-{formatCurrency(cart.discount_amount)}</span>
                    </div>
                  )}
                </div>

                <Separator />

                <div className="flex justify-between text-lg font-bold">
                  <span>Total</span>
                  <span>{formatCurrency(cart.total || totalPrice)}</span>
                </div>

                <Button
                  className="w-full"
                  size="lg"
                  onClick={handleCheckout}
                  disabled={isProcessing || !selectedShippingAddress || !selectedPaymentMethod || !selectedGateway}
                >
                  {isProcessing ? (
                    <>
                      <Loader2 className="h-5 w-5 mr-2 animate-spin" />
                      Processing...
                    </>
                  ) : (
                    <>
                      <Check className="h-5 w-5 mr-2" />
                      Place Order
                    </>
                  )}
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}
