-- Drop triggers
DROP TRIGGER IF EXISTS update_blog_posts_updated_at ON blog_posts;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;

-- Drop indexes
DROP INDEX IF EXISTS idx_blog_posts_tags;
DROP INDEX IF EXISTS idx_blog_posts_view_count;
DROP INDEX IF EXISTS idx_blog_posts_created_at;
DROP INDEX IF EXISTS idx_blog_posts_published_at;
DROP INDEX IF EXISTS idx_blog_posts_category_id;
DROP INDEX IF EXISTS idx_blog_posts_author_id;
DROP INDEX IF EXISTS idx_blog_posts_status;
DROP INDEX IF EXISTS idx_blog_posts_slug;
DROP INDEX IF EXISTS idx_categories_slug;

-- Drop tables
DROP TABLE IF EXISTS blog_posts;
DROP TABLE IF EXISTS categories;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();
