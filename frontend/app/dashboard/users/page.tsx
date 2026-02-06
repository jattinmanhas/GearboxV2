"use client"

import { useState, useEffect, useCallback } from "react"
import { userApi, roleApi } from "@/lib/apiFunctions"
import { User, Role, UserFilters, UpdateUserRequest, ChangePasswordRequest } from "@/lib/types"
import { UserTable } from "./components/user-table"
import { UserForm } from "./components/user-form"
import { UserFilters as UserFiltersComponent } from "./components/user-filters"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Plus, Users, UserCheck, UserX, Eraser } from "lucide-react"
import { toast } from "sonner"

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(true)
  const [editingUser, setEditingUser] = useState<User | null>(null)
  const [filters, setFilters] = useState<UserFilters>({
    page: 1,
    limit: 10,
  })
  const [pagination, setPagination] = useState({
    total: 0,
    page: 1,
    limit: 10,
    total_pages: 0,
  })

  // Load users and roles
  const loadUsers = useCallback(async () => {
    try {
      setLoading(true)
      const response = await userApi.getUsers(filters)

      // Ensure we have valid data
      const usersData = response.data?.users || []
      setUsers(usersData)
      setPagination({
        total: response.data?.total || 0,
        page: response.data?.page || 1,
        limit: response.data?.limit || 10,
        total_pages: response.data?.total_pages || 0,
      })
    } catch (err: unknown) {
      const error = err as Error
      console.error('Error loading users:', error)
      // The improved error handling in api.ts will handle 401s appropriately
      toast.error(error.message || "Failed to load users")
      setUsers([]) // Set empty array on error
    } finally {
      setLoading(false)
    }
  }, [filters])

  const loadRoles = useCallback(async () => {
    try {
      const response = await roleApi.getRoles()

      // Ensure we have an array
      const rolesData = Array.isArray(response) ? response : []
      setRoles(rolesData)
    } catch (err: unknown) {
      const error = err as Error
      console.error("Failed to load roles:", error)
      setRoles([]) // Set empty array on error
    }
  }, [])

  useEffect(() => {
    loadUsers()
    loadRoles()
  }, [loadUsers, loadRoles])


  const handleUpdateUser = async (id: number, userData: UpdateUserRequest) => {
    try {
      await userApi.updateUser(id, userData)
      toast.success("User updated successfully")
      setEditingUser(null)
      loadUsers()
    } catch (err: unknown) {
      const error = err as Error
      toast.error(error.message || "Failed to update user")
    }
  }

  const handleDeleteUser = async (id: number) => {
    if (!confirm("Are you sure you want to delete this user?")) return

    try {
      await userApi.deleteUser(id)
      toast.success("User deleted successfully")
      loadUsers()
    } catch (err: unknown) {
      const error = err as Error
      toast.error(error.message || "Failed to delete user")
    }
  }

  const handleChangePassword = async (id: number, passwordData: ChangePasswordRequest) => {
    try {
      await userApi.changePassword(id, passwordData)
      toast.success("Password changed successfully")
    } catch (err: unknown) {
      const error = err as Error
      toast.error(error.message || "Failed to change password")
    }
  }

  const handleLogoutAll = async (id: number) => {
    if (!confirm("Are you sure you want to logout this user from all devices?")) return

    try {
      await userApi.logoutAll(id)
      toast.success("User logged out from all devices")
    } catch (err: unknown) {
      const error = err as Error
      toast.error(error.message || "Failed to logout user")
    }
  }

  const handleCleanupTokens = async () => {
    const promise = userApi.cleanupExpiredTokens()
    toast.promise(promise, {
      loading: 'Cleaning up expired tokens...',
      success: (data) => {
        return data.message || 'Tokens cleaned up successfully'
      },
      error: 'Failed to cleanup tokens'
    })
  }

  const handleFiltersChange = (newFilters: UserFilters) => {
    setFilters({ ...filters, ...newFilters, page: 1 })
  }

  const handlePageChange = (page: number) => {
    setFilters({ ...filters, page })
  }

  const activeUsers = Array.isArray(users) ? users.filter(user => user.is_active).length : 0
  const inactiveUsers = Array.isArray(users) ? users.filter(user => !user.is_active).length : 0

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">User Management</h1>
          <p className="text-muted-foreground">
            Manage users, roles, and permissions
          </p>
        </div>
        <Button variant="outline" onClick={handleCleanupTokens}>
          <Eraser className="h-4 w-4 mr-2" />
          Cleanup Tokens
        </Button>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Users</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{pagination.total}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Users</CardTitle>
            <UserCheck className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{activeUsers}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Inactive Users</CardTitle>
            <UserX className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{inactiveUsers}</div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardHeader>
          <CardTitle>Filters</CardTitle>
          <CardDescription>Filter users by various criteria</CardDescription>
        </CardHeader>
        <CardContent>
          <UserFiltersComponent
            filters={filters}
            roles={roles}
            onFiltersChange={handleFiltersChange}
          />
        </CardContent>
      </Card>


      {/* User Table */}
      <Card>
        <CardHeader>
          <CardTitle>Users</CardTitle>
          <CardDescription>
            Manage user accounts and their permissions
          </CardDescription>
        </CardHeader>
        <CardContent>
          <UserTable
            users={users}
            roles={roles}
            loading={loading}
            pagination={pagination}
            onEdit={setEditingUser}
            onDelete={handleDeleteUser}
            onChangePassword={handleChangePassword}
            onLogoutAll={handleLogoutAll}
            onPageChange={handlePageChange}
          />
        </CardContent>
      </Card>

      {/* Edit User Modal */}
      {editingUser && (
        <UserForm
          user={editingUser}
          roles={roles}
          onSubmit={(userData: UpdateUserRequest) => handleUpdateUser(editingUser.id, userData)}
          onClose={() => setEditingUser(null)}
        />
      )}
    </div>
  )
}
