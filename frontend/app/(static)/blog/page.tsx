'use client';

import { useState, useEffect } from 'react';
import { useBlogStore } from '@/lib/stores/blog-store';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Search, User, Eye, Clock, Tag } from 'lucide-react';
import Link from 'next/link';
import Image from 'next/image';

export default function BlogPage() {
  const {
    posts,
    categories,
    currentPage,
    totalPages,
    totalPosts,
    filters,
    isLoading,
    error,
    fetchPosts,
    fetchAllCategories,
    setFilters,
    setCurrentPage,
    clearError,
  } = useBlogStore();

  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [sortBy, setSortBy] = useState<string>('createdAt');
  const [sortOrder, setSortOrder] = useState<string>('desc');

  useEffect(() => {
    fetchAllCategories();
    fetchPosts();
  }, [fetchAllCategories, fetchPosts]);

  useEffect(() => {
    const newFilters = {
      ...filters,
      search: searchQuery || undefined,
      categoryId: selectedCategory !== 'all' ? selectedCategory : undefined,
      sortBy: sortBy as 'createdAt' | 'updatedAt' | 'publishedAt' | 'viewCount' | 'title',
      sortOrder: sortOrder as 'asc' | 'desc',
      page: 1,
    };
    setFilters(newFilters);
    fetchPosts(newFilters);
  }, [searchQuery, selectedCategory, sortBy, sortOrder, fetchPosts]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    const newFilters = { ...filters, page };
    setFilters(newFilters);
    fetchPosts(newFilters);
  };

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
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center max-w-md mx-auto px-4">
          <div className="w-24 h-24 mx-auto mb-6 bg-gradient-to-br from-red-100 to-red-200 dark:from-red-900/50 dark:to-red-800/50 rounded-full flex items-center justify-center">
            <div className="text-red-500 text-4xl">⚠️</div>
          </div>
          <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-100 mb-4">Oops! Something went wrong</h1>
          <p className="text-slate-600 dark:text-slate-400 mb-8">{error}</p>
          <Button onClick={clearError} className="bg-blue-600 hover:bg-blue-700 text-white">
            Try Again
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <div className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-600/10 via-purple-600/10 to-pink-600/10 dark:from-blue-400/20 dark:via-purple-400/20 dark:to-pink-400/20"></div>
        <div className="relative container mx-auto px-4 py-16">
          <div className="text-center max-w-4xl mx-auto">
            <h1 className="text-5xl md:text-6xl font-bold bg-gradient-to-r from-slate-900 via-slate-800 to-slate-700 dark:from-slate-100 dark:via-slate-200 dark:to-slate-300 bg-clip-text text-transparent mb-6">
              Our Blog
            </h1>
            <p className="text-xl md:text-2xl text-slate-600 dark:text-slate-300 mb-8 leading-relaxed">
              Discover insights, tutorials, and stories from our community of creators and innovators
            </p>
            <div className="flex flex-wrap justify-center gap-4 text-sm text-slate-500 dark:text-slate-400">
              <span className="flex items-center gap-2">
                <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                {totalPosts} Articles
              </span>
              <span className="flex items-center gap-2">
                <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                {categories.length} Categories
              </span>
              <span className="flex items-center gap-2">
                <div className="w-2 h-2 bg-purple-500 rounded-full"></div>
                Updated Weekly
              </span>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 -mt-8 relative z-10">
        {/* Filters Card */}
        <Card className="shadow-xl border-0 bg-card/80 dark:bg-card/80 backdrop-blur-sm mb-12">
          <CardContent className="p-8">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              {/* Search */}
              <div className="relative group">
                <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-slate-400 group-focus-within:text-blue-500 h-5 w-5 transition-colors" />
                <Input
                  placeholder="Search articles..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-12 h-12 border-slate-200 dark:border-slate-700 focus:border-blue-500 focus:ring-blue-500/20 transition-all duration-200"
                />
              </div>

              {/* Category Filter */}
              <Select value={selectedCategory} onValueChange={setSelectedCategory}>
                <SelectTrigger className="h-12 border-slate-200 dark:border-slate-700 focus:border-blue-500 focus:ring-blue-500/20">
                  <SelectValue placeholder="All Categories" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Categories</SelectItem>
                  {categories.map((category) => (
                    <SelectItem key={category.id} value={category.id}>
                      {category.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>

              {/* Sort By */}
              <Select value={sortBy} onValueChange={setSortBy}>
                <SelectTrigger className="h-12 border-slate-200 dark:border-slate-700 focus:border-blue-500 focus:ring-blue-500/20">
                  <SelectValue placeholder="Sort by" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="createdAt">Date Created</SelectItem>
                  <SelectItem value="updatedAt">Last Updated</SelectItem>
                  <SelectItem value="publishedAt">Published Date</SelectItem>
                  <SelectItem value="viewCount">Views</SelectItem>
                  <SelectItem value="title">Title</SelectItem>
                </SelectContent>
              </Select>

              {/* Sort Order */}
              <Select value={sortOrder} onValueChange={setSortOrder}>
                <SelectTrigger className="h-12 border-slate-200 dark:border-slate-700 focus:border-blue-500 focus:ring-blue-500/20">
                  <SelectValue placeholder="Sort order" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="desc">Newest First</SelectItem>
                  <SelectItem value="asc">Oldest First</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>

        {/* Results Summary */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <h2 className="text-2xl font-bold text-slate-900 dark:text-slate-100">
                Latest Articles
              </h2>
              <Badge variant="secondary" className="px-3 py-1">
                {posts.length} of {totalPosts}
              </Badge>
            </div>
            <div className="hidden md:flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
              <span>Published Articles</span>
            </div>
          </div>
        </div>

        {/* Blog Posts Grid */}
        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {Array.from({ length: 6 }).map((_, i) => (
              <Card key={i} className="overflow-hidden border-0 shadow-lg bg-card/50 dark:bg-card/50 backdrop-blur-sm">
                <Skeleton className="h-48 w-full" />
                <CardHeader className="p-6">
                  <Skeleton className="h-6 w-3/4 mb-2" />
                  <Skeleton className="h-4 w-1/2" />
                </CardHeader>
                <CardContent className="p-6 pt-0">
                  <Skeleton className="h-4 w-full mb-2" />
                  <Skeleton className="h-4 w-2/3" />
                </CardContent>
              </Card>
            ))}
          </div>
        ) : posts.length === 0 ? (
          <div className="text-center py-16">
            <div className="max-w-md mx-auto">
              <div className="w-24 h-24 mx-auto mb-6 bg-gradient-to-br from-slate-100 to-slate-200 dark:from-slate-700 dark:to-slate-600 rounded-full flex items-center justify-center">
                <Search className="h-12 w-12 text-slate-400" />
              </div>
              <h3 className="text-2xl font-bold text-slate-900 dark:text-slate-100 mb-3">No articles found</h3>
              <p className="text-slate-600 dark:text-slate-400 mb-6">Try adjusting your search terms or filters to find what you're looking for.</p>
              <Button 
                onClick={() => {
                  setSearchQuery('');
                  setSelectedCategory('all');
                }}
                variant="outline"
                className="border-slate-300 dark:border-slate-600"
              >
                Clear Filters
              </Button>
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {posts.map((post) => (
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
                      {formatDate(post.createdAt)}
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
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-center mt-16 mb-8">
            <Card className="border-0 shadow-lg bg-card/80 dark:bg-card/80 backdrop-blur-sm">
              <CardContent className="p-4">
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    onClick={() => handlePageChange(currentPage - 1)}
                    disabled={currentPage === 1}
                    className="border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700"
                  >
                    Previous
                  </Button>
                  
                  {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                    <Button
                      key={page}
                      variant={currentPage === page ? "default" : "outline"}
                      onClick={() => handlePageChange(page)}
                      className={`w-10 ${currentPage === page 
                        ? 'bg-blue-600 hover:bg-blue-700 text-white' 
                        : 'border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700'
                      }`}
                    >
                      {page}
                    </Button>
                  ))}
                  
                  <Button
                    variant="outline"
                    onClick={() => handlePageChange(currentPage + 1)}
                    disabled={currentPage === totalPages}
                    className="border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700"
                  >
                    Next
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </div>
  );
}
