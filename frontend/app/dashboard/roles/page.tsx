"use client"

import { useState, useEffect, useCallback } from "react"
import { roleApi } from "@/lib/api"
import { Role, AssignRoleRequest, RemoveRoleRequest } from "@/lib/types"
import { RoleTable } from "./components/role-table"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Plus, Shield, ShieldCheck, ShieldX } from "lucide-react"
import { AlertMessage } from "@/components/ui/alert-message"

export default function RolesPage() {
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  const [editingRole, setEditingRole] = useState<Role | null>(null)

  // Load roles
  const loadRoles = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await roleApi.getRoles()
      // Ensure we have an array
      const rolesData = Array.isArray(response) ? response : []
      setRoles(rolesData)
    } catch (err: unknown) {
      const error = err as Error
      setError(error.message || "Failed to load roles")
      setRoles([]) // Set empty array on error
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadRoles()
  }, [loadRoles])

  const handleAssignRole = async (assignmentData: AssignRoleRequest) => {
    try {
      setError(null)
      await roleApi.assignRole(assignmentData)
      setSuccess("Role assigned successfully")
      loadRoles()
    } catch (err: unknown) {
      const error = err as Error
      setError(error.message || "Failed to assign role")
    }
  }

  const handleRemoveRole = async (removalData: RemoveRoleRequest) => {
    try {
      setError(null)
      await roleApi.removeRole(removalData)
      setSuccess("Role removed successfully")
      loadRoles()
    } catch (err: unknown) {
      const error = err as Error
      setError(error.message || "Failed to remove role")
    }
  }

  const activeRoles = Array.isArray(roles) ? roles.filter(role => role.is_active).length : 0
  const inactiveRoles = Array.isArray(roles) ? roles.filter(role => !role.is_active).length : 0
  const totalRoles = Array.isArray(roles) ? roles.length : 0

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Role Management</h1>
          <p className="text-muted-foreground">
            Manage roles and permissions for users
          </p>
        </div>
        <div className="text-sm text-muted-foreground">
          Roles are fixed and cannot be created or modified
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Roles</CardTitle>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{totalRoles}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Roles</CardTitle>
            <ShieldCheck className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{activeRoles}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Inactive Roles</CardTitle>
            <ShieldX className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{inactiveRoles}</div>
          </CardContent>
        </Card>
      </div>


      {/* Alerts */}
      {error && (
        <AlertMessage
          type="error"
          message={error}
        />
      )}
      {success && (
        <AlertMessage
          type="success"
          message={success}
        />
      )}

      {/* Role Table */}
      <Card>
        <CardHeader>
          <CardTitle>Roles</CardTitle>
          <CardDescription>
            Manage user roles and their permissions
          </CardDescription>
        </CardHeader>
        <CardContent>
          <RoleTable
            roles={roles}
            loading={loading}
            onAssignRole={handleAssignRole}
            onRemoveRole={handleRemoveRole}
          />
        </CardContent>
      </Card>

    </div>
  )
}
