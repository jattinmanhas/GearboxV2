# Cloudinary Setup Guide

This guide will help you set up Cloudinary for image uploads in the Gearbox frontend.

## 1. Create a Cloudinary Account

1. Go to [Cloudinary](https://cloudinary.com/) and sign up for a free account
2. Once logged in, go to your [Dashboard](https://cloudinary.com/console)
3. Note down your Cloud Name, API Key, and API Secret

## 2. Environment Variables

Create a `.env.local` file in the frontend directory with the following variables:

```env
# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
```

## 3. Cloudinary Configuration

The Cloudinary configuration is handled in `lib/cloudinary.ts`. The default settings include:

- **Folder**: `gearbox-uploads` - All uploads will be organized in this folder
- **Quality**: `auto` - Automatically optimizes quality based on content
- **Format**: `auto` - Automatically serves the best format (WebP, AVIF, etc.)
- **Overwrite**: `true` - Allows overwriting files with the same name

## 4. Features

### Automatic Optimization
- Images are automatically optimized for web delivery
- Multiple formats supported (WebP, AVIF, JPEG, PNG)
- Quality optimization based on content analysis

### Thumbnail Generation
- Automatic thumbnail generation in multiple sizes:
  - `thumbnail`: 150x150px
  - `small`: 300x300px
  - `medium`: 600x600px
  - `large`: 1200x1200px

### Transformations
- Automatic cropping and resizing
- Format conversion
- Quality optimization
- Responsive image delivery

## 5. Usage

The Cloudinary integration is transparent to your existing code. The `EnhancedImageUpload` component will automatically use Cloudinary when the environment variables are set.

### API Routes
- `POST /api/v1/upload/image` - Upload single image
- `DELETE /api/v1/upload/image/[id]` - Delete image
- `GET /api/v1/upload/images` - List user images

### Frontend Components
- `EnhancedImageUpload` - Main upload component
- `ImageUpload` - Basic upload component

## 6. Migration from Local Storage

When you set up Cloudinary environment variables, the system will automatically:
1. Upload new images to Cloudinary
2. Generate optimized URLs
3. Create thumbnails automatically
4. Handle image deletion

Existing local images will continue to work, but new uploads will use Cloudinary.

## 7. Cost Considerations

Cloudinary's free tier includes:
- 25 GB storage
- 25 GB bandwidth per month
- 25,000 transformations per month

For most development and small production sites, this is sufficient.

## 8. Security

- API keys are server-side only
- Images are served over HTTPS
- Public IDs are generated securely
- Folder organization prevents unauthorized access
