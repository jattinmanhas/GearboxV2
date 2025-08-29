import React from 'react'
import { CpuIcon } from 'lucide-react'
import { ThemeToggle } from '@/components/ui/theme-toggle'
import { RegisterForm } from './components/RegisterForm'
import Link from 'next/link'

const RegisterPage = () => {
  return (
    <div className="min-h-svh bg-background flex flex-col items-center justify-center">
      <div className="flex items-center justify-between w-full max-w-md mb-8">
        <Link href="/" className="flex items-center gap-2 font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <CpuIcon className="size-4" />
          </div>
          <span className="text-secondary-foreground text-xl">GearBox</span>
          </Link>
        <ThemeToggle />
      </div>
      <div className="w-full max-w-md">
        <RegisterForm />
      </div>
    </div>
  )
}

export default RegisterPage