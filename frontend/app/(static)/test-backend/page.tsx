"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export default function TestBackendPage() {
  const [result, setResult] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const testBackendConnection = async () => {
    setLoading(true)
    setError(null)
    try {
      // Test direct backend connection
      const response = await fetch('http://localhost:8080/health', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })
      
      const data = await response.json()
      setResult({
        status: response.status,
        data: data,
        url: 'http://localhost:8080/health'
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
      console.error('Backend test error:', err)
    } finally {
      setLoading(false)
    }
  }

  const testProductsAPI = async () => {
    setLoading(true)
    setError(null)
    try {
      // Test through Next.js API proxy
      const response = await fetch('/api/v1/products?page=1&limit=5', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })
      
      const data = await response.json()
      setResult({
        status: response.status,
        data: data,
        url: '/api/v1/products'
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
      console.error('Products API test error:', err)
    } finally {
      setLoading(false)
    }
  }

  const testWishlistsAPI = async () => {
    setLoading(true)
    setError(null)
    try {
      // Test wishlist API
      const response = await fetch('/api/v1/wishlists', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include'
      })
      
      const data = await response.json()
      setResult({
        status: response.status,
        data: data,
        url: '/api/v1/wishlists'
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
      console.error('Wishlist API test error:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="container mx-auto max-w-4xl">
        <h1 className="text-3xl font-bold mb-6">Backend Connection Test</h1>
        
        <div className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Test Backend Services</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex gap-4">
                <Button 
                  onClick={testBackendConnection}
                  disabled={loading}
                >
                  Test Direct Backend (Port 8080)
                </Button>
                <Button 
                  onClick={testProductsAPI}
                  disabled={loading}
                  variant="outline"
                >
                  Test Products API
                </Button>
                <Button 
                  onClick={testWishlistsAPI}
                  disabled={loading}
                  variant="outline"
                >
                  Test Wishlists API
                </Button>
              </div>
              
              {loading && (
                <div className="text-blue-600">Testing connection...</div>
              )}
              
              {error && (
                <div className="text-red-600 bg-red-50 p-4 rounded">
                  <strong>Error:</strong> {error}
                </div>
              )}
              
              {result && (
                <div className="bg-green-50 p-4 rounded">
                  <strong>Response from {result.url}:</strong>
                  <div className="mt-2">
                    <strong>Status:</strong> {result.status}
                  </div>
                  <pre className="mt-2 text-sm overflow-auto">
                    {JSON.stringify(result.data, null, 2)}
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
