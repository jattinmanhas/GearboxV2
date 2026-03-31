"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { useState, useEffect } from "react"
import { useUserStore } from "@/lib/stores/user-store"
import {
  LayoutDashboard,
  ShoppingCart,
  Package,
  Users,
  User,
  ChevronDown,
  Shield,
  FileText,
  FolderTree,
  Warehouse,
  Ticket,
  CreditCard
} from "lucide-react"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

const navigationSections = [
  {
    title: "Overview",
    items: [
      { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
    ]
  },
  {
    title: "Commerce",
    items: [
      { name: "Products", href: "/dashboard/products", icon: Package },
      { name: "Product Categories", href: "/dashboard/categories", icon: FolderTree },
      { name: "Inventory", href: "/dashboard/inventory", icon: Warehouse },
      { name: "Orders", href: "/dashboard/orders", icon: ShoppingCart },
      { name: "Payments", href: "/dashboard/payments", icon: CreditCard },
      { name: "Coupons", href: "/dashboard/coupons", icon: Ticket },
    ]
  },
  {
    title: "Content",
    items: [
      { name: "Blog Posts", href: "/dashboard/blog", icon: FileText },
      { name: "Blog Categories", href: "/dashboard/blog/categories", icon: FolderTree },
    ]
  },
  {
    title: "Access",
    items: [
      { name: "Users", href: "/dashboard/users", icon: Users },
      { name: "Roles", href: "/dashboard/roles", icon: Shield },
    ]
  },
]

export function DashboardSidebar() {
  const pathname = usePathname()
  const { user } = useUserStore()
  const [mounted, setMounted] = useState(false)
  const [expandedSections, setExpandedSections] = useState<Set<string>>(
    new Set(["Overview", "Commerce", "Content"]) // Start with key sections expanded
  )

  useEffect(() => {
    setMounted(true)
  }, [])

  const toggleSection = (sectionTitle: string) => {
    setExpandedSections(prev => {
      const newSet = new Set(prev)
      if (newSet.has(sectionTitle)) {
        newSet.delete(sectionTitle)
      } else {
        newSet.add(sectionTitle)
      }
      return newSet
    })
  }


  return (
    <Sidebar className="bg-sidebar/95">
      <SidebarHeader className="px-4 py-3 border-b border-sidebar-border/60">
        <div className="flex items-center gap-3">
          <div className="relative flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-primary/90 to-primary/40 text-primary-foreground shadow-sm ring-1 ring-primary/20">
            <Package className="h-4 w-4" />
            <span className="absolute -bottom-1 -right-1 h-2.5 w-2.5 rounded-full bg-emerald-400 ring-2 ring-sidebar" />
          </div>
          <div className="grid flex-1 text-left leading-tight">
            <span className="truncate text-sm font-semibold tracking-wide">Gearbox</span>
            <span className="truncate text-[11px] text-muted-foreground">Operations Console</span>
          </div>
        </div>
      </SidebarHeader>

      <SidebarContent className="px-3 py-2">
        {navigationSections.map((section) => {
          const isExpanded = expandedSections.has(section.title)
          const hasActiveItem = section.items.some(item => pathname === item.href)

          return (
            <SidebarGroup key={section.title}>
              <SidebarGroupLabel
                className={`cursor-pointer rounded-lg px-3 transition-all duration-200 text-[11px] tracking-[0.18em] uppercase ${hasActiveItem
                    ? "text-sidebar-foreground"
                    : "text-sidebar-foreground/50 hover:text-sidebar-foreground"
                  }`}
                onClick={() => toggleSection(section.title)}
              >
                <div className="flex items-center justify-between w-full">
                  <span>{section.title}</span>
                  <div className={`transition-transform duration-200 ${isExpanded ? 'rotate-0' : '-rotate-90'}`}>
                    <ChevronDown className="h-4 w-4" />
                  </div>
                </div>
              </SidebarGroupLabel>

              <div className={`overflow-hidden transition-all duration-300 ease-in-out ${isExpanded ? 'max-h-[420px] opacity-100' : 'max-h-0 opacity-0'
                }`}>
                <SidebarGroupContent className="pt-2">
                  <SidebarMenu className="gap-1.5">
                    {section.items.map((item) => {
                      const isActive = pathname === item.href
                      return (
                        <SidebarMenuItem key={item.name}>
                          <SidebarMenuButton
                            asChild
                            isActive={isActive}
                            className="relative h-9 text-[13px] font-medium ml-1 rounded-xl px-3 text-sidebar-foreground/80 transition-all duration-200 hover:bg-sidebar-accent/60 hover:text-sidebar-foreground data-[active=true]:bg-gradient-to-r data-[active=true]:from-primary/15 data-[active=true]:to-transparent data-[active=true]:text-sidebar-foreground data-[active=true]:shadow-[inset_0_0_0_1px_hsl(var(--sidebar-border))]"
                          >
                            <Link href={item.href} className="flex items-center gap-2">
                              <item.icon className="h-4 w-4 flex-shrink-0 opacity-80" />
                              <span className="flex-1">{item.name}</span>
                            </Link>
                          </SidebarMenuButton>
                        </SidebarMenuItem>
                      )
                    })}
                  </SidebarMenu>
                </SidebarGroupContent>
              </div>
            </SidebarGroup>
          )
        })}
      </SidebarContent>

      <SidebarFooter className="px-3 py-2 border-t border-sidebar-border/60">
        <SidebarMenu className="space-y-1">
          <SidebarMenuItem>
            <div className="flex items-center gap-3 px-3 py-2 rounded-xl bg-gradient-to-r from-muted/40 to-transparent">
              <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-muted">
                <User className="h-3 w-3" />
              </div>
              <div className="grid flex-1 text-left leading-tight">
                {mounted ? (
                  <>
                    <span className="truncate text-xs font-semibold">
                      {user?.firstName && user?.lastName
                        ? `${user.firstName} ${user.lastName}`
                        : user?.username || 'User'
                      }
                    </span>
                    <span className="truncate text-xs text-muted-foreground">
                      {user?.email || user?.username || 'No email'}
                    </span>
                  </>
                ) : (
                  // Skeleton placeholder — matches SSR output so no hydration mismatch
                  <>
                    <span className="h-3 w-24 rounded bg-muted animate-pulse" />
                    <span className="h-2.5 w-32 rounded bg-muted animate-pulse mt-1" />
                  </>
                )}
              </div>
            </div>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
