"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import { 
  Filter, 
  X, 
  ChevronDown,
  ChevronUp
} from "lucide-react"
import { Category } from "@/lib/types"

interface ProductFiltersProps {
  categories: Category[]
  selectedCategory: number | undefined
  onCategoryChange: (categoryId: number | undefined) => void
  priceRange: [number, number]
  onPriceChange: (range: [number, number]) => void
  inStock?: boolean | undefined
  onInStockChange: (value: boolean | undefined) => void
  onSale?: boolean | undefined
  onOnSaleChange: (value: boolean | undefined) => void
  isDigital?: boolean | undefined
  onIsDigitalChange: (value: boolean | undefined) => void
}

export function ProductFilters({
  categories,
  selectedCategory,
  onCategoryChange,
  priceRange,
  onPriceChange,
  inStock,
  onInStockChange,
  onSale,
  onOnSaleChange,
  isDigital,
  onIsDigitalChange
}: ProductFiltersProps) {
  const [expandedSections, setExpandedSections] = useState({
    categories: true,
    price: true,
    features: true
  })

  const toggleSection = (section: keyof typeof expandedSections) => {
    setExpandedSections(prev => ({
      ...prev,
      [section]: !prev[section]
    }))
  }

  const handlePriceChange = (type: 'min' | 'max', value: string) => {
    const numValue = parseFloat(value) || 0
    if (type === 'min') {
      onPriceChange([numValue, priceRange[1]])
    } else {
      onPriceChange([priceRange[0], numValue])
    }
  }

  const clearFilters = () => {
    onCategoryChange(undefined)
    onPriceChange([0, 1000])
    onInStockChange(undefined)
    onOnSaleChange(undefined)
    onIsDigitalChange(undefined)
  }

  const hasActiveFilters = selectedCategory !== undefined || priceRange[0] > 0 || priceRange[1] < 1000 || inStock !== undefined || onSale !== undefined || isDigital !== undefined

  return (
    <div className="space-y-4">
      {/* Filter Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold flex items-center gap-2">
          <Filter className="h-5 w-5" />
          Filters
        </h2>
        {hasActiveFilters && (
          <Button variant="ghost" size="sm" onClick={clearFilters}>
            <X className="h-4 w-4 mr-1" />
            Clear
          </Button>
        )}
      </div>

      {/* Categories Filter */}
      <Card>
        <CardHeader 
          className="cursor-pointer"
          onClick={() => toggleSection('categories')}
        >
          <div className="flex items-center justify-between">
            <CardTitle className="text-base">Categories</CardTitle>
            {expandedSections.categories ? (
              <ChevronUp className="h-4 w-4" />
            ) : (
              <ChevronDown className="h-4 w-4" />
            )}
          </div>
        </CardHeader>
        {expandedSections.categories && (
          <CardContent className="pt-0">
            <div className="space-y-2">
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="all-categories"
                  checked={selectedCategory === undefined}
                  onCheckedChange={() => onCategoryChange(undefined)}
                />
                <Label htmlFor="all-categories" className="text-sm">
                  All Categories
                </Label>
              </div>
              {categories.map((category) => (
                <div key={category.id} className="flex items-center space-x-2">
                  <Checkbox
                    id={`category-${category.id}`}
                    checked={selectedCategory === category.id}
                    onCheckedChange={() => 
                      onCategoryChange(selectedCategory === category.id ? undefined : category.id)
                    }
                  />
                  <Label htmlFor={`category-${category.id}`} className="text-sm">
                    {category.name}
                  </Label>
                </div>
              ))}
            </div>
          </CardContent>
        )}
      </Card>

      {/* Price Filter */}
      <Card>
        <CardHeader 
          className="cursor-pointer"
          onClick={() => toggleSection('price')}
        >
          <div className="flex items-center justify-between">
            <CardTitle className="text-base">Price Range</CardTitle>
            {expandedSections.price ? (
              <ChevronUp className="h-4 w-4" />
            ) : (
              <ChevronDown className="h-4 w-4" />
            )}
          </div>
        </CardHeader>
        {expandedSections.price && (
          <CardContent className="pt-0">
            <div className="space-y-4">
              <div className="flex items-center space-x-2">
                <input
                  type="number"
                  placeholder="Min"
                  value={priceRange[0]}
                  onChange={(e) => handlePriceChange('min', e.target.value)}
                  className="w-full px-3 py-2 border border-input bg-background rounded-md text-sm"
                />
                <span className="text-muted-foreground">to</span>
                <input
                  type="number"
                  placeholder="Max"
                  value={priceRange[1]}
                  onChange={(e) => handlePriceChange('max', e.target.value)}
                  className="w-full px-3 py-2 border border-input bg-background rounded-md text-sm"
                />
              </div>
              
              {/* Quick Price Ranges */}
              <div className="space-y-2">
                <div className="text-sm font-medium">Quick Select</div>
                <div className="grid grid-cols-2 gap-2">
                  {[
                    { label: "Under $25", range: [0, 25] as [number, number] },
                    { label: "$25 - $50", range: [25, 50] as [number, number] },
                    { label: "$50 - $100", range: [50, 100] as [number, number] },
                    { label: "Over $100", range: [100, 1000] as [number, number] }
                  ].map((item, index) => (
                    <Button
                      key={index}
                      variant="outline"
                      size="sm"
                      className="text-xs"
                      onClick={() => onPriceChange(item.range)}
                    >
                      {item.label}
                    </Button>
                  ))}
                </div>
              </div>
            </div>
          </CardContent>
        )}
      </Card>

      {/* Features Filter */}
      <Card>
        <CardHeader 
          className="cursor-pointer"
          onClick={() => toggleSection('features')}
        >
          <div className="flex items-center justify-between">
            <CardTitle className="text-base">Features</CardTitle>
            {expandedSections.features ? (
              <ChevronUp className="h-4 w-4" />
            ) : (
              <ChevronDown className="h-4 w-4" />
            )}
          </div>
        </CardHeader>
        {expandedSections.features && (
          <CardContent className="pt-0">
            <div className="space-y-2">
              <div className="flex items-center space-x-2">
                <Checkbox 
                  id="in-stock" 
                  checked={inStock === true}
                  onCheckedChange={(checked) => onInStockChange(checked ? true : undefined)}
                />
                <Label htmlFor="in-stock" className="text-sm">
                  In Stock
                </Label>
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox 
                  id="on-sale" 
                  checked={onSale === true}
                  onCheckedChange={(checked) => onOnSaleChange(checked ? true : undefined)}
                />
                <Label htmlFor="on-sale" className="text-sm">
                  On Sale
                </Label>
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox 
                  id="digital" 
                  checked={isDigital === true}
                  onCheckedChange={(checked) => onIsDigitalChange(checked ? true : undefined)}
                />
                <Label htmlFor="digital" className="text-sm">
                  Digital Products
                </Label>
              </div>
            </div>
          </CardContent>
        )}
      </Card>

      {/* Active Filters */}
      {hasActiveFilters && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Active Filters</CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex flex-wrap gap-2">
              {selectedCategory && (
                <Badge variant="secondary" className="flex items-center gap-1">
                  {categories.find(c => c.id === selectedCategory)?.name}
                  <X 
                    className="h-3 w-3 cursor-pointer" 
                    onClick={() => onCategoryChange(undefined)}
                  />
                </Badge>
              )}
              {(priceRange[0] > 0 || priceRange[1] < 1000) && (
                <Badge variant="secondary" className="flex items-center gap-1">
                  ${priceRange[0]} - ${priceRange[1]}
                  <X 
                    className="h-3 w-3 cursor-pointer" 
                    onClick={() => onPriceChange([0, 1000])}
                  />
                </Badge>
              )}
              {inStock === true && (
                <Badge variant="secondary" className="flex items-center gap-1">
                  In Stock
                  <X 
                    className="h-3 w-3 cursor-pointer" 
                    onClick={() => onInStockChange(undefined)}
                  />
                </Badge>
              )}
              {onSale === true && (
                <Badge variant="secondary" className="flex items-center gap-1">
                  On Sale
                  <X 
                    className="h-3 w-3 cursor-pointer" 
                    onClick={() => onOnSaleChange(undefined)}
                  />
                </Badge>
              )}
              {isDigital === true && (
                <Badge variant="secondary" className="flex items-center gap-1">
                  Digital
                  <X 
                    className="h-3 w-3 cursor-pointer" 
                    onClick={() => onIsDigitalChange(undefined)}
                  />
                </Badge>
              )}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
