# Markdown Editor Guide

## Overview

The blog system now includes a rich markdown editor that provides a user-friendly interface for creating blog posts with markdown support. Users can write content using either the visual editor or directly in markdown format.

## Features

### 🎨 Visual Editor
- **WYSIWYG Interface**: Write content in a visual editor that automatically converts to markdown
- **Live Preview**: Switch between edit and preview modes to see how content will appear
- **Toolbar**: Quick access to formatting options like bold, italic, lists, and quotes

### 📝 Markdown Support
- **Full Markdown Syntax**: Supports all standard markdown features including:
  - Headers (# ## ###)
  - Bold (**text**) and italic (*text*)
  - Lists (ordered and unordered)
  - Links [text](url)
  - Images ![alt](url)
  - Code blocks and inline code
  - Blockquotes
  - Tables
  - Horizontal rules

### 🖼️ Image Management
- **Drag & Drop**: Simply drag images into the editor
- **URL Input**: Add images by pasting URLs
- **Alt Text**: Add accessibility descriptions for images
- **Preview**: See image previews before inserting

### 🔗 Link Management
- **Easy Link Creation**: Add links with custom text and URLs
- **External Links**: Automatically opens in new tabs for security

## How to Use

### Creating a New Blog Post

1. **Navigate** to Dashboard → Blog → Create Post
2. **Fill in** the basic information (title, excerpt, author details)
3. **Write Content** using the markdown editor:
   - Use the toolbar buttons for quick formatting
   - Switch to preview mode to see the final result
   - Add images using the image button or drag & drop
4. **Preview** your post using the "Show Preview" button
5. **Publish** when ready

### Editor Toolbar

| Button | Function | Markdown Output |
|--------|----------|----------------|
| **B** | Bold | `**text**` |
| *I* | Italic | `*text*` |
| • | List | `- item` |
| " | Quote | `> text` |
| 🖼️ | Image | `![alt](url)` |
| 🔗 | Link | `[text](url)` |

### Markdown Examples

#### Headers
```markdown
# Main Title
## Section Title
### Subsection Title
```

#### Text Formatting
```markdown
**Bold text** and *italic text*
~~Strikethrough~~ and `inline code`
```

#### Lists
```markdown
- Unordered list item
- Another item
  - Nested item

1. Ordered list item
2. Another item
```

#### Links and Images
```markdown
[Link text](https://example.com)
![Image alt text](https://example.com/image.jpg)
```

#### Code Blocks
```markdown
```javascript
function hello() {
  console.log("Hello, World!");
}
```
```

#### Tables
```markdown
| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Data 1   | Data 2   | Data 3   |
| Data 4   | Data 5   | Data 6   |
```

## Technical Details

### Components Used
- **MarkdownEditor**: Main editor component with toolbar and tabs
- **MarkdownRenderer**: Displays markdown content with proper styling
- **ImageUpload**: Handles image insertion with drag & drop support

### Dependencies
- `@uiw/react-md-editor`: Core markdown editor
- `react-markdown`: Markdown rendering
- `remark-gfm`: GitHub Flavored Markdown support
- `rehype-highlight`: Syntax highlighting for code blocks

### Styling
- Uses Tailwind CSS for consistent styling
- Dark mode support
- Responsive design
- Custom syntax highlighting themes

## Best Practices

### Content Creation
1. **Use Headers** to structure your content logically
2. **Add Alt Text** to all images for accessibility
3. **Preview Regularly** to ensure formatting looks correct
4. **Use Lists** to break up long paragraphs
5. **Add Code Blocks** for technical content

### Image Guidelines
- **Optimize Images**: Compress images before uploading
- **Use Descriptive Alt Text**: Help screen readers understand images
- **Consistent Sizing**: Use similar aspect ratios for better layout
- **High Quality**: Use high-resolution images for better display

### Performance Tips
- **Large Images**: Consider using external image hosting services
- **Code Blocks**: Use appropriate language tags for syntax highlighting
- **Tables**: Keep tables simple and mobile-friendly

## Troubleshooting

### Common Issues

**Images not displaying?**
- Check that image URLs are accessible
- Ensure URLs are complete (include https://)
- Try using a different image hosting service

**Formatting not working?**
- Make sure you're using proper markdown syntax
- Check for extra spaces or characters
- Use the preview mode to verify formatting

**Editor not loading?**
- Refresh the page
- Check browser console for errors
- Ensure all dependencies are installed

### Getting Help

If you encounter issues:
1. Check the browser console for error messages
2. Verify all required fields are filled
3. Try creating a simple post first
4. Contact the development team for technical support

## Future Enhancements

Planned features for future releases:
- **Cloud Image Storage**: Direct integration with cloud storage services
- **Collaborative Editing**: Real-time collaboration features
- **Version History**: Track changes and revert to previous versions
- **Advanced Tables**: Rich table editing capabilities
- **Math Support**: LaTeX math equation rendering
- **Custom Shortcodes**: Add custom content blocks
