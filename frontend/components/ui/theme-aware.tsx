"use client"

import { useTheme } from "next-themes"
import { ReactNode } from "react"

interface ThemeAwareProps {
  children: ReactNode
  light?: ReactNode
  dark?: ReactNode
  system?: ReactNode
  fallback?: ReactNode
}

export function ThemeAware({ 
  children, 
  light, 
  dark, 
  system, 
  fallback 
}: ThemeAwareProps) {
  const { theme, systemTheme } = useTheme()
  
  const currentTheme = theme === "system" ? systemTheme : theme
  
  if (currentTheme === "light" && light !== undefined) {
    return <>{light}</>
  }
  
  if (currentTheme === "dark" && dark !== undefined) {
    return <>{dark}</>
  }
  
  if (theme === "system" && system !== undefined) {
    return <>{system}</>
  }
  
  if (fallback !== undefined) {
    return <>{fallback}</>
  }
  
  return <>{children}</>
}

// Convenience components for common use cases
export function LightOnly({ children }: { children: ReactNode }) {
  return <ThemeAware light={children} children={null} />
}

export function DarkOnly({ children }: { children: ReactNode }) {
  return <ThemeAware dark={children} children={null} />
}

export function SystemOnly({ children }: { children: ReactNode }) {
  return <ThemeAware system={children} children={null} />
}
