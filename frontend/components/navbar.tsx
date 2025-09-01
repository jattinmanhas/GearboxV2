"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { ThemeToggle } from "@/components/ui/theme-toggle"
import { CpuIcon, Menu, X, User, Settings, LogOut } from "lucide-react"
import { useUserStore } from "@/lib/stores/user-store"

export type NavItems = {
  label: string
  link: string
}

export const navLinks: NavItems[] = [
  {
    label: 'Blogs',
    link: '/blogs'
  },
  {
    label: 'Shop',
    link: '/shop'
  },
  {
    label: 'Contact',
    link: '/contact'
  },
  {
    label: 'About',
    link: '/about'
  },
]

export function Navbar() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [isNavbarVisible, setIsNavbarVisible] = useState(true)
  const [lastScrollY, setLastScrollY] = useState(0)
  const pathname = usePathname()
  
  // Get user data from Zustand store
  const { user, isAuthenticated, logout } = useUserStore()
  
  // Debug: Log user data
  useEffect(() => {
    console.log("Navbar - User data:", user)
    console.log("Navbar - Is authenticated:", isAuthenticated)
  }, [user, isAuthenticated])

  const isActiveLink = (link: string) => {
    return pathname === link
  }

  const handleLogout = async () => {
    try {
      await logout()
      // Redirect to home page after logout
      window.location.href = '/'
    } catch (error) {
      console.error("Logout error:", error)
      // Still redirect even if logout fails
      window.location.href = '/'
    }
  }

  // Scroll behavior for navbar hide/show
  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY
      
      // Show navbar when scrolling up or at the top
      if (currentScrollY < lastScrollY || currentScrollY < 10) {
        setIsNavbarVisible(true)
      } 
      // Hide navbar when scrolling down (but only after scrolling down a bit)
      else if (currentScrollY > lastScrollY && currentScrollY > 100) {
        setIsNavbarVisible(false)
      }
      
      setLastScrollY(currentScrollY)
    }

    // Throttle scroll events for better performance
    let ticking = false
    const throttledHandleScroll = () => {
      if (!ticking) {
        requestAnimationFrame(() => {
          handleScroll()
          ticking = false
        })
        ticking = true
      }
    }

    window.addEventListener('scroll', throttledHandleScroll, { passive: true })
    
    return () => {
      window.removeEventListener('scroll', throttledHandleScroll)
    }
  }, [lastScrollY])

  return (
    <nav className={cn(
      "fixed top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 transition-transform duration-300 ease-in-out",
      isNavbarVisible ? "translate-y-0" : "-translate-y-full"
    )}>
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2 font-medium">
            <div className="bg-primary text-primary-foreground flex size-8 items-center justify-center rounded-md">
              <CpuIcon className="size-5" />
            </div>
            <span className="text-xl font-bold text-foreground">GearBox</span>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-8">
            {navLinks.map((link) => (
              <Link
                key={link.label}
                href={link.link}
                className={cn(
                  "text-sm font-medium transition-colors hover:text-primary",
                  isActiveLink(link.link)
                    ? "text-primary"
                    : "text-muted-foreground"
                )}
              >
                {link.label}
              </Link>
            ))}
          </div>

          {/* Right side - User menu and theme toggle */}
          <div className="flex items-center gap-4">
            {/* Theme Toggle */}
            <ThemeToggle />

            {/* User Menu */}
            {isAuthenticated && user ? (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                    <Avatar className="h-8 w-8">
                      <AvatarImage src={user.avatar} alt={`${user.firstName} ${user.lastName}`} />
                      <AvatarFallback>
                        {user.firstName && user.lastName 
                          ? `${user.firstName} ${user.lastName}`.split(' ').map(n => n[0]).join('').toUpperCase()
                          : user.username ? user.username.substring(0, 2).toUpperCase() : 'U'
                        }
                      </AvatarFallback>
                    </Avatar>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-56" align="end" forceMount>
                  <DropdownMenuLabel className="font-normal">
                    <div className="flex flex-col space-y-1">
                      <p className="text-sm font-medium leading-none">
                        {user.firstName && user.lastName 
                          ? `${user.firstName} ${user.lastName}`
                          : user.username || 'User'
                        }
                      </p>
                      <p className="text-xs leading-none text-muted-foreground">
                        {user.email || user.username}
                      </p>
                    </div>
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem asChild>
                    <Link href="/profile" className="flex items-center">
                      <User className="mr-2 h-4 w-4" />
                      <span>Profile</span>
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem asChild>
                    <Link href="/settings" className="flex items-center">
                      <Settings className="mr-2 h-4 w-4" />
                      <span>Settings</span>
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={handleLogout} className="text-red-600">
                    <LogOut className="mr-2 h-4 w-4" />
                    <span>Log out</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            ) : (
              <div className="hidden md:flex items-center gap-2">
                <Button variant="ghost" asChild>
                  <Link href="/login">Sign in</Link>
                </Button>
                <Button asChild>
                  <Link href="/register">Sign up</Link>
                </Button>
              </div>
            )}

            {/* Mobile menu button */}
            <Button
              variant="ghost"
              size="sm"
              className="md:hidden"
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            >
              {isMobileMenuOpen ? (
                <X className="h-5 w-5" />
              ) : (
                <Menu className="h-5 w-5" />
              )}
            </Button>
          </div>
        </div>

        {/* Mobile Navigation */}
        {isMobileMenuOpen && (
          <div className="md:hidden">
            <div className="px-2 pt-2 pb-3 space-y-1 border-t">
              {navLinks.map((link) => (
                <Link
                  key={link.label}
                  href={link.link}
                  className={cn(
                    "block px-3 py-2 text-base font-medium rounded-md transition-colors",
                    isActiveLink(link.link)
                      ? "text-primary bg-primary/10"
                      : "text-muted-foreground hover:text-primary hover:bg-primary/10"
                  )}
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  {link.label}
                </Link>
              ))}
              
              {/* Mobile auth buttons */}
              {!isAuthenticated && (
                <div className="pt-4 space-y-2">
                  <Button variant="ghost" asChild className="w-full justify-start">
                    <Link href="/login" onClick={() => setIsMobileMenuOpen(false)}>
                      Sign in
                    </Link>
                  </Button>
                  <Button asChild className="w-full">
                    <Link href="/register" onClick={() => setIsMobileMenuOpen(false)}>
                      Sign up
                    </Link>
                  </Button>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </nav>
  )
}
