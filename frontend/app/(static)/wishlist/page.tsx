"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { 
  Heart, 
  ShoppingCart, 
  Trash2, 
  Plus,
  Package,
  Eye,
  Star,
  Loader2
} from "lucide-react"
import { useWishlistStore } from "@/lib/stores/wishlist-store"
import { useCartStore } from "@/lib/stores/cart-store"
import { formatCurrency } from "@/lib/currency"
import Link from "next/link"

export default function WishlistPage() {
  const [newWishlistName, setNewWishlistName] = useState("")
  const [isCreatingWishlist, setIsCreatingWishlist] = useState(false)
  const [showCreateForm, setShowCreateForm] = useState(false)

  const { 
    wishlists, 
    isLoading, 
    error, 
    loadWishlists, 
    createWishlist, 
    deleteWishlist,
    removeItemFromWishlist,
    moveItemToCart 
  } = useWishlistStore()

  const { addItem } = useCartStore()

  useEffect(() => {
    loadWishlists()
  }, [loadWishlists])

  const handleCreateWishlist = async () => {
    if (!newWishlistName.trim()) return
    
    setIsCreatingWishlist(true)
    try {
      await createWishlist({
        name: newWishlistName.trim(),
        description: "",
        is_public: false
      })
      setNewWishlistName("")
      setShowCreateForm(false)
    } catch (error) {
      console.error("Failed to create wishlist:", error)
    } finally {
      setIsCreatingWishlist(false)
    }
  }

  const handleDeleteWishlist = async (wishlistId: number) => {
    if (confirm("Are you sure you want to delete this wishlist?")) {
      await deleteWishlist(wishlistId)
    }
  }

  const handleRemoveItem = async (wishlistId: number, itemId: number) => {
    await removeItemFromWishlist(wishlistId, itemId)
  }

  const handleMoveToCart = async (itemId: number) => {
    try {
      await moveItemToCart(itemId)
    } catch (error) {
      console.error("Failed to move item to cart:", error)
    }
  }

  const handleAddToCart = async (productId: number) => {
    try {
      await addItem({
        product_id: productId,
        quantity: 1
      })
    } catch (error) {
      console.error("Failed to add to cart:", error)
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-center h-64">
            <Loader2 className="h-6 w-6 animate-spin" />
            <span className="ml-2">Loading wishlists...</span>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-6">
          <div className="text-center py-8">
            <div className="text-destructive mb-4">{error}</div>
            <Button onClick={loadWishlists}>Try Again</Button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">My Wishlists</h1>
              <p className="text-muted-foreground">
                Save your favorite products for later
              </p>
            </div>
            <Button onClick={() => setShowCreateForm(!showCreateForm)}>
              <Plus className="h-4 w-4 mr-2" />
              New Wishlist
            </Button>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6">
        {/* Create Wishlist Form */}
        {showCreateForm && (
          <Card className="mb-6">
            <CardHeader>
              <CardTitle>Create New Wishlist</CardTitle>
              <CardDescription>
                Give your wishlist a name to get started
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex gap-2">
                <Input
                  placeholder="Wishlist name"
                  value={newWishlistName}
                  onChange={(e) => setNewWishlistName(e.target.value)}
                  className="flex-1"
                />
                <Button 
                  onClick={handleCreateWishlist}
                  disabled={!newWishlistName.trim() || isCreatingWishlist}
                >
                  {isCreatingWishlist ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    "Create"
                  )}
                </Button>
                <Button 
                  variant="outline"
                  onClick={() => setShowCreateForm(false)}
                >
                  Cancel
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Wishlists */}
        {wishlists.length === 0 ? (
          <div className="text-center py-12">
            <Heart className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
            <h3 className="text-lg font-semibold mb-2">No wishlists yet</h3>
            <p className="text-muted-foreground mb-4">
              Create your first wishlist to start saving products
            </p>
            <Button onClick={() => setShowCreateForm(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Create Wishlist
            </Button>
          </div>
        ) : (
          <div className="space-y-6">
            {wishlists.map((wishlist) => (
              <Card key={wishlist.id}>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle className="flex items-center gap-2">
                        <Heart className="h-5 w-5 text-red-500" />
                        {wishlist.name}
                      </CardTitle>
                      {wishlist.description && (
                        <CardDescription className="mt-1">
                          {wishlist.description}
                        </CardDescription>
                      )}
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant="outline">
                        {wishlist.items.length} items
                      </Badge>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleDeleteWishlist(wishlist.id)}
                        className="text-destructive hover:text-destructive"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  {wishlist.items.length === 0 ? (
                    <div className="text-center py-8 text-muted-foreground">
                      <Package className="h-8 w-8 mx-auto mb-2" />
                      <p>No items in this wishlist</p>
                    </div>
                  ) : (
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                      {wishlist.items.map((item) => (
                        <Card key={item.id} className="group">
                          <CardContent className="p-4">
                            {/* Product Image */}
                            <div className="aspect-square bg-muted rounded-md flex items-center justify-center mb-3">
                              <Package className="h-8 w-8 text-muted-foreground" />
                            </div>
                            
                            {/* Product Info */}
                            <div className="space-y-2">
                              <h4 className="font-medium text-sm line-clamp-2">
                                {item.product_name}
                              </h4>
                              <p className="text-xs text-muted-foreground">
                                SKU: {item.product_sku}
                              </p>
                              
                              {/* Price */}
                              <div className="text-lg font-bold">
                                {formatCurrency(item.price)}
                              </div>
                              
                              {/* Actions */}
                              <div className="flex gap-1">
                                <Button
                                  size="sm"
                                  className="flex-1"
                                  onClick={() => handleAddToCart(item.product_id)}
                                >
                                  <ShoppingCart className="h-3 w-3 mr-1" />
                                  Add to Cart
                                </Button>
                                <Button
                                  variant="outline"
                                  size="sm"
                                  onClick={() => handleMoveToCart(item.id)}
                                >
                                  <Eye className="h-3 w-3" />
                                </Button>
                                <Button
                                  variant="outline"
                                  size="sm"
                                  onClick={() => handleRemoveItem(wishlist.id, item.id)}
                                  className="text-destructive hover:text-destructive"
                                >
                                  <Trash2 className="h-3 w-3" />
                                </Button>
                              </div>
                            </div>
                          </CardContent>
                        </Card>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
