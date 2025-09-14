"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { wishlistApi } from "@/lib/api"

export default function DebugWishlistPage() {
  const [result, setResult] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const testGetWishlists = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await wishlistApi.getWishlists()
      setResult(response)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
      console.error('Wishlist test error:', err)
    } finally {
      setLoading(false)
    }
  }

  const testCreateWishlist = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await wishlistApi.createWishlist({
        name: "Test Wishlist",
        description: "Test description",
        is_public: false
      })
      setResult(response)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
      console.error('Create wishlist test error:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="container mx-auto max-w-4xl">
        <h1 className="text-3xl font-bold mb-6">Wishlist API Debug</h1>
        
        <div className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Test Wishlist API</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex gap-4">
                <Button 
                  onClick={testGetWishlists}
                  disabled={loading}
                >
                  Test Get Wishlists
                </Button>
                <Button 
                  onClick={testCreateWishlist}
                  disabled={loading}
                  variant="outline"
                >
                  Test Create Wishlist
                </Button>
              </div>
              
              {loading && (
                <div className="text-blue-600">Loading...</div>
              )}
              
              {error && (
                <div className="text-red-600 bg-red-50 p-4 rounded">
                  <strong>Error:</strong> {error}
                </div>
              )}
              
              {result && (
                <div className="bg-green-50 p-4 rounded">
                  <strong>Success:</strong>
                  <pre className="mt-2 text-sm overflow-auto">
                    {JSON.stringify(result, null, 2)}
                  </pre>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
