import React from "react"

export default function ShopPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-4xl font-bold mb-6">Tech Shop</h1>
        
        <div className="space-y-6">
          <div className="p-6 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-3">Featured Products</h2>
            <p className="text-muted-foreground">
              Discover the latest and greatest tech products with expert reviews and competitive prices.
            </p>
          </div>

          <div className="grid gap-6 md:grid-cols-3">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div key={i} className="p-6 border rounded-lg">
                <div className="h-48 bg-muted rounded-lg mb-4 flex items-center justify-center">
                  <span className="text-muted-foreground">Product Image {i}</span>
                </div>
                <h3 className="text-xl font-semibold mb-2">Product {i}</h3>
                <p className="text-muted-foreground mb-4">
                  High-quality tech product with excellent features and performance.
                </p>
                <div className="flex items-center justify-between">
                  <span className="text-2xl font-bold">$299.99</span>
                  <button className="bg-primary text-primary-foreground px-4 py-2 rounded-md hover:bg-primary/90">
                    Add to Cart
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
