import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";

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
import type { StatusTimeDataPoint } from "@/types/analytics";

interface StatusOverTimeChartProps {
  data: StatusTimeDataPoint[];
}

const chartConfig = {
  open: {
    label: "Opened",
    color: "hsl(var(--primary))",
  },
  closed: {
    label: "Closed",
    color: "hsl(142, 76%, 36%)",
  },
} satisfies ChartConfig;

export function StatusOverTimeChart({ data }: StatusOverTimeChartProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Opened vs Closed Over Time</CardTitle>
        <CardDescription>Risks opened and closed per period</CardDescription>
      </CardHeader>
      <CardContent>
        {data.length === 0 ? (
          <div className="flex h-[200px] items-center justify-center text-muted-foreground italic text-sm">
            No data available for this period
          </div>
        ) : (
          <ChartContainer config={chartConfig} className="h-[200px] w-full">
            <BarChart data={data} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} opacity={0.3} />
              <XAxis
                dataKey="period"
                tickLine={false}
                axisLine={false}
                tickMargin={8}
              />
              <YAxis tickLine={false} axisLine={false} tickMargin={8} />
              <ChartTooltip content={<ChartTooltipContent />} />
              <Bar
                dataKey="open"
                fill="hsl(var(--primary))"
                radius={[4, 4, 0, 0]}
              />
              <Bar
                dataKey="closed"
                fill="hsl(142, 76%, 36%)"
                radius={[4, 4, 0, 0]}
              />
            </BarChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  );
}
