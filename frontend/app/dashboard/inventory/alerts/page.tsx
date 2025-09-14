"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { 
  AlertTriangle, 
  TrendingDown, 
  TrendingUp,
  RefreshCw,
  CheckCircle,
  XCircle
} from "lucide-react"
import { InventoryAlert } from "@/lib/types"

export default function InventoryAlertsPage() {
  const [alerts, setAlerts] = useState<InventoryAlert[]>([])
  const [loading, setLoading] = useState(true)
  const [resolvedFilter, setResolvedFilter] = useState<boolean | undefined>(undefined)

  // Fetch inventory alerts
  const fetchAlerts = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      if (resolvedFilter !== undefined) {
        params.append('resolved', resolvedFilter.toString())
      }

      const response = await fetch(`/api/v1/inventory/alerts?${params}`, {
        credentials: 'include'
      })
      const data = await response.json()
      
      if (response.ok) {
        setAlerts(data.data || [])
      }
    } catch (error) {
      console.error('Error fetching inventory alerts:', error)
    } finally {
      setLoading(false)
    }
  }

  // Resolve alert
  const resolveAlert = async (alertId: number) => {
    try {
      const response = await fetch(`/api/v1/inventory/alerts/${alertId}/resolve`, {
        method: 'PUT',
        credentials: 'include'
      })

      if (response.ok) {
        fetchAlerts()
      }
    } catch (error) {
      console.error('Error resolving alert:', error)
    }
  }

  useEffect(() => {
    fetchAlerts()
  }, [resolvedFilter])

  const getAlertIcon = (alertType: string) => {
    switch (alertType) {
      case 'low_stock':
        return <TrendingDown className="h-5 w-5 text-yellow-500" />
      case 'out_of_stock':
        return <XCircle className="h-5 w-5 text-red-500" />
      case 'overstock':
        return <TrendingUp className="h-5 w-5 text-blue-500" />
      default:
        return <AlertTriangle className="h-5 w-5 text-gray-500" />
    }
  }

  const getAlertColor = (alertType: string) => {
    switch (alertType) {
      case 'low_stock':
        return 'bg-yellow-50 border-yellow-200'
      case 'out_of_stock':
        return 'bg-red-50 border-red-200'
      case 'overstock':
        return 'bg-blue-50 border-blue-200'
      default:
        return 'bg-gray-50 border-gray-200'
    }
  }

  const getAlertBadgeColor = (alertType: string) => {
    switch (alertType) {
      case 'low_stock':
        return 'bg-yellow-100 text-yellow-800'
      case 'out_of_stock':
        return 'bg-red-100 text-red-800'
      case 'overstock':
        return 'bg-blue-100 text-blue-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const filteredAlerts = alerts.filter(alert => {
    if (resolvedFilter === undefined) return true
    return alert.is_resolved === resolvedFilter
  })

  const activeAlerts = alerts.filter(alert => !alert.is_resolved)
  const resolvedAlerts = alerts.filter(alert => alert.is_resolved)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Inventory Alerts</h1>
          <p className="text-muted-foreground">
            Monitor and manage inventory alerts and warnings
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchAlerts}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Alerts</CardTitle>
            <AlertTriangle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{alerts.length}</div>
            <p className="text-xs text-muted-foreground">
              All inventory alerts
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Alerts</CardTitle>
            <XCircle className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{activeAlerts.length}</div>
            <p className="text-xs text-muted-foreground">
              Requiring attention
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Resolved Alerts</CardTitle>
            <CheckCircle className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{resolvedAlerts.length}</div>
            <p className="text-xs text-muted-foreground">
              Successfully resolved
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardHeader>
          <CardTitle>Filter Alerts</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2">
            <Button
              variant={resolvedFilter === undefined ? "default" : "outline"}
              onClick={() => setResolvedFilter(undefined)}
            >
              All Alerts
            </Button>
            <Button
              variant={resolvedFilter === false ? "default" : "outline"}
              onClick={() => setResolvedFilter(false)}
            >
              Active Only
            </Button>
            <Button
              variant={resolvedFilter === true ? "default" : "outline"}
              onClick={() => setResolvedFilter(true)}
            >
              Resolved Only
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Alerts List */}
      <Card>
        <CardHeader>
          <CardTitle>Inventory Alerts</CardTitle>
          <CardDescription>
            {filteredAlerts.length} alert(s) found
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="text-center py-8">
              <RefreshCw className="h-6 w-6 animate-spin mx-auto mb-2" />
              Loading alerts...
            </div>
          ) : filteredAlerts.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              No alerts found
            </div>
          ) : (
            <div className="space-y-4">
              {filteredAlerts.map((alert) => (
                <div
                  key={alert.id}
                  className={`p-4 border rounded-lg ${getAlertColor(alert.alert_type)}`}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex items-start gap-3">
                      {getAlertIcon(alert.alert_type)}
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          <h3 className="font-medium">
                            {alert.product_name} {alert.variant_name && `(${alert.variant_name})`}
                          </h3>
                          <Badge className={getAlertBadgeColor(alert.alert_type)}>
                            {alert.alert_type.replace('_', ' ').toUpperCase()}
                          </Badge>
                          <Badge variant={alert.is_resolved ? "secondary" : "destructive"}>
                            {alert.is_resolved ? "Resolved" : "Active"}
                          </Badge>
                        </div>
                        <div className="text-sm text-muted-foreground space-y-1">
                          <p>
                            <strong>Current Stock:</strong> {alert.current_quantity} units
                          </p>
                          <p>
                            <strong>Threshold:</strong> {alert.threshold_quantity} units
                          </p>
                          <p>
                            <strong>SKU:</strong> {alert.product_sku} {alert.variant_sku && `(${alert.variant_sku})`}
                          </p>
                          <p>
                            <strong>Created:</strong> {new Date(alert.created_at).toLocaleString()}
                          </p>
                          {alert.resolved_at && (
                            <p>
                              <strong>Resolved:</strong> {new Date(alert.resolved_at).toLocaleString()}
                            </p>
                          )}
                        </div>
                      </div>
                    </div>
                    {!alert.is_resolved && (
                      <Button
                        size="sm"
                        onClick={() => resolveAlert(alert.id)}
                        className="ml-4"
                      >
                        <CheckCircle className="h-4 w-4 mr-2" />
                        Resolve
                      </Button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
