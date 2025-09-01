import React from "react"

export default function BlogsPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold mb-6">Tech Blogs</h1>
        
        <div className="space-y-6">
          <div className="p-6 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-3">Latest Tech Insights</h2>
            <p className="text-muted-foreground">
              Stay updated with the latest technology trends, reviews, and insights from our expert team.
            </p>
          </div>

          <div className="grid gap-6 md:grid-cols-2">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="p-6 border rounded-lg">
                <h3 className="text-xl font-semibold mb-2">Blog Post {i}</h3>
                <p className="text-muted-foreground mb-4">
                  This is a sample blog post about the latest in technology and gadgets.
                </p>
                <button className="text-primary hover:underline">Read More</button>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
