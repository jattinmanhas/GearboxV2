"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  CreateInventoryRequest,
  UpdateInventoryRequest,
  Inventory,
  Product,
} from "@/lib/types";
import { productApi } from "@/lib/apiFunctions";
import { Search, X, Loader2 } from "lucide-react";

interface InventoryFormProps {
  inventory?: Inventory | null;
  // products prop removed as we fetch dynamically
  onSave: (data: any) => void;
  onCancel: () => void;
}

export function InventoryForm({
  inventory,
  onSave,
  onCancel,
}: InventoryFormProps) {
  const [formData, setFormData] = useState({
    product_id: 0,
    product_variant_id: 0,
    quantity: 0,
    min_stock_level: 0,
    max_stock_level: 0,
    reorder_point: 0,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [variants, setVariants] = useState<any[]>([]);
  const [loadingVariants, setLoadingVariants] = useState(false);

  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState<Product[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [showResults, setShowResults] = useState(false);

  useEffect(() => {
    if (inventory) {
      setFormData({
        product_id: inventory.product_id,
        product_variant_id: inventory.product_variant_id || 0,
        quantity: inventory.quantity,
        min_stock_level: inventory.min_stock_level,
        max_stock_level: inventory.max_stock_level || 0,
        reorder_point: inventory.reorder_point,
      });

      // If editing, we might need to fetch the product details if we want to show it nicely
      // For now, we can rely on the inventory item having product_name if available
      if (inventory.product_name) {
        setSearchTerm(inventory.product_name);
        // We can optionally fetch the full product here to populate selectedProduct
        productApi
          .getProduct(inventory.product_id)
          .then((p) => setSelectedProduct(p))
          .catch(console.error);
      }
    }
  }, [inventory]);

  // Debounced search
  useEffect(() => {
    const timer = setTimeout(async () => {
      if (searchTerm.trim() && !selectedProduct && searchTerm.length >= 2) {
        setIsSearching(true);
        try {
          console.log("Searching for:", searchTerm);
          const results = await productApi.searchProducts(searchTerm);
          console.log("Search results:", results);
          setSearchResults(results);
          setShowResults(true);
        } catch (error) {
          console.error("Search failed:", error);
        } finally {
          setIsSearching(false);
        }
      } else if (!searchTerm.trim()) {
        setSearchResults([]);
        setShowResults(false);
      }
    }, 500); // 500ms debounce

    return () => clearTimeout(timer);
  }, [searchTerm, selectedProduct]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.product_id) {
      newErrors.product_id = "Product is required";
    }

    if (formData.quantity < 0) {
      newErrors.quantity = "Quantity cannot be negative";
    }

    if (formData.min_stock_level < 0) {
      newErrors.min_stock_level = "Minimum stock level cannot be negative";
    }

    if (
      formData.max_stock_level &&
      formData.max_stock_level < formData.min_stock_level
    ) {
      newErrors.max_stock_level =
        "Maximum stock level must be greater than minimum stock level";
    }

    if (formData.reorder_point < 0) {
      newErrors.reorder_point = "Reorder point cannot be negative";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    setIsSubmitting(true);
    try {
      const submitData = {
        ...formData,
        product_variant_id: formData.product_variant_id || undefined,
        max_stock_level: formData.max_stock_level || undefined,
      };
      await onSave(submitData);
    } catch (error) {
      console.error("Error saving inventory:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleProductSelect = async (product: Product) => {
    setSelectedProduct(product);
    setSearchTerm(product.name);
    setSearchResults([]);
    setShowResults(false);
    setFormData((prev) => ({
      ...prev,
      product_id: product.id,
      product_variant_id: 0,
    }));

    // Fetch variants for the selected product
    setLoadingVariants(true);
    try {
      const variants = await productApi.getProductVariants(product.id);
      setVariants(variants || []);
    } catch (error) {
      console.error("Error fetching variants:", error);
      setVariants([]);
    } finally {
      setLoadingVariants(false);
    }
  };

  const clearSelection = () => {
    setSelectedProduct(null);
    setSearchTerm("");
    setFormData((prev) => ({ ...prev, product_id: 0, product_variant_id: 0 }));
    setVariants([]);
  };

  return (
    <Dialog open={true} onOpenChange={onCancel}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>
            {inventory ? "Edit Inventory" : "Create New Inventory"}
          </DialogTitle>
          <DialogDescription>
            {inventory
              ? "Update the inventory information below."
              : "Fill in the details to create a new inventory record."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Product Selection */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Product Selection</h3>
            <div className="space-y-2 relative">
              <Label htmlFor="product_search">Product *</Label>
              <div className="relative">
                <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  id="product_search"
                  placeholder="Search for a product..."
                  value={searchTerm}
                  onChange={(e) => {
                    setSearchTerm(e.target.value);
                    if (selectedProduct) {
                      setSelectedProduct(null); // Clear selection if user types
                      setFormData((prev) => ({
                        ...prev,
                        product_id: 0,
                        product_variant_id: 0,
                      }));
                    }
                  }}
                  className="pl-8 pr-8"
                  autoComplete="off"
                />
                {isSearching && (
                  <Loader2 className="absolute right-2 top-2.5 h-4 w-4 animate-spin text-muted-foreground" />
                )}
                {searchTerm && !isSearching && (
                  <button
                    type="button"
                    onClick={clearSelection}
                    className="absolute right-2 top-2.5 text-muted-foreground hover:text-foreground"
                  >
                    <X className="h-4 w-4" />
                  </button>
                )}
              </div>

              {/* Search Results Dropdown */}
              {showResults && searchResults.length > 0 && (
                <div className="absolute z-10 w-full mt-1 bg-popover border rounded-md shadow-md max-h-60 overflow-auto">
                  {searchResults.map((product) => (
                    <div
                      key={product.id}
                      className="p-2 hover:bg-accent hover:text-accent-foreground cursor-pointer flex flex-col"
                      onClick={() => handleProductSelect(product)}
                    >
                      <span className="font-medium">{product.name}</span>
                      <span className="text-xs text-muted-foreground">
                        SKU: {product.sku}
                      </span>
                    </div>
                  ))}
                </div>
              )}

              {showResults &&
                searchResults.length === 0 &&
                searchTerm.length >= 2 &&
                !isSearching && (
                  <div className="absolute z-10 w-full mt-1 bg-popover border rounded-md shadow-md p-4 text-center text-muted-foreground">
                    No products found
                  </div>
                )}

              {errors.product_id && (
                <p className="text-sm text-destructive">{errors.product_id}</p>
              )}
            </div>

            {selectedProduct && (
              <div className="space-y-2">
                <Label htmlFor="product_variant_id">
                  Product Variant (Optional)
                </Label>
                <select
                  id="product_variant_id"
                  value={formData.product_variant_id}
                  onChange={(e) =>
                    setFormData((prev) => ({
                      ...prev,
                      product_variant_id: Number(e.target.value),
                    }))
                  }
                  className="w-full p-2 border border-input rounded-md bg-background"
                  disabled={loadingVariants}
                >
                  <option value={0}>
                    {loadingVariants
                      ? "Loading variants..."
                      : "No variant (base product)"}
                  </option>
                  {variants.map((variant) => (
                    <option key={variant.id} value={variant.id}>
                      {variant.name} ({variant.sku})
                    </option>
                  ))}
                </select>
                {loadingVariants && (
                  <p className="text-sm text-muted-foreground">
                    Loading product variants...
                  </p>
                )}
              </div>
            )}
          </div>

          {/* Stock Levels */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Stock Levels</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="quantity">Current Quantity *</Label>
                <Input
                  id="quantity"
                  type="number"
                  value={formData.quantity}
                  onChange={(e) =>
                    setFormData((prev) => ({
                      ...prev,
                      quantity: parseInt(e.target.value) || 0,
                    }))
                  }
                  className={errors.quantity ? "border-destructive" : ""}
                />
                {errors.quantity && (
                  <p className="text-sm text-destructive">{errors.quantity}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="min_stock_level">Minimum Stock Level *</Label>
                <Input
                  id="min_stock_level"
                  type="number"
                  value={formData.min_stock_level}
                  onChange={(e) =>
                    setFormData((prev) => ({
                      ...prev,
                      min_stock_level: parseInt(e.target.value) || 0,
                    }))
                  }
                  className={errors.min_stock_level ? "border-destructive" : ""}
                />
                {errors.min_stock_level && (
                  <p className="text-sm text-destructive">
                    {errors.min_stock_level}
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="max_stock_level">
                  Maximum Stock Level (Optional)
                </Label>
                <Input
                  id="max_stock_level"
                  type="number"
                  value={formData.max_stock_level}
                  onChange={(e) =>
                    setFormData((prev) => ({
                      ...prev,
                      max_stock_level: parseInt(e.target.value) || 0,
                    }))
                  }
                  className={errors.max_stock_level ? "border-destructive" : ""}
                />
                {errors.max_stock_level && (
                  <p className="text-sm text-destructive">
                    {errors.max_stock_level}
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="reorder_point">Reorder Point *</Label>
                <Input
                  id="reorder_point"
                  type="number"
                  value={formData.reorder_point}
                  onChange={(e) =>
                    setFormData((prev) => ({
                      ...prev,
                      reorder_point: parseInt(e.target.value) || 0,
                    }))
                  }
                  className={errors.reorder_point ? "border-destructive" : ""}
                />
                {errors.reorder_point && (
                  <p className="text-sm text-destructive">
                    {errors.reorder_point}
                  </p>
                )}
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting
                ? "Saving..."
                : inventory
                  ? "Update Inventory"
                  : "Create Inventory"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}