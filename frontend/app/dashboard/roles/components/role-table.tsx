"use client"

import { useState } from "react"
import { Role, AssignRoleRequest, RemoveRoleRequest } from "@/lib/types"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { 
  MoreHorizontal, 
  Edit, 
  Trash2, 
  Shield
} from "lucide-react"
import { format } from "date-fns"

interface RoleTableProps {
  roles: Role[]
  loading: boolean
  onAssignRole: (assignmentData: AssignRoleRequest) => void
  onRemoveRole: (removalData: RemoveRoleRequest) => void
}

export function RoleTable({
  roles,
  loading,
  onAssignRole,
  onRemoveRole,
}: RoleTableProps) {

  if (loading) {
    return (
      <div className="flex items-center justify-center h-32">
        <div className="text-muted-foreground">Loading roles...</div>
      </div>
    )
  }

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
            <TableRow>
              <TableHead>Role</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Created</TableHead>
              <TableHead className="w-[70px]">Actions</TableHead>
            </TableRow>
        </TableHeader>
        <TableBody>
          {roles.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                No roles found
              </TableCell>
            </TableRow>
          ) : (
            Array.isArray(roles) && roles.map((role) => (
              <TableRow key={role.id}>
                <TableCell>
                  <div className="flex items-center space-x-2">
                    <Shield className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <div className="font-medium">{role.name}</div>
                      <div className="text-sm text-muted-foreground">
                        Role ID: {role.id}
                      </div>
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  <div className="max-w-xs truncate">
                    {role.description || "No description"}
                  </div>
                </TableCell>
                <TableCell>
                  <Badge variant={role.is_active ? "default" : "secondary"}>
                    {role.is_active ? "Active" : "Inactive"}
                  </Badge>
                </TableCell>
                <TableCell>
                  <div className="text-sm">
                    {format(new Date(role.created_at), "MMM d, yyyy")}
                  </div>
                </TableCell>
                <TableCell>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">Open menu</span>
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuLabel>Actions</DropdownMenuLabel>
                        <DropdownMenuItem onClick={() => {/* TODO: Implement role assignment */}}>
                          <Edit className="mr-2 h-4 w-4" />
                          Assign to User
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => {/* TODO: Implement role removal */}}>
                          <Trash2 className="mr-2 h-4 w-4" />
                          Remove from User
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>

    </div>
  )
}
