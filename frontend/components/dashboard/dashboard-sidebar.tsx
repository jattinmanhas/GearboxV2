"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { useState } from "react"
import { 
  LayoutDashboard, 
  ShoppingCart, 
  Package, 
  Users, 
  Settings,
  LogOut,
  User,
  ChevronDown,
  ChevronRight,
  Shield,
  FileText,
  BarChart3
} from "lucide-react"
import { ThemeToggle } from "@/components/ui/theme-toggle"
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
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar"

const navigationSections = [
  {
    title: "Overview",
    items: [
      { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
      { name: "Analytics", href: "/dashboard/analytics", icon: BarChart3 },
    ]
  },
  {
    title: "Authentication",
    items: [
      { name: "Users", href: "/dashboard/users", icon: Users },
      { name: "Roles", href: "/dashboard/roles", icon: Shield },
    ]
  },
  {
    title: "E-commerce",
    items: [
      { name: "Products", href: "/dashboard/products", icon: Package },
      { name: "Orders", href: "/dashboard/orders", icon: ShoppingCart },
    ]
  },
  {
    title: "Content",
    items: [
      { name: "Blog Posts", href: "/dashboard/blog", icon: FileText },
    ]
  },
  {
    title: "System",
    items: [
      { name: "Settings", href: "/dashboard/settings", icon: Settings },
    ]
  }
]

export function DashboardSidebar() {
  const pathname = usePathname()
  const [expandedSections, setExpandedSections] = useState<Set<string>>(
    new Set(["Overview"]) // Start with Overview expanded
  )

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
    <Sidebar>
      <SidebarHeader className="px-4 py-2">
        <div className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <Package className="h-4 w-4" />
          </div>
          <div className="grid flex-1 text-left leading-tight">
            <span className="truncate text-sm font-semibold">GearboxV2</span>
            <span className="truncate text-xs text-muted-foreground">Admin Dashboard</span>
          </div>
        </div>
      </SidebarHeader>
      
      <SidebarContent className="px-2">
        {navigationSections.map((section) => {
          const isExpanded = expandedSections.has(section.title)
          const hasActiveItem = section.items.some(item => pathname === item.href)
          
          return (
            <SidebarGroup key={section.title}>
              <SidebarGroupLabel 
                className="cursor-pointer hover:bg-sidebar-accent hover:text-sidebar-accent-foreground rounded-lg px-3 transition-all duration-200 text-sm font-medium"
                onClick={() => toggleSection(section.title)}
              >
                <div className="flex items-center justify-between w-full">
                  <span>{section.title}</span>
                  <div className={`transition-transform duration-200 ${isExpanded ? 'rotate-0' : '-rotate-90'}`}>
                    <ChevronDown className="h-4 w-4" />
                  </div>
                </div>
              </SidebarGroupLabel>
              
              <div className={`overflow-hidden transition-all duration-300 ease-in-out ${
                isExpanded ? 'max-h-96 opacity-100' : 'max-h-0 opacity-0'
              }`}>
                <SidebarGroupContent className="pt-2">
                  <SidebarMenu className="space-y-1">
                    {section.items.map((item) => {
                      const isActive = pathname === item.href
                      return (
                        <SidebarMenuItem key={item.name}>
                          <SidebarMenuButton 
                            asChild 
                            isActive={isActive}
                            className="h-9 text-sm font-normal ml-4 rounded-lg transition-all duration-200 hover:bg-sidebar-accent/50"
                          >
                            <Link href={item.href} className="flex items-center gap-2">
                              <item.icon className="h-4 w-4 flex-shrink-0" />
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
      
      <SidebarFooter className="px-3 py-2">
        <SidebarMenu className="space-y-1">
          <SidebarMenuItem>
            <div className="flex items-center gap-2 px-2 py-2 rounded-md bg-muted/30">
              <div className="flex h-6 w-6 items-center justify-center rounded-md bg-muted">
                <User className="h-3 w-3" />
              </div>
              <div className="grid flex-1 text-left leading-tight">
                <span className="truncate text-xs font-medium">Admin User</span>
                <span className="truncate text-xs text-muted-foreground">admin@gearbox.com</span>
              </div>
            </div>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <div className="flex items-center justify-between px-2 py-1">
              <ThemeToggle />
              <SidebarMenuButton size="sm" className="h-7 w-7 p-0 rounded-md">
                <LogOut className="h-3 w-3" />
                <span className="sr-only">Logout</span>
              </SidebarMenuButton>
            </div>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
