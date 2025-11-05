"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { 
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { UserProfile, UpdateProfileRequest } from "@/lib/types"
import { showError, showSuccess, NotificationMessages } from "@/lib/notifications"
import { EnhancedImageUpload } from "@/components/ui/enhanced-image-upload"
import { UploadedImage } from "@/lib/image-upload"
import { X } from "lucide-react"

interface ProfileEditFormProps {
  profile: UserProfile
  onSave: (data: UpdateProfileRequest) => void
  onCancel: () => void
}

export function ProfileEditForm({ profile, onSave, onCancel }: ProfileEditFormProps) {
  const [formData, setFormData] = useState({
    first_name: "",
    middle_name: "",
    last_name: "",
    phone_number: "",
    date_of_birth: "",
    avatar: "",
  })

  const [errors, setErrors] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [uploadedAvatar, setUploadedAvatar] = useState<UploadedImage | null>(null)
  const [avatarUrl, setAvatarUrl] = useState("")

  // Convert RFC3339 date (2006-01-02T15:04:05Z07:00) to YYYY-MM-DD for date input
  const formatDateForInput = (dateString: string): string => {
    if (!dateString || dateString === '0001-01-01T00:00:00Z' || dateString === '') {
      return ''
    }
    try {
      const date = new Date(dateString)
      if (isNaN(date.getTime())) {
        return ''
      }
      // Format as YYYY-MM-DD for date input
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      return `${year}-${month}-${day}`
    } catch {
      return ''
    }
  }

  useEffect(() => {
    setFormData({
      first_name: profile.first_name,
      middle_name: profile.middle_name || "",
      last_name: profile.last_name,
      phone_number: profile.phone_number || "",
      date_of_birth: formatDateForInput(profile.date_of_birth || ""),
      avatar: profile.avatar || "",
    })
    setAvatarUrl(profile.avatar || "")
    setUploadedAvatar(null)
  }, [profile])

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.first_name.trim()) {
      newErrors.first_name = "First name is required"
    }

    if (!formData.last_name.trim()) {
      newErrors.last_name = "Last name is required"
    }

    if (formData.phone_number && !/^\+?[\d\s\-\(\)]+$/.test(formData.phone_number)) {
      newErrors.phone_number = "Please enter a valid phone number"
    }

    if (formData.date_of_birth) {
      const date = new Date(formData.date_of_birth)
      const today = new Date()
      if (date >= today) {
        newErrors.date_of_birth = "Date of birth must be in the past"
      }
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  // Convert YYYY-MM-DD to RFC3339 format (2006-01-02T15:04:05Z07:00)
  const formatDateForBackend = (dateString: string): string | undefined => {
    if (!dateString || dateString.trim() === '') {
      return undefined
    }
    try {
      // Create date from YYYY-MM-DD format (local date, no time)
      const date = new Date(dateString + 'T00:00:00')
      if (isNaN(date.getTime())) {
        return undefined
      }
      // Format as RFC3339 with timezone
      return date.toISOString()
    } catch {
      return undefined
    }
  }

  const handleAvatarUpload = (image: UploadedImage) => {
    setUploadedAvatar(image)
    // Use secure URL if available, otherwise use regular URL
    const imageUrl = image.secureUrl || image.url
    setAvatarUrl(imageUrl)
    setFormData(prev => ({ ...prev, avatar: imageUrl }))
  }

  const handleAvatarRemove = () => {
    setUploadedAvatar(null)
    setAvatarUrl("")
    setFormData(prev => ({ ...prev, avatar: "" }))
  }

  const handleAvatarUrlChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const url = e.target.value
    setAvatarUrl(url)
    setFormData(prev => ({ ...prev, avatar: url }))
    // Clear uploaded avatar if URL is manually entered
    if (url && uploadedAvatar) {
      setUploadedAvatar(null)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) return

    setIsSubmitting(true)
    try {
      // Convert date format before sending
      // Use uploaded avatar URL if available, otherwise use manually entered URL
      const finalAvatarUrl = uploadedAvatar ? (uploadedAvatar.secureUrl || uploadedAvatar.url) : avatarUrl
      
      const profileData: UpdateProfileRequest = {
        ...formData,
        avatar: finalAvatarUrl,
        date_of_birth: formatDateForBackend(formData.date_of_birth),
      }
      await onSave(profileData)
      // Success notification is shown in the parent component (handleProfileUpdate)
    } catch (error) {
      console.error("Error saving profile:", error)
      const errorMessage = error instanceof Error ? error.message : "Failed to save profile"
      showError(NotificationMessages.user.profileUpdateError, {
        description: errorMessage
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Dialog open={true} onOpenChange={onCancel}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Edit Profile</DialogTitle>
          <DialogDescription>
            Update your personal information below.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* First Name */}
            <div className="space-y-2">
              <Label htmlFor="first_name">First Name *</Label>
              <Input
                id="first_name"
                value={formData.first_name}
                onChange={(e) => setFormData(prev => ({ ...prev, first_name: e.target.value }))}
                className={errors.first_name ? "border-destructive" : ""}
              />
              {errors.first_name && (
                <p className="text-sm text-destructive">{errors.first_name}</p>
              )}
            </div>

            {/* Last Name */}
            <div className="space-y-2">
              <Label htmlFor="last_name">Last Name *</Label>
              <Input
                id="last_name"
                value={formData.last_name}
                onChange={(e) => setFormData(prev => ({ ...prev, last_name: e.target.value }))}
                className={errors.last_name ? "border-destructive" : ""}
              />
              {errors.last_name && (
                <p className="text-sm text-destructive">{errors.last_name}</p>
              )}
            </div>
          </div>

          {/* Middle Name */}
          <div className="space-y-2">
            <Label htmlFor="middle_name">Middle Name</Label>
            <Input
              id="middle_name"
              value={formData.middle_name}
              onChange={(e) => setFormData(prev => ({ ...prev, middle_name: e.target.value }))}
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Phone Number */}
            <div className="space-y-2">
              <Label htmlFor="phone_number">Phone Number</Label>
              <Input
                id="phone_number"
                value={formData.phone_number}
                onChange={(e) => setFormData(prev => ({ ...prev, phone_number: e.target.value }))}
                placeholder="+1 (555) 123-4567"
                className={errors.phone_number ? "border-destructive" : ""}
              />
              {errors.phone_number && (
                <p className="text-sm text-destructive">{errors.phone_number}</p>
              )}
            </div>

            {/* Date of Birth */}
            <div className="space-y-2">
              <Label htmlFor="date_of_birth">Date of Birth</Label>
              <Input
                id="date_of_birth"
                type="date"
                value={formData.date_of_birth}
                onChange={(e) => setFormData(prev => ({ ...prev, date_of_birth: e.target.value }))}
                className={errors.date_of_birth ? "border-destructive" : ""}
              />
              {errors.date_of_birth && (
                <p className="text-sm text-destructive">{errors.date_of_birth}</p>
              )}
            </div>
          </div>

          {/* Avatar Upload */}
          <div className="space-y-4">
            <Label>Profile Picture</Label>
            
            {/* Current Avatar Preview */}
            {(avatarUrl || uploadedAvatar) && (
              <div className="relative inline-block">
                <Avatar className="h-24 w-24 border-2">
                  <AvatarImage 
                    src={uploadedAvatar ? (uploadedAvatar.secureUrl || uploadedAvatar.url) : avatarUrl} 
                    alt="Profile avatar" 
                  />
                  <AvatarFallback className="text-lg">
                    {(formData.first_name?.[0] || '').toUpperCase()}{(formData.last_name?.[0] || '').toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                {uploadedAvatar && (
                  <Button
                    type="button"
                    variant="destructive"
                    size="sm"
                    className="absolute -top-2 -right-2 h-6 w-6 p-0 rounded-full"
                    onClick={handleAvatarRemove}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                )}
              </div>
            )}

            {/* Image Upload Component */}
            {!uploadedAvatar && (
              <div className="border rounded-lg p-4">
                <EnhancedImageUpload
                  onImageSelect={handleAvatarUpload}
                  selectedImages={uploadedAvatar ? [uploadedAvatar] : []}
                  multiple={false}
                  maxImages={1}
                  showPreview={false}
                  showThumbnails={false}
                  config={{
                    maxFileSize: 5 * 1024 * 1024, // 5MB
                    allowedTypes: ['image/jpeg', 'image/png', 'image/webp', 'image/gif'],
                    maxWidth: 512,
                    maxHeight: 512,
                    quality: 85,
                    generateThumbnails: false,
                    thumbnailSizes: [],
                    useCloudinary: true
                  }}
                />
              </div>
            )}

            {/* Manual URL Input (Alternative) */}
            <div className="space-y-2">
              <Label htmlFor="avatar-url">Or enter avatar URL manually</Label>
              <Input
                id="avatar-url"
                type="url"
                value={avatarUrl}
                onChange={handleAvatarUrlChange}
                placeholder="https://example.com/avatar.jpg"
                disabled={!!uploadedAvatar}
              />
              <p className="text-xs text-muted-foreground">
                {uploadedAvatar 
                  ? "Remove the uploaded image to enter a URL manually"
                  : "Upload an image above or enter a direct URL to an image"}
              </p>
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Saving..." : "Save Changes"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
