"use client"

import { useState, useEffect } from "react"
import { type UserFilters, Role } from "@/lib/types"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Card, CardContent } from "@/components/ui/card"
import { Search, Filter, X } from "lucide-react"

interface UserFiltersProps {
  filters: UserFilters
  roles: Role[]
  onFiltersChange: (filters: UserFilters) => void
}

export function UserFilters({ filters, roles, onFiltersChange }: UserFiltersProps) {
  const [localFilters, setLocalFilters] = useState<UserFilters>(filters)

  // Sync local filters with parent filters when they change
  useEffect(() => {
    setLocalFilters(filters)
  }, [filters])

  const handleFilterChange = (field: keyof UserFilters, value: string | number | boolean | undefined) => {
    const newFilters = { ...localFilters, [field]: value }
    setLocalFilters(newFilters)
  }

  const handleApplyFilters = () => {
    console.log('Applying filters:', localFilters)
    onFiltersChange(localFilters)
  }

  const handleClearFilters = () => {
    const clearedFilters: UserFilters = {
      page: 1,
      limit: 10,
    }
    setLocalFilters(clearedFilters)
    onFiltersChange(clearedFilters)
  }

  const hasActiveFilters = Object.keys(filters).some(
    key => key !== 'page' && key !== 'limit' && filters[key as keyof UserFilters]
  )

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Search */}
        <div className="space-y-2">
          <Label htmlFor="search">Search</Label>
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              id="search"
              placeholder="Search users..."
              value={localFilters.search || ""}
              onChange={(e) => handleFilterChange("search", e.target.value)}
              className="pl-10"
            />
          </div>
        </div>

        {/* Status Filter */}
        <div className="space-y-2">
          <Label htmlFor="status">Status</Label>
          <Select
            value={localFilters.is_active?.toString() || "all"}
            onValueChange={(value) => 
              handleFilterChange("is_active", value === "all" ? undefined : value === "true")
            }
          >
            <SelectTrigger>
              <SelectValue placeholder="All statuses" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All statuses</SelectItem>
              <SelectItem value="true">Active</SelectItem>
              <SelectItem value="false">Inactive</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Role Filter */}
        <div className="space-y-2">
          <Label htmlFor="role">Role</Label>
          <Select
            value={localFilters.role_id?.toString() || "all"}
            onValueChange={(value) => 
              handleFilterChange("role_id", value === "all" ? undefined : parseInt(value))
            }
          >
            <SelectTrigger>
              <SelectValue placeholder="All roles" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All roles</SelectItem>
              {Array.isArray(roles) && roles.map((role) => (
                <SelectItem key={role.id} value={role.id.toString()}>
                  {role.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Results per page */}
      <div className="flex items-center space-x-4">
        <div className="flex items-center space-x-2">
          <Label htmlFor="limit">Results per page:</Label>
          <Select
            value={localFilters.limit?.toString() || "10"}
            onValueChange={(value) => handleFilterChange("limit", parseInt(value))}
          >
            <SelectTrigger className="w-20">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="5">5</SelectItem>
              <SelectItem value="10">10</SelectItem>
              <SelectItem value="25">25</SelectItem>
              <SelectItem value="50">50</SelectItem>
              <SelectItem value="100">100</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Action Buttons */}
      <div className="flex items-center space-x-2">
        <Button onClick={handleApplyFilters} size="sm">
          <Filter className="h-4 w-4 mr-2" />
          Apply Filters
        </Button>
        {hasActiveFilters && (
          <Button onClick={handleClearFilters} variant="outline" size="sm">
            <X className="h-4 w-4 mr-2" />
            Clear Filters
          </Button>
        )}
      </div>

      {/* Active Filters Display */}
      {hasActiveFilters && (
        <Card>
          <CardContent className="pt-4">
            <div className="flex flex-wrap gap-2">
              <span className="text-sm font-medium text-muted-foreground">Active filters:</span>
              {filters.search && (
                <div className="flex items-center space-x-1 bg-primary/10 text-primary px-2 py-1 rounded-md text-sm">
                  <span>Search: &quot;{filters.search}&quot;</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-4 w-4 p-0 hover:bg-primary/20"
                    onClick={() => {
                      const newFilters = { ...localFilters, search: undefined }
                      setLocalFilters(newFilters)
                      onFiltersChange(newFilters)
                    }}
                  >
                    <X className="h-3 w-3" />
                  </Button>
                </div>
              )}
              {filters.is_active !== undefined && (
                <div className="flex items-center space-x-1 bg-primary/10 text-primary px-2 py-1 rounded-md text-sm">
                  <span>Status: {filters.is_active ? "Active" : "Inactive"}</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-4 w-4 p-0 hover:bg-primary/20"
                    onClick={() => {
                      const newFilters = { ...localFilters, is_active: undefined }
                      setLocalFilters(newFilters)
                      onFiltersChange(newFilters)
                    }}
                  >
                    <X className="h-3 w-3" />
                  </Button>
                </div>
              )}
              {filters.role_id && (
                <div className="flex items-center space-x-1 bg-primary/10 text-primary px-2 py-1 rounded-md text-sm">
                  <span>Role: {roles.find(r => r.id === filters.role_id)?.name}</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-4 w-4 p-0 hover:bg-primary/20"
                    onClick={() => {
                      const newFilters = { ...localFilters, role_id: undefined }
                      setLocalFilters(newFilters)
                      onFiltersChange(newFilters)
                    }}
                  >
                    <X className="h-3 w-3" />
                  </Button>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
