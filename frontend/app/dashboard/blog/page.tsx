'use client';

import { useState, useEffect } from 'react';
import { useBlogStore } from '@/lib/stores/blog-store';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table';
import { 
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { 
  Search, 
  Plus, 
  MoreHorizontal, 
  Edit, 
  Trash2, 
  Eye, 
  Calendar,
  User,
  Tag
} from 'lucide-react';
import Link from 'next/link';
import Image from 'next/image';

export default function BlogManagementPage() {
  const {
    posts,
    currentPage,
    totalPages,
    totalPosts,
    filters,
    isLoading,
    error,
    fetchAllPosts,
    deletePost,
    setFilters,
    setCurrentPage,
    clearError,
  } = useBlogStore();

  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');

  useEffect(() => {
    fetchAllPosts();
  }, [fetchAllPosts]);

  useEffect(() => {
    const newFilters = {
      ...filters,
      search: searchQuery || undefined,
      status: statusFilter !== 'all' ? statusFilter as 'draft' | 'published' | 'archived' : undefined,
      page: 1,
    };
    setFilters(newFilters);
    fetchAllPosts(newFilters);
  }, [searchQuery, statusFilter, fetchAllPosts]);

  const handleDeletePost = async (id: string, title: string) => {
    if (window.confirm(`Are you sure you want to delete "${title}"?`)) {
      try {
        await deletePost(id);
      } catch (error) {
        console.error('Error deleting post:', error);
      }
    }
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    const newFilters = { ...filters, page };
    setFilters(newFilters);
    fetchAllPosts(newFilters);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
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
      <div className="w-full px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-4">Error Loading Posts</h1>
          <p className="text-gray-600 mb-4">{error}</p>
          <Button onClick={clearError}>Try Again</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold">Blog Management</h1>
          <p className="text-gray-600">Manage your blog posts and content</p>
        </div>
        <Button asChild>
          <Link href="/dashboard/blog/create">
            <Plus className="h-4 w-4 mr-2" />
            Create Post
          </Link>
        </Button>
      </div>

      {/* Filters */}
      <Card className="mb-6">
        <CardContent className="pt-6">
          <div className="flex gap-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
              <Input
                placeholder="Search posts..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Status</option>
              <option value="published">Published</option>
              <option value="draft">Draft</option>
              <option value="archived">Archived</option>
            </select>
          </div>
        </CardContent>
      </Card>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <Card>
          <CardContent className="pt-6">
            <div className="text-2xl font-bold">{totalPosts}</div>
            <p className="text-sm text-gray-600">Total Posts</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="text-2xl font-bold">
              {posts.filter(post => post.status === 'published').length}
            </div>
            <p className="text-sm text-gray-600">Published</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="text-2xl font-bold">
              {posts.filter(post => post.status === 'draft').length}
            </div>
            <p className="text-sm text-gray-600">Drafts</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="text-2xl font-bold">
              {posts.reduce((sum, post) => sum + (post.viewCount || 0), 0)}
            </div>
            <p className="text-sm text-gray-600">Total Views</p>
          </CardContent>
        </Card>
      </div>

      {/* Posts Table */}
      <Card>
        <CardHeader>
          <CardTitle>Blog Posts</CardTitle>
          <CardDescription>
            Showing {posts.length} of {totalPosts} posts
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="flex items-center space-x-4">
                  <Skeleton className="h-12 w-12 rounded" />
                  <div className="space-y-2 flex-1">
                    <Skeleton className="h-4 w-3/4" />
                    <Skeleton className="h-4 w-1/2" />
                  </div>
                </div>
              ))}
            </div>
          ) : posts.length === 0 ? (
            <div className="text-center py-8">
              <h3 className="text-lg font-semibold mb-2">No posts found</h3>
              <p className="text-gray-600 mb-4">Get started by creating your first blog post</p>
              <Button asChild>
                <Link href="/dashboard/blog/create">Create Post</Link>
              </Button>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <div className="min-w-[500px]">
                {/* Header */}
                <div className="grid grid-cols-12 gap-2 py-2 px-3 text-sm font-medium text-gray-500 border-b">
                  <div className="col-span-4">Post</div>
                  <div className="col-span-2">Status</div>
                  <div className="col-span-2 hidden sm:block">Author</div>
                  <div className="col-span-1 hidden md:block">Views</div>
                  <div className="col-span-2 hidden lg:block">Created</div>
                  <div className="col-span-1">Actions</div>
                </div>
                
                {/* Rows */}
                {posts.map((post) => (
                  <div key={post.id} className="grid grid-cols-12 gap-2 py-3 px-3 text-sm border-b hover:bg-muted/50">
                    {/* Post Column */}
                    <div className="col-span-4 min-w-0">
                      <div className="flex items-center space-x-2">
                        {post.featuredImage && (
                          <div className="relative h-6 w-6 rounded overflow-hidden flex-shrink-0">
                            <Image
                              src={post.featuredImage}
                              alt={post.title}
                              fill
                              className="object-cover"
                            />
                          </div>
                        )}
                        <div className="min-w-0 flex-1">
                          <h3 className="font-medium text-xs truncate leading-tight">{post.title}</h3>
                          <p className="text-xs text-gray-500 truncate mt-0.5">
                            {post.excerpt || 'No excerpt'}
                          </p>
                          {post.tags && post.tags.length > 0 && (
                            <div className="flex items-center gap-1 mt-0.5">
                              <Tag className="h-2 w-2 text-gray-400" />
                              <span className="text-xs text-gray-500">
                                {post.tags.slice(0, 1).join(', ')}
                                {post.tags.length > 1 && ` +${post.tags.length - 1}`}
                              </span>
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                    
                    {/* Status Column */}
                    <div className="col-span-2 flex items-center">
                      <Badge className={`${getStatusColor(post.status)} text-xs px-1.5 py-0.5`}>
                        {post.status}
                      </Badge>
                    </div>
                    
                    {/* Author Column */}
                    <div className="col-span-2 hidden sm:flex items-center">
                      <div className="flex items-center gap-1">
                        <User className="h-3 w-3 text-gray-400" />
                        <span className="text-xs truncate">{post.authorName}</span>
                      </div>
                    </div>
                    
                    {/* Views Column */}
                    <div className="col-span-1 hidden md:flex items-center">
                      <div className="flex items-center gap-1">
                        <Eye className="h-3 w-3 text-gray-400" />
                        <span className="text-xs">{post.viewCount || 0}</span>
                      </div>
                    </div>
                    
                    {/* Created Column */}
                    <div className="col-span-2 hidden lg:flex items-center">
                      <div className="flex items-center gap-1">
                        <Calendar className="h-3 w-3 text-gray-400" />
                        <span className="text-xs">{formatDate(post.createdAt)}</span>
                      </div>
                    </div>
                    
                    {/* Actions Column */}
                    <div className="col-span-1 flex items-center justify-end">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem asChild>
                            <Link href={`/blog/${post.slug}`}>
                              <Eye className="h-4 w-4 mr-2" />
                              View
                            </Link>
                          </DropdownMenuItem>
                          <DropdownMenuItem asChild>
                            <Link href={`/dashboard/blog/edit/${post.id}`}>
                              <Edit className="h-4 w-4 mr-2" />
                              Edit
                            </Link>
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onClick={() => handleDeletePost(post.id, post.title)}
                            className="text-red-600"
                          >
                            <Trash2 className="h-4 w-4 mr-2" />
                            Delete
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex justify-center mt-8">
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() => handlePageChange(currentPage - 1)}
              disabled={currentPage === 1}
            >
              Previous
            </Button>
            
            {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
              <Button
                key={page}
                variant={currentPage === page ? "default" : "outline"}
                onClick={() => handlePageChange(page)}
                className="w-10"
              >
                {page}
              </Button>
            ))}
            
            <Button
              variant="outline"
              onClick={() => handlePageChange(currentPage + 1)}
              disabled={currentPage === totalPages}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
