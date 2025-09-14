"use client"

import { useState, useEffect } from "react"
import { useCouponStore } from "@/lib/stores/coupon-store"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Switch } from "@/components/ui/switch"
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { X, Save } from "lucide-react"
import { CreateCouponRequest, UpdateCouponRequest } from "@/lib/types/coupon"

interface CouponFormProps {
  couponId?: number | null
  onClose: () => void
  onSuccess: () => void
}

export function CouponForm({ couponId, onClose, onSuccess }: CouponFormProps) {
  const { 
    currentCoupon, 
    isLoading, 
    error, 
    loadCoupon, 
    createCoupon, 
    updateCoupon 
  } = useCouponStore()

  const [formData, setFormData] = useState<CreateCouponRequest>({
    code: "",
    name: "",
    description: "",
    type: "percentage",
    value: 0,
    min_order_amount: 0,
    max_discount_amount: 0,
    usage_limit: 0,
    is_active: true,
    starts_at: new Date().toISOString(),
    expires_at: "",
  })

  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (couponId) {
      loadCoupon(couponId.toString())
    }
  }, [couponId, loadCoupon])

  useEffect(() => {
    if (currentCoupon && couponId) {
      setFormData({
        code: currentCoupon.code,
        name: currentCoupon.name,
        description: currentCoupon.description || "",
        type: currentCoupon.type,
        value: currentCoupon.value,
        min_order_amount: currentCoupon.min_order_amount || 0,
        max_discount_amount: currentCoupon.max_discount_amount || 0,
        usage_limit: currentCoupon.usage_limit || 0,
        is_active: currentCoupon.is_active,
        starts_at: currentCoupon.starts_at,
        expires_at: currentCoupon.expires_at ? currentCoupon.expires_at.slice(0, 16) : "",
      })
    }
  }, [currentCoupon, couponId])

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.code.trim()) {
      newErrors.code = "Coupon code is required"
    } else if (!/^[A-Z0-9_-]+$/.test(formData.code)) {
      newErrors.code = "Coupon code must contain only uppercase letters, numbers, hyphens, and underscores"
    }

    if (!formData.name.trim()) {
      newErrors.name = "Coupon name is required"
    }

    if (formData.value <= 0) {
      newErrors.value = "Value must be greater than 0"
    }

    if (formData.type === "percentage" && formData.value > 100) {
      newErrors.value = "Percentage value cannot exceed 100"
    }

    if ((formData.min_order_amount || 0) < 0) {
      newErrors.min_order_amount = "Minimum order amount cannot be negative"
    }

    if ((formData.max_discount_amount || 0) < 0) {
      newErrors.max_discount_amount = "Maximum discount amount cannot be negative"
    }

    if ((formData.usage_limit || 0) < 0) {
      newErrors.usage_limit = "Usage limit cannot be negative"
    }

    if (!formData.starts_at) {
      newErrors.starts_at = "Start date is required"
    }

    if (formData.expires_at && formData.expires_at <= formData.starts_at) {
      newErrors.expires_at = "Expiry date must be after start date"
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) {
      return
    }

    try {
      if (couponId) {
        const updateData: UpdateCouponRequest = { ...formData }
        await updateCoupon(couponId.toString(), updateData)
      } else {
        await createCoupon(formData)
      }
      onSuccess()
    } catch (error) {
      console.error("Error saving coupon:", error)
    }
  }

  const handleInputChange = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }))
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: "" }))
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            {couponId ? "Edit Coupon" : "Create Coupon"}
          </h1>
          <p className="text-muted-foreground">
            {couponId ? "Update coupon details" : "Create a new discount coupon"}
          </p>
        </div>
        <Button variant="outline" onClick={onClose}>
          <X className="mr-2 h-4 w-4" />
          Cancel
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Coupon Details</CardTitle>
          <CardDescription>
            Fill in the details for your discount coupon
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label htmlFor="code">Coupon Code *</Label>
                <Input
                  id="code"
                  value={formData.code}
                  onChange={(e) => handleInputChange("code", e.target.value.toUpperCase())}
                  placeholder="SAVE20"
                  className={errors.code ? "border-destructive" : ""}
                />
                {errors.code && (
                  <p className="text-sm text-destructive">{errors.code}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="name">Coupon Name *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => handleInputChange("name", e.target.value)}
                  placeholder="20% Off Sale"
                  className={errors.name ? "border-destructive" : ""}
                />
                {errors.name && (
                  <p className="text-sm text-destructive">{errors.name}</p>
                )}
              </div>

              <div className="space-y-2 md:col-span-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => handleInputChange("description", e.target.value)}
                  placeholder="Get 20% off on all items"
                  rows={3}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="type">Discount Type *</Label>
                <Select
                  value={formData.type}
                  onValueChange={(value) => handleInputChange("type", value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="percentage">Percentage</SelectItem>
                    <SelectItem value="fixed_amount">Fixed Amount</SelectItem>
                    <SelectItem value="free_shipping">Free Shipping</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="value">Value *</Label>
                <Input
                  id="value"
                  type="number"
                  value={formData.value}
                  onChange={(e) => handleInputChange("value", parseFloat(e.target.value) || 0)}
                  placeholder={formData.type === "percentage" ? "20" : "100"}
                  min="0"
                  max={formData.type === "percentage" ? "100" : undefined}
                  className={errors.value ? "border-destructive" : ""}
                />
                {errors.value && (
                  <p className="text-sm text-destructive">{errors.value}</p>
                )}
                <p className="text-xs text-muted-foreground">
                  {formData.type === "percentage" 
                    ? "Enter percentage (0-100)" 
                    : formData.type === "fixed_amount" 
                    ? "Enter amount in ₹" 
                    : "Free shipping (value ignored)"
                  }
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="min_order_amount">Minimum Order Amount</Label>
                <Input
                  id="min_order_amount"
                  type="number"
                  value={formData.min_order_amount}
                  onChange={(e) => handleInputChange("min_order_amount", parseFloat(e.target.value) || 0)}
                  placeholder="500"
                  min="0"
                  className={errors.min_order_amount ? "border-destructive" : ""}
                />
                {errors.min_order_amount && (
                  <p className="text-sm text-destructive">{errors.min_order_amount}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="max_discount_amount">Maximum Discount Amount</Label>
                <Input
                  id="max_discount_amount"
                  type="number"
                  value={formData.max_discount_amount}
                  onChange={(e) => handleInputChange("max_discount_amount", parseFloat(e.target.value) || 0)}
                  placeholder="1000"
                  min="0"
                  className={errors.max_discount_amount ? "border-destructive" : ""}
                />
                {errors.max_discount_amount && (
                  <p className="text-sm text-destructive">{errors.max_discount_amount}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="usage_limit">Usage Limit</Label>
                <Input
                  id="usage_limit"
                  type="number"
                  value={formData.usage_limit}
                  onChange={(e) => handleInputChange("usage_limit", parseInt(e.target.value) || 0)}
                  placeholder="100"
                  min="0"
                  className={errors.usage_limit ? "border-destructive" : ""}
                />
                {errors.usage_limit && (
                  <p className="text-sm text-destructive">{errors.usage_limit}</p>
                )}
                <p className="text-xs text-muted-foreground">
                  Leave 0 for unlimited usage
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="starts_at">Start Date *</Label>
                <Input
                  id="starts_at"
                  type="datetime-local"
                  value={formData.starts_at ? new Date(formData.starts_at).toISOString().slice(0, 16) : ""}
                  onChange={(e) => handleInputChange("starts_at", e.target.value ? new Date(e.target.value).toISOString() : "")}
                  className={errors.starts_at ? "border-destructive" : ""}
                />
                {errors.starts_at && (
                  <p className="text-sm text-destructive">{errors.starts_at}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="expires_at">Expiry Date</Label>
                <Input
                  id="expires_at"
                  type="datetime-local"
                  value={formData.expires_at ? new Date(formData.expires_at).toISOString().slice(0, 16) : ""}
                  onChange={(e) => handleInputChange("expires_at", e.target.value ? new Date(e.target.value).toISOString() : "")}
                  className={errors.expires_at ? "border-destructive" : ""}
                />
                {errors.expires_at && (
                  <p className="text-sm text-destructive">{errors.expires_at}</p>
                )}
                <p className="text-xs text-muted-foreground">
                  Leave empty for no expiry
                </p>
              </div>
            </div>

            <div className="flex items-center space-x-2">
              <Switch
                id="is_active"
                checked={formData.is_active}
                onCheckedChange={(checked) => handleInputChange("is_active", checked)}
              />
              <Label htmlFor="is_active">Active</Label>
            </div>

            {error && (
              <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-md">
                {error}
              </div>
            )}

            <div className="flex justify-end space-x-2">
              <Button type="button" variant="outline" onClick={onClose}>
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                <Save className="mr-2 h-4 w-4" />
                {isLoading ? "Saving..." : couponId ? "Update Coupon" : "Create Coupon"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
