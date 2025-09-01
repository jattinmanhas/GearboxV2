"use client"

import { useUserStore } from "@/lib/stores/user-store"
import { Button } from "@/components/ui/button"

export default function DebugUserPage() {
  const { user, isAuthenticated, logout } = useUserStore()

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
          
          <div className="flex gap-4">
            <Button onClick={testLogin}>
              Test Login
            </Button>
            <Button onClick={() => logout()} variant="outline">
              Logout
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
