# Loading Components

A collection of reusable loading components for consistent loading states across the application.

## Components

### LoadingSpinner
A simple spinner with optional text.

```tsx
import { LoadingSpinner } from "@/components/ui/loading"

// Basic spinner
<LoadingSpinner />

// With custom size and text
<LoadingSpinner size="lg" text="Loading data..." />

// With custom styling
<LoadingSpinner className="my-4" />
```

**Props:**
- `size?: "sm" | "md" | "lg" | "xl"` - Size of the spinner (default: "md")
- `className?: string` - Additional CSS classes
- `text?: string` - Optional text to display below the spinner

### Skeleton
A simple skeleton placeholder.

```tsx
import { Skeleton } from "@/components/ui/loading"

// Basic skeleton
<Skeleton className="h-4 w-32" />

// Card skeleton
<Skeleton className="h-20 w-full rounded-lg" />
```

**Props:**
- `className?: string` - CSS classes for sizing and styling

### SkeletonCard
A pre-built skeleton for card layouts.

```tsx
import { SkeletonCard } from "@/components/ui/loading"

<SkeletonCard />
```

### SkeletonRow
A pre-built skeleton for table/list rows.

```tsx
import { SkeletonRow } from "@/components/ui/loading"

<SkeletonRow />
```

### LoadingState
A comprehensive loading component that handles different states.

```tsx
import { LoadingState } from "@/components/ui/loading"

// Skeleton loading for list view
<LoadingState
  type="skeleton"
  viewMode="list"
  itemCount={5}
  text="Loading Categories"
/>

// Skeleton loading for grid view
<LoadingState
  type="skeleton"
  viewMode="grid"
  itemCount={6}
  text="Loading Products"
/>

// Empty state
<LoadingState
  type="empty"
  emptyText="No items found. Create your first item to get started."
  emptyIcon={<Folder className="h-8 w-8" />}
/>

// Simple spinner
<LoadingState
  type="spinner"
  text="Loading data..."
/>
```

**Props:**
- `type?: "spinner" | "skeleton" | "empty"` - Type of loading state (default: "spinner")
- `viewMode?: "grid" | "list"` - Layout mode for skeleton (default: "list")
- `itemCount?: number` - Number of skeleton items to show (default: 5)
- `text?: string` - Loading text (default: "Loading...")
- `emptyText?: string` - Text for empty state
- `emptyIcon?: React.ReactNode` - Icon for empty state
- `className?: string` - Additional CSS classes

## Usage Examples

### In a Data Table
```tsx
{loading ? (
  <LoadingState
    type="skeleton"
    viewMode="list"
    itemCount={10}
    text="Loading Data"
  />
) : (
  <DataTable data={data} />
)}
```

### In a Card Grid
```tsx
{loading ? (
  <LoadingState
    type="skeleton"
    viewMode="grid"
    itemCount={8}
    text="Loading Cards"
  />
) : (
  <CardGrid cards={cards} />
)}
```

### Empty State
```tsx
{data.length === 0 ? (
  <LoadingState
    type="empty"
    emptyText="No data available. Click the button above to add some."
    emptyIcon={<Plus className="h-8 w-8" />}
  />
) : (
  <DataList data={data} />
)}
```

### Simple Loading
```tsx
{loading && (
  <LoadingSpinner size="lg" text="Saving changes..." />
)}
```

## Benefits

- **Consistency**: All loading states look and behave the same across the app
- **Reusability**: Easy to use in any component
- **Customizable**: Flexible props for different use cases
- **Accessible**: Proper ARIA labels and semantic HTML
- **Performance**: Lightweight and optimized
- **Maintainable**: Centralized loading logic
