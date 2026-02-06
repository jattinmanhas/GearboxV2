"use client"

import { useState } from "react"
import { User, Role, ChangePasswordRequest } from "@/lib/types"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  MoreHorizontal,
  Edit,
  Trash2,
  Key,
  LogOut,
  Eye,
  EyeOff,
  ChevronLeft,
  ChevronRight
} from "lucide-react"
import { format } from "date-fns"

interface UserTableProps {
  users: User[]
  roles: Role[]
  loading: boolean
  pagination: {
    total: number
    page: number
    limit: number
    total_pages: number
  }
  onEdit: (user: User) => void
  onDelete: (id: number) => void
  onChangePassword: (id: number, passwordData: ChangePasswordRequest) => void
  onLogoutAll: (id: number) => void
  onPageChange: (page: number) => void
}

export function UserTable({
  users,
  roles,
  loading,
  pagination,
  onEdit,
  onDelete,
  onChangePassword,
  onLogoutAll,
  onPageChange,
}: UserTableProps) {
  const [changePasswordUser, setChangePasswordUser] = useState<User | null>(null)
  const [passwordData, setPasswordData] = useState({
    current_password: "",
    new_password: "",
    confirm_password: "",
  })
  const [showPassword, setShowPassword] = useState(false)

  const handleChangePassword = () => {
    if (!changePasswordUser) return

    if (passwordData.new_password !== passwordData.confirm_password) {
      toast.error("New passwords do not match")
      return
    }

    onChangePassword(changePasswordUser.id, passwordData)
    setChangePasswordUser(null)
    setPasswordData({
      current_password: "",
      new_password: "",
      confirm_password: "",
    })
  }

  const getRoleName = (user: User) => {
    // First try to get role from the role field (string)
    if (user.role) return user.role

    // Fallback to role_id lookup
    if (user.role_id && Array.isArray(roles)) {
      const role = roles.find(r => r.id === user.role_id)
      return role?.name || "Unknown Role"
    }

    return "No Role"
  }

  const getRoleColor = (user: User) => {
    const roleName = getRoleName(user).toLowerCase()

    switch (roleName) {
      case "admin":
        return "destructive"
      case "editor":
        return "default"
      case "moderator":
        return "secondary"
      case "viewer":
        return "outline"
      case "user":
        return "secondary"
      default:
        return "secondary"
    }
  }

  const getInitials = (user: User) => {
    const firstChar = (val: any) => (typeof val === 'string' && val.length > 0) ? val.charAt(0) : ""
    const first = (firstChar(user.first_name) || firstChar(user.username) || "U").toUpperCase()
    const last = firstChar(user.last_name).toUpperCase()
    return `${first}${last}`
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-32">
        <div className="text-muted-foreground">Loading users...</div>
      </div>
    )
  }

  return (
    <>
      <div className="rounded-md border overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>User</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Role</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="hidden md:table-cell">Last Login</TableHead>
              <TableHead className="hidden lg:table-cell">Created</TableHead>
              <TableHead className="w-[70px]">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {users.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                  No users found
                </TableCell>
              </TableRow>
            ) : (
              users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>
                    <div className="flex items-center space-x-3">
                      <Avatar className="h-8 w-8">
                        <AvatarImage src={user.avatar} />
                        <AvatarFallback>{getInitials(user)}</AvatarFallback>
                      </Avatar>
                      <div>
                        <div className="font-medium">
                          {user.first_name} {user.last_name}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          @{user.username}
                        </div>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell className="max-w-[200px] truncate" title={user.email}>{user.email}</TableCell>
                  <TableCell>
                    <Badge variant={getRoleColor(user)}>
                      {getRoleName(user)}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge variant={user.is_active ? "default" : "secondary"}>
                      {user.is_active ? "Active" : "Inactive"}
                    </Badge>
                  </TableCell>
                  <TableCell className="hidden md:table-cell">
                    {user.last_login && user.last_login !== "0001-01-01T00:00:00Z" ? (
                      <div className="text-sm">
                        {format(new Date(user.last_login), "MMM d, yyyy")}
                        <div className="text-xs text-muted-foreground">
                          {format(new Date(user.last_login), "h:mm a")}
                        </div>
                      </div>
                    ) : (
                      <span className="text-muted-foreground">Never</span>
                    )}
                  </TableCell>
                  <TableCell className="hidden lg:table-cell">
                    <div className="text-sm">
                      {format(new Date(user.created_at), "MMM d, yyyy")}
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
                        <DropdownMenuItem onClick={() => onEdit(user)}>
                          <Edit className="mr-2 h-4 w-4" />
                          Edit User
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => setChangePasswordUser(user)}>
                          <Key className="mr-2 h-4 w-4" />
                          Change Password
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => onLogoutAll(user.id)}>
                          <LogOut className="mr-2 h-4 w-4" />
                          Logout All Devices
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          onClick={() => onDelete(user.id)}
                          className="text-destructive"
                        >
                          <Trash2 className="mr-2 h-4 w-4" />
                          Delete User
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

      {/* Pagination */}
      {pagination.total_pages > 1 && (
        <div className="flex items-center justify-between px-2 py-4">
          <div className="text-sm text-muted-foreground">
            Showing {((pagination.page - 1) * pagination.limit) + 1} to{" "}
            {Math.min(pagination.page * pagination.limit, pagination.total)} of{" "}
            {pagination.total} results
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => onPageChange(pagination.page - 1)}
              disabled={pagination.page <= 1}
            >
              <ChevronLeft className="h-4 w-4" />
              Previous
            </Button>
            <div className="flex items-center space-x-1">
              {Array.from({ length: pagination.total_pages }, (_, i) => i + 1).map((page) => (
                <Button
                  key={page}
                  variant={page === pagination.page ? "default" : "outline"}
                  size="sm"
                  onClick={() => onPageChange(page)}
                  className="w-8 h-8 p-0"
                >
                  {page}
                </Button>
              ))}
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => onPageChange(pagination.page + 1)}
              disabled={pagination.page >= pagination.total_pages}
            >
              Next
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
        </div>
      )}

      {/* Change Password Dialog */}
      <Dialog open={!!changePasswordUser} onOpenChange={() => setChangePasswordUser(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Change Password</DialogTitle>
            <DialogDescription>
              Change password for {changePasswordUser?.first_name} {changePasswordUser?.last_name}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label htmlFor="current_password">Current Password</Label>
              <div className="relative">
                <Input
                  id="current_password"
                  type={showPassword ? "text" : "password"}
                  value={passwordData.current_password}
                  onChange={(e) => setPasswordData({ ...passwordData, current_password: e.target.value })}
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? (
                    <EyeOff className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>
            <div>
              <Label htmlFor="new_password">New Password</Label>
              <Input
                id="new_password"
                type={showPassword ? "text" : "password"}
                value={passwordData.new_password}
                onChange={(e) => setPasswordData({ ...passwordData, new_password: e.target.value })}
              />
            </div>
            <div>
              <Label htmlFor="confirm_password">Confirm New Password</Label>
              <Input
                id="confirm_password"
                type={showPassword ? "text" : "password"}
                value={passwordData.confirm_password}
                onChange={(e) => setPasswordData({ ...passwordData, confirm_password: e.target.value })}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setChangePasswordUser(null)}>
              Cancel
            </Button>
            <Button onClick={handleChangePassword}>
              Change Password
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
