"use client"

import { useState } from "react"
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
import { Product, Category } from "@/lib/types"
import { useCartStore } from "@/lib/stores/cart-store"

interface ProductCardProps {
  product: Product
  viewMode: "grid" | "list"
  categories: Category[]
}

export function ProductCard({ product, viewMode, categories }: ProductCardProps) {
  const [isWishlisted, setIsWishlisted] = useState(false)
  const [isAddingToCart, setIsAddingToCart] = useState(false)
  
  const { addItem } = useCartStore()

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(price)
  }

  const getCategoryNames = (categoryIds: number[]) => {
    return categoryIds
      .map(id => categories.find(cat => cat.id === id)?.name)
      .filter(Boolean)
      .slice(0, 2) // Show only first 2 categories
  }

  const handleAddToCart = async () => {
    if (!product.is_active) return
    
    setIsAddingToCart(true)
    try {
      // Add to cart using the store
      addItem({
        id: product.id,
        name: product.name,
        price: product.price,
        sku: product.sku,
        maxQuantity: product.max_quantity || 999
      })
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 500))
    } catch (error) {
      console.error("Failed to add to cart:", error)
    } finally {
      setIsAddingToCart(false)
    }
  }

  const handleWishlist = () => {
    setIsWishlisted(!isWishlisted)
    // TODO: Implement wishlist functionality
  }

  if (viewMode === "list") {
    return (
      <Card className="overflow-hidden">
        <CardContent className="p-0">
          <div className="flex">
            {/* Product Image */}
            <div className="w-48 h-48 bg-muted flex items-center justify-center">
              <Package className="h-12 w-12 text-muted-foreground" />
            </div>
            
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
                  
                  {getCategoryNames(product.category_ids).length > 0 && (
                    <div className="flex gap-1 mb-3">
                      {getCategoryNames(product.category_ids).map((name, index) => (
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
                      {formatPrice(product.price)}
                    </div>
                    {product.compare_price > 0 && (
                      <div className="text-sm text-muted-foreground line-through">
                        {formatPrice(product.compare_price)}
                      </div>
                    )}
                  </div>
                  
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handleWishlist}
                      className={isWishlisted ? "text-red-500" : ""}
                    >
                      <Heart className={`h-4 w-4 ${isWishlisted ? "fill-current" : ""}`} />
                    </Button>
                    <Button
                      size="sm"
                      onClick={handleAddToCart}
                      disabled={!product.is_active || isAddingToCart}
                    >
                      <ShoppingCart className="h-4 w-4 mr-1" />
                      {isAddingToCart ? "Adding..." : "Add to Cart"}
                    </Button>
                    <Button variant="outline" size="sm" asChild>
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
    <Card className="group overflow-hidden hover:shadow-lg transition-shadow">
      <CardContent className="p-0">
        {/* Product Image */}
        <div className="relative aspect-square bg-muted flex items-center justify-center">
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
              onClick={handleWishlist}
              className={`h-8 w-8 p-0 ${isWishlisted ? "text-red-500" : ""}`}
            >
              <Heart className={`h-4 w-4 ${isWishlisted ? "fill-current" : ""}`} />
            </Button>
          </div>
        </div>
        
        {/* Product Info */}
        <div className="p-4">
          <div className="mb-2">
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
          {getCategoryNames(product.category_ids).length > 0 && (
            <div className="flex gap-1 mb-3">
              {getCategoryNames(product.category_ids).map((name, index) => (
                <Badge key={index} variant="outline" className="text-xs">
                  {name}
                </Badge>
              ))}
            </div>
          )}
          
          {/* Price */}
          <div className="mb-3">
            <div className="text-lg font-bold">
              {formatPrice(product.price)}
            </div>
            {product.compare_price > 0 && (
              <div className="text-sm text-muted-foreground line-through">
                {formatPrice(product.compare_price)}
              </div>
            )}
          </div>
          
          {/* Actions */}
          <div className="flex gap-2">
            <Button
              className="flex-1"
              onClick={handleAddToCart}
              disabled={!product.is_active || isAddingToCart}
            >
              <ShoppingCart className="h-4 w-4 mr-1" />
              {isAddingToCart ? "Adding..." : "Add to Cart"}
            </Button>
            <Button variant="outline" size="sm" asChild>
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
