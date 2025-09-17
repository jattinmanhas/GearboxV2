'use client';

import { useState, useEffect } from 'react';
import { useBlogStore } from '@/lib/stores/blog-store';
import { Category } from '@/lib/types/blog';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { 
  Plus, 
  MoreHorizontal, 
  Edit, 
  Trash2, 
  Calendar,
  Tag
} from 'lucide-react';

export default function CategoryManagementPage() {
  const {
    categories,
    isLoading,
    error,
    fetchAllCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    clearError,
  } = useBlogStore();

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    color: '#3B82F6',
  });
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    fetchAllCategories();
  }, [fetchAllCategories]);


  const handleCreateCategory = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createCategory(formData);
      setFormData({ name: '', description: '', color: '#3B82F6' });
      setIsCreateDialogOpen(false);
    } catch (error) {
      console.error('Error creating category:', error);
    }
  };

  const handleEditCategory = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingCategory) return;
    
    try {
      await updateCategory(editingCategory.id, formData);
      setFormData({ name: '', description: '', color: '#3B82F6' });
      setEditingCategory(null);
      setIsEditDialogOpen(false);
    } catch (error) {
      console.error('Error updating category:', error);
    }
  };

  const handleDeleteCategory = async (id: string, name: string) => {
    if (window.confirm(`Are you sure you want to delete "${name}"?`)) {
      try {
        await deleteCategory(id);
      } catch (error) {
        console.error('Error deleting category:', error);
      }
    }
  };

  const openEditDialog = (category: Category) => {
    setEditingCategory(category);
    setFormData({
      name: category.name,
      description: category.description || '',
      color: category.color || '#3B82F6',
    });
    setIsEditDialogOpen(true);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  // Filter categories based on search query
  const filteredCategories = categories.filter(category =>
    category.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    category.description?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    category.slug.toLowerCase().includes(searchQuery.toLowerCase())
  );

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-4">Error Loading Categories</h1>
          <p className="text-gray-600 mb-4">{error}</p>
          <Button onClick={clearError}>Try Again</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold">Category Management</h1>
          <p className="text-gray-600">Manage blog post categories</p>
        </div>
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              Create Category
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create New Category</DialogTitle>
              <DialogDescription>
                Add a new category to organize your blog posts.
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleCreateCategory}>
              <div className="space-y-4 py-4">
                <div>
                  <Label htmlFor="name">Name *</Label>
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                    placeholder="Category name"
                    required
                  />
                </div>
                <div>
                  <Label htmlFor="description">Description</Label>
                  <Textarea
                    id="description"
                    value={formData.description}
                    onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                    placeholder="Category description"
                    rows={3}
                  />
                </div>
                <div>
                  <Label htmlFor="color">Color</Label>
                  <div className="flex items-center gap-2">
                    <Input
                      id="color"
                      type="color"
                      value={formData.color}
                      onChange={(e) => setFormData(prev => ({ ...prev, color: e.target.value }))}
                      className="w-16 h-10 p-1"
                    />
                    <Input
                      value={formData.color}
                      onChange={(e) => setFormData(prev => ({ ...prev, color: e.target.value }))}
                      placeholder="#3B82F6"
                      className="flex-1"
                    />
                  </div>
                </div>
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                  Cancel
                </Button>
                <Button type="submit">Create Category</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {/* Categories Table */}
      <Card>
        <CardHeader>
          <CardTitle>Categories</CardTitle>
          <CardDescription>
            {categories.length} categories available
          </CardDescription>
        </CardHeader>
        <CardContent>
          {/* Search Input */}
          <div className="mb-6">
            <Input
              placeholder="Search categories..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="max-w-sm"
            />
          </div>
          
          {isLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="flex items-center space-x-4">
                  <Skeleton className="h-8 w-8 rounded" />
                  <div className="space-y-2 flex-1">
                    <Skeleton className="h-4 w-1/4" />
                    <Skeleton className="h-4 w-1/2" />
                  </div>
                </div>
              ))}
            </div>
          ) : filteredCategories.length === 0 ? (
            <div className="text-center py-8">
              <h3 className="text-lg font-semibold mb-2">
                {searchQuery ? 'No categories found' : 'No categories found'}
              </h3>
              <p className="text-gray-600 mb-4">
                {searchQuery 
                  ? `No categories match "${searchQuery}". Try a different search term.`
                  : 'Get started by creating your first category'
                }
              </p>
              {!searchQuery && (
                <Button onClick={() => setIsCreateDialogOpen(true)}>
                  Create Category
                </Button>
              )}
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Category</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="w-12"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredCategories.map((category, index) => (
                  <TableRow key={category.id || `category-${index}`}>
                    <TableCell>
                      <div className="flex items-center space-x-3">
                        <div
                          className="w-4 h-4 rounded-full"
                          style={{ backgroundColor: category.color || '#3B82F6' }}
                        />
                        <div>
                          <div className="font-medium">{category.name}</div>
                          <div className="text-sm text-gray-500">
                            Slug: {category.slug}
                          </div>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="max-w-xs">
                        <p className="text-sm text-gray-600 truncate">
                          {category.description || 'No description'}
                        </p>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Calendar className="h-4 w-4 text-gray-400" />
                        {formatDate(category.createdAt)}
                      </div>
                    </TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => openEditDialog(category)}>
                            <Edit className="h-4 w-4 mr-2" />
                            Edit
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onClick={() => handleDeleteCategory(category.id, category.name)}
                            className="text-red-600"
                          >
                            <Trash2 className="h-4 w-4 mr-2" />
                            Delete
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Category</DialogTitle>
            <DialogDescription>
              Update the category information.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleEditCategory}>
            <div className="space-y-4 py-4">
              <div>
                <Label htmlFor="edit-name">Name *</Label>
                <Input
                  id="edit-name"
                  value={formData.name}
                  onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                  placeholder="Category name"
                  required
                />
              </div>
              <div>
                <Label htmlFor="edit-description">Description</Label>
                <Textarea
                  id="edit-description"
                  value={formData.description}
                  onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                  placeholder="Category description"
                  rows={3}
                />
              </div>
              <div>
                <Label htmlFor="edit-color">Color</Label>
                <div className="flex items-center gap-2">
                  <Input
                    id="edit-color"
                    type="color"
                    value={formData.color}
                    onChange={(e) => setFormData(prev => ({ ...prev, color: e.target.value }))}
                    className="w-16 h-10 p-1"
                  />
                  <Input
                    value={formData.color}
                    onChange={(e) => setFormData(prev => ({ ...prev, color: e.target.value }))}
                    placeholder="#3B82F6"
                    className="flex-1"
                  />
                </div>
              </div>
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setIsEditDialogOpen(false)}>
                Cancel
              </Button>
              <Button type="submit">Update Category</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
