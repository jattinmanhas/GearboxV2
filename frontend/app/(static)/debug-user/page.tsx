"use client"

import { useUserStore } from "@/lib/stores/user-store"
import { Button } from "@/components/ui/button"
import { useState } from "react"

export default function DebugUserPage() {
  const { user, isAuthenticated, logout, validateToken, checkAuthStatus, clearUser } = useUserStore()
  const [validationResult, setValidationResult] = useState<string>("")
  const [isValidating, setIsValidating] = useState(false)

  const testLogin = () => {
    const testUser = {
      id: 1,
      username: "testuser",
      email: "test@example.com",
      firstName: "Test",
      lastName: "User",
      avatar: "",
      role: "user",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }
    
    useUserStore.getState().login(testUser)
  }

  const handleValidateToken = async () => {
    setIsValidating(true)
    setValidationResult("")
    
    try {
      const isValid = await validateToken()
      setValidationResult(`Token validation result: ${isValid ? "Valid" : "Invalid"}`)
    } catch (error) {
      setValidationResult(`Token validation error: ${error}`)
    } finally {
      setIsValidating(false)
    }
  }

  const handleCheckAuthStatus = async () => {
    setIsValidating(true)
    setValidationResult("")
    
    try {
      const isAuthValid = await checkAuthStatus()
      setValidationResult(`Auth status check result: ${isAuthValid ? "Authenticated and valid" : "Not authenticated or invalid"}`)
    } catch (error) {
      setValidationResult(`Auth status check error: ${error}`)
    } finally {
      setIsValidating(false)
    }
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-2xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">Debug User State</h1>
        
        <div className="space-y-4 p-6 border rounded-lg">
          <div>
            <h2 className="text-xl font-semibold mb-2">Current State</h2>
            <p><strong>Is Authenticated:</strong> {isAuthenticated ? "Yes" : "No"}</p>
            <p><strong>User:</strong> {user ? JSON.stringify(user, null, 2) : "null"}</p>
          </div>
          
          <div className="space-y-4">
            <div className="flex gap-4 flex-wrap">
              <Button onClick={testLogin}>
                Test Login
              </Button>
              <Button onClick={() => logout()} variant="outline">
                Logout
              </Button>
              <Button onClick={handleValidateToken} disabled={isValidating}>
                {isValidating ? "Validating..." : "Validate Token"}
              </Button>
              <Button onClick={handleCheckAuthStatus} disabled={isValidating}>
                {isValidating ? "Checking..." : "Check Auth Status"}
              </Button>
              <Button onClick={() => clearUser()} variant="destructive">
                Clear User
              </Button>
            </div>
            
            {validationResult && (
              <div className="p-4 bg-gray-100 rounded-lg">
                <p className="text-sm font-mono">{validationResult}</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
