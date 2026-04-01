import * as React from "react"
import { Progress } from "@/components/ui/progress"
import { cn } from "@/lib/utils"

interface ProgressReportProps {
  label: string
  value: number
  total: number
  color?: string
  className?: string
}

export function ProgressReport({ 
  label, 
  value, 
  total, 
  color = "bg-primary",
  className 
}: ProgressReportProps) {
  const percentage = total > 0 ? Math.round((value / total) * 100) : 0

  return (
    <div className={cn("space-y-2", className)}>
      <div className="flex items-center justify-between text-sm">
        <span className="font-medium">{label}</span>
        <span className="text-muted-foreground">{value} / {total} ({percentage}%)</span>
      </div>
      <Progress value={percentage} className="h-2" />
    </div>
  )
}
