import { Bar, BarChart, XAxis, YAxis } from "recharts";

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
import type { CategoryCount } from "@/types/dashboard";

interface CategoryBarChartProps {
  data: CategoryCount[];
}

const chartConfig = {
  count: {
    label: "Risks",
    color: "hsl(var(--primary))",
  },
} satisfies ChartConfig;

export function CategoryBarChart({ data }: CategoryBarChartProps) {
  const chartData = data.map((cat) => ({
    name: cat.category_name,
    count: cat.count,
  }));

  return (
    <Card className="flex flex-col">
      <CardHeader>
        <CardTitle>Risks by Category</CardTitle>
        <CardDescription>Number of risks in each category</CardDescription>
      </CardHeader>
      <CardContent className="flex-1">
        {chartData.length === 0 ? (
          <div className="flex h-[200px] items-center justify-center text-muted-foreground">
            No categories available
          </div>
        ) : (
          <ChartContainer config={chartConfig} className="min-h-[200px] w-full">
            <BarChart
              layout="vertical"
              data={chartData}
              margin={{ left: 80 }}
            >
              <XAxis type="number" />
              <YAxis
                dataKey="name"
                type="category"
                tickLine={false}
                axisLine={false}
                width={70}
              />
              <ChartTooltip content={<ChartTooltipContent />} />
              <Bar dataKey="count" fill="hsl(var(--primary))" radius={4} />
            </BarChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  );
}
