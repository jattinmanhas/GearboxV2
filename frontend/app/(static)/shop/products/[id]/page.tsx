"use client"

import { useState, useEffect } from "react"
import { useParams } from "next/navigation"
import Image from "next/image"
import Link from "next/link"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { 
  ArrowLeft, 
  Heart, 
  ShoppingCart, 
  Star,
  Package,
  Truck,
  Shield,
  RotateCcw
} from "lucide-react"
import { productApi } from "@/lib/api"
import { Product, Category, ProductVariant } from "@/lib/types"
import { useCartStore } from "@/lib/stores/cart-store"
import { useWishlistStore } from "@/lib/stores/wishlist-store"
import { formatPrice } from "@/lib/currency"

export default function ProductDetailPage() {
  const params = useParams()
  const productId = params.id as string
  
  const [product, setProduct] = useState<Product | null>(null)
  const [categories, setCategories] = useState<Category[]>([])
  const [variants, setVariants] = useState<ProductVariant[]>([])
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [quantity, setQuantity] = useState(1)
  const [isAddingToCart, setIsAddingToCart] = useState(false)
  const [isWishlisting, setIsWishlisting] = useState(false)

  const { addItem, isLoading: cartLoading } = useCartStore()
  const { 
    addItemToWishlist, 
    removeItemFromWishlist, 
    isProductInWishlist,
    wishlists,
    loadWishlists 
  } = useWishlistStore()

  // Check if product is in any wishlist
  const isWishlisted = product ? wishlists.some(wishlist => 
    isProductInWishlist(product.id, wishlist.id)
  ) : false

  useEffect(() => {
    const loadProduct = async () => {
      try {
        setLoading(true)
        setError(null)
        const [productData, categoriesData, variantsData] = await Promise.all([
          productApi.getProduct(parseInt(productId)),
          productApi.getCategories({
            page: 1,
            limit: 10
          }),
          productApi.getProductVariants(parseInt(productId))
        ])
        setProduct(productData)
        setCategories(categoriesData.data.categories)
        setVariants(variantsData)
        
        // Set the first active variant as selected, or the first variant if none are active
        const activeVariants = variantsData.filter(v => v.is_active)
        if (activeVariants.length > 0) {
          setSelectedVariant(activeVariants[0])
        } else if (variantsData.length > 0) {
          setSelectedVariant(variantsData[0])
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load product")
      } finally {
        setLoading(false)
      }
    }

    if (productId) {
      loadProduct()
    }
  }, [productId])


  // No longer needed - using category_names from API response

  const handleAddToCart = async () => {
    if (!product) return
    
    setIsAddingToCart(true)
    try {
      await addItem({
        product_id: product.id,
        product_variant_id: selectedVariant?.id,
        quantity: quantity
      })
    } catch (error) {
      console.error("Failed to add to cart:", error)
    } finally {
      setIsAddingToCart(false)
    }
  }

  const handleWishlist = async () => {
    if (!product || isWishlisting) return
    
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
          // Create a default wishlist first
          // TODO: Implement default wishlist creation
          console.log("No wishlists available, need to create one first")
        }
      }
    } catch (error) {
      console.error("Failed to update wishlist:", error)
    } finally {
      setIsWishlisting(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">Loading product...</div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !product) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-6">
          <div className="text-center py-8">
            <div className="text-destructive mb-4">
              {error || "Product not found"}
            </div>
            <Button asChild>
              <Link href="/shop">Back to Shop</Link>
            </Button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Breadcrumb */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Link href="/shop" className="hover:text-foreground">
              Shop
            </Link>
            <span>/</span>
            <span className="text-foreground">{product.name}</span>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Product Images */}
          <div className="space-y-4">
            <div className="aspect-square bg-muted rounded-lg flex items-center justify-center">
              <Package className="h-24 w-24 text-muted-foreground" />
            </div>
            
            {/* Thumbnail images would go here */}
            <div className="flex gap-2">
              <div className="w-16 h-16 bg-muted rounded-md flex items-center justify-center">
                <Package className="h-6 w-6 text-muted-foreground" />
              </div>
              <div className="w-16 h-16 bg-muted rounded-md flex items-center justify-center">
                <Package className="h-6 w-6 text-muted-foreground" />
              </div>
              <div className="w-16 h-16 bg-muted rounded-md flex items-center justify-center">
                <Package className="h-6 w-6 text-muted-foreground" />
              </div>
            </div>
          </div>

          {/* Product Info */}
          <div className="space-y-6">
            {/* Header */}
            <div>
              <div className="flex items-center gap-2 mb-2">
                <h1 className="text-3xl font-bold">{product.name}</h1>
                {!product.is_active && (
                  <Badge variant="secondary">Out of Stock</Badge>
                )}
              </div>
              
              <div className="flex items-center gap-4 mb-4">
                <div className="flex items-center">
                  {[...Array(5)].map((_, i) => (
                    <Star key={i} className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                  ))}
                  <span className="text-sm text-muted-foreground ml-2">(4.5) 23 reviews</span>
                </div>
                <span className="text-sm text-muted-foreground">•</span>
                <span className="text-sm text-muted-foreground">SKU: {product.sku}</span>
              </div>

              {/* Price */}
              <div className="mb-4">
                <div className="text-3xl font-bold">
                  {selectedVariant ? formatPrice(selectedVariant.price) : formatPrice(product.price)}
                </div>
                {(selectedVariant ? selectedVariant.compare_price : product.compare_price) > 0 && (
                  <div className="text-lg text-muted-foreground line-through">
                    {formatPrice(selectedVariant ? selectedVariant.compare_price : product.compare_price)}
                  </div>
                )}
              </div>
            </div>

            {/* Categories */}
            {product.category_names && product.category_names.length > 0 && (
              <div>
                <h3 className="text-sm font-medium mb-2">Categories</h3>
                <div className="flex gap-2">
                  {product.category_names.map((name, index) => (
                    <Badge key={index} variant="outline">
                      {name}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            {/* Variants */}
            {variants.length > 0 && (
              <div>
                <h3 className="text-sm font-medium mb-2">Options</h3>
                <div className="flex flex-wrap gap-2">
                  {variants
                    .filter(variant => variant.is_active)
                    .sort((a, b) => a.position - b.position)
                    .map((variant) => (
                    <Button
                      key={variant.id}
                      variant={selectedVariant?.id === variant.id ? "default" : "outline"}
                      size="sm"
                      onClick={() => setSelectedVariant(variant)}
                      className="flex items-center gap-2"
                    >
                      {variant.name}
                      <span className="text-xs text-muted-foreground">
                        ({formatPrice(variant.price)})
                      </span>
                    </Button>
                  ))}
                </div>
                {selectedVariant && (
                  <div className="mt-2 text-sm text-muted-foreground">
                    Selected: <span className="font-medium">{selectedVariant.name}</span>
                    <span className="ml-2">•</span>
                    <span className="ml-2">SKU: {selectedVariant.sku}</span>
                  </div>
                )}
              </div>
            )}

            {/* Description */}
            <div>
              <h3 className="text-sm font-medium mb-2">Description</h3>
              <p className="text-muted-foreground">
                {product.description || product.short_description}
              </p>
            </div>

            {/* Quantity and Add to Cart */}
            <div className="space-y-4">
              <div>
                <h3 className="text-sm font-medium mb-2">Quantity</h3>
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setQuantity(Math.max(1, quantity - 1))}
                    className="h-10 w-10 p-0"
                  >
                    -
                  </Button>
                  <span className="text-lg font-medium min-w-[3rem] text-center">
                    {quantity}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setQuantity(quantity + 1)}
                    className="h-10 w-10 p-0"
                  >
                    +
                  </Button>
                </div>
              </div>

              <div className="flex gap-3">
                <Button
                  className="flex-1"
                  size="lg"
                  onClick={handleAddToCart}
                  disabled={!product.is_active || isAddingToCart || cartLoading}
                >
                  <ShoppingCart className="h-4 w-4 mr-2" />
                  {isAddingToCart ? "Adding..." : "Add to Cart"}
                </Button>
                <Button
                  variant="outline"
                  size="lg"
                  onClick={handleWishlist}
                  disabled={isWishlisting}
                  className={isWishlisted ? "text-red-500" : ""}
                >
                  <Heart className={`h-4 w-4 ${isWishlisted ? "fill-current" : ""}`} />
                </Button>
              </div>
            </div>

            {/* Features */}
            <div className="space-y-3">
              <div className="flex items-center gap-3 text-sm">
                <Truck className="h-4 w-4 text-muted-foreground" />
                <span>Free shipping on orders over $50</span>
              </div>
              <div className="flex items-center gap-3 text-sm">
                <RotateCcw className="h-4 w-4 text-muted-foreground" />
                <span>30-day return policy</span>
              </div>
              <div className="flex items-center gap-3 text-sm">
                <Shield className="h-4 w-4 text-muted-foreground" />
                <span>2-year warranty</span>
              </div>
            </div>
          </div>
        </div>

        {/* Product Details Tabs */}
        <div className="mt-12">
          <Card>
            <CardContent className="p-6">
              <div className="space-y-6">
                <div>
                  <h3 className="text-lg font-semibold mb-3">Product Details</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="font-medium">SKU:</span> {product.sku}
                    </div>
                    <div>
                      <span className="font-medium">Weight:</span> {product.weight} lbs
                    </div>
                    {product.dimensions && (
                      <div>
                        <span className="font-medium">Dimensions:</span> {product.dimensions}
                      </div>
                    )}
                    <div>
                      <span className="font-medium">Type:</span> {product.is_digital ? "Digital" : "Physical"}
                    </div>
                    <div>
                      <span className="font-medium">Shipping:</span> {product.requires_shipping ? "Required" : "Not Required"}
                    </div>
                    <div>
                      <span className="font-medium">Taxable:</span> {product.taxable ? "Yes" : "No"}
                    </div>
                  </div>
                </div>

                {product.tags && (
                  <div>
                    <h3 className="text-lg font-semibold mb-3">Tags</h3>
                    <div className="flex flex-wrap gap-2">
                      {product.tags.split(',').map((tag, index) => (
                        <Badge key={index} variant="secondary">
                          {tag.trim()}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Back to Shop */}
        <div className="mt-8">
          <Button variant="outline" asChild>
            <Link href="/shop">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Shop
            </Link>
          </Button>
        </div>
      </div>
    </div>
  )
}
