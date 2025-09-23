'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useBlogStore } from '@/lib/stores/blog-store';
import { CreateBlogPostRequest } from '@/lib/types/blog';
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
import { EnhancedImageUpload } from '@/components/ui/enhanced-image-upload';
import { UploadedImage } from '@/lib/image-upload';
import { X, Plus, Save, Eye, Settings, ArrowLeft, Image as ImageIcon, Tag, User, Calendar } from 'lucide-react';

export default function CreateBlogPostPage() {
  const router = useRouter();
  const {
    categories,
    isCreating,
    error,
    fetchAllCategories,
    createPost,
    clearError,
  } = useBlogStore();

  const [formData, setFormData] = useState<CreateBlogPostRequest>({
    title: '',
    content: '',
    excerpt: '',
    authorId: '',
    authorName: '',
    authorEmail: '',
    status: 'draft',
    featuredImage: '',
    tags: [],
    categoryId: '',
  });

  const [tagInput, setTagInput] = useState('');
  const [showPreview, setShowPreview] = useState(false);
  const [showSettings, setShowSettings] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [selectedImages, setSelectedImages] = useState<UploadedImage[]>([]);

  useEffect(() => {
    fetchAllCategories();
  }, [fetchAllCategories]);

  const handleInputChange = (field: keyof CreateBlogPostRequest, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleAddTag = () => {
    const trimmedTag = tagInput.trim();
    console.log('Adding tag:', trimmedTag, 'Current tags:', formData.tags);
    if (trimmedTag && !formData.tags?.includes(trimmedTag)) {
      setFormData(prev => ({
        ...prev,
        tags: [...(prev.tags || []), trimmedTag]
      }));
      setTagInput('');
    } else if (trimmedTag && formData.tags?.includes(trimmedTag)) {
      console.log('Tag already exists:', trimmedTag);
    } else {
      console.log('Empty tag input');
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

  const handleImageSelect = (image: UploadedImage) => {
    setSelectedImages(prev => [...prev, image])
    // Set the first image as featured image
    if (selectedImages.length === 0) {
      setFormData(prev => ({ ...prev, featuredImage: image.url }))
    }
  }

  const handleImageRemove = (imageId: string) => {
    setSelectedImages(prev => {
      const newImages = prev.filter(img => img.id !== imageId)
      // Update featured image if the removed image was featured
      if (formData.featuredImage && prev.find(img => img.id === imageId)?.url === formData.featuredImage) {
        setFormData(prevForm => ({ 
          ...prevForm, 
          featuredImage: newImages.length > 0 ? newImages[0].url : '' 
        }))
      }
      return newImages
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    console.log('Form submitted with data:', formData);
    
    if (!formData.title || !formData.content || !formData.authorId || !formData.authorName || !formData.authorEmail) {
      console.log('Missing required fields:', {
        title: !!formData.title,
        content: !!formData.content,
        authorId: !!formData.authorId,
        authorName: !!formData.authorName,
        authorEmail: !!formData.authorEmail
      });
      alert('Please fill in all required fields (Title, Content, Author ID, Author Name, Author Email)');
      return;
    }

    try {
      console.log('Creating post...');
      const result = await createPost(formData);
      console.log('Post created successfully:', result);
      
      // Only redirect on successful creation
      router.push('/dashboard/blog');
    } catch (error) {
      console.error('Error creating post:', error);
      
      // Don't redirect on error - stay on the form
      if (error && typeof error === 'object' && 'message' in error) {
        alert(`Failed to create post: ${error.message}`);
      } else {
        alert('Failed to create post. Please check your connection and try again.');
      }
    }
  };

  const generateExcerpt = () => {
    const plainText = formData.content.replace(/<[^>]*>/g, '');
    const excerpt = plainText.length > 160 
      ? plainText.substring(0, 160).replace(/\s+\S*$/, '') + '...'
      : plainText;
    setFormData(prev => ({ ...prev, excerpt }));
  };

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

  return (
    <div className={`${isFullscreen ? 'fixed inset-0 z-50 bg-white dark:bg-neutral-900 p-4' : 'container mx-auto px-4 py-8'}`}>
      <div className={`${isFullscreen ? 'h-full overflow-auto' : 'max-w-6xl mx-auto'}`}>
        {/* Header with modern styling but consistent with app design */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold mb-2">Create New Blog Post</h1>
              <p className="text-gray-600">Write and publish a new blog post</p>
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
            {/* Main Content - More space when preview is off */}
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
                    <Label>Featured Image & Media</Label>
                    <EnhancedImageUpload
                      onImageSelect={handleImageSelect}
                      onImageRemove={handleImageRemove}
                      selectedImages={selectedImages}
                      multiple={true}
                      maxImages={5}
                      showPreview={true}
                      showThumbnails={true}
                      config={{
                        maxFileSize: 5 * 1024 * 1024, // 5MB for blog images
                        allowedTypes: ['image/jpeg', 'image/png', 'image/webp', 'image/gif'],
                        maxWidth: 1920,
                        maxHeight: 1080,
                        quality: 90,
                        generateThumbnails: true,
                        thumbnailSizes: [
                          { width: 150, height: 150, name: 'thumbnail' },
                          { width: 300, height: 300, name: 'small' },
                          { width: 600, height: 600, name: 'medium' },
                          { width: 1200, height: 1200, name: 'large' }
                        ]
                      }}
                    />
                    {formData.featuredImage && (
                      <div className="mt-2 text-sm text-muted-foreground">
                        Featured image: {formData.featuredImage}
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
                    value={formData.content}
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

              {/* Author Information */}
              <Card>
                <CardHeader>
                  <CardTitle>Author Information</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div>
                    <Label htmlFor="authorId">Author ID *</Label>
                    <Input
                      id="authorId"
                      value={formData.authorId}
                      onChange={(e) => handleInputChange('authorId', e.target.value)}
                      placeholder="Author ID"
                      required
                    />
                  </div>

                  <div>
                    <Label htmlFor="authorName">Author Name *</Label>
                    <Input
                      id="authorName"
                      value={formData.authorName}
                      onChange={(e) => handleInputChange('authorName', e.target.value)}
                      placeholder="Author Name"
                      required
                    />
                  </div>

                  <div>
                    <Label htmlFor="authorEmail">Author Email *</Label>
                    <Input
                      id="authorEmail"
                      type="email"
                      value={formData.authorEmail}
                      onChange={(e) => handleInputChange('authorEmail', e.target.value)}
                      placeholder="author@example.com"
                      required
                    />
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
                      disabled={isCreating}
                    >
                      <Save className="h-4 w-4 mr-2" />
                      {isCreating ? 'Creating...' : 'Create Post'}
                    </Button>
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

                  {/* Author Information in Preview Mode */}
                  <Card>
                    <CardHeader>
                      <CardTitle>Author Information</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div>
                        <Label htmlFor="authorId-preview">Author ID *</Label>
                        <Input
                          id="authorId-preview"
                          value={formData.authorId}
                          onChange={(e) => handleInputChange('authorId', e.target.value)}
                          placeholder="Author ID"
                          required
                        />
                      </div>

                      <div>
                        <Label htmlFor="authorName-preview">Author Name *</Label>
                        <Input
                          id="authorName-preview"
                          value={formData.authorName}
                          onChange={(e) => handleInputChange('authorName', e.target.value)}
                          placeholder="Author Name"
                          required
                        />
                      </div>

                      <div>
                        <Label htmlFor="authorEmail-preview">Author Email *</Label>
                        <Input
                          id="authorEmail-preview"
                          type="email"
                          value={formData.authorEmail}
                          onChange={(e) => handleInputChange('authorEmail', e.target.value)}
                          placeholder="author@example.com"
                          required
                        />
                      </div>
                    </CardContent>
                  </Card>

                  <Card>
                    <CardContent className="pt-6">
                      <Button 
                        type="submit" 
                        className="w-full" 
                        disabled={isCreating}
                      >
                        <Save className="h-4 w-4 mr-2" />
                        {isCreating ? 'Creating...' : 'Create Post'}
                      </Button>
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
