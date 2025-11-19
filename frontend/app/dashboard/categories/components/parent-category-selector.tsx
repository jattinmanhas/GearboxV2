"use client"

import * as React from "react"
import { Check, ChevronDown, Search } from "lucide-react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Category } from "@/lib/types"

interface ParentCategorySelectorProps {
    value?: number
    onChange: (value?: number) => void
    onSearch: (query: string) => void
    options: Category[]
    placeholder?: string
}

export function ParentCategorySelector({
    value,
    onChange,
    onSearch,
    options,
    placeholder = "Select parent category"
}: ParentCategorySelectorProps) {
    const [isOpen, setIsOpen] = React.useState(false)
    const [searchTerm, setSearchTerm] = React.useState("")

    // Debounce search
    React.useEffect(() => {
        const timer = setTimeout(() => {
            onSearch(searchTerm)
        }, 300)
        return () => clearTimeout(timer)
    }, [searchTerm, onSearch])

    const selectedLabel = React.useMemo(() => {
        if (value === 0) return "Root Categories Only"
        if (value === undefined) return "All Categories"
        return options.find(o => o.id === value)?.name || "Select parent"
    }, [value, options])

    return (
        <DropdownMenu open={isOpen} onOpenChange={setIsOpen}>
            <DropdownMenuTrigger asChild>
                <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={isOpen}
                    className="w-full justify-between font-normal"
                >
                    <span className="truncate">{selectedLabel}</span>
                    <ChevronDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-[300px] p-0" align="start">
                <div className="flex items-center border-b px-3">
                    <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
                    <Input
                        placeholder="Search categories..."
                        className="flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50 border-0 focus-visible:ring-0 px-0"
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>
                <div className="max-h-[300px] overflow-y-auto p-1">
                    <div
                        className={cn(
                            "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
                            value === undefined && "bg-accent text-accent-foreground"
                        )}
                        onClick={() => {
                            onChange(undefined)
                            setIsOpen(false)
                        }}
                    >
                        <Check
                            className={cn(
                                "mr-2 h-4 w-4",
                                value === undefined ? "opacity-100" : "opacity-0"
                            )}
                        />
                        All Categories
                    </div>
                    <div
                        className={cn(
                            "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
                            value === 0 && "bg-accent text-accent-foreground"
                        )}
                        onClick={() => {
                            onChange(0)
                            setIsOpen(false)
                        }}
                    >
                        <Check
                            className={cn(
                                "mr-2 h-4 w-4",
                                value === 0 ? "opacity-100" : "opacity-0"
                            )}
                        />
                        Root Categories Only
                    </div>
                    {options.length > 0 && <div className="my-1 h-px bg-muted" />}
                    {options.map((category) => (
                        <div
                            key={category.id}
                            className={cn(
                                "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
                                value === category.id && "bg-accent text-accent-foreground"
                            )}
                            onClick={() => {
                                onChange(category.id)
                                setIsOpen(false)
                            }}
                        >
                            <Check
                                className={cn(
                                    "mr-2 h-4 w-4",
                                    value === category.id ? "opacity-100" : "opacity-0"
                                )}
                            />
                            {category.name}
                        </div>
                    ))}
                    {options.length === 0 && searchTerm && (
                        <div className="py-6 text-center text-sm text-muted-foreground">
                            No categories found.
                        </div>
                    )}
                </div>
            </DropdownMenuContent>
        </DropdownMenu>
    )
}
