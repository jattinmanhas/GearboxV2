"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { Home } from "lucide-react"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { SidebarTrigger } from "@/components/ui/sidebar"
import { ThemeToggle } from "@/components/ui/theme-toggle"
import { LogOut } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useUserStore } from "@/lib/stores/user-store"
import { useAuth } from "@/lib/contexts/auth-context"

// Breadcrumb mapping
const breadcrumbMap: Record<string, string> = {
  "/dashboard": "Dashboard",
  "/dashboard/analytics": "Analytics",
  "/dashboard/users": "Users",
  "/dashboard/roles": "Roles",
  "/dashboard/products": "Products",
  "/dashboard/orders": "Orders",
  "/dashboard/blog": "Blog Posts",
  "/dashboard/settings": "Settings",
}

function getBreadcrumbs(pathname: string) {
  const segments = pathname.split('/').filter(Boolean)
  const breadcrumbs = [
    { label: "Home", href: "/" }
  ]

  let currentPath = ""
  segments.forEach((segment) => {
    currentPath += `/${segment}`
    const label = breadcrumbMap[currentPath] || segment.charAt(0).toUpperCase() + segment.slice(1)
    breadcrumbs.push({ label, href: currentPath })
  })

  return breadcrumbs
}

export function DashboardBreadcrumbs() {
  const pathname = usePathname()
  const breadcrumbs = getBreadcrumbs(pathname)
  const { logout } = useUserStore()
  const { clearAuth } = useAuth()

  const handleLogout = async () => {
    try {
      await logout()
      clearAuth()
      window.location.href = '/'
    } catch (error) {
      console.error("Logout error:", error)
      clearAuth()
      window.location.href = '/'
    }
  }

  return (
    <header className="flex h-12 shrink-0 items-center gap-2 border-b px-4">
      <SidebarTrigger className="-ml-1" />

      <div className="flex items-center gap-2">
        <Breadcrumb>
          <BreadcrumbList>
            {breadcrumbs.map((breadcrumb, index) => (
              <div key={breadcrumb.href} className="flex items-center">
                <BreadcrumbItem>
                  {index === breadcrumbs.length - 1 ? (
                    <BreadcrumbPage className="text-sm">{breadcrumb.label}</BreadcrumbPage>
                  ) : (
                    <BreadcrumbLink asChild>
                      <Link href={breadcrumb.href} className="text-sm">
                        {breadcrumb.label === "Home" ? <Home className="h-3 w-3" /> : breadcrumb.label}
                      </Link>
                    </BreadcrumbLink>
                  )}
                </BreadcrumbItem>
                {index < breadcrumbs.length - 1 && <BreadcrumbSeparator />}
              </div>
            ))}
          </BreadcrumbList>
        </Breadcrumb>
      </div>

      <div className="ml-auto flex items-center gap-2">
        <ThemeToggle />
        <Button
          variant="ghost"
          size="icon"
          onClick={handleLogout}
          className="h-9 w-9 text-muted-foreground hover:text-foreground"
          title="Logout"
        >
          <LogOut className="h-4 w-4" />
          <span className="sr-only">Logout</span>
        </Button>
      </div>
    </header>
  )
}
