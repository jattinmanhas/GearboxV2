'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { CheckCircle2, Loader2 } from 'lucide-react'
import { useUserStore } from '@/lib/stores/user-store'
import { profileApi } from '@/lib/api'
import { showError } from '@/lib/notifications'

export default function OAuthSuccessPage() {
  const router = useRouter()
  const [loading, setLoading] = useState(true)
  const [mounted, setMounted] = useState(false)
  const { login } = useUserStore()

  // Prevent hydration mismatch by ensuring component is mounted
  useEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    const handleSuccess = async () => {
      try {
        // Cookies are already set by the backend
        // Fetch user profile to update the store
        const profileResponse: any = await profileApi.getProfile()
        
        // Handle both response formats: direct profile or wrapped in data
        const profile = profileResponse.data || profileResponse
        
        if (profile) {
          // Convert profile to user format expected by store
          const userData = {
            id: profile.id,
            username: profile.username,
            email: profile.email,
            firstName: profile.first_name,
            middleName: profile.middle_name || '',
            lastName: profile.last_name || '',
            avatar: profile.avatar || '',
            role: 'user', // Profile doesn't include role, will be fetched separately if needed
            createdAt: profile.created_at,
            updatedAt: profile.updated_at,
          }
          // Use the same login flow as normal login to ensure cart merging and state setup
          await login(userData)
        }

        // Wait a moment to show success message
        setTimeout(() => {
          // Redirect to dashboard
          router.push('/')
        }, 500)
      } catch (error) {
        showError('Failed to fetch user profile', {
          description: error instanceof Error ? error.message : 'Failed to fetch user profile',
          duration: 5000,
        })
        // Even if profile fetch fails, redirect to dashboard
        // as auth cookies should still be valid
        setTimeout(() => {
          router.push('/')
        }, 500)
      } finally {
        setLoading(false)
      }
    }

    handleSuccess()
  }, [router, login])

  // Prevent hydration mismatch - don't render until mounted
  if (!mounted) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center px-4">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-green-100 dark:bg-green-900">
              <Loader2 className="h-6 w-6 animate-spin text-green-600 dark:text-green-400" />
            </div>
            <CardTitle className="text-2xl font-bold">Authentication Successful!</CardTitle>
            <CardDescription>
              Setting up your account...
            </CardDescription>
          </CardHeader>
          <CardContent className="text-center text-sm text-muted-foreground">
            <p>You have successfully authenticated with your OAuth provider.</p>
            <p className="mt-2">Please wait while we redirect you...</p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-green-100 dark:bg-green-900">
            {loading ? (
              <Loader2 className="h-6 w-6 animate-spin text-green-600 dark:text-green-400" />
            ) : (
              <CheckCircle2 className="h-6 w-6 text-green-600 dark:text-green-400" />
            )}
          </div>
          <CardTitle className="text-2xl font-bold">Authentication Successful!</CardTitle>
          <CardDescription>
            {loading 
              ? 'Setting up your account...'
              : 'Redirecting to your dashboard...'}
          </CardDescription>
        </CardHeader>
        <CardContent className="text-center text-sm text-muted-foreground">
          <p>You have successfully authenticated with your OAuth provider.</p>
          <p className="mt-2">Please wait while we redirect you...</p>
        </CardContent>
      </Card>
    </div>
  )
}