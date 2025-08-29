import React from 'react'
import { LoginForm } from './components/LoginForm'
import { CpuIcon } from 'lucide-react'
import { ThemeToggle } from '@/components/ui/theme-toggle'

const LoginPage = () => {
  return (
    <div className="min-h-svh bg-background flex flex-col items-center justify-center">
      <div className="flex items-center justify-between w-full max-w-xs mb-8">
        <a href="#" className="flex items-center gap-2 font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <CpuIcon className="size-4" />
          </div>
         <span className="text-secondary-foreground text-xl">GearBox</span>
        </a>
        <ThemeToggle />
      </div>
      <div className="w-full max-w-xs">
        <LoginForm />
      </div>
    </div>
  )
}

export default LoginPage