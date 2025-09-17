'use client';

import { useEffect } from 'react';
import { useParams } from 'next/navigation';
import { useBlogStore } from '@/lib/stores/blog-store';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { MarkdownRenderer } from '@/components/ui/markdown-renderer';
import { Calendar, User, Eye, Clock, Tag, ArrowLeft, Share2, Bookmark } from 'lucide-react';
import Link from 'next/link';
import Image from 'next/image';

export default function BlogPostPage() {
  const params = useParams();
  const slug = params.slug as string;
  
  const {
    currentPost,
    relatedPosts,
    isLoading,
    error,
    fetchPostBySlug,
    fetchRelatedPosts,
    clearError,
  } = useBlogStore();

  useEffect(() => {
    if (slug) {
      fetchPostBySlug(slug);
    }
  }, [slug, fetchPostBySlug]);

  useEffect(() => {
    if (currentPost) {
      fetchRelatedPosts(currentPost.id);
    }
  }, [currentPost, fetchRelatedPosts]);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'published':
        return 'bg-green-100 text-green-800';
      case 'draft':
        return 'bg-yellow-100 text-yellow-800';
      case 'archived':
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-4">Post Not Found</h1>
          <p className="text-gray-600 mb-4">{error}</p>
          <div className="space-x-4">
            <Button onClick={clearError}>Try Again</Button>
            <Button variant="outline" asChild>
              <Link href="/blog">Back to Blog</Link>
            </Button>
          </div>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <Skeleton className="h-8 w-3/4 mb-4" />
          <Skeleton className="h-4 w-1/2 mb-8" />
          <Skeleton className="h-64 w-full mb-8" />
          <div className="space-y-4">
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-3/4" />
          </div>
        </div>
      </div>
    );
  }

  if (!currentPost) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Post Not Found</h1>
          <p className="text-gray-600 mb-4">The blog post you&apos;re looking for doesn&apos;t exist.</p>
          <Button asChild>
            <Link href="/blog">Back to Blog</Link>
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        {/* Back Button */}
        <Button variant="ghost" asChild className="mb-6">
          <Link href="/blog" className="flex items-center gap-2">
            <ArrowLeft className="h-4 w-4" />
            Back to Blog
          </Link>
        </Button>

        {/* Article Header */}
        <article className="mb-8">
          <header className="mb-8">
            {currentPost.categoryName && (
              <div className="mb-4">
                <Badge variant="outline" className="text-sm px-3 py-1">
                  {currentPost.categoryName}
                </Badge>
              </div>
            )}
            
            <h1 className="text-4xl font-bold mb-4">{currentPost.title}</h1>
            
            <div className="flex items-center gap-6 text-sm text-gray-600 mb-6">
              <div className="flex items-center gap-2">
                <User className="h-4 w-4" />
                {currentPost.authorName}
              </div>
              <div className="flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                {formatDate(currentPost.publishedAt || currentPost.createdAt)}
              </div>
              <div className="flex items-center gap-2">
                <Eye className="h-4 w-4" />
                {currentPost.viewCount} views
              </div>
              <div className="flex items-center gap-2">
                <Clock className="h-4 w-4" />
                {currentPost.readTime} min read
              </div>
            </div>

            <div className="flex items-center gap-4">
              <Button variant="outline" size="sm">
                <Share2 className="h-4 w-4 mr-2" />
                Share
              </Button>
              <Button variant="outline" size="sm">
                <Bookmark className="h-4 w-4 mr-2" />
                Save
              </Button>
            </div>
          </header>

          {/* Featured Image */}
          {currentPost.featuredImage && (
            <div className="relative h-64 md:h-96 w-full mb-8 rounded-lg overflow-hidden">
              <Image
                src={currentPost.featuredImage}
                alt={currentPost.title}
                fill
                className="object-cover"
              />
            </div>
          )}

          {/* Article Content */}
          <MarkdownRenderer content={currentPost.content} />

          {/* Tags */}
          {currentPost.tags && currentPost.tags.length > 0 && (
            <div className="mt-8 pt-6 border-t">
              <h3 className="text-lg font-semibold mb-4">Tags</h3>
              <div className="flex flex-wrap gap-2">
                {currentPost.tags.map((tag) => (
                  <Badge key={tag} variant="secondary" className="flex items-center gap-1">
                    <Tag className="h-3 w-3" />
                    {tag}
                  </Badge>
                ))}
              </div>
            </div>
          )}
        </article>

        {/* Related Posts */}
        {relatedPosts.length > 0 && (
          <section className="mt-16 pt-8 border-t border-slate-200 dark:border-slate-700">
            <div className="text-center mb-12">
              <h2 className="text-3xl font-bold text-slate-900 dark:text-slate-100 mb-4">
                You might also like
              </h2>
              <p className="text-slate-600 dark:text-slate-400">
                Discover more articles that might interest you
              </p>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {relatedPosts.map((post) => (
                <Card key={post.id} className="group overflow-hidden border-0 shadow-lg hover:shadow-2xl bg-card/80 dark:bg-card/80 backdrop-blur-sm transition-all duration-300 hover:-translate-y-1">
                  {post.featuredImage && (
                    <div className="relative h-48 w-full overflow-hidden">
                      <Image
                        src={post.featuredImage}
                        alt={post.title}
                        fill
                        className="object-cover group-hover:scale-105 transition-transform duration-300"
                      />
                      <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent"></div>
                    </div>
                  )}
                  <CardHeader className="p-6">
                    <div className="flex items-center justify-between mb-3">
                      <span className="text-sm text-slate-500 dark:text-slate-400 flex items-center gap-1">
                        <Clock className="h-4 w-4" />
                        {formatDate(post.publishedAt || post.createdAt)}
                      </span>
                      <div className="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
                        <div className="flex items-center gap-1">
                          <Eye className="h-4 w-4" />
                          {post.viewCount}
                        </div>
                      </div>
                    </div>
                    <CardTitle className="line-clamp-2 text-xl group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                      <Link 
                        href={`/blog/${post.slug}`}
                        className="hover:no-underline"
                      >
                        {post.title}
                      </Link>
                    </CardTitle>
                    <CardDescription className="line-clamp-3 text-slate-600 dark:text-slate-300 mt-2">
                      {post.excerpt}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="p-6 pt-0">
                    <div className="flex items-center gap-3 text-sm text-slate-500 dark:text-slate-400 mb-4">
                      <div className="flex items-center gap-1">
                        <User className="h-4 w-4" />
                        {post.authorName}
                      </div>
                      <div className="flex items-center gap-1">
                        <Clock className="h-4 w-4" />
                        {post.readTime} min read
                      </div>
                    </div>
                    
                    {post.tags && post.tags.length > 0 && (
                      <div className="flex flex-wrap gap-2">
                        {post.tags.slice(0, 3).map((tag) => (
                          <Badge key={tag} variant="secondary" className="text-xs bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-600 transition-colors">
                            <Tag className="h-3 w-3 mr-1" />
                            {tag}
                          </Badge>
                        ))}
                        {post.tags.length > 3 && (
                          <Badge variant="secondary" className="text-xs bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-300">
                            +{post.tags.length - 3} more
                          </Badge>
                        )}
                      </div>
                    )}
                  </CardContent>
                </Card>
              ))}
            </div>
          </section>
        )}
      </div>
    </div>
  );
}
