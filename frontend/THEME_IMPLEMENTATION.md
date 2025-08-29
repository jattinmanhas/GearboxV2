# Theme Implementation Guide

This document explains how the light/dark mode theme system is implemented in the GearboxV2 frontend application.

## Overview

The theme system uses:
- **next-themes**: For theme management and persistence
- **CSS Variables**: For theme-specific colors and styles
- **Tailwind CSS**: For utility classes that automatically adapt to themes
- **React Context**: For providing theme state throughout the app

## Components

### 1. ThemeProvider (`lib/theme-provider.tsx`)
Wraps the entire application and provides theme context. Configured in `app/layout.tsx`.

**Features:**
- Automatic theme detection
- System theme support
- Theme persistence in localStorage
- Hydration-safe rendering

### 2. ThemeToggle (`components/ui/theme-toggle.tsx`)
A dropdown button that allows users to switch between:
- Light theme
- Dark theme  
- System theme (follows OS preference)

**Usage:**
```tsx
import { ThemeToggle } from '@/components/ui/theme-toggle'

// Add to any page or component
<ThemeToggle />
```

### 3. ThemeAware (`components/ui/theme-aware.tsx`)
Conditionally renders content based on the current theme.

**Usage:**
```tsx
import { ThemeAware, LightOnly, DarkOnly } from '@/components/ui/theme-aware'

// Conditional rendering
<ThemeAware 
  light={<SunIcon />}
  dark={<MoonIcon />}
  fallback={<DefaultIcon />}
/>

// Or use convenience components
<LightOnly>Only visible in light mode</LightOnly>
<DarkOnly>Only visible in dark mode</DarkOnly>
```

## Hooks

### useThemeToggle (`hooks/use-theme.ts`)
Custom hook providing theme state and functions.

**Usage:**
```tsx
import { useThemeToggle } from '@/hooks/use-theme'

function MyComponent() {
  const { theme, toggleTheme, isDark, isLight } = useThemeToggle()
  
  return (
    <button onClick={toggleTheme}>
      Current theme: {theme}
    </button>
  )
}
```

## CSS Variables

The theme system uses CSS custom properties defined in `app/globals.css`:

**Light Theme Variables:**
- `--background`: White backgrounds
- `--foreground`: Dark text
- `--card`: Card backgrounds
- `--primary`: Primary brand colors

**Dark Theme Variables:**
- `--background`: Dark backgrounds  
- `--foreground`: Light text
- `--card`: Dark card backgrounds
- `--primary`: Adjusted primary colors

## Tailwind Integration

Use Tailwind's semantic color classes that automatically adapt to themes:

```tsx
// These automatically use the correct theme colors
<div className="bg-background text-foreground">
  <div className="bg-card text-card-foreground border-border">
    <p className="text-muted-foreground">Content</p>
  </div>
</div>
```

## Adding Theme Toggle to New Pages

1. **Import the component:**
```tsx
import { ThemeToggle } from '@/components/ui/theme-toggle'
```

2. **Add to your layout:**
```tsx
<div className="flex justify-between items-center">
  <h1>Page Title</h1>
  <ThemeToggle />
</div>
```

## Best Practices

1. **Use semantic color classes** instead of hardcoded colors
2. **Test both themes** during development
3. **Use the ThemeAware component** for theme-specific content
4. **Avoid mixing theme systems** - stick to the established pattern

## Troubleshooting

**Theme not persisting:**
- Check that ThemeProvider is wrapping your app
- Verify localStorage is available

**Hydration warnings:**
- The `suppressHydrationWarning` attribute is already added to `<html>`
- This is normal for theme systems

**Colors not changing:**
- Ensure you're using CSS variable-based classes
- Check that the theme toggle is working
- Verify CSS variables are properly defined

## Example Implementation

See these files for complete examples:
- `app/page.tsx` - Main page with theme toggle
- `app/(auth)/login/page.tsx` - Login page with theme toggle
- `app/(auth)/register/page.tsx` - Register page with theme toggle
