"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  AlertTriangle,
  TrendingDown,
  TrendingUp,
  RefreshCw,
  CheckCircle,
  XCircle,
} from "lucide-react";
import { InventoryAlert } from "@/lib/types";
import { inventoryApi } from "@/lib/apiFunctions";

export default function InventoryAlertsPage() {
  const [alerts, setAlerts] = useState<InventoryAlert[]>([]);
  const [loading, setLoading] = useState(true);
  const [resolvedFilter, setResolvedFilter] = useState<boolean | undefined>(
    undefined
  );
  const [isCheckingAlerts, setIsCheckingAlerts] = useState(false);

  // Fetch inventory alerts
  const fetchAlerts = async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams();
      if (resolvedFilter !== undefined) {
        params.append("resolved", resolvedFilter.toString());
      }

      const response = await inventoryApi.getInventoryAlerts(params.toString());
      console.log(response);
      if (response.success) {
        setAlerts(response.data || []);
      } else {
        const msg = response.message || "Failed to fetch inventory alerts";
        console.error("Inventory alerts fetch failed:", response);
      }
    } catch (error) {
      console.error("Error fetching inventory alerts:", error);
    } finally {
      setLoading(false);
    }
  };

  // Check inventory alerts
  const checkAlerts = async () => {
    try {
      setIsCheckingAlerts(true);
      const response = await inventoryApi.checkInventoryAlerts();

      if (response.success) {
        // Refresh the alerts list after checking
        await fetchAlerts();
      } else {
        console.error("Failed to check inventory alerts:", response.message);
      }
    } catch (error) {
      console.error("Error checking inventory alerts:", error);
    } finally {
      setIsCheckingAlerts(false);
    }
  };

  // Resolve alert
  const resolveAlert = async (alertId: number) => {
    try {
      const response = await inventoryApi.resolveInventoryAlert(alertId);

      if (response.success) {
        fetchAlerts();
      }
    } catch (error) {
      console.error("Error resolving alert:", error);
    }
  };

  useEffect(() => {
    fetchAlerts();
  }, [resolvedFilter]);

  const getAlertIcon = (alertType: string) => {
    switch (alertType) {
      case "low_stock":
        return <TrendingDown className="h-5 w-5 text-yellow-500" />;
      case "out_of_stock":
        return <XCircle className="h-5 w-5 text-red-500" />;
      case "overstock":
        return <TrendingUp className="h-5 w-5 text-blue-500" />;
      default:
        return <AlertTriangle className="h-5 w-5 text-gray-500" />;
    }
  };

  const getAlertColor = (alertType: string) => {
    switch (alertType) {
      case "low_stock":
        return "bg-yellow-50 border-yellow-200";
      case "out_of_stock":
        return "bg-red-50 border-red-200";
      case "overstock":
        return "bg-blue-50 border-blue-200";
      default:
        return "bg-gray-50 border-gray-200";
    }
  };

  const getAlertBadgeColor = (alertType: string) => {
    switch (alertType) {
      case "low_stock":
        return "bg-yellow-100 text-yellow-800";
      case "out_of_stock":
        return "bg-red-100 text-red-800";
      case "overstock":
        return "bg-blue-100 text-blue-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const filteredAlerts = alerts.filter((alert) => {
    if (resolvedFilter === undefined) return true;
    return alert.is_resolved === resolvedFilter;
  });

  const activeAlerts = alerts.filter((alert) => !alert.is_resolved);
  const resolvedAlerts = alerts.filter((alert) => alert.is_resolved);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            Inventory Alerts
          </h1>
          <p className="text-muted-foreground">
            Monitor and manage inventory alerts and warnings
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="default"
            size="sm"
            onClick={checkAlerts}
            disabled={isCheckingAlerts}
          >
            <AlertTriangle className="h-4 w-4 mr-2" />
            {isCheckingAlerts ? "Checking..." : "Check Inventory Alerts"}
          </Button>
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
            <div className="text-2xl font-bold text-red-600">
              {activeAlerts.length}
            </div>
            <p className="text-xs text-muted-foreground">Requiring attention</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Resolved Alerts
            </CardTitle>
            <CheckCircle className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {resolvedAlerts.length}
            </div>
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
                  className="p-4 border rounded-lg bg-card hover:bg-accent/5 transition-colors"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex items-start gap-3">
                      {getAlertIcon(alert.alert_type)}
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          <h3 className="font-semibold text-foreground">
                            {alert.product_name}{" "}
                            {alert.variant_name && `(${alert.variant_name})`}
                          </h3>
                          <Badge
                            className={getAlertBadgeColor(alert.alert_type)}
                          >
                            {alert.alert_type.replace("_", " ").toUpperCase()}
                          </Badge>
                          <Badge
                            variant={
                              alert.is_resolved ? "secondary" : "destructive"
                            }
                          >
                            {alert.is_resolved ? "Resolved" : "Active"}
                          </Badge>
                        </div>
                        <div className="text-sm space-y-1.5">
                          <p className="text-foreground/80">
                            <span className="font-medium text-foreground">
                              Current Stock:
                            </span>{" "}
                            {alert.current_quantity} units
                          </p>
                          <p className="text-foreground/80">
                            <span className="font-medium text-foreground">
                              Threshold:
                            </span>{" "}
                            {alert.threshold_quantity} units
                          </p>
                          <p className="text-foreground/80">
                            <span className="font-medium text-foreground">
                              SKU:
                            </span>{" "}
                            {alert.product_sku}{" "}
                            {alert.variant_sku && `(${alert.variant_sku})`}
                          </p>
                          <p className="text-foreground/70 text-xs">
                            <span className="font-medium text-foreground/80">
                              Created:
                            </span>{" "}
                            {new Date(alert.created_at).toLocaleString()}
                          </p>
                          {alert.resolved_at && (
                            <p className="text-foreground/70 text-xs">
                              <span className="font-medium text-foreground/80">
                                Resolved:
                              </span>{" "}
                              {new Date(alert.resolved_at).toLocaleString()}
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
  );
}