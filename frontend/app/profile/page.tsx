"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Separator } from "@/components/ui/separator"
import {
  User,
  Mail,
  Phone,
  Calendar,
  MapPin,
  Edit,
  Plus,
} from "lucide-react"
import { profileApi } from "@/lib/apiFunctions"
import { UserProfile, Address } from "@/lib/types"
import { ProfileEditForm } from "./components/profile-edit-form"
import { AddressForm } from "./components/address-form"
import { AddressCard } from "./components/address-card"
import { useUserStore } from "@/lib/stores/user-store"
import { showSuccess, showError, NotificationMessages } from "@/lib/notifications"

export default function ProfilePage() {
  const [profile, setProfile] = useState<UserProfile | null>(null)
  const [addresses, setAddresses] = useState<Address[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showProfileEdit, setShowProfileEdit] = useState(false)
  const [showAddressForm, setShowAddressForm] = useState(false)
  const [editingAddress, setEditingAddress] = useState<Address | null>(null)

  const { updateProfile } = useUserStore()

  const loadProfile = async () => {
    try {
      setLoading(true)
      setError(null)
      const [profileData, addressesData] = await Promise.all([
        profileApi.getProfile(),
        profileApi.getAddresses()
      ])
      setProfile(profileData)
      setAddresses(Array.isArray(addressesData) ? addressesData : [])
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to load profile"
      setError(errorMessage)
      showError("Failed to load profile", {
        description: errorMessage
      })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadProfile()
  }, [])

  const handleProfileUpdate = async (profileData: any) => {
    try {
      const updatedProfile = await profileApi.updateProfile(profileData)
      setProfile(updatedProfile)

      // Update user store with new profile data
      updateProfile({
        firstName: updatedProfile.first_name,
        middleName: updatedProfile.middle_name,
        lastName: updatedProfile.last_name,
        avatar: updatedProfile.avatar,
      })

      showSuccess(NotificationMessages.user.profileUpdated)
      setShowProfileEdit(false)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to update profile"
      setError(errorMessage)
      showError(NotificationMessages.user.profileUpdateError, {
        description: errorMessage
      })
    }
  }

  const handleAddressCreate = async (addressData: any) => {
    try {
      const newAddress = await profileApi.createAddress(addressData)
      setAddresses(prev => [...prev, newAddress])
      showSuccess("Address created successfully")
      setShowAddressForm(false)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to create address"
      setError(errorMessage)
      showError("Failed to create address", {
        description: errorMessage
      })
    }
  }

  const handleAddressUpdate = async (id: number, addressData: any) => {
    try {
      const updatedAddress = await profileApi.updateAddress(id, addressData)
      setAddresses(prev => prev.map(addr => addr.id === id ? updatedAddress : addr))
      showSuccess("Address updated successfully")
      setEditingAddress(null)
      setShowAddressForm(false)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to update address"
      setError(errorMessage)
      showError("Failed to update address", {
        description: errorMessage
      })
    }
  }

  const handleAddressDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this address?")) return

    try {
      await profileApi.deleteAddress(id)
      setAddresses(prev => prev.filter(addr => addr.id !== id))
      showSuccess("Address deleted successfully")
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to delete address"
      setError(errorMessage)
      showError("Failed to delete address", {
        description: errorMessage
      })
    }
  }

  const handleEditAddress = (address: Address) => {
    setEditingAddress(address)
    setShowAddressForm(true)
  }

  const handleAddAddress = () => {
    setEditingAddress(null)
    setShowAddressForm(true)
  }

  const formatDate = (dateString: string) => {
    if (!dateString || dateString === '0001-01-01T00:00:00Z' || dateString === '') {
      return 'Not provided'
    }
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-center h-64">
            <div className="text-muted-foreground">Loading profile...</div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !profile) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-6">
          <div className="text-center py-8">
            <div className="text-destructive mb-4">
              {error || "Failed to load profile"}
            </div>
            <Button onClick={loadProfile}>Try Again</Button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto px-4 py-6">
        <div className="max-w-4xl mx-auto space-y-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">Profile</h1>
              <p className="text-muted-foreground">
                Manage your account information and addresses
              </p>
            </div>
          </div>

          {/* Profile Information */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Profile Information</CardTitle>
                  <CardDescription>
                    Your personal account details
                  </CardDescription>
                </div>
                <Button onClick={() => setShowProfileEdit(true)}>
                  <Edit className="h-4 w-4 mr-2" />
                  Edit Profile
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-start space-x-6">
                <Avatar className="h-20 w-20">
                  <AvatarImage src={profile.avatar || ''} alt={`${profile.first_name || ''} ${profile.last_name || ''}`} />
                  <AvatarFallback className="text-lg">
                    {(profile.first_name?.[0] || '').toUpperCase()}{(profile.last_name?.[0] || '').toUpperCase()}
                  </AvatarFallback>
                </Avatar>

                <div className="flex-1 space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="flex items-center space-x-3">
                      <User className="h-4 w-4 text-muted-foreground" />
                      <div>
                        <p className="text-sm font-medium">Full Name</p>
                        <p className="text-sm text-muted-foreground">
                          {profile.first_name || 'N/A'} {profile.middle_name && `${profile.middle_name} `}{profile.last_name || 'N/A'}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center space-x-3">
                      <Mail className="h-4 w-4 text-muted-foreground" />
                      <div>
                        <p className="text-sm font-medium">Email</p>
                        <p className="text-sm text-muted-foreground">{profile.email}</p>
                      </div>
                    </div>

                    <div className="flex items-center space-x-3">
                      <User className="h-4 w-4 text-muted-foreground" />
                      <div>
                        <p className="text-sm font-medium">Username</p>
                        <p className="text-sm text-muted-foreground">{profile.username}</p>
                      </div>
                    </div>

                    {profile.phone_number && (
                      <div className="flex items-center space-x-3">
                        <Phone className="h-4 w-4 text-muted-foreground" />
                        <div>
                          <p className="text-sm font-medium">Phone</p>
                          <p className="text-sm text-muted-foreground">{profile.phone_number}</p>
                        </div>
                      </div>
                    )}

                    {profile.date_of_birth && (
                      <div className="flex items-center space-x-3">
                        <Calendar className="h-4 w-4 text-muted-foreground" />
                        <div>
                          <p className="text-sm font-medium">Date of Birth</p>
                          <p className="text-sm text-muted-foreground">
                            {formatDate(profile.date_of_birth)}
                          </p>
                        </div>
                      </div>
                    )}
                  </div>

                  <Separator />

                  <div className="text-sm text-muted-foreground">
                    <p>Member since {formatDate(profile.created_at)}</p>
                    <p>Last updated {formatDate(profile.updated_at)}</p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Addresses */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="flex items-center gap-2">
                    <MapPin className="h-5 w-5" />
                    Addresses
                  </CardTitle>
                  <CardDescription>
                    Manage your addresses
                  </CardDescription>
                </div>
                <Button onClick={handleAddAddress}>
                  <Plus className="h-4 w-4 mr-2" />
                  Add Address
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {(addresses || []).length === 0 ? (
                <div className="text-center py-8">
                  <MapPin className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
                  <h3 className="text-lg font-semibold mb-2">No addresses yet</h3>
                  <p className="text-muted-foreground mb-4">
                    Add your first address to get started
                  </p>
                  <Button onClick={handleAddAddress}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Address
                  </Button>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {(addresses || []).map((address) => (
                    <AddressCard
                      key={address.id}
                      address={address}
                      onEdit={() => handleEditAddress(address)}
                      onDelete={() => handleAddressDelete(address.id)}
                    />
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Profile Edit Modal */}
      {showProfileEdit && (
        <ProfileEditForm
          profile={profile}
          onSave={handleProfileUpdate}
          onCancel={() => setShowProfileEdit(false)}
        />
      )}

      {/* Address Form Modal */}
      {showAddressForm && (
        <AddressForm
          address={editingAddress}
          onSave={editingAddress ?
            (data) => handleAddressUpdate(editingAddress.id, data) :
            handleAddressCreate
          }
          onCancel={() => {
            setShowAddressForm(false)
            setEditingAddress(null)
          }}
        />
      )}
    </div>
  )
}
