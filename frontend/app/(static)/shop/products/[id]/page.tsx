"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import Image from "next/image";
import Link from "next/link";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import {
  ArrowLeft,
  Heart,
  ShoppingCart,
  Star,
  Package,
  Truck,
  Shield,
  RotateCcw,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { productApi } from "@/lib/apiFunctions";
import { Product, Category, ProductVariant } from "@/lib/types";
import { useCartStore } from "@/lib/stores/cart-store";
import { useWishlistStore } from "@/lib/stores/wishlist-store";
import { formatPrice } from "@/lib/currency";

export default function ProductDetailPage() {
  const params = useParams();
  const productId = params.id as string;

  const [product, setProduct] = useState<Product | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);
  const [variants, setVariants] = useState<ProductVariant[]>([]);
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(
    null
  );
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);

  // Debug: Log selectedVariant changes
  useEffect(() => {
    console.log("Selected variant changed:", selectedVariant);
    if (selectedVariant) {
      console.log("Selected variant is_in_stock:", selectedVariant.is_in_stock);
      console.log("Selected variant is_active:", selectedVariant.is_active);
    }
  }, [selectedVariant]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [isAddingToCart, setIsAddingToCart] = useState(false);
  const [isWishlisting, setIsWishlisting] = useState(false);

  const { addItem, isLoading: cartLoading } = useCartStore();
  const {
    addItemToWishlist,
    removeItemFromWishlist,
    isProductInWishlist,
    wishlists,
    loadWishlists,
  } = useWishlistStore();

  // Check if product is in any wishlist
  const isWishlisted = product
    ? wishlists.some((wishlist) => isProductInWishlist(product.id, wishlist.id))
    : false;

  useEffect(() => {
    const loadProduct = async () => {
      try {
        setLoading(true);
        setError(null);
        // Fetch product and categories first
        const [productData, categoriesData] = await Promise.all([
          productApi.getProduct(parseInt(productId)),
          productApi.getCategories({
            page: 1,
            limit: 10,
          }),
        ]);
        setProduct(productData);
        setCategories(categoriesData.data.categories);
        setSelectedImageIndex(0);

        // Fetch variants separately; ignore errors that may arise from zero inventory
        let variantsData: ProductVariant[] = [];
        try {
          variantsData = await productApi.getProductVariantsWithInventory(
            parseInt(productId)
          );
        } catch (variantErr) {
          console.error(
            "Failed to fetch variants (may be zero inventory):",
            variantErr
          );
          // Keep variantsData as empty array
        }
        setVariants(variantsData);

        // Debug: Log the variants data
        console.log("Variants data:", variantsData);
        console.log("Product data:", productData);

        // Set the first in-stock and active variant as selected, or the first active variant, or the first variant
        const inStockActiveVariants = variantsData.filter(
          (v) => v.is_active && v.is_in_stock
        );
        const activeVariants = variantsData.filter((v) => v.is_active);

        console.log("In stock active variants:", inStockActiveVariants);
        console.log("Active variants:", activeVariants);

        if (inStockActiveVariants.length > 0) {
          console.log(
            "Setting selected variant to in-stock variant:",
            inStockActiveVariants[0]
          );
          setSelectedVariant(inStockActiveVariants[0]);
        } else if (activeVariants.length > 0) {
          console.log(
            "Setting selected variant to active variant:",
            activeVariants[0]
          );
          setSelectedVariant(activeVariants[0]);
        } else if (variantsData.length > 0) {
          console.log(
            "Setting selected variant to first variant:",
            variantsData[0]
          );
          setSelectedVariant(variantsData[0]);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load product");
      } finally {
        setLoading(false);
      }
    };

    if (productId) {
      loadProduct();
    }
  }, [productId]);

  // No longer needed - using category_names from API response

  const handleAddToCart = async () => {
    if (!product) return;

    setIsAddingToCart(true);
    try {
      await addItem({
        product_id: product.id,
        product_variant_id: selectedVariant?.id,
        quantity: quantity,
      });
    } catch (error) {
      console.error("Failed to add to cart:", error);
    } finally {
      setIsAddingToCart(false);
    }
  };

  const handleWishlist = async () => {
    if (!product || isWishlisting) return;

    setIsWishlisting(true);
    try {
      // Load wishlists if not already loaded
      if (wishlists.length === 0) {
        await loadWishlists();
      }

      if (isWishlisted) {
        // Find the wishlist containing this product and remove it
        const wishlistWithProduct = wishlists.find((wishlist) =>
          isProductInWishlist(product.id, wishlist.id)
        );
        if (wishlistWithProduct) {
          const item = wishlistWithProduct.items.find(
            (item) => item.product_id === product.id
          );
          if (item) {
            await removeItemFromWishlist(wishlistWithProduct.id, item.id);
          }
        }
      } else {
        // Add to the first wishlist (or create a default one)
        if (wishlists.length > 0) {
          await addItemToWishlist(wishlists[0].id, product.id);
        } else {
          // Create a default wishlist first
          // TODO: Implement default wishlist creation
          console.log("No wishlists available, need to create one first");
        }
      }
    } catch (error) {
      console.error("Failed to update wishlist:", error);
    } finally {
      setIsWishlisting(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">Loading product...</div>
          </div>
        </div>
      </div>
    );
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
    );
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
            <div className="relative aspect-square bg-muted rounded-lg flex items-center justify-center overflow-hidden group">
              {product.images && product.images.length > 0 ? (
                <Image
                  src={product.images[selectedImageIndex].url}
                  alt={product.images[selectedImageIndex].alt || product.name}
                  width={600}
                  height={600}
                  className="w-full h-full object-cover"
                />
              ) : (
                <Package className="h-24 w-24 text-muted-foreground" />
              )}

              {/* Navigation arrows */}
              {product.images && product.images.length > 1 && (
                <>
                  <Button
                    variant="outline"
                    size="sm"
                    className="absolute left-2 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-opacity bg-background/90 hover:bg-background border-border shadow-lg"
                    onClick={() =>
                      setSelectedImageIndex(
                        selectedImageIndex === 0
                          ? (product.images?.length || 1) - 1
                          : selectedImageIndex - 1
                      )
                    }
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    className="absolute right-2 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-opacity bg-background/90 hover:bg-background border-border shadow-lg"
                    onClick={() =>
                      setSelectedImageIndex(
                        selectedImageIndex === (product.images?.length || 1) - 1
                          ? 0
                          : selectedImageIndex + 1
                      )
                    }
                  >
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </>
              )}

              {/* Image counter */}
              {product.images && product.images.length > 1 && (
                <div className="absolute bottom-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
                  {selectedImageIndex + 1} / {product.images.length}
                </div>
              )}
            </div>

            {/* Thumbnail images */}
            {product.images && product.images.length > 1 && (
              <div className="flex gap-2">
                {product.images.slice(0, 4).map((image, index) => (
                  <div
                    key={image.id}
                    className={`w-16 h-16 bg-muted rounded-md flex items-center justify-center overflow-hidden cursor-pointer transition-all duration-200 ${selectedImageIndex === index
                        ? "ring-2 ring-primary ring-offset-2"
                        : "hover:ring-2 hover:ring-primary"
                      }`}
                    onClick={() => setSelectedImageIndex(index)}
                  >
                    <Image
                      src={image.url}
                      alt={image.alt || product.name}
                      width={64}
                      height={64}
                      className="w-full h-full object-cover"
                    />
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Product Info */}
          <div className="space-y-6">
            {/* Header */}
            <div>
              <div className="flex items-center gap-2 mb-2">
                <h1 className="text-3xl font-bold">{product.name}</h1>
                {!product.is_in_stock && (
                  <Badge variant="destructive">Out of Stock</Badge>
                )}
              </div>

              <div className="flex items-center gap-4 mb-4">
                <div className="flex items-center">
                  {[...Array(5)].map((_, i) => (
                    <Star
                      key={i}
                      className="h-4 w-4 fill-yellow-400 text-yellow-400"
                    />
                  ))}
                  <span className="text-sm text-muted-foreground ml-2">
                    (4.5) 23 reviews
                  </span>
                </div>
                <span className="text-sm text-muted-foreground">•</span>
                <span className="text-sm text-muted-foreground">
                  SKU: {product.sku}
                </span>
              </div>

              {/* Price */}
              <div className="mb-4">
                <div className="text-3xl font-bold">
                  {selectedVariant
                    ? formatPrice(selectedVariant.price)
                    : formatPrice(product.price)}
                </div>
                {(selectedVariant
                  ? selectedVariant.compare_price
                  : product.compare_price) > 0 && (
                    <div className="text-lg text-muted-foreground line-through">
                      {formatPrice(
                        selectedVariant
                          ? selectedVariant.compare_price
                          : product.compare_price
                      )}
                    </div>
                  )}
              </div>

              {/* Stock Information */}
              <div className="mb-4">
                <div className="flex items-center gap-2">
                  {selectedVariant ? (
                    <>
                      <Badge
                        variant={
                          selectedVariant.is_in_stock
                            ? "default"
                            : "destructive"
                        }
                      >
                        {selectedVariant.is_in_stock
                          ? "In Stock"
                          : "Out of Stock"}
                      </Badge>
                      {selectedVariant.is_in_stock && (
                        <span className="text-sm text-muted-foreground">
                          {selectedVariant.available_quantity} available
                        </span>
                      )}
                    </>
                  ) : (
                    <>
                      <Badge
                        variant={
                          product.is_in_stock ? "default" : "destructive"
                        }
                      >
                        {product.is_in_stock ? "In Stock" : "Out of Stock"}
                      </Badge>
                      {product.is_in_stock && (
                        <span className="text-sm text-muted-foreground">
                          {product.available_quantity} available
                        </span>
                      )}
                    </>
                  )}
                </div>
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
                    .filter((variant) => variant.is_active)
                    .sort((a, b) => a.position - b.position)
                    .map((variant) => (
                      <Button
                        key={variant.id}
                        variant={
                          selectedVariant?.id === variant.id
                            ? "default"
                            : "outline"
                        }
                        size="sm"
                        onClick={() => setSelectedVariant(variant)}
                        className={`flex items-center gap-2 ${!variant.is_in_stock
                            ? "opacity-60 cursor-not-allowed"
                            : ""
                          }`}
                        disabled={!variant.is_in_stock}
                      >
                        {variant.name}
                        <span className="text-xs text-muted-foreground">
                          ({formatPrice(variant.price)})
                        </span>
                        {!variant.is_in_stock && (
                          <span className="text-xs text-red-500 ml-1">
                            (Out of Stock)
                          </span>
                        )}
                      </Button>
                    ))}
                </div>
                {selectedVariant && (
                  <div className="mt-2 text-sm text-muted-foreground">
                    Selected:{" "}
                    <span className="font-medium">{selectedVariant.name}</span>
                    <span className="ml-2">•</span>
                    <span className="ml-2">SKU: {selectedVariant.sku}</span>
                    {!selectedVariant.is_in_stock && (
                      <span className="ml-2 text-red-500">• Out of Stock</span>
                    )}
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

            {/* Stock Status */}
            <div className="space-y-2">
              {(() => {
                const isInStock =
                  product.is_digital ||
                  (variants.length > 0
                    ? selectedVariant?.is_in_stock ?? false
                    : product.is_in_stock);

                return (
                  <div
                    className={`flex items-center gap-2 p-3 rounded-lg ${isInStock
                        ? "bg-green-50 dark:bg-green-950/20 border border-green-200 dark:border-green-800"
                        : "bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-800"
                      }`}
                  >
                    <div
                      className={`w-2 h-2 rounded-full ${isInStock ? "bg-green-500" : "bg-red-500"
                        }`}
                    />
                    <span
                      className={`text-sm font-medium ${isInStock
                          ? "text-green-700 dark:text-green-300"
                          : "text-red-700 dark:text-red-300"
                        }`}
                    >
                      {isInStock ? "In Stock" : "Out of Stock"}
                    </span>
                    {isInStock && product.available_quantity > 0 && (
                      <span className="text-xs text-muted-foreground">
                        ({product.available_quantity} available)
                      </span>
                    )}
                  </div>
                );
              })()}
            </div>

            {/* Quantity and Add to Cart */}
            <div className="space-y-4">
              {(() => {
                const isInStock =
                  product.is_digital ||
                  (variants.length > 0
                    ? selectedVariant?.is_in_stock ?? false
                    : product.is_in_stock);

                if (!isInStock) {
                  return (
                    <div className="text-center py-4">
                      <p className="text-muted-foreground mb-4">
                        This product is currently out of stock.
                      </p>
                      <Button variant="outline" disabled>
                        <ShoppingCart className="h-4 w-4 mr-2" />
                        Out of Stock
                      </Button>
                    </div>
                  );
                }

                return (
                  <>
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
                        disabled={
                          !isInStock ||
                          !product.is_active ||
                          isAddingToCart ||
                          cartLoading
                        }
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
                        <Heart
                          className={`h-4 w-4 ${isWishlisted ? "fill-current" : ""
                            }`}
                        />
                      </Button>
                    </div>
                  </>
                );
              })()}
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
                  <h3 className="text-lg font-semibold mb-3">
                    Product Details
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="font-medium">SKU:</span> {product.sku}
                    </div>
                    <div>
                      <span className="font-medium">Weight:</span>{" "}
                      {product.weight} lbs
                    </div>
                    {product.dimensions && (
                      <div>
                        <span className="font-medium">Dimensions:</span>{" "}
                        {product.dimensions}
                      </div>
                    )}
                    <div>
                      <span className="font-medium">Type:</span>{" "}
                      {product.is_digital ? "Digital" : "Physical"}
                    </div>
                    <div>
                      <span className="font-medium">Shipping:</span>{" "}
                      {product.requires_shipping ? "Required" : "Not Required"}
                    </div>
                    <div>
                      <span className="font-medium">Taxable:</span>{" "}
                      {product.taxable ? "Yes" : "No"}
                    </div>
                  </div>
                </div>

                {product.tags && (
                  <div>
                    <h3 className="text-lg font-semibold mb-3">Tags</h3>
                    <div className="flex flex-wrap gap-2">
                      {product.tags.split(",").map((tag, index) => (
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
  );
}