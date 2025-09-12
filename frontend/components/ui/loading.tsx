"use client"

import { Loader2 } from "lucide-react"
import { cn } from "@/lib/utils"

interface LoadingSpinnerProps {
  size?: "sm" | "md" | "lg" | "xl"
  className?: string
  text?: string
}

export function LoadingSpinner({ 
  size = "md", 
  className,
  text 
}: LoadingSpinnerProps) {
  const sizeClasses = {
    sm: "w-4 h-4",
    md: "w-6 h-6", 
    lg: "w-8 h-8",
    xl: "w-12 h-12"
  }

  return (
    <div className={cn("flex items-center justify-center", className)}>
      <div className="flex flex-col items-center space-y-2">
        <div className="relative">
          <div className={cn(
            "border-4 border-primary/20 border-t-primary rounded-full animate-spin",
            sizeClasses[size]
          )}></div>
          <Loader2 className={cn(
            "text-primary animate-spin absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2",
            size === "sm" ? "w-2 h-2" : size === "md" ? "w-3 h-3" : size === "lg" ? "w-4 h-4" : "w-6 h-6"
          )} />
        </div>
        {text && (
          <div className="text-center">
            <p className="text-sm text-muted-foreground">{text}</p>
          </div>
        )}
      </div>
    </div>
  )
}

interface SkeletonProps {
  className?: string
}

export function Skeleton({ className }: SkeletonProps) {
  return (
    <div className={cn("animate-pulse rounded-md bg-muted", className)} />
  )
}

interface SkeletonCardProps {
  className?: string
}

export function SkeletonCard({ className }: SkeletonCardProps) {
  return (
    <div className={cn("border rounded-lg p-4 space-y-3", className)}>
      <div className="flex items-start justify-between">
        <div className="flex items-center space-x-3">
          <Skeleton className="w-12 h-12 rounded-lg" />
          <div className="space-y-2">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-3 w-16" />
          </div>
        </div>
        <Skeleton className="h-6 w-16" />
      </div>
      <div className="space-y-2">
        <Skeleton className="h-3 w-full" />
        <Skeleton className="h-3 w-3/4" />
      </div>
      <div className="space-y-2">
        <Skeleton className="h-3 w-1/2" />
        <Skeleton className="h-3 w-1/3" />
      </div>
      <div className="flex space-x-2">
        <Skeleton className="h-8 w-16" />
        <Skeleton className="h-8 w-16" />
      </div>
    </div>
  )
}

interface SkeletonRowProps {
  className?: string
}

export function SkeletonRow({ className }: SkeletonRowProps) {
  return (
    <div className={cn("flex items-center space-x-4 p-4 border rounded-lg", className)}>
      <Skeleton className="w-10 h-10 rounded" />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-4 w-1/4" />
        <Skeleton className="h-3 w-1/6" />
      </div>
      <Skeleton className="h-6 w-16" />
      <Skeleton className="h-4 w-12" />
      <Skeleton className="h-4 w-20" />
      <Skeleton className="w-8 h-8 rounded" />
    </div>
  )
}

interface LoadingStateProps {
  type?: "spinner" | "skeleton" | "empty"
  viewMode?: "grid" | "list"
  itemCount?: number
  text?: string
  emptyText?: string
  emptyIcon?: React.ReactNode
  className?: string
}

export function LoadingState({
  type = "spinner",
  viewMode = "list",
  itemCount = 5,
  text = "Loading...",
  emptyText = "No items found",
  emptyIcon,
  className
}: LoadingStateProps) {
  if (type === "empty") {
    return (
      <div className={cn("flex flex-col items-center justify-center py-12", className)}>
        <div className="w-16 h-16 bg-muted rounded-full flex items-center justify-center mb-4">
          {emptyIcon || <div className="w-8 h-8 bg-muted-foreground/20 rounded" />}
        </div>
        <h3 className="text-lg font-semibold text-foreground mb-2">No Items Found</h3>
        <p className="text-muted-foreground text-center max-w-md mb-6">
          {emptyText}
        </p>
      </div>
    )
  }

  if (type === "skeleton") {
    return (
      <div className={cn("space-y-4", className)}>
        {viewMode === "list" ? (
          <div className="space-y-3">
            {Array.from({ length: itemCount }).map((_, i) => (
              <SkeletonRow key={i} />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {Array.from({ length: itemCount }).map((_, i) => (
              <SkeletonCard key={i} />
            ))}
          </div>
        )}
      </div>
    )
  }

  return (
    <div className={cn("space-y-4", className)}>
      <div className="flex items-center justify-center py-8">
        <div className="flex flex-col items-center space-y-4">
          <div className="relative">
            <div className="w-12 h-12 border-4 border-primary/20 border-t-primary rounded-full animate-spin"></div>
            <Loader2 className="w-6 h-6 text-primary animate-spin absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2" />
          </div>
          <div className="text-center">
            <h3 className="font-medium text-foreground">{text}</h3>
            <p className="text-sm text-muted-foreground">Please wait while we fetch your data...</p>
          </div>
        </div>
      </div>
    </div>
  )
}
