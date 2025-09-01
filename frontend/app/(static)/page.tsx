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
  Users,
  TrendingUp,
} from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

const MainLandingPage = async () => {
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
            {[
              {
                title: "iPhone 15 Pro",
                description: "Latest flagship with titanium design and A17 Pro chip",
                price: "$999",
                rating: 4.8,
                reviews: 1247,
                image: "📱"
              },
              {
                title: "MacBook Pro M3",
                description: "Powerful laptop for professionals and creators",
                price: "$1,999",
                rating: 4.9,
                reviews: 892,
                image: "💻"
              },
              {
                title: "AirPods Pro 2",
                description: "Premium wireless earbuds with active noise cancellation",
                price: "$249",
                rating: 4.7,
                reviews: 2156,
                image: "🎧"
              }
            ].map((product, index) => (
              <Card key={index} className="group hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="text-4xl mb-2">{product.image}</div>
                  <CardTitle className="text-xl">{product.title}</CardTitle>
                  <CardDescription>{product.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between mb-4">
                    <span className="text-2xl font-bold text-primary">{product.price}</span>
                    <div className="flex items-center gap-1">
                      <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                      <span className="text-sm text-muted-foreground">
                        {product.rating} ({product.reviews})
                      </span>
                    </div>
                  </div>
                  <Button className="w-full" asChild>
                    <Link href="/shop">
                      View Details
                      <ArrowRight className="ml-2 h-4 w-4" />
                    </Link>
                  </Button>
                </CardContent>
              </Card>
            ))}
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
            {[
              {
                title: "iPhone 15 Pro Review: Is It Worth the Upgrade?",
                excerpt: "We dive deep into Apple's latest flagship to see if the titanium design and A17 Pro chip justify the premium price tag.",
                author: "Sarah Chen",
                date: "2 days ago",
                readTime: "8 min read",
                category: "Reviews",
                image: "📱"
              },
              {
                title: "Best Laptops for Developers in 2024",
                excerpt: "Our comprehensive guide to finding the perfect development machine, from budget-friendly options to high-end workstations.",
                author: "Mike Rodriguez",
                date: "5 days ago",
                readTime: "12 min read",
                category: "Guides",
                image: "💻"
              },
              {
                title: "The Future of AI in Smartphones",
                excerpt: "Exploring how artificial intelligence is revolutionizing mobile technology and what to expect in the coming years.",
                author: "Alex Kim",
                date: "1 week ago",
                readTime: "6 min read",
                category: "Tech News",
                image: "🤖"
              }
            ].map((article, index) => (
              <Card key={index} className="group hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="text-4xl mb-2">{article.image}</div>
                  <Badge variant="secondary" className="w-fit mb-2">{article.category}</Badge>
                  <CardTitle className="text-xl group-hover:text-primary transition-colors">
                    {article.title}
                  </CardTitle>
                  <CardDescription className="line-clamp-3">{article.excerpt}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between text-sm text-muted-foreground mb-4">
                    <span>By {article.author}</span>
                    <span>{article.readTime}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">{article.date}</span>
                    <Button variant="ghost" size="sm" asChild>
                      <Link href="/blogs">
                        Read More
                        <ArrowRight className="ml-1 h-3 w-3" />
                      </Link>
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
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
