# Image Cleanup System

## Overview

The image cleanup system automatically removes images from Cloudinary when entities (products, categories, blog posts) are deleted. This prevents orphaned images from accumulating in your Cloudinary account and helps manage storage costs.

## How It Works

### 1. Automatic Cleanup Process

When you delete an entity (product, category, or blog post), the system:

1. **Fetches the entity data** before deletion to extract image URLs
2. **Deletes the entity** from the backend service
3. **If deletion is successful**, automatically cleans up associated images from Cloudinary
4. **Logs the cleanup results** for monitoring

### 2. Supported Entities

The cleanup system works with:

- **Products**: Main product images, variant images, gallery images
- **Categories**: Category images, logos
- **Blog Posts**: Featured images, content images

### 3. Image Detection

The system automatically detects images from various field names:

- `image_url`, `imageUrl`
- `featured_image`, `featuredImage`
- `thumbnail_url`, `thumbnailUrl`
- `cover_image`, `coverImage`
- `logo_url`, `logoUrl`
- `images` array
- `variants` array (for products)

## Implementation Details

### Core Files

- **`lib/cloudinary-cleanup.ts`**: Core cleanup utilities
- **`app/api/v1/products/[id]/route.ts`**: Product deletion with cleanup
- **`app/api/v1/products/categories/[id]/route.ts`**: Category deletion with cleanup
- **`app/api/v1/blog/posts/[id]/route.ts`**: Blog post deletion with cleanup

### Key Functions

#### `extractCloudinaryPublicId(url: string)`
Extracts the Cloudinary public ID from a URL.

```typescript
// Examples:
extractCloudinaryPublicId('https://res.cloudinary.com/dtah8vzac/image/upload/v123/gearbox-uploads/image.jpg')
// Returns: 'gearbox-uploads/image'

extractCloudinaryPublicId('http://localhost:3000/uploads/image.jpg')
// Returns: null (not a Cloudinary image)
```

#### `cleanupEntityImages(entity: any)`
Cleans up all images associated with an entity.

```typescript
const result = await cleanupEntityImages(productData)
// Returns: { success: boolean, deleted: number, errors: string[] }
```

#### `deleteCloudinaryImages(imageUrls: string[])`
Deletes multiple images from Cloudinary.

```typescript
const result = await deleteCloudinaryImages(['url1', 'url2', 'url3'])
// Returns: { success: boolean, deleted: number, errors: string[] }
```

## Error Handling

### Graceful Degradation

- **Image cleanup failures don't affect entity deletion**
- **Non-Cloudinary images are ignored** (no errors thrown)
- **Partial cleanup is supported** (some images may fail, others succeed)
- **Detailed error logging** for debugging

### Error Types

1. **Invalid URLs**: Non-Cloudinary URLs are silently ignored
2. **Cloudinary API Errors**: Logged but don't fail the operation
3. **Network Issues**: Retry logic could be added in the future

## Monitoring and Logging

### Console Logs

The system provides detailed logging:

```
✅ Success: "Cleaned up 3 images for deleted product 123"
⚠️  Warning: "Some images could not be deleted: [error details]"
❌ Error: "Failed to cleanup images for deleted product: [error details]"
```

### Cloudinary Dashboard

- **Deleted images** are removed from your Cloudinary dashboard
- **Storage usage** decreases automatically
- **Bandwidth costs** are reduced

## Configuration

### Environment Variables

The cleanup system uses the same Cloudinary configuration:

```env
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
```

### Customization

You can customize the cleanup behavior by modifying:

1. **Image field detection** in `extractImageUrls()`
2. **URL pattern matching** in `extractCloudinaryPublicId()`
3. **Error handling** in the delete endpoints

## Testing

### Manual Testing

1. **Create an entity** with Cloudinary images
2. **Delete the entity** through the API
3. **Check Cloudinary dashboard** - images should be removed
4. **Check server logs** - cleanup should be logged

### API Testing

```bash
# Test category deletion with cleanup
curl -X DELETE http://localhost:3000/api/v1/products/categories/123

# Test product deletion with cleanup
curl -X DELETE http://localhost:3000/api/v1/products/456

# Test blog post deletion with cleanup
curl -X DELETE http://localhost:3000/api/v1/blog/posts/789
```

## Benefits

### Cost Savings
- **Reduced storage costs** by removing unused images
- **Lower bandwidth costs** from fewer stored images
- **Automatic cleanup** prevents manual maintenance

### Data Integrity
- **No orphaned images** in Cloudinary
- **Consistent state** between database and image storage
- **Cleaner dashboard** with only active images

### Performance
- **Faster Cloudinary operations** with fewer images
- **Reduced API calls** for image management
- **Better organization** of image assets

## Future Enhancements

### Potential Improvements

1. **Batch cleanup**: Process multiple entities at once
2. **Retry logic**: Retry failed image deletions
3. **Cleanup scheduling**: Periodic cleanup of orphaned images
4. **Analytics**: Track cleanup statistics and savings
5. **Webhook integration**: Real-time cleanup notifications

### Advanced Features

1. **Soft delete support**: Cleanup images on soft-deleted entities
2. **Image usage tracking**: Only delete images not used elsewhere
3. **Archive mode**: Move images to archive folder instead of deletion
4. **Cleanup policies**: Configurable cleanup rules per entity type

## Troubleshooting

### Common Issues

1. **Images not being deleted**: Check Cloudinary configuration and permissions
2. **Partial cleanup**: Some images may fail due to network issues
3. **Performance impact**: Cleanup adds overhead to delete operations

### Debug Steps

1. **Check server logs** for cleanup messages
2. **Verify Cloudinary configuration** is correct
3. **Test with a single image** to isolate issues
4. **Check Cloudinary dashboard** for deletion status

## Security Considerations

### API Security
- **Authentication required** for all delete operations
- **Authorization checks** before cleanup
- **Rate limiting** to prevent abuse

### Data Protection
- **No sensitive data** in image URLs
- **Secure deletion** from Cloudinary
- **Audit logging** for compliance

The image cleanup system ensures your Cloudinary account stays clean and cost-effective while maintaining data integrity across your application.
