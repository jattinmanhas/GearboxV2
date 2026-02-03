import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Package } from "lucide-react"
import { BarList } from "@/components/dashboard/charts"
import { formatCurrency } from "@/lib/currency"
import { httpClient } from "@/lib/apiFunctions/http-client"

interface TopProduct {
  product_id: number
  product_name: string
  sku: string
  total_quantity: number
  total_revenue: number
  order_count: number
}

export function TopProductsChart() {
  const [products, setProducts] = useState<TopProduct[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchTopProducts = async () => {
      try {
        setLoading(true)
        const result = await httpClient.get<{
          success: boolean
          message: string
          data: TopProduct[]
        }>('/dashboard/top-products?limit=5')

        if (result.success) {
          setProducts(result.data || [])
        } else {
          throw new Error(result.message || 'Failed to load top products')
        }
      } catch (error) {
        setError(error instanceof Error ? error.message : 'Unknown error')
      } finally {
        setLoading(false)
      }
    }

    fetchTopProducts()
  }, [])

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Top Products</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-24 bg-muted animate-pulse rounded" />
        </CardContent>
      </Card>
    )
  }

  if (error) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Top Products</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-destructive">{error}</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0">
        <CardTitle className="flex items-center gap-2">
          <Package className="h-5 w-5" />
          Top Products
        </CardTitle>
      </CardHeader>
      <CardContent>
        {products.length === 0 ? (
          <p className="text-sm text-muted-foreground">No product data available.</p>
        ) : (
          <BarList
            items={products.map((product) => ({
              label: product.product_name,
              value: product.total_quantity,
              description: `${formatCurrency(product.total_revenue)} • ${product.order_count} orders`,
            }))}
          />
        )}
      </CardContent>
    </Card>
  )
}
