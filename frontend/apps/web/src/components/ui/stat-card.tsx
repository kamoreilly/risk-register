import * as React from "react"
import { cn } from "@/lib/utils"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

interface StatCardProps {
  title: string
  value: string | number
  icon?: React.ReactNode
  description?: string
  trend?: string
  isUrgent?: boolean
  className?: string
}

function StatCard({
  title,
  value,
  icon,
  description,
  trend,
  isUrgent,
  className,
}: StatCardProps) {
  return (
    <Card
      className={cn(
        "transition-all hover:shadow-md overflow-hidden",
        isUrgent
          ? "border-l-4 border-l-destructive shadow-destructive/10"
          : "border-l-4 border-l-primary shadow-primary/5",
        className
      )}
    >
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icon && (
          <div className="bg-primary/10 text-primary p-2 rounded-lg">
            {icon}
          </div>
        )}
      </CardHeader>
      <CardContent>
        <div className="flex items-baseline gap-2">
          <div className="text-2xl font-bold tracking-tight">{value}</div>
          {trend && (
            <Badge
              variant={trend === "High" ? "destructive" : "secondary"}
              className="text-[10px] h-4"
            >
              {trend}
            </Badge>
          )}
        </div>
        {description && (
          <p className="text-xs text-muted-foreground mt-1">{description}</p>
        )}
      </CardContent>
    </Card>
  )
}

export { StatCard }
export type { StatCardProps }
