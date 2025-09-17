'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useBlogStore } from '@/lib/stores/blog-store';
import { UpdateBlogPostRequest, BlogPost } from '@/lib/types/blog';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { MarkdownEditor } from '@/components/ui/markdown-editor';
import { MarkdownRenderer } from '@/components/ui/markdown-renderer';
import { 
  Save, 
  Eye, 
  ArrowLeft, 
  Plus, 
  X, 
  Calendar,
  User,
  Tag,
  Globe,
  FileText
} from 'lucide-react';

export default function EditBlogPostPage() {
  const router = useRouter();
  const params = useParams();
  const postId = params.id as string;
  
  const {
    currentPost,
    categories,
    isLoading,
    isUpdating,
    error,
    fetchPostById,
    fetchAllCategories,
    updatePost,
    clearError,
  } = useBlogStore();

  const [formData, setFormData] = useState<UpdateBlogPostRequest>({
    title: '',
    content: '',
    excerpt: '',
    status: 'draft',
    featuredImage: '',
    tags: [],
    categoryId: '',
  });

  const [originalPost, setOriginalPost] = useState<BlogPost | null>(null);
  const [tagInput, setTagInput] = useState('');
  const [showPreview, setShowPreview] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [hasChanges, setHasChanges] = useState(false);

  useEffect(() => {
    if (postId) {
      fetchPostById(postId);
      fetchAllCategories();
    }
  }, [postId, fetchPostById, fetchAllCategories]);

  // Populate form when post data is loaded
  useEffect(() => {
    if (currentPost) {
      setOriginalPost(currentPost);
      setFormData({
        title: currentPost.title || '',
        content: currentPost.content || '',
        excerpt: currentPost.excerpt || '',
        status: currentPost.status || 'draft',
        featuredImage: currentPost.featuredImage || '',
        tags: currentPost.tags || [],
        categoryId: currentPost.categoryId || '',
      });
    }
  }, [currentPost]);

  // Track changes
  useEffect(() => {
    if (originalPost) {
      const hasFormChanges = 
        formData.title !== (originalPost.title || '') ||
        formData.content !== (originalPost.content || '') ||
        formData.excerpt !== (originalPost.excerpt || '') ||
        formData.status !== (originalPost.status || 'draft') ||
        formData.featuredImage !== (originalPost.featuredImage || '') ||
        formData.categoryId !== (originalPost.categoryId || '') ||
        JSON.stringify(formData.tags || []) !== JSON.stringify(originalPost.tags || []);
      
      setHasChanges(hasFormChanges);
    }
  }, [formData, originalPost]);

  const handleInputChange = (field: keyof UpdateBlogPostRequest, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleAddTag = () => {
    const trimmedTag = tagInput.trim();
    if (trimmedTag && !formData.tags?.includes(trimmedTag)) {
      setFormData(prev => ({
        ...prev,
        tags: [...(prev.tags || []), trimmedTag]
      }));
      setTagInput('');
    }
  };

  const handleRemoveTag = (tagToRemove: string) => {
    setFormData(prev => ({
      ...prev,
      tags: prev.tags?.filter(tag => tag !== tagToRemove) || []
    }));
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddTag();
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.title || !formData.content) {
      alert('Please fill in all required fields (Title, Content)');
      return;
    }

    try {
      await updatePost(postId, formData);
      router.push('/dashboard/blog');
    } catch (error) {
      console.error('Error updating post:', error);
      if (error && typeof error === 'object' && 'message' in error) {
        alert(`Failed to update post: ${error.message}`);
      } else {
        alert('Failed to update post. Please try again.');
      }
    }
  };

  const handlePublish = async () => {
    if (!formData.title || !formData.content) {
      alert('Please fill in all required fields before publishing');
      return;
    }

    try {
      await updatePost(postId, { ...formData, status: 'published' });
      router.push('/dashboard/blog');
    } catch (error) {
      console.error('Error publishing post:', error);
      if (error && typeof error === 'object' && 'message' in error) {
        alert(`Failed to publish post: ${error.message}`);
      } else {
        alert('Failed to publish post. Please try again.');
      }
    }
  };

  const generateExcerpt = () => {
    const content = formData.content || '';
    const plainText = content.replace(/<[^>]*>/g, '');
    const excerpt = plainText.length > 160 
      ? plainText.substring(0, 160).replace(/\s+\S*$/, '') + '...'
      : plainText;
    setFormData(prev => ({ ...prev, excerpt }));
  };

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="animate-pulse space-y-6">
            <div className="h-8 bg-gray-200 rounded w-1/3"></div>
            <div className="h-4 bg-gray-200 rounded w-1/2"></div>
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-2 space-y-6">
                <div className="h-64 bg-gray-200 rounded"></div>
                <div className="h-32 bg-gray-200 rounded"></div>
              </div>
              <div className="space-y-6">
                <div className="h-32 bg-gray-200 rounded"></div>
                <div className="h-24 bg-gray-200 rounded"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Alert className="mb-4">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
        <Button onClick={clearError}>Try Again</Button>
      </div>
    );
  }

  if (!currentPost) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Alert className="mb-4">
          <AlertDescription>Post not found</AlertDescription>
        </Alert>
        <Button onClick={() => router.push('/dashboard/blog')}>
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Blog
        </Button>
      </div>
    );
  }

  return (
    <div className={`${isFullscreen ? 'fixed inset-0 z-50 bg-white dark:bg-neutral-900 p-4' : 'container mx-auto px-4 py-8'}`}>
      <div className={`${isFullscreen ? 'h-full overflow-auto' : 'max-w-6xl mx-auto'}`}>
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold mb-2">Edit Blog Post</h1>
              <p className="text-gray-600">Update your blog post content and settings</p>
              {originalPost && (
                <div className="flex items-center gap-4 mt-2 text-sm text-gray-500">
                  <span className="flex items-center gap-1">
                    <Calendar className="h-4 w-4" />
                    Created: {new Date(originalPost.createdAt).toLocaleDateString()}
                  </span>
                  <span className="flex items-center gap-1">
                    <User className="h-4 w-4" />
                    {originalPost.authorName}
                  </span>
                  <Badge variant={originalPost.status === 'published' ? 'default' : 'secondary'}>
                    {originalPost.status === 'published' ? 'Published' : 'Draft'}
                  </Badge>
                </div>
              )}
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                onClick={() => setShowPreview(!showPreview)}
              >
                <Eye className="h-4 w-4 mr-2" />
                {showPreview ? 'Hide Preview' : 'Show Preview'}
              </Button>
              <Button
                variant="outline"
                onClick={() => setIsFullscreen(!isFullscreen)}
              >
                {isFullscreen ? 'Exit Fullscreen' : 'Fullscreen'}
              </Button>
              {isFullscreen && (
                <Button
                  variant="outline"
                  onClick={() => router.push('/dashboard/blog')}
                >
                  <ArrowLeft className="h-4 w-4 mr-2" />
                  Back to Blog
                </Button>
              )}
            </div>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className={`grid grid-cols-1 gap-6 ${showPreview ? 'lg:grid-cols-3' : 'lg:grid-cols-3'}`}>
            {/* Main Content */}
            <div className={`${showPreview ? 'lg:col-span-2' : 'lg:col-span-2'} space-y-6`}>
              {/* Post Details Card */}
              <Card>
                <CardHeader>
                  <CardTitle>Post Details</CardTitle>
                  <CardDescription>Basic information about your blog post</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div>
                    <Label htmlFor="title">Title *</Label>
                    <Input
                      id="title"
                      value={formData.title}
                      onChange={(e) => handleInputChange('title', e.target.value)}
                      placeholder="Enter post title"
                      required
                    />
                  </div>

                  <div>
                    <Label htmlFor="excerpt">Excerpt</Label>
                    <Textarea
                      id="excerpt"
                      value={formData.excerpt}
                      onChange={(e) => handleInputChange('excerpt', e.target.value)}
                      placeholder="Brief description of the post"
                      rows={3}
                    />
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={generateExcerpt}
                      className="mt-2"
                    >
                      Generate from content
                    </Button>
                  </div>

                  <div>
                    <Label htmlFor="featuredImage">Featured Image URL</Label>
                    <Input
                      id="featuredImage"
                      value={formData.featuredImage}
                      onChange={(e) => handleInputChange('featuredImage', e.target.value)}
                      placeholder="https://example.com/image.jpg"
                    />
                    {formData.featuredImage && (
                      <div className="relative h-48 w-full mt-2 rounded-lg overflow-hidden">
                        <img
                          src={formData.featuredImage}
                          alt="Featured image preview"
                          className="object-cover w-full h-full"
                        />
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>

              {/* Content Editor Card */}
              <Card>
                <CardHeader>
                  <CardTitle>Content *</CardTitle>
                  <CardDescription>Write your blog post content using the markdown editor</CardDescription>
                </CardHeader>
                <CardContent>
                  <MarkdownEditor
                    value={formData.content || ''}
                    onChange={(value) => handleInputChange('content', value)}
                    height={500}
                  />
                </CardContent>
              </Card>

              {/* Tags Card */}
              <Card>
                <CardHeader>
                  <CardTitle>Tags</CardTitle>
                  <CardDescription>Add tags to help categorize your post</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex gap-2 mb-4">
                    <Input
                      value={tagInput}
                      onChange={(e) => setTagInput(e.target.value)}
                      onKeyPress={handleKeyPress}
                      placeholder="Enter a tag"
                      className="flex-1"
                    />
                    <Button
                      type="button"
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        handleAddTag();
                      }}
                      disabled={!tagInput.trim()}
                      variant="outline"
                    >
                      <Plus className="h-4 w-4" />
                    </Button>
                  </div>

                  {formData.tags && formData.tags.length > 0 ? (
                    <div className="flex flex-wrap gap-2">
                      {formData.tags.map((tag) => (
                        <Badge key={tag} variant="secondary" className="flex items-center gap-1">
                          {tag}
                          <X
                            className="h-3 w-3 cursor-pointer hover:text-red-500"
                            onClick={() => handleRemoveTag(tag)}
                          />
                        </Badge>
                      ))}
                    </div>
                  ) : (
                    <p className="text-sm text-gray-500 italic">No tags added yet. Type a tag and click the + button or press Enter.</p>
                  )}
                </CardContent>
              </Card>
            </div>

            {/* Sidebar - Only visible when preview is off */}
            {!showPreview && (
              <div className="space-y-6">
                {/* Publish Settings */}
                <Card>
                  <CardHeader>
                    <CardTitle>Publish Settings</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <Label htmlFor="status">Status</Label>
                      <Select
                        value={formData.status}
                        onValueChange={(value) => handleInputChange('status', value)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="draft">Draft</SelectItem>
                          <SelectItem value="published">Published</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>

                    <div>
                      <Label htmlFor="category">Category</Label>
                      <Select
                        value={formData.categoryId}
                        onValueChange={(value) => handleInputChange('categoryId', value)}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Select category" />
                        </SelectTrigger>
                        <SelectContent>
                          {categories.map((category) => (
                            <SelectItem key={category.id} value={category.id}>
                              {category.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  </CardContent>
                </Card>

                {/* Author Information - Read Only */}
                <Card>
                  <CardHeader>
                    <CardTitle>Author Information</CardTitle>
                    <CardDescription>Author details cannot be changed after post creation</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <Label>Author Name</Label>
                      <div className="p-2 bg-gray-50 dark:bg-gray-800 rounded-md text-sm">
                        {originalPost?.authorName || 'Unknown'}
                      </div>
                    </div>

                    <div>
                      <Label>Author Email</Label>
                      <div className="p-2 bg-gray-50 dark:bg-gray-800 rounded-md text-sm">
                        {originalPost?.authorEmail || 'Unknown'}
                      </div>
                    </div>
                  </CardContent>
                </Card>

                {/* Actions */}
                <Card>
                  <CardContent className="pt-6">
                    <div className="space-y-2">
                      <Button
                        type="submit"
                        className="w-full"
                        disabled={isUpdating || !hasChanges}
                      >
                        <Save className="h-4 w-4 mr-2" />
                        {isUpdating ? 'Updating...' : 'Update Post'}
                      </Button>
                      
                      {originalPost?.status === 'draft' && (
                        <Button
                          type="button"
                          onClick={handlePublish}
                          className="w-full"
                          disabled={isUpdating}
                          variant="default"
                        >
                          <Globe className="h-4 w-4 mr-2" />
                          {isUpdating ? 'Publishing...' : 'Publish Now'}
                        </Button>
                      )}
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Preview Panel */}
            {showPreview && (
              <div className="lg:col-span-1 space-y-6">
                {/* Settings in Preview Mode */}
                <div className="space-y-4">
                  <Card>
                    <CardHeader>
                      <CardTitle>Publish Settings</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div>
                        <Label htmlFor="status-preview">Status</Label>
                        <Select
                          value={formData.status}
                          onValueChange={(value) => handleInputChange('status', value)}
                        >
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="draft">Draft</SelectItem>
                            <SelectItem value="published">Published</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>

                      <div>
                        <Label htmlFor="category-preview">Category</Label>
                        <Select
                          value={formData.categoryId}
                          onValueChange={(value) => handleInputChange('categoryId', value)}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Select category" />
                          </SelectTrigger>
                          <SelectContent>
                            {categories.map((category) => (
                              <SelectItem key={category.id} value={category.id}>
                                {category.name}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                    </CardContent>
                  </Card>

                  {/* Author Information in Preview Mode - Read Only */}
                  <Card>
                    <CardHeader>
                      <CardTitle>Author Information</CardTitle>
                      <CardDescription>Author details cannot be changed after post creation</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div>
                        <Label>Author Name</Label>
                        <div className="p-2 bg-gray-50 dark:bg-gray-800 rounded-md text-sm">
                          {originalPost?.authorName || 'Unknown'}
                        </div>
                      </div>

                      <div>
                        <Label>Author Email</Label>
                        <div className="p-2 bg-gray-50 dark:bg-gray-800 rounded-md text-sm">
                          {originalPost?.authorEmail || 'Unknown'}
                        </div>
                      </div>
                    </CardContent>
                  </Card>

                  <Card>
                    <CardContent className="pt-6">
                      <div className="space-y-2">
                        <Button
                          type="submit"
                          className="w-full"
                          disabled={isUpdating || !hasChanges}
                        >
                          <Save className="h-4 w-4 mr-2" />
                          {isUpdating ? 'Updating...' : 'Update Post'}
                        </Button>
                        
                        {originalPost?.status === 'draft' && (
                          <Button
                            type="button"
                            onClick={handlePublish}
                            className="w-full"
                            disabled={isUpdating}
                            variant="default"
                          >
                            <Globe className="h-4 w-4 mr-2" />
                            {isUpdating ? 'Publishing...' : 'Publish Now'}
                          </Button>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                </div>

                {/* Preview Content */}
                <Card>
                  <CardHeader>
                    <CardTitle>Preview</CardTitle>
                    <CardDescription>This is how your blog post will appear to readers</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="prose max-w-none">
                      <article>
                        <header className="mb-8">
                          {formData.categoryId && (
                            <div className="mb-4">
                              <Badge variant="outline" className="text-sm px-3 py-1">
                                {categories.find(cat => cat.id === formData.categoryId)?.name || 'Category'}
                              </Badge>
                            </div>
                          )}

                          <h1 className="text-4xl font-bold mb-4">{formData.title || 'Untitled'}</h1>

                          {formData.excerpt && (
                            <p className="text-lg text-gray-600 dark:text-gray-400 mb-6">{formData.excerpt}</p>
                          )}

                          {formData.featuredImage && (
                            <div className="relative h-64 md:h-96 w-full mb-8 rounded-lg overflow-hidden">
                              <img
                                src={formData.featuredImage}
                                alt={formData.title || 'Featured image'}
                                className="object-cover w-full h-full"
                              />
                            </div>
                          )}
                        </header>

                        <MarkdownRenderer content={formData.content || 'No content yet...'} />

                        {formData.tags && formData.tags.length > 0 && (
                          <div className="mt-8 pt-6 border-t">
                            <h3 className="text-lg font-semibold mb-4">Tags</h3>
                            <div className="flex flex-wrap gap-2">
                              {formData.tags.map((tag) => (
                                <Badge key={tag} variant="secondary" className="flex items-center gap-1">
                                  {tag}
                                </Badge>
                              ))}
                            </div>
                          </div>
                        )}
                      </article>
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}
          </div>
        </form>
      </div>
    </div>
  );
}
