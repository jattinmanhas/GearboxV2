import React from "react";
import {
  ChevronRight,
  Cpu,
  Smartphone,
  Laptop,
  HeadphonesIcon,
  TabletIcon,
  WatchIcon,
  ArrowRight,
  Star,
  TrendingUp,
} from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { headers } from "next/headers";

const MainLandingPage = async () => {
  // Fetch recent products and blogs
  let products: any[] = [];
  let blogs: any[] = [];
  
  try {
    // Get the base URL for API calls
    const headersList = await headers();
    const host = headersList.get('host') || 'localhost:3000';
    const protocol = process.env.NODE_ENV === 'production' ? 'https' : 'http';
    const baseUrl = `${protocol}://${host}`;
    
    // Fetch products - get 3 most recent
    const params = new URLSearchParams({
      limit: '3',
      is_active: 'true',
      sort_by: 'created_at',
      sort_order: 'desc'
    });
    
    const productsResponse = await fetch(`${baseUrl}/api/v1/products?${params}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    if (productsResponse.ok) {
      const data = await productsResponse.json();
      products = data.data?.products || data.data || [];
    }
    
    // Fetch blogs - get 3 most recent
    const blogsResponse = await fetch(`${baseUrl}/api/v1/blog/posts/recent?limit=3`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    if (blogsResponse.ok) {
      const blogsData = await blogsResponse.json();
      blogs = blogsData.data || blogsData || [];
    }
  } catch (error) {
    console.error("Error fetching data for landing page:", error);
    // Use empty arrays if fetch fails - the page will still render
  }

  // Helper function to format price
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
    }).format(price);
  };

  // Helper function to get primary image URL
  const getProductImage = (product: any) => {
    if (product.images && product.images.length > 0) {
      const primaryImage = product.images.find((img: any) => img.is_primary) || product.images[0];
      return primaryImage.url;
    }
    return null;
  };

  // Helper function to format date for blogs
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);
    
    if (diffInSeconds < 60) return 'just now';
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} minutes ago`;
    if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} hours ago`;
    if (diffInSeconds < 604800) return `${Math.floor(diffInSeconds / 86400)} days ago`;
    if (diffInSeconds < 2592000) return `${Math.floor(diffInSeconds / 604800)} weeks ago`;
    if (diffInSeconds < 31536000) return `${Math.floor(diffInSeconds / 2592000)} months ago`;
    return `${Math.floor(diffInSeconds / 31536000)} years ago`;
  };
  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <section className="relative overflow-hidden py-20 lg:py-32">
        <div className="container mx-auto px-4">
          <div className="text-center">
            <Badge variant="secondary" className="mb-6">
              <TrendingUp className="mr-2 h-4 w-4" />
              Trusted by 10,000+ Tech Enthusiasts
            </Badge>
            
            <h1 className="mb-6 text-4xl font-bold tracking-tight text-foreground md:text-6xl lg:text-7xl">
              Your Ultimate
              <span className="bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                {" "}Tech{" "}
              </span>
              Destination
            </h1>
            
            <p className="mx-auto mb-8 max-w-2xl text-lg text-muted-foreground">
              Discover cutting-edge gadgets and stay informed with expert
              reviews, guides, and the latest tech insights.
            </p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12">
              <Button size="lg" asChild>
                <Link href="/shop">
                  Explore Products
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
              <Button size="lg" variant="outline" asChild>
                <Link href="/blogs">
                  Read Reviews
                </Link>
              </Button>
            </div>

            {/* Quick Categories */}
            <div className="flex flex-wrap justify-center gap-3">
              {[
                { icon: Smartphone, label: "Phones" },
                { icon: Laptop, label: "Laptops" },
                { icon: HeadphonesIcon, label: "Audio" },
                { icon: WatchIcon, label: "Wearables" },
                { icon: TabletIcon, label: "Tablets" },
                { icon: Cpu, label: "Components" },
              ].map((category, index) => (
                <Button
                  key={index}
                  variant="secondary"
                  size="sm"
                  className="flex items-center space-x-2"
                >
                  <category.icon className="h-4 w-4" />
                  <span>{category.label}</span>
                </Button>
              ))}
            </div>
          </div>
        </div>
      </section>
      {/* Featured Products Section */}
      <section className="py-16 bg-muted/30">
        <div className="container mx-auto px-4">
          <div className="mb-12 text-center">
            <h2 className="text-3xl font-bold text-foreground mb-4">Trending Products</h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              Discover the most popular tech products that everyone's talking about
            </p>
          </div>

          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 mb-8">
            {products.length > 0 ? (
              products.map((product) => {
                const imageUrl = getProductImage(product);
                const price = product.variants && product.variants.length > 0 
                  ? product.variants[0].price 
                  : product.price || 0;
                
                return (
                  <Card key={product.id} className="group hover:shadow-lg transition-shadow flex flex-col">
                    <CardHeader>
                      <div className="w-full h-48 mb-2 overflow-hidden rounded-lg bg-muted flex items-center justify-center">
                        {imageUrl ? (
                          <img 
                            src={imageUrl} 
                            alt={product.name}
                            className="w-full h-full object-cover"
                          />
                        ) : (
                          <div className="text-4xl">📦</div>
                        )}
                      </div>
                      <CardTitle className="text-xl">{product.name}</CardTitle>
                      <CardDescription className="line-clamp-2">
                        {product.description || 'Check out this amazing product!'}
                      </CardDescription>
                    </CardHeader>
                    <CardContent className="flex-1 flex flex-col">
                      <div className="flex items-center justify-between mb-4">
                        <span className="text-2xl font-bold text-primary">{formatPrice(price)}</span>
                        {product.rating && (
                          <div className="flex items-center gap-1">
                            <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                            <span className="text-sm text-muted-foreground">
                              {product.rating.toFixed(1)}
                            </span>
                          </div>
                        )}
                      </div>
                      <Button className="w-full mt-auto" asChild>
                        <Link href={`/shop/products/${product.id}`}>
                          View Details
                          <ArrowRight className="ml-2 h-4 w-4" />
                        </Link>
                      </Button>
                    </CardContent>
                  </Card>
                );
              })
            ) : (
              // Fallback message if no products
              <div className="col-span-full text-center text-muted-foreground py-12">
                No products available yet. Check back soon!
              </div>
            )}
          </div>

          <div className="text-center">
            <Button variant="outline" size="lg" asChild>
              <Link href="/shop">
                View All Products
                <ChevronRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </section>
      {/* Latest Blog Posts Section */}
      <section className="py-16">
        <div className="container mx-auto px-4">
          <div className="mb-12 text-center">
            <h2 className="text-3xl font-bold text-foreground mb-4">Latest Tech Insights</h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              Stay ahead with expert reviews, guides, and the latest tech news
            </p>
          </div>

          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 mb-8">
            {blogs.length > 0 ? (
              blogs.map((blog) => (
                <Card key={blog.id} className="group hover:shadow-lg transition-shadow flex flex-col">
                  <CardHeader>
                    <div className="w-full h-48 mb-2 overflow-hidden rounded-lg bg-muted flex items-center justify-center">
                      {blog.featuredImage ? (
                        <img 
                          src={blog.featuredImage} 
                          alt={blog.title}
                          className="w-full h-full object-cover"
                        />
                      ) : (
                        <div className="text-4xl">📝</div>
                      )}
                    </div>
                    {blog.categoryName && (
                      <Badge variant="secondary" className="w-fit mb-2">{blog.categoryName}</Badge>
                    )}
                    <CardTitle className="text-xl group-hover:text-primary transition-colors">
                      {blog.title}
                    </CardTitle>
                    <CardDescription className="line-clamp-3">
                      {blog.excerpt || blog.content.substring(0, 150) + '...'}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="flex-1 flex flex-col">
                    <div className="flex items-center justify-between text-sm text-muted-foreground mb-4">
                      <span>By {blog.authorName}</span>
                      <span>{blog.readTime || Math.ceil(blog.content.length / 1000)} min read</span>
                    </div>
                    <div className="flex items-center justify-between mt-auto">
                      <span className="text-sm text-muted-foreground">
                        {blog.publishedAt ? formatDate(blog.publishedAt) : formatDate(blog.createdAt)}
                      </span>
                      <Button variant="ghost" size="sm" asChild>
                        <Link href={`/blog/${blog.slug}`}>
                          Read More
                          <ArrowRight className="ml-1 h-3 w-3" />
                        </Link>
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              // Fallback message if no blogs
              <div className="col-span-full text-center text-muted-foreground py-12">
                No blog posts available yet. Check back soon!
              </div>
            )}
          </div>

          <div className="text-center">
            <Button variant="outline" size="lg" asChild>
              <Link href="/blogs">
                View All Articles
                <ChevronRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="py-16 bg-muted/30">
        <div className="container mx-auto px-4">
          <div className="grid gap-8 md:grid-cols-3 text-center">
            <div>
              <div className="text-4xl font-bold text-primary mb-2">10,000+</div>
              <div className="text-muted-foreground">Happy Customers</div>
            </div>
            <div>
              <div className="text-4xl font-bold text-primary mb-2">500+</div>
              <div className="text-muted-foreground">Product Reviews</div>
            </div>
            <div>
              <div className="text-4xl font-bold text-primary mb-2">50+</div>
              <div className="text-muted-foreground">Expert Authors</div>
            </div>
          </div>
        </div>
      </section>

      {/* Newsletter Section */}
      <section className="py-16">
        <div className="container mx-auto px-4">
          <Card className="max-w-4xl mx-auto">
            <CardHeader className="text-center">
              <CardTitle className="text-3xl">Stay Updated</CardTitle>
              <CardDescription className="text-lg">
                Get the latest tech news, exclusive deals, and product launches delivered to your inbox.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex flex-col sm:flex-row gap-4 max-w-md mx-auto">
                <Input
                  type="email"
                  placeholder="Enter your email"
                  className="flex-1"
                />
                <Button className="whitespace-nowrap">
                  Subscribe
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </div>
              <p className="text-center text-sm text-muted-foreground mt-4">
                Join 5,000+ subscribers. Unsubscribe at any time.
              </p>
            </CardContent>
          </Card>
        </div>
      </section>
    </div>
  );
};

export default MainLandingPage;
