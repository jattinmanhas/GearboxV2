import { cn } from "@/lib/utils"

interface PasswordStrengthProps {
  password: string
  className?: string
}

export function PasswordStrength({ password, className }: PasswordStrengthProps) {
  if (!password) return null

  const checks = [
    {
      label: "At least 6 characters",
      valid: password.length >= 6,
    },
    {
      label: "Contains uppercase letter",
      valid: /[A-Z]/.test(password),
    },
    {
      label: "Contains lowercase letter", 
      valid: /[a-z]/.test(password),
    },
    {
      label: "Contains number",
      valid: /\d/.test(password),
    },
  ]

  const validChecks = checks.filter(check => check.valid).length
  const strength = validChecks / checks.length

  return (
    <div className={cn("space-y-2", className)}>
      <div className="space-y-1">
        {checks.map((check, index) => (
          <div key={index} className="flex items-center gap-2 text-xs">
            <div className={cn(
              "w-1.5 h-1.5 rounded-full",
              check.valid ? "bg-green-500" : "bg-gray-300"
            )} />
            <span className={cn(
              check.valid ? "text-green-600" : "text-gray-500"
            )}>
              {check.label}
            </span>
          </div>
        ))}
      </div>
      <div className="w-full bg-gray-200 rounded-full h-1.5">
        <div 
          className={cn(
            "h-1.5 rounded-full transition-all duration-300",
            strength === 1 ? "bg-green-500" : 
            strength >= 0.75 ? "bg-yellow-500" :
            strength >= 0.5 ? "bg-orange-500" : "bg-red-500"
          )}
          style={{ width: `${strength * 100}%` }}
        />
      </div>
    </div>
  )
}
