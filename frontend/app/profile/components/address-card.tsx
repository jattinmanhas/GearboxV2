"use client"

import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { 
  MapPin, 
  Edit, 
  Trash2, 
  Check,
  Building2,
  Phone
} from "lucide-react"
import { Address } from "@/lib/types"

interface AddressCardProps {
  address: Address
  onEdit: () => void
  onDelete: () => void
}

export function AddressCard({ address, onEdit, onDelete }: AddressCardProps) {
  const formatAddress = () => {
    const parts = [
      address.address_line_1,
      address.address_line_2,
      address.city,
      address.state,
      address.postal_code,
      address.country
    ].filter(Boolean)
    
    return parts.join(", ")
  }

  return (
    <Card className="relative">
      <CardContent className="p-4">
        <div className="flex items-start justify-between mb-3">
          <div className="flex items-center gap-2">
            <MapPin className="h-4 w-4 text-muted-foreground" />
            <span className="font-medium capitalize">{address.type} Address</span>
            {address.is_default && (
              <Badge variant="default" className="text-xs">
                <Check className="h-3 w-3 mr-1" />
                Default
              </Badge>
            )}
          </div>
          
          <div className="flex gap-1">
            <Button
              variant="ghost"
              size="sm"
              onClick={onEdit}
              className="h-8 w-8 p-0"
            >
              <Edit className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={onDelete}
              className="h-8 w-8 p-0 text-destructive hover:text-destructive"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </div>

        <div className="space-y-2">
          <div className="font-medium">
            {address.first_name} {address.last_name}
          </div>
          
          {address.company && (
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Building2 className="h-3 w-3" />
              {address.company}
            </div>
          )}
          
          <div className="text-sm text-muted-foreground">
            {formatAddress()}
          </div>
          
          {address.phone_number && (
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Phone className="h-3 w-3" />
              {address.phone_number}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
