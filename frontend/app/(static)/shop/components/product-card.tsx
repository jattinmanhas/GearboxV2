"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import Image from "next/image"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { 
  Heart, 
  ShoppingCart, 
  Star,
  Package,
  Eye
} from "lucide-react"
import { Product, Category, ProductVariant } from "@/lib/types"
import { useCartStore } from "@/lib/stores/cart-store"
import { useWishlistStore } from "@/lib/stores/wishlist-store"
import { productApi } from "@/lib/api"
import { formatPrice } from "@/lib/currency"

interface ProductCardProps {
  product: Product
  viewMode: "grid" | "list"
  categories: Category[]
}

export function ProductCard({ product, viewMode, categories }: ProductCardProps) {
  const [isAddingToCart, setIsAddingToCart] = useState(false)
  const [isWishlisting, setIsWishlisting] = useState(false)
  const [variants, setVariants] = useState<ProductVariant[]>([])
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(null)
  const [variantsLoading, setVariantsLoading] = useState(false)
  
  const { addItem, isLoading: cartLoading } = useCartStore()
  const { 
    addItemToWishlist, 
    removeItemFromWishlist, 
    isProductInWishlist,
    wishlists,
    loadWishlists,
    createWishlist,
    getWishlistItemCount
  } = useWishlistStore()

  // Load variants for this product
  useEffect(() => {
    const loadVariants = async () => {
      try {
        setVariantsLoading(true)
        const variantsData = await productApi.getProductVariants(product.id)
        setVariants(variantsData)
        
        // Set the first active variant as selected, or the first variant if none are active
        const activeVariants = variantsData.filter(v => v.is_active)
        if (activeVariants.length > 0) {
          setSelectedVariant(activeVariants[0])
        } else if (variantsData.length > 0) {
          setSelectedVariant(variantsData[0])
        }
      } catch (error) {
        console.error("Failed to load variants:", error)
      } finally {
        setVariantsLoading(false)
      }
    }

    loadVariants()
  }, [product.id])

  // Check if product is in any wishlist
  const isWishlisted = wishlists.some(wishlist => 
    isProductInWishlist(product.id, wishlist.id)
  )

  const handleAddToCart = async () => {
    if (!product.is_active) return
    
    setIsAddingToCart(true)
    try {
      await addItem({
        product_id: product.id,
        product_variant_id: selectedVariant?.id,
        quantity: 1
      })
    } catch (error) {
      console.error("Failed to add to cart:", error)
    } finally {
      setIsAddingToCart(false)
    }
  }

  const handleWishlist = async () => {
    if (isWishlisting) return
    
    setIsWishlisting(true)
    try {
      // Load wishlists if not already loaded
      if (wishlists.length === 0) {
        await loadWishlists()
      }
      
      if (isWishlisted) {
        // Find the wishlist containing this product and remove it
        const wishlistWithProduct = wishlists.find(wishlist => 
          isProductInWishlist(product.id, wishlist.id)
        )
        if (wishlistWithProduct) {
          const item = wishlistWithProduct.items.find(item => item.product_id === product.id)
          if (item) {
            await removeItemFromWishlist(wishlistWithProduct.id, item.id)
          }
        }
      } else {
        // Add to the first wishlist (or create a default one)
        if (wishlists.length > 0) {
          await addItemToWishlist(wishlists[0].id, product.id)
        } else {
          // Try to create a default wishlist first
          try {
            await createWishlist({
              name: "My Wishlist",
              description: "Default wishlist",
              is_public: false
            })
            // After creating, reload wishlists and then add the item
            await loadWishlists()
            // Get the updated wishlists from the store
            const store = useWishlistStore.getState()
            if (store.wishlists.length > 0) {
              await addItemToWishlist(store.wishlists[0].id, product.id)
            }
          } catch (createError) {
            console.error("Failed to create default wishlist:", createError)
            // Show a user-friendly message
            alert("Wishlist functionality is currently unavailable. Please try again later.")
          }
        }
      }
    } catch (error) {
      console.error("Failed to update wishlist:", error)
      // Show a user-friendly message
      alert("Failed to update wishlist. Please try again later.")
    } finally {
      setIsWishlisting(false)
    }
  }

  if (viewMode === "list") {
    return (
      <Card className="overflow-hidden">
        <CardContent className="p-0">
          <div className="flex">
            {/* Product Image - Clickable */}
            <Link href={`/shop/products/${product.id}`} className="block">
              <div className="w-48 h-48 bg-muted flex items-center justify-center cursor-pointer hover:bg-muted/80 transition-colors">
                <Package className="h-12 w-12 text-muted-foreground" />
              </div>
            </Link>
            
            {/* Product Info */}
            <div className="flex-1 p-6">
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <h3 className="text-lg font-semibold line-clamp-1">
                      {product.name}
                    </h3>
                    {!product.is_active && (
                      <Badge variant="secondary">Out of Stock</Badge>
                    )}
                  </div>
                  
                  <p className="text-sm text-muted-foreground mb-2 line-clamp-2">
                    {product.short_description || product.description}
                  </p>
                  
                  <div className="flex items-center gap-2 mb-3">
                    <div className="flex items-center">
                      {[...Array(5)].map((_, i) => (
                        <Star key={i} className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                      ))}
                      <span className="text-sm text-muted-foreground ml-1">(4.5)</span>
                    </div>
                    <span className="text-sm text-muted-foreground">•</span>
                    <span className="text-sm text-muted-foreground">SKU: {product.sku}</span>
                  </div>
                  
                  {product.category_names && product.category_names.length > 0 && (
                    <div className="flex gap-1 mb-3">
                      {product.category_names.slice(0, 2).map((name, index) => (
                        <Badge key={index} variant="outline" className="text-xs">
                          {name}
                        </Badge>
                      ))}
                    </div>
                  )}
                </div>
                
                {/* Price and Actions */}
                <div className="text-right ml-4">
                  <div className="mb-3">
                    <div className="text-2xl font-bold">
                      {selectedVariant ? formatPrice(selectedVariant.price) : formatPrice(product.price)}
                    </div>
                    {(selectedVariant ? selectedVariant.compare_price : product.compare_price) > 0 && (
                      <div className="text-sm text-muted-foreground line-through">
                        {formatPrice(selectedVariant ? selectedVariant.compare_price : product.compare_price)}
                      </div>
                    )}
                  </div>
                  
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={(e) => {
                        e.preventDefault()
                        e.stopPropagation()
                        handleWishlist()
                      }}
                      disabled={isWishlisting}
                      className={isWishlisted ? "text-red-500" : ""}
                    >
                      <Heart className={`h-4 w-4 ${isWishlisted ? "fill-current" : ""}`} />
                    </Button>
                    <Button
                      size="sm"
                      onClick={(e) => {
                        e.preventDefault()
                        e.stopPropagation()
                        handleAddToCart()
                      }}
                      disabled={!product.is_active || isAddingToCart || cartLoading || variantsLoading}
                    >
                      <ShoppingCart className="h-4 w-4 mr-1" />
                      {isAddingToCart ? "Adding..." : "Add to Cart"}
                    </Button>
                    <Button 
                      variant="outline" 
                      size="sm" 
                      asChild
                      onClick={(e) => {
                        e.preventDefault()
                        e.stopPropagation()
                      }}
                    >
                      <Link href={`/shop/products/${product.id}`}>
                        <Eye className="h-4 w-4" />
                      </Link>
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    )
  }

  // Grid view
  return (
    <Card className="group overflow-hidden hover:shadow-lg transition-shadow h-full flex flex-col">
      <CardContent className="p-0 h-full flex flex-col">
        {/* Product Image - Clickable */}
        <Link href={`/shop/products/${product.id}`} className="block">
          <div className="relative aspect-square bg-muted flex items-center justify-center cursor-pointer">
            <Package className="h-16 w-16 text-muted-foreground" />
            {!product.is_active && (
              <div className="absolute inset-0 bg-background/80 flex items-center justify-center">
                <Badge variant="secondary">Out of Stock</Badge>
              </div>
            )}
            
            {/* Quick Actions */}
            <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <Button
                variant="outline"
                size="sm"
                onClick={(e) => {
                  e.preventDefault()
                  e.stopPropagation()
                  handleWishlist()
                }}
                disabled={isWishlisting}
                className={`h-8 w-8 p-0 ${isWishlisted ? "text-red-500" : ""}`}
              >
                <Heart className={`h-4 w-4 ${isWishlisted ? "fill-current" : ""}`} />
              </Button>
            </div>
          </div>
        </Link>
        
        {/* Product Info - Flex grow to push buttons to bottom */}
        <div className="p-4 flex flex-col flex-grow">
          <div className="mb-2 flex-grow">
            <h3 className="font-semibold line-clamp-2 mb-1">
              {product.name}
            </h3>
            <p className="text-sm text-muted-foreground line-clamp-2">
              {product.short_description}
            </p>
          </div>
          
          {/* Rating */}
          <div className="flex items-center gap-1 mb-2">
            <div className="flex items-center">
              {[...Array(5)].map((_, i) => (
                <Star key={i} className="h-3 w-3 fill-yellow-400 text-yellow-400" />
              ))}
            </div>
            <span className="text-xs text-muted-foreground">(4.5)</span>
          </div>
          
          {/* Categories */}
          {product.category_names && product.category_names.length > 0 && (
            <div className="flex gap-1 mb-3">
              {product.category_names.slice(0, 2).map((name, index) => (
                <Badge key={index} variant="outline" className="text-xs">
                  {name}
                </Badge>
              ))}
            </div>
          )}
          
          {/* Price */}
          <div className="mb-3">
            <div className="text-lg font-bold">
              {selectedVariant ? formatPrice(selectedVariant.price) : formatPrice(product.price)}
            </div>
            {(selectedVariant ? selectedVariant.compare_price : product.compare_price) > 0 && (
              <div className="text-sm text-muted-foreground line-through">
                {formatPrice(selectedVariant ? selectedVariant.compare_price : product.compare_price)}
              </div>
            )}
          </div>
          
          {/* Actions - Always at bottom */}
          <div className="flex gap-2 mt-auto">
            <Button
              className="flex-1"
              onClick={(e) => {
                e.preventDefault()
                e.stopPropagation()
                handleAddToCart()
              }}
              disabled={!product.is_active || isAddingToCart || cartLoading || variantsLoading}
            >
              <ShoppingCart className="h-4 w-4 mr-1" />
              {isAddingToCart ? "Adding..." : "Add to Cart"}
            </Button>
            <Button 
              variant="outline" 
              size="sm" 
              asChild
              onClick={(e) => {
                e.preventDefault()
                e.stopPropagation()
              }}
            >
              <Link href={`/shop/products/${product.id}`}>
                <Eye className="h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
