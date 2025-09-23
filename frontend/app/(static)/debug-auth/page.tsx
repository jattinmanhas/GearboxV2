"use client"

import { useState, useEffect } from 'react'
import { useUserStore } from '@/lib/stores/user-store'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'

export default function DebugAuthPage() {
  const { user, isAuthenticated, debugAuthState, validateToken, checkAuthStatus } = useUserStore()
  const [authStatus, setAuthStatus] = useState<boolean | null>(null)
  const [tokenValid, setTokenValid] = useState<boolean | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const handleDebugAuth = () => {
    debugAuthState()
  }

  const handleValidateToken = async () => {
    setIsLoading(true)
    try {
      const isValid = await validateToken()
      setTokenValid(isValid)
    } catch (error) {
      console.error('Token validation error:', error)
      setTokenValid(false)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCheckAuthStatus = async () => {
    setIsLoading(true)
    try {
      const status = await checkAuthStatus()
      setAuthStatus(status)
    } catch (error) {
      console.error('Auth status check error:', error)
      setAuthStatus(false)
    } finally {
      setIsLoading(false)
    }
  }

  const testApiCall = async () => {
    setIsLoading(true)
    try {
      const response = await fetch('/api/v1/auth/profile', {
        method: 'GET',
        credentials: 'include'
      })
      
      console.log('API Call Result:', {
        status: response.status,
        ok: response.ok,
        statusText: response.statusText
      })
      
      if (response.ok) {
        const data = await response.json()
        console.log('Profile data:', data)
      } else {
        const errorData = await response.json()
        console.log('Error data:', errorData)
      }
    } catch (error) {
      console.error('API call error:', error)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="container mx-auto py-8 px-4">
      <div className="max-w-4xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Authentication Debug Page</h1>
          <p className="text-muted-foreground mt-2">
            This page helps debug authentication issues and test the improved 401 error handling.
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          {/* Current Auth State */}
          <Card>
            <CardHeader>
              <CardTitle>Current Authentication State</CardTitle>
              <CardDescription>
                Real-time view of the authentication state
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="font-medium">Authenticated:</span>
                <Badge variant={isAuthenticated ? "default" : "secondary"}>
                  {isAuthenticated ? "Yes" : "No"}
                </Badge>
              </div>
              
              {user && (
                <>
                  <Separator />
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="font-medium">User ID:</span>
                      <span className="text-sm text-muted-foreground">{user.id}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="font-medium">Username:</span>
                      <span className="text-sm text-muted-foreground">{user.username}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="font-medium">Email:</span>
                      <span className="text-sm text-muted-foreground">{user.email}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="font-medium">Role:</span>
                      <span className="text-sm text-muted-foreground">{user.role || 'user'}</span>
                    </div>
                  </div>
                </>
              )}
            </CardContent>
          </Card>

          {/* Token Validation */}
          <Card>
            <CardHeader>
              <CardTitle>Token Validation</CardTitle>
              <CardDescription>
                Test token validity and authentication status
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="font-medium">Token Valid:</span>
                <Badge variant={tokenValid === null ? "outline" : tokenValid ? "default" : "destructive"}>
                  {tokenValid === null ? "Not tested" : tokenValid ? "Valid" : "Invalid"}
                </Badge>
              </div>
              
              <div className="flex items-center justify-between">
                <span className="font-medium">Auth Status:</span>
                <Badge variant={authStatus === null ? "outline" : authStatus ? "default" : "destructive"}>
                  {authStatus === null ? "Not tested" : authStatus ? "Valid" : "Invalid"}
                </Badge>
              </div>

              <Separator />

              <div className="space-y-2">
                <Button 
                  onClick={handleValidateToken} 
                  disabled={isLoading}
                  className="w-full"
                >
                  {isLoading ? "Validating..." : "Validate Token"}
                </Button>
                
                <Button 
                  onClick={handleCheckAuthStatus} 
                  disabled={isLoading}
                  variant="outline"
                  className="w-full"
                >
                  {isLoading ? "Checking..." : "Check Auth Status"}
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Debug Actions */}
        <Card>
          <CardHeader>
            <CardTitle>Debug Actions</CardTitle>
            <CardDescription>
              Actions to help debug authentication issues
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-2 md:grid-cols-3">
              <Button onClick={handleDebugAuth} variant="outline">
                Log Auth State
              </Button>
              
              <Button onClick={testApiCall} disabled={isLoading} variant="outline">
                {isLoading ? "Testing..." : "Test API Call"}
              </Button>
              
              <Button 
                onClick={() => window.location.reload()} 
                variant="outline"
              >
                Reload Page
              </Button>
            </div>
            
            <div className="text-sm text-muted-foreground">
              <p><strong>Instructions:</strong></p>
              <ul className="list-disc list-inside mt-2 space-y-1">
                <li>Check the browser console for detailed logs</li>
                <li>Test API calls to see how 401 errors are handled</li>
                <li>Verify that the navbar shows the correct authentication state</li>
                <li>Try accessing protected routes to test the flow</li>
              </ul>
            </div>
          </CardContent>
        </Card>

        {/* Test Scenarios */}
        <Card>
          <CardHeader>
            <CardTitle>Test Scenarios</CardTitle>
            <CardDescription>
              Scenarios to test the improved 401 error handling
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="p-4 border rounded-lg">
                <h4 className="font-medium mb-2">Scenario 1: Valid Authentication</h4>
                <p className="text-sm text-muted-foreground mb-2">
                  If you're logged in, the navbar should show your user menu and the debug page should show "Authenticated: Yes"
                </p>
                <Badge variant="outline">Expected: User menu visible, isAuthenticated = true</Badge>
              </div>
              
              <div className="p-4 border rounded-lg">
                <h4 className="font-medium mb-2">Scenario 2: 401 Error Handling</h4>
                <p className="text-sm text-muted-foreground mb-2">
                  When a 401 error occurs, the system should only clear user state if it's a genuine authentication failure, not for other 401 scenarios
                </p>
                <Badge variant="outline">Expected: Smart 401 handling, no unnecessary logout</Badge>
              </div>
              
              <div className="p-4 border rounded-lg">
                <h4 className="font-medium mb-2">Scenario 3: Token Expiration</h4>
                <p className="text-sm text-muted-foreground mb-2">
                  If your token is actually expired, the system should clear your state and redirect to login
                </p>
                <Badge variant="outline">Expected: Clear state only on genuine auth failure</Badge>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
