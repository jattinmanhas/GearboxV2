'use client';

import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkBreaks from 'remark-breaks';
import rehypeHighlight from 'rehype-highlight';
import rehypeRaw from 'rehype-raw';
import Image from 'next/image';
import { cn } from '@/lib/utils';

interface MarkdownRendererProps {
  content: string;
  className?: string;
}

export function MarkdownRenderer({ content, className }: MarkdownRendererProps) {
  return (
    <div className={cn("prose prose-lg max-w-none dark:prose-invert", className)}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm, remarkBreaks]}
        rehypePlugins={[rehypeHighlight, rehypeRaw]}
        components={{
          // Custom image component with Next.js Image optimization
          img: ({ src, alt, ...props }) => {
            if (!src || typeof src !== 'string') return null;
            
            return (
              <Image
                src={src}
                alt={alt || ''}
                width={800}
                height={400}
                className="rounded-lg object-cover w-full h-auto"
                {...(props as any)}
              />
            );
          },
          
          // Custom heading styles
          h1: ({ children, ...props }) => (
            <h1 className="text-4xl font-bold mb-6 mt-8 text-slate-900 dark:text-slate-100" {...props}>
              {children}
            </h1>
          ),
          h2: ({ children, ...props }) => (
            <h2 className="text-3xl font-bold mb-4 mt-6 text-slate-900 dark:text-slate-100" {...props}>
              {children}
            </h2>
          ),
          h3: ({ children, ...props }) => (
            <h3 className="text-2xl font-semibold mb-3 mt-5 text-slate-900 dark:text-slate-100" {...props}>
              {children}
            </h3>
          ),
          h4: ({ children, ...props }) => (
            <h4 className="text-xl font-semibold mb-2 mt-4 text-slate-900 dark:text-slate-100" {...props}>
              {children}
            </h4>
          ),
          
          // Custom paragraph styles
          p: ({ children, ...props }) => {
            // Check if the paragraph contains only an image
            const hasOnlyImage = React.Children.count(children) === 1 && 
              React.isValidElement(children) && 
              children.type === 'img';
            
            if (hasOnlyImage) {
              // If paragraph contains only an image, render it as a div to avoid nesting issues
              return (
                <div className="my-6" {...props}>
                  {children}
                </div>
              );
            }
            
            return (
              <p className="mb-4 text-slate-700 dark:text-slate-300 leading-relaxed" {...props}>
                {children}
              </p>
            );
          },
          
          // Custom list styles
          ul: ({ children, ...props }) => (
            <ul className="mb-4 ml-6 space-y-2 list-disc text-slate-700 dark:text-slate-300" {...props}>
              {children}
            </ul>
          ),
          ol: ({ children, ...props }) => (
            <ol className="mb-4 ml-6 space-y-2 list-decimal text-slate-700 dark:text-slate-300" {...props}>
              {children}
            </ol>
          ),
          li: ({ children, ...props }) => (
            <li className="leading-relaxed" {...props}>
              {children}
            </li>
          ),
          
          // Custom blockquote styles
          blockquote: ({ children, ...props }) => (
            <blockquote 
              className="border-l-4 border-blue-500 pl-4 py-2 my-4 bg-slate-50 dark:bg-slate-800 rounded-r-lg italic text-slate-600 dark:text-slate-400" 
              {...props}
            >
              {children}
            </blockquote>
          ),
          
          // Custom code styles
          code: ({ children, className, ...props }) => {
            const isInline = !className;
            if (isInline) {
              return (
                <code 
                  className="bg-slate-100 dark:bg-slate-800 px-1.5 py-0.5 rounded text-sm font-mono text-slate-800 dark:text-slate-200" 
                  {...props}
                >
                  {children}
                </code>
              );
            }
            return (
              <code className={className} {...props}>
                {children}
              </code>
            );
          },
          
          pre: ({ children, ...props }) => (
            <pre 
              className="bg-slate-900 dark:bg-slate-950 text-slate-100 p-4 rounded-lg overflow-x-auto my-4 border" 
              {...props}
            >
              {children}
            </pre>
          ),
          
          // Custom link styles
          a: ({ children, href, ...props }) => (
            <a 
              href={href}
              className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 underline decoration-2 underline-offset-2 hover:decoration-blue-800 dark:hover:decoration-blue-300 transition-colors" 
              target="_blank"
              rel="noopener noreferrer"
              {...props}
            >
              {children}
            </a>
          ),
          
          // Custom table styles
          table: ({ children, ...props }) => (
            <div className="overflow-x-auto my-6">
              <table className="min-w-full border-collapse border border-slate-300 dark:border-slate-600" {...props}>
                {children}
              </table>
            </div>
          ),
          thead: ({ children, ...props }) => (
            <thead className="bg-slate-50 dark:bg-slate-800" {...props}>
              {children}
            </thead>
          ),
          tbody: ({ children, ...props }) => (
            <tbody className="divide-y divide-slate-200 dark:divide-slate-700" {...props}>
              {children}
            </tbody>
          ),
          tr: ({ children, ...props }) => (
            <tr className="hover:bg-slate-50 dark:hover:bg-slate-800" {...props}>
              {children}
            </tr>
          ),
          th: ({ children, ...props }) => (
            <th className="px-4 py-2 text-left font-semibold text-slate-900 dark:text-slate-100 border border-slate-300 dark:border-slate-600" {...props}>
              {children}
            </th>
          ),
          td: ({ children, ...props }) => (
            <td className="px-4 py-2 text-slate-700 dark:text-slate-300 border border-slate-300 dark:border-slate-600" {...props}>
              {children}
            </td>
          ),
          
          // Custom horizontal rule
          hr: ({ ...props }) => (
            <hr className="my-8 border-0 border-t border-slate-200 dark:border-slate-700" {...props} />
          ),
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}
