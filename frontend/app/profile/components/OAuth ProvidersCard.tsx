'use client'

import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { oauthApi } from '@/lib/api'
import { OAuthProviderLink, OAuthProvider } from '@/lib/types'
import { Loader2, Link2, Unlink, AlertCircle, CheckCircle2 } from 'lucide-react'

// Provider display names and colors
const providerConfig = {
  google: {
    name: 'Google',
    color: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
  },
  github: {
    name: 'GitHub',
    color: 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200',
  },
}

export function OAuthProvidersCard() {
  const [providers, setProviders] = useState<OAuthProviderLink[]>([])
  const [loading, setLoading] = useState(true)
  const [unlinking, setUnlinking] = useState<OAuthProvider | null>(null)
  const [error, setError] = useState<string>('')
  const [success, setSuccess] = useState<string>('')

  useEffect(() => {
    fetchLinkedProviders()
  }, [])

  const fetchLinkedProviders = async () => {
    try {
      setLoading(true)
      setError('')
      const response = await oauthApi.getLinkedProviders()
      setProviders(response.providers || [])
    } catch (err: any) {
      console.error('Failed to fetch linked providers:', err)
      setError(err.message || 'Failed to load linked providers')
    } finally {
      setLoading(false)
    }
  }

  const handleUnlink = async (provider: OAuthProvider) => {
    if (!confirm(`Are you sure you want to unlink your ${providerConfig[provider].name} account?`)) {
      return
    }

    try {
      setUnlinking(provider)
      setError('')
      setSuccess('')
      await oauthApi.unlinkProvider(provider)
      setSuccess(`${providerConfig[provider].name} account unlinked successfully`)
      // Refresh the list
      await fetchLinkedProviders()
    } catch (err: any) {
      console.error('Failed to unlink provider:', err)
      setError(err.message || 'Failed to unlink provider')
    } finally {
      setUnlinking(null)
    }
  }

  const handleLink = async (provider: OAuthProvider) => {
    try {
      setError('')
      setSuccess('')
      const { auth_url } = await oauthApi.initiateOAuth(provider)
      // For linking, we'll redirect but user should come back to profile
      // Store return URL in session storage
      sessionStorage.setItem('oauth_return_url', '/profile')
      window.location.href = auth_url
    } catch (err: any) {
      console.error('Failed to initiate OAuth:', err)
      setError(err.message || 'Failed to initiate OAuth')
    }
  }

  const isProviderLinked = (provider: OAuthProvider) => {
    return providers.some((p) => p.provider === provider)
  }

  const getProviderLink = (provider: OAuthProvider) => {
    return providers.find((p) => p.provider === provider)
  }

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Connected Accounts</CardTitle>
          <CardDescription>Manage your OAuth provider connections</CardDescription>
        </CardHeader>
        <CardContent className="flex items-center justify-center py-8">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Connected Accounts</CardTitle>
        <CardDescription>
          Link your accounts from other services to enable quick login
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {success && (
          <Alert className="border-green-500 bg-green-50 dark:bg-green-950">
            <CheckCircle2 className="h-4 w-4 text-green-600 dark:text-green-400" />
            <AlertDescription className="text-green-800 dark:text-green-200">
              {success}
            </AlertDescription>
          </Alert>
        )}

        <div className="space-y-3">
          {(['google', 'github'] as OAuthProvider[]).map((provider) => {
            const linked = isProviderLinked(provider)
            const providerLink = getProviderLink(provider)
            const config = providerConfig[provider]

            return (
              <div
                key={provider}
                className="flex items-center justify-between rounded-lg border p-4"
              >
                <div className="flex items-center gap-4">
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="font-medium">{config.name}</span>
                      {linked && (
                        <Badge className={config.color} variant="secondary">
                          <Link2 className="mr-1 h-3 w-3" />
                          Connected
                        </Badge>
                      )}
                    </div>
                    {linked && providerLink && (
                      <p className="text-sm text-muted-foreground mt-1">
                        {providerLink.email}
                      </p>
                    )}
                    {!linked && (
                      <p className="text-sm text-muted-foreground mt-1">
                        Not connected
                      </p>
                    )}
                  </div>
                </div>

                {linked ? (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleUnlink(provider)}
                    disabled={unlinking === provider}
                  >
                    {unlinking === provider ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Unlinking...
                      </>
                    ) : (
                      <>
                        <Unlink className="mr-2 h-4 w-4" />
                        Unlink
                      </>
                    )}
                  </Button>
                ) : (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleLink(provider)}
                  >
                    <Link2 className="mr-2 h-4 w-4" />
                    Connect
                  </Button>
                )}
              </div>
            )
          })}
        </div>

        {providers.length === 0 && (
          <div className="text-center py-6 text-muted-foreground">
            <p className="text-sm">No accounts connected yet.</p>
            <p className="text-xs mt-1">Connect an account to enable quick login.</p>
          </div>
        )}

        <div className="mt-4 rounded-lg bg-muted p-3">
          <p className="text-xs text-muted-foreground">
            <strong>Note:</strong> You must have at least one authentication method (password or OAuth provider) 
            to access your account. Set a password before unlinking your last OAuth provider.
          </p>
        </div>
      </CardContent>
    </Card>
  )
}

