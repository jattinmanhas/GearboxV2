"use client"

import React from "react"

interface DonutSegment {
  label: string
  value: number
  color: string
}

export function DonutChart({
  segments,
  size = 140,
  thickness = 18,
}: {
  segments: DonutSegment[]
  size?: number
  thickness?: number
}) {
  const radius = (size - thickness) / 2
  const circumference = 2 * Math.PI * radius
  const total = segments.reduce((sum, segment) => sum + segment.value, 0) || 1
  let offset = 0

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
      <circle
        cx={size / 2}
        cy={size / 2}
        r={radius}
        fill="none"
        stroke="hsl(var(--muted))"
        strokeWidth={thickness}
      />
      {segments.map((segment, index) => {
        const value = Math.max(segment.value, 0)
        const length = (value / total) * circumference
        const dashOffset = circumference - offset
        offset += length

        return (
          <circle
            key={`${segment.label}-${index}`}
            cx={size / 2}
            cy={size / 2}
            r={radius}
            fill="none"
            stroke={segment.color}
            strokeWidth={thickness}
            strokeDasharray={`${length} ${circumference - length}`}
            strokeDashoffset={dashOffset}
            strokeLinecap="round"
            transform={`rotate(-90 ${size / 2} ${size / 2})`}
          />
        )
      })}
    </svg>
  )
}

export function LineChart({
  data,
  height = 120,
}: {
  data: Array<{ label: string; value: number }>
  height?: number
}) {
  const width = 360
  const maxValue = Math.max(...data.map((point) => point.value), 1)
  const points = data.map((point, index) => {
    const x = (index / Math.max(data.length - 1, 1)) * width
    const y = height - (point.value / maxValue) * height
    return `${x},${y}`
  })

  return (
    <svg width="100%" height={height} viewBox={`0 0 ${width} ${height}`}>
      <polyline
        fill="none"
        stroke="hsl(var(--primary))"
        strokeWidth="3"
        points={points.join(" ")}
      />
      {data.map((point, index) => {
        const x = (index / Math.max(data.length - 1, 1)) * width
        const y = height - (point.value / maxValue) * height
        return (
          <circle
            key={`${point.label}-${index}`}
            cx={x}
            cy={y}
            r="3"
            fill="hsl(var(--primary))"
          />
        )
      })}
    </svg>
  )
}

export function BarList({
  items,
}: {
  items: Array<{ label: string; value: number; description?: string }>
}) {
  const maxValue = Math.max(...items.map((item) => item.value), 1)

  return (
    <div className="space-y-3">
      {items.map((item, index) => (
        <div key={`${item.label}-${index}`} className="space-y-1">
          <div className="flex items-center justify-between text-sm">
            <span className="font-medium">{item.label}</span>
            <span className="text-muted-foreground">{item.value.toLocaleString()}</span>
          </div>
          {item.description ? (
            <div className="text-xs text-muted-foreground">{item.description}</div>
          ) : null}
          <div className="h-2 rounded-full bg-muted">
            <div
              className="h-2 rounded-full bg-primary"
              style={{ width: `${(item.value / maxValue) * 100}%` }}
            />
          </div>
        </div>
      ))}
    </div>
  )
}
