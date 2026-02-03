import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { FileText, Eye } from "lucide-react"
import { httpClient } from "@/lib/apiFunctions/http-client"

interface BlogAnalytics {
  total_posts: number
  published_posts: number
  draft_posts: number
  archived_posts: number
  total_views: number
  average_read_time: number
  top_posts: Array<{
    id: string
    title: string
    slug: string
    viewCount: number
    authorName: string
    publishedAt: string | null
  }>
  recent_posts: Array<{
    id: string
    title: string
    slug: string
    viewCount: number
    authorName: string
    publishedAt: string | null
  }>
}

export function BlogInsights() {
  const [analytics, setAnalytics] = useState<BlogAnalytics | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchBlogAnalytics = async () => {
      try {
        setLoading(true)
        const result = await httpClient.get<{
          success: boolean
          message: string
          data: BlogAnalytics
        }>('/dashboard/blog')

        if (result.success) {
          setAnalytics(result.data)
        } else {
          throw new Error(result.message || 'Failed to load blog analytics')
        }
      } catch (error) {
        setError(error instanceof Error ? error.message : 'Unknown error')
      } finally {
        setLoading(false)
      }
    }

    fetchBlogAnalytics()
  }, [])

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Blog Insights</CardTitle>
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
          <CardTitle>Blog Insights</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-destructive">{error}</p>
        </CardContent>
      </Card>
    )
  }

  if (!analytics) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Blog Insights</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">No blog data available.</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0">
        <CardTitle className="flex items-center gap-2">
          <FileText className="h-5 w-5" />
          Blog Insights
        </CardTitle>
        <Badge variant="secondary">Content</Badge>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-2 gap-3 text-sm">
          <div>
            <p className="text-muted-foreground">Total Posts</p>
            <p className="font-semibold">{analytics.total_posts.toLocaleString()}</p>
          </div>
          <div>
            <p className="text-muted-foreground">Published</p>
            <p className="font-semibold">{analytics.published_posts.toLocaleString()}</p>
          </div>
          <div>
            <p className="text-muted-foreground">Drafts</p>
            <p className="font-semibold">{analytics.draft_posts.toLocaleString()}</p>
          </div>
          <div>
            <p className="text-muted-foreground">Views</p>
            <p className="font-semibold flex items-center gap-1">
              <Eye className="h-4 w-4" />
              {analytics.total_views.toLocaleString()}
            </p>
          </div>
        </div>

        <div className="space-y-2">
          <p className="text-sm font-medium">Top Posts</p>
          {analytics.top_posts.length === 0 ? (
            <p className="text-sm text-muted-foreground">No published posts yet.</p>
          ) : (
            <div className="space-y-2">
              {analytics.top_posts.map((post) => (
                <div key={post.id} className="flex items-center justify-between text-sm">
                  <div>
                    <p className="font-medium">{post.title}</p>
                    <p className="text-xs text-muted-foreground">{post.authorName}</p>
                  </div>
                  <span className="text-xs text-muted-foreground">
                    {post.viewCount.toLocaleString()} views
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
