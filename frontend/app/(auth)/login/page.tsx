import React from 'react'
import { LoginForm } from './components/LoginForm'
import { CpuIcon } from 'lucide-react'
import Link from 'next/link'
import { ThemeToggle } from '@/components/ui/theme-toggle'

export default function LoginPage() {
  return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center px-4">
      <div className="flex items-center justify-between w-full max-w-sm md:max-w-lg mb-8">
        <Link href="/" className="flex items-center gap-2 font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <CpuIcon className="size-4" />
          </div>
         <span className="text-secondary-foreground text-xl">GearBox</span>
        </Link>
        <ThemeToggle />
      </div>
      <div className="w-full max-w-sm md:max-w-lg">
        <LoginForm />
      </div>
    </div>
  )
}