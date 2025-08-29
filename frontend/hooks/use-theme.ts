"use client"

import { useTheme } from "next-themes"

export function useThemeToggle() {
  const { theme, setTheme, systemTheme } = useTheme()
  
  const toggleTheme = () => {
    if (theme === "dark") {
      setTheme("light")
    } else {
      setTheme("dark")
    }
  }
  
  const currentTheme = theme === "system" ? systemTheme : theme
  
  return {
    theme,
    systemTheme,
    currentTheme,
    setTheme,
    toggleTheme,
    isDark: currentTheme === "dark",
    isLight: currentTheme === "light",
    isSystem: theme === "system"
  }
}
