"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/components/ui/table"
import { 
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { 
  Edit, 
  Trash2, 
  Eye, 
  EyeOff, 
  MoreHorizontal,
  ChevronLeft,
  ChevronRight,
  Package,
  TrendingUp
} from "lucide-react"
import { Product, Category } from "@/lib/types"
import { formatPrice } from "@/lib/currency"

interface ProductTableProps {
  products: Product[]
  categories: Category[]
  onEdit: (product: Product) => void
  onDelete: (id: number) => void
  onManageVariants: (product: Product) => void
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
}

export function ProductTable({ 
  products, 
  categories, 
  onEdit, 
  onDelete, 
  onManageVariants,
  currentPage, 
  totalPages, 
  onPageChange 
}: ProductTableProps) {
  const [deletingId, setDeletingId] = useState<number | null>(null)

  const handleDelete = async (id: number) => {
    setDeletingId(id)
    try {
      await onDelete(id)
    } finally {
      setDeletingId(null)
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  }



  if (!products || products.length === 0) {
    return (
      <div className="text-center py-8">
        <div className="text-muted-foreground mb-4">
          No products found. Create your first product to get started.
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="rounded-md border overflow-x-auto">
        <Table className="min-w-full">
          <TableHeader>
            <TableRow className="h-8">
              <TableHead className="w-[160px] text-xs">Product</TableHead>
              <TableHead className="w-[50px] hidden sm:table-cell text-xs">SKU</TableHead>
              <TableHead className="w-[50px] text-xs">Price</TableHead>
              <TableHead className="w-[60px] hidden md:table-cell text-xs">Status</TableHead>
              <TableHead className="w-[40px] text-xs">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {products.map((product) => (
              <TableRow key={product.id} className="h-10">
                <TableCell className="font-medium">
                  <div className="flex items-center space-x-2">
                    <div className="w-8 h-8 bg-gradient-to-br from-blue-100 to-blue-200 rounded flex items-center justify-center flex-shrink-0 overflow-hidden border border-blue-300 relative">
                      {product.images && product.images.length > 0 ? (
                        <img 
                          src={product.images[0].url} 
                          alt={product.images[0].alt || product.name}
                          className="w-full h-full object-cover"
                        />
                      ) : (
                        <div className="flex items-center justify-center w-full h-full">
                          <Package className="h-4 w-4 text-blue-600" />
                        </div>
                      )}
                      {/* Add a small indicator for products without images */}
                      {(!product.images || product.images.length === 0) && (
                        <div className="absolute -top-1 -right-1 w-3 h-3 bg-orange-400 rounded-full border border-white"></div>
                      )}
                    </div>
                    <div className="min-w-0 flex-1">
                      <div className="font-medium truncate text-xs">{product.name}</div>
                    </div>
                  </div>
                </TableCell>
                <TableCell className="hidden sm:table-cell">
                  <code className="text-xs bg-muted px-1 py-0.5 rounded truncate block">
                    {product.sku}
                  </code>
                </TableCell>
                <TableCell>
                  <div className="text-xs font-medium">{formatPrice(product.price)}</div>
                </TableCell>
                <TableCell className="hidden md:table-cell">
                  <Badge variant={(product.is_active ?? true) ? "default" : "secondary"} className="text-xs px-1 py-0">
                    {(product.is_active ?? true) ? "✓" : "✗"}
                  </Badge>
                </TableCell>
                <TableCell>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" className="h-6 w-6 p-0">
                        <span className="sr-only">Open menu</span>
                        <MoreHorizontal className="h-3 w-3" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => onEdit(product)}>
                        <Edit className="mr-2 h-4 w-4" />
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={() => onManageVariants(product)}>
                        <Package className="mr-2 h-4 w-4" />
                        Manage Variants
                      </DropdownMenuItem>
                      <DropdownMenuItem 
                        onClick={() => handleDelete(product.id)}
                        disabled={deletingId === product.id}
                        className="text-destructive"
                      >
                        <Trash2 className="mr-2 h-4 w-4" />
                        {deletingId === product.id ? "Deleting..." : "Delete"}
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* Compact Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(currentPage - 1)}
            disabled={currentPage <= 1}
            className="h-6 px-2"
          >
            <ChevronLeft className="h-3 w-3" />
          </Button>
          <span className="px-2 py-1 text-xs font-medium">
            {currentPage} / {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(currentPage + 1)}
            disabled={currentPage >= totalPages}
            className="h-6 px-2"
          >
            <ChevronRight className="h-3 w-3" />
          </Button>
        </div>
      )}
    </div>
  )
}
