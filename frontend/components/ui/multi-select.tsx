"use client"

import * as React from "react"
import { Check, ChevronDown, X } from "lucide-react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

export interface MultiSelectOption {
  value: string | number
  label: string
}

export interface MultiSelectProps {
  options: MultiSelectOption[]
  selectedValues: (string | number)[]
  onSelectionChange: (selectedValues: (string | number)[]) => void
  placeholder?: string
  searchPlaceholder?: string
  className?: string
  disabled?: boolean
  maxDisplayed?: number
}

export function MultiSelect({
  options,
  selectedValues,
  onSelectionChange,
  placeholder = "Select items...",
  searchPlaceholder = "Search...",
  className,
  disabled = false,
  maxDisplayed = 3,
}: MultiSelectProps) {
  const [isOpen, setIsOpen] = React.useState(false)
  const [searchTerm, setSearchTerm] = React.useState("")
  const [focusedIndex, setFocusedIndex] = React.useState(-1)

  const filteredOptions = React.useMemo(() => {
    if (!searchTerm) return options
    return options.filter(option =>
      option.label.toLowerCase().includes(searchTerm.toLowerCase())
    )
  }, [options, searchTerm])

  const selectedOptions = React.useMemo(() => {
    return options.filter(option => selectedValues.includes(option.value))
  }, [options, selectedValues])

  const handleToggleOption = (value: string | number) => {
    const newSelection = selectedValues.includes(value)
      ? selectedValues.filter(v => v !== value)
      : [...selectedValues, value]
    onSelectionChange(newSelection)
  }

  const handleRemoveOption = (value: string | number) => {
    const newSelection = selectedValues.filter(v => v !== value)
    onSelectionChange(newSelection)
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!isOpen) {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault()
        setIsOpen(true)
      }
      return
    }

    switch (e.key) {
      case "ArrowDown":
        e.preventDefault()
        setFocusedIndex(prev => 
          prev < filteredOptions.length - 1 ? prev + 1 : 0
        )
        break
      case "ArrowUp":
        e.preventDefault()
        setFocusedIndex(prev => 
          prev > 0 ? prev - 1 : filteredOptions.length - 1
        )
        break
      case "Enter":
        e.preventDefault()
        if (focusedIndex >= 0 && focusedIndex < filteredOptions.length) {
          handleToggleOption(filteredOptions[focusedIndex].value)
        }
        break
      case "Escape":
        setIsOpen(false)
        setFocusedIndex(-1)
        break
    }
  }

  const displayText = React.useMemo(() => {
    if (selectedValues.length === 0) return placeholder
    
    if (selectedValues.length <= maxDisplayed) {
      return selectedOptions.map(option => option.label).join(", ")
    }
    
    return `${selectedValues.length} items selected`
  }, [selectedValues, selectedOptions, placeholder, maxDisplayed])

  return (
    <div className={cn("relative", className)}>
      <DropdownMenu open={isOpen} onOpenChange={setIsOpen}>
        <DropdownMenuTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={isOpen}
            className={cn(
              "w-full justify-between text-left font-normal",
              !selectedValues.length && "text-muted-foreground"
            )}
            disabled={disabled}
            onKeyDown={handleKeyDown}
          >
            <span className="truncate">{displayText}</span>
            <ChevronDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </DropdownMenuTrigger>
        
        <DropdownMenuContent className="w-[--radix-dropdown-menu-trigger-width] min-w-0 p-0" align="start" sideOffset={4}>
          <div className="p-2">
            <Input
              placeholder={searchPlaceholder}
              value={searchTerm}
              onChange={(e) => {
                setSearchTerm(e.target.value)
                setFocusedIndex(-1)
              }}
              className="h-8"
              onKeyDown={(e) => {
                if (e.key === "ArrowDown") {
                  e.preventDefault()
                  setFocusedIndex(0)
                }
              }}
            />
          </div>
          
          <div className="max-h-60 overflow-y-auto">
            {filteredOptions.length === 0 ? (
              <div className="px-2 py-1.5 text-sm text-muted-foreground">
                No options found
              </div>
            ) : (
              filteredOptions.map((option, index) => (
                <div
                  key={option.value}
                  className={cn(
                    "relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground",
                    focusedIndex === index && "bg-accent text-accent-foreground",
                    selectedValues.includes(option.value) && "bg-accent/50"
                  )}
                  onClick={() => handleToggleOption(option.value)}
                  onMouseEnter={() => setFocusedIndex(index)}
                >
                  <span className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
                    {selectedValues.includes(option.value) && (
                      <Check className="h-4 w-4" />
                    )}
                  </span>
                  <span className="pl-8">{option.label}</span>
                </div>
              ))
            )}
          </div>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Selected items display */}
      {selectedValues.length > 0 && (
        <div className="mt-2 flex flex-wrap gap-1">
          {selectedOptions.slice(0, maxDisplayed).map((option) => (
            <Badge
              key={option.value}
              variant="secondary"
              className="flex items-center gap-1 pr-1"
            >
              <span className="truncate max-w-[120px]">{option.label}</span>
              <Button
                variant="ghost"
                size="sm"
                className="h-4 w-4 p-0 hover:bg-transparent"
                onClick={() => handleRemoveOption(option.value)}
              >
                <X className="h-3 w-3" />
              </Button>
            </Badge>
          ))}
          {selectedValues.length > maxDisplayed && (
            <Badge variant="outline" className="text-xs">
              +{selectedValues.length - maxDisplayed} more
            </Badge>
          )}
        </div>
      )}
    </div>
  )
}
