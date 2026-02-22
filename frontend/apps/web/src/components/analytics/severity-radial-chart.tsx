import { Pie, PieChart } from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";

interface SeverityRadialChartProps {
  data: Record<string, number>;
}

const chartConfig = {
  critical: {
    label: "Critical",
    color: "hsl(var(--destructive))",
  },
  high: {
    label: "High",
    color: "hsl(25, 95%, 53%)",
  },
  medium: {
    label: "Medium",
    color: "hsl(45, 93%, 47%)",
  },
  low: {
    label: "Low",
    color: "hsl(var(--muted-foreground))",
  },
} satisfies ChartConfig;

export function SeverityRadialChart({ data }: SeverityRadialChartProps) {
  const chartData = Object.entries(data).map(([severity, count]) => ({
    severity,
    count,
    fill: chartConfig[severity as keyof typeof chartConfig]?.color ?? "hsl(var(--muted))",
  }));

  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Risk Severity</CardTitle>
        <CardDescription>Distribution by severity level</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-0">
        <ChartContainer
          config={chartConfig}
          className="mx-auto aspect-square max-h-[250px]"
        >
          <PieChart>
            <ChartTooltip
              content={<ChartTooltipContent nameKey="severity" hideLabel />}
            />
            <Pie
              data={chartData}
              dataKey="count"
              nameKey="severity"
              innerRadius={60}
              strokeWidth={5}
            />
          </PieChart>
        </ChartContainer>
        <div className="mt-4 flex flex-wrap justify-center gap-4 text-xs">
          {Object.entries(chartConfig).map(([key, config]) => (
            <div key={key} className="flex items-center gap-1.5">
              <div
                className="size-2.5 shrink-0 rounded-[2px]"
                style={{ backgroundColor: config.color }}
              />
              <span className="text-muted-foreground">{config.label}</span>
              <span className="font-medium">{data[key] ?? 0}</span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
