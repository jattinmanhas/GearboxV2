import React from "react"

export default function AboutPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold mb-6">About GearBox</h1>
        
        <div className="space-y-8">
          <div className="p-6 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-4">Our Mission</h2>
            <p className="text-muted-foreground leading-relaxed">
              At GearBox, we're passionate about technology and committed to bringing you the latest 
              and greatest tech products with expert reviews, competitive prices, and exceptional 
              customer service. Our mission is to make technology accessible and understandable for everyone.
            </p>
          </div>

          <div className="p-6 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-4">What We Do</h2>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <h3 className="font-semibold mb-2">Product Reviews</h3>
                <p className="text-muted-foreground text-sm">
                  In-depth reviews of the latest tech products to help you make informed decisions.
                </p>
              </div>
              <div>
                <h3 className="font-semibold mb-2">Tech News</h3>
                <p className="text-muted-foreground text-sm">
                  Stay updated with the latest technology trends and industry news.
                </p>
              </div>
              <div>
                <h3 className="font-semibold mb-2">Expert Guides</h3>
                <p className="text-muted-foreground text-sm">
                  Comprehensive guides to help you get the most out of your tech products.
                </p>
              </div>
              <div>
                <h3 className="font-semibold mb-2">Community</h3>
                <p className="text-muted-foreground text-sm">
                  Join our community of tech enthusiasts and share your experiences.
                </p>
              </div>
            </div>
          </div>

          <div className="p-6 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-4">Our Team</h2>
            <p className="text-muted-foreground leading-relaxed">
              Our team consists of passionate tech experts, reviewers, and writers who are dedicated 
              to providing you with accurate, unbiased information about the latest technology. 
              We believe in transparency, quality, and putting our customers first.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
