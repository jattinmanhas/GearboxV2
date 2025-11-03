'use client'

import { useSearchParams } from 'next/navigation'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { AlertCircle } from 'lucide-react'
import Link from 'next/link'

export default function OAuthErrorPage() {
  const searchParams = useSearchParams()
  const error = searchParams.get('message') || 'An unknown error occurred'

  // Parse the error message to make it more user-friendly
  const getErrorMessage = (errorStr: string) => {
    const lowerError = errorStr.toLowerCase()
    
    if (lowerError.includes('invalid_provider')) {
      return 'The authentication provider is not supported.'
    }
    if (lowerError.includes('missing_parameters')) {
      return 'Required authentication parameters are missing.'
    }
    if (lowerError.includes('invalid_state')) {
      return 'Security validation failed. Please try again.'
    }
    if (lowerError.includes('email')) {
      return 'Unable to retrieve your email address. Please ensure your email is verified with the provider.'
    }
    if (lowerError.includes('already linked')) {
      return 'This account is already linked to another user.'
    }
    
    // Default: return the original error message
    return errorStr.replace(/%20/g, ' ')
  }

  const friendlyError = getErrorMessage(error)

  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4">
      <Card className="w-full max-w-md border-destructive">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900">
            <AlertCircle className="h-6 w-6 text-red-600 dark:text-red-400" />
          </div>
          <CardTitle className="text-2xl font-bold text-destructive">Authentication Failed</CardTitle>
          <CardDescription>
            There was a problem authenticating with your OAuth provider.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="rounded-lg bg-muted p-4">
            <p className="text-sm text-muted-foreground">
              <strong>Error:</strong> {friendlyError}
            </p>
          </div>

          <div className="space-y-2">
            <Button asChild className="w-full">
              <Link href="/login">
                Try Again
              </Link>
            </Button>
            <Button asChild variant="outline" className="w-full">
              <Link href="/">
                Go Home
              </Link>
            </Button>
          </div>

          <div className="text-center text-xs text-muted-foreground">
            <p>If this problem persists, please contact support.</p>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

