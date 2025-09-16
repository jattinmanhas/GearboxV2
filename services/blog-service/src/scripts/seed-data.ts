import { db } from '../config/database';
import { blogPosts, categories } from '../config/schema';
import { BlogPostDomain } from '../domain/blog';
import postgres from 'postgres';

const connectionString = process.env.DATABASE_URL || 'postgresql://localhost:5432/blog_service';
const client = postgres(connectionString, {
  max: 10,
  idle_timeout: 20,
  connect_timeout: 10,
});

async function seedData() {
  try {
    console.log('🌱 Starting to seed data...');

    // Create sample categories
    console.log('📁 Creating categories...');
    const [techCategory, lifestyleCategory, businessCategory] = await db
      .insert(categories)
      .values([
        {
          name: 'Technology',
          slug: 'technology',
          description: 'Latest trends and insights in technology',
          color: '#3B82F6'
        },
        {
          name: 'Lifestyle',
          slug: 'lifestyle',
          description: 'Tips and stories about modern living',
          color: '#10B981'
        },
        {
          name: 'Business',
          slug: 'business',
          description: 'Entrepreneurship and business insights',
          color: '#F59E0B'
        }
      ])
      .returning();

    console.log('✅ Categories created successfully');

    // Create sample blog posts
    console.log('📝 Creating blog posts...');
    
    const samplePosts = [
      {
        title: 'Getting Started with TypeScript: A Comprehensive Guide',
        content: `
# Getting Started with TypeScript: A Comprehensive Guide

TypeScript has become one of the most popular programming languages in recent years, and for good reason. It brings static typing to JavaScript, making your code more robust, maintainable, and less prone to runtime errors.

## What is TypeScript?

TypeScript is a strongly typed programming language that builds on JavaScript, giving you better tooling at any scale. It's developed and maintained by Microsoft and is open source.

## Why Use TypeScript?

1. **Static Typing**: Catch errors at compile time rather than runtime
2. **Better IDE Support**: Enhanced autocomplete, refactoring, and navigation
3. **Improved Code Quality**: Self-documenting code with type annotations
4. **Large Scale Development**: Better suited for large applications and teams

## Setting Up TypeScript

To get started with TypeScript, you'll need to install it globally or in your project:

\`\`\`bash
npm install -g typescript
# or
npm install --save-dev typescript
\`\`\`

## Basic Types

TypeScript provides several basic types:

- \`string\`: Text data
- \`number\`: Numeric data
- \`boolean\`: True/false values
- \`array\`: Collections of items
- \`object\`: Complex data structures

## Conclusion

TypeScript is a powerful tool that can significantly improve your JavaScript development experience. While there's a learning curve, the benefits far outweigh the initial investment in time and effort.

Start small, gradually add types to your existing JavaScript code, and you'll soon see the benefits of this amazing language.
        `,
        authorId: 'user-1',
        authorName: 'John Doe',
        authorEmail: 'john.doe@example.com',
        status: 'published' as const,
        featuredImage: 'https://images.unsplash.com/photo-1516116216621-301e3777632f?w=800&h=400&fit=crop',
        tags: ['typescript', 'javascript', 'programming', 'web-development'],
        categoryId: techCategory.id,
        publishedAt: new Date('2024-01-15T10:00:00Z')
      },
      {
        title: 'The Art of Minimalist Living: Less is More',
        content: `
# The Art of Minimalist Living: Less is More

In our fast-paced, consumer-driven world, the concept of minimalist living has gained significant traction. But what does it really mean to live minimally, and how can it benefit your life?

## What is Minimalism?

Minimalism is a lifestyle that focuses on living with less. It's about intentionally promoting what you value most and removing everything that distracts you from it.

## The Benefits of Minimalist Living

### 1. Reduced Stress
When you have fewer possessions, you have less to worry about, maintain, and organize.

### 2. More Time
Less time spent cleaning, organizing, and managing stuff means more time for what truly matters.

### 3. Financial Freedom
Buying less means spending less, which can lead to significant financial benefits.

### 4. Environmental Impact
Consuming less reduces your environmental footprint and contributes to sustainability.

## How to Start Your Minimalist Journey

1. **Start Small**: Begin with one room or category of items
2. **Ask the Right Questions**: Does this item add value to my life?
3. **One In, One Out**: For every new item you bring in, remove one
4. **Digital Minimalism**: Don't forget about digital clutter

## Common Misconceptions

- Minimalism doesn't mean living in an empty white room
- It's not about being cheap or depriving yourself
- It's not a one-size-fits-all approach

## Conclusion

Minimalist living is about finding freedom through less. It's a personal journey that looks different for everyone. Start where you are, use what you have, and do what you can.

Remember, the goal isn't to have as little as possible, but to have exactly what you need to live a fulfilling life.
        `,
        authorId: 'user-2',
        authorName: 'Jane Smith',
        authorEmail: 'jane.smith@example.com',
        status: 'published' as const,
        featuredImage: 'https://images.unsplash.com/photo-1484101403633-562f0dc0d60a?w=800&h=400&fit=crop',
        tags: ['minimalism', 'lifestyle', 'organization', 'wellness'],
        categoryId: lifestyleCategory.id,
        publishedAt: new Date('2024-01-20T14:30:00Z')
      },
      {
        title: 'Building a Successful Startup: Lessons from the Trenches',
        content: `
# Building a Successful Startup: Lessons from the Trenches

Starting a business is one of the most challenging yet rewarding experiences you can have. After building and scaling multiple startups, I've learned valuable lessons that I wish I knew when I started.

## The Reality of Entrepreneurship

Building a startup is not glamorous. It's filled with uncertainty, long hours, and constant problem-solving. But it's also incredibly rewarding when you see your vision come to life.

## Key Lessons Learned

### 1. Start with a Problem, Not a Solution
The most successful startups solve real problems that people are willing to pay for. Don't fall in love with your solution; fall in love with the problem.

### 2. Validate Early and Often
Don't spend months building something nobody wants. Get out of the building and talk to potential customers.

### 3. Focus on One Thing
It's tempting to try to do everything, but focus is crucial. Do one thing exceptionally well before expanding.

### 4. Build a Strong Team
Your team is your most valuable asset. Hire people who are smarter than you and share your vision.

### 5. Cash Flow is King
Revenue solves most problems. Focus on generating revenue early and often.

## Common Mistakes to Avoid

- Building in isolation without customer feedback
- Trying to be everything to everyone
- Not having a clear value proposition
- Ignoring unit economics
- Scaling too fast without proper systems

## The Importance of Persistence

Success in startups is often about persistence. Most successful entrepreneurs have failed multiple times before finding success. The key is to learn from each failure and keep moving forward.

## Conclusion

Building a startup is a marathon, not a sprint. It requires patience, persistence, and a willingness to learn and adapt. Focus on solving real problems, validate your ideas, and build a strong team.

Remember, every successful entrepreneur was once a beginner. The only difference is they didn't give up.
        `,
        authorId: 'user-3',
        authorName: 'Mike Johnson',
        authorEmail: 'mike.johnson@example.com',
        status: 'published' as const,
        featuredImage: 'https://images.unsplash.com/photo-1551434678-e076c223a692?w=800&h=400&fit=crop',
        tags: ['startup', 'entrepreneurship', 'business', 'leadership'],
        categoryId: businessCategory.id,
        publishedAt: new Date('2024-01-25T09:15:00Z')
      },
      {
        title: 'The Future of Web Development: Trends to Watch in 2024',
        content: `
# The Future of Web Development: Trends to Watch in 2024

The web development landscape is constantly evolving, and 2024 promises to bring exciting new trends and technologies. Let's explore what's on the horizon.

## Emerging Technologies

### 1. WebAssembly (WASM)
WebAssembly is revolutionizing web performance by allowing code written in languages like C, C++, and Rust to run in browsers at near-native speed.

### 2. Edge Computing
Moving computation closer to users reduces latency and improves performance, making edge computing a key trend.

### 3. Progressive Web Apps (PWAs)
PWAs continue to bridge the gap between web and native apps, offering app-like experiences in browsers.

## Framework Evolution

### React 18+ Features
- Concurrent rendering
- Suspense for data fetching
- Server components

### Vue 3 Composition API
- Better TypeScript support
- Improved performance
- More flexible component composition

### SvelteKit
- Compile-time optimizations
- Smaller bundle sizes
- Better developer experience

## Development Tools

### Vite
Lightning-fast build tool that's becoming the standard for modern web development.

### Turbopack
Next.js's new bundler promises 10x faster builds than Webpack.

## Best Practices

1. **Performance First**: Core Web Vitals are more important than ever
2. **Accessibility**: Building inclusive web experiences
3. **Security**: Protecting user data and preventing attacks
4. **Sustainability**: Building eco-friendly web applications

## Conclusion

The future of web development is bright, with new technologies and tools making it easier to build fast, secure, and accessible web applications. Stay curious, keep learning, and embrace the change.

The best developers are those who continuously adapt and learn new technologies while maintaining a solid foundation in the fundamentals.
        `,
        authorId: 'user-1',
        authorName: 'John Doe',
        authorEmail: 'john.doe@example.com',
        status: 'published' as const,
        featuredImage: 'https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=800&h=400&fit=crop',
        tags: ['web-development', 'javascript', 'react', 'vue', 'svelte'],
        categoryId: techCategory.id,
        publishedAt: new Date('2024-02-01T11:45:00Z')
      },
      {
        title: 'Draft: The Psychology of Color in Web Design',
        content: `
# The Psychology of Color in Web Design

Color is one of the most powerful tools in a designer's arsenal. It can evoke emotions, influence decisions, and create memorable experiences. Understanding color psychology is crucial for effective web design.

## The Science Behind Color Psychology

Colors have the ability to trigger specific emotional responses in users. This is why choosing the right color palette is so important for your website's success.

## Color Meanings and Associations

### Red
- Energy, passion, urgency
- Often used for call-to-action buttons
- Can also represent danger or warning

### Blue
- Trust, stability, professionalism
- Popular choice for corporate websites
- Associated with technology and finance

### Green
- Growth, nature, harmony
- Commonly used for health and environmental sites
- Represents money and prosperity

## Best Practices for Color Selection

1. **Consider Your Brand**: Colors should align with your brand personality
2. **Think About Your Audience**: Different cultures may interpret colors differently
3. **Use Color Theory**: Complementary and analogous color schemes work well
4. **Test Accessibility**: Ensure sufficient contrast ratios
5. **Limit Your Palette**: Too many colors can be overwhelming

## Common Mistakes to Avoid

- Using too many colors
- Ignoring accessibility guidelines
- Not considering cultural differences
- Choosing colors based on personal preference only

## Tools for Color Selection

- Adobe Color
- Coolors.co
- Material Design Color Tool
- WebAIM Contrast Checker

## Conclusion

Color psychology is a complex field, but understanding the basics can significantly improve your web design. Take time to research your target audience and test different color combinations to find what works best for your specific use case.

Remember, the best color choices are those that support your content and enhance the user experience.
        `,
        authorId: 'user-2',
        authorName: 'Jane Smith',
        authorEmail: 'jane.smith@example.com',
        status: 'draft' as const,
        featuredImage: 'https://images.unsplash.com/photo-1541701494587-cb58502866ab?w=800&h=400&fit=crop',
        tags: ['design', 'color-theory', 'psychology', 'ux'],
        categoryId: techCategory.id
      }
    ];

    // Process each post to add calculated fields
    const processedPosts = samplePosts.map(post => {
      const processedPost = BlogPostDomain.prepareBlogPostForCreation(post);
      return {
        ...processedPost,
        publishedAt: post.publishedAt
      };
    });

    await db.insert(blogPosts).values(processedPosts);

    console.log('✅ Blog posts created successfully');
    console.log('🎉 Data seeding completed!');
    
    // Display summary
    console.log('\n📊 Summary:');
    console.log(`- Categories: 3`);
    console.log(`- Blog Posts: ${processedPosts.length}`);
    console.log(`- Published: ${processedPosts.filter(p => p.status === 'published').length}`);
    console.log(`- Drafts: ${processedPosts.filter(p => p.status === 'draft').length}`);

  } catch (error) {
    console.error('❌ Error seeding data:', error);
    process.exit(1);
  } finally {
    await client.end();
  }
}

// Run the seed function
if (require.main === module) {
  seedData();
}

export { seedData };
