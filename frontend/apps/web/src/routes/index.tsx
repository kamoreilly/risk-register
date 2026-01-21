import { createFileRoute, Link } from "@tanstack/react-router";
import * as React from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  XAxis,
} from "recharts";
import {
  BotIcon,
  BrainCircuitIcon,
  CalendarClockIcon,
  FingerprintIcon,
  LineChartIcon,
  SearchIcon,
  SparklesIcon,
  WorkflowIcon,
} from "lucide-react";

import { HeaderNav } from "@/components/header-nav";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/")({
  component: HomeComponent,
});

function HomeComponent() {
  const chartData = React.useMemo(
    () => [
      { date: "2025-08-01", critical: 4, high: 12, medium: 18, low: 25 },
      { date: "2025-09-01", critical: 3, high: 15, medium: 22, low: 20 },
      { date: "2025-10-01", critical: 5, high: 10, medium: 25, low: 18 },
      { date: "2025-11-01", critical: 2, high: 18, medium: 20, low: 15 },
      { date: "2025-12-01", critical: 6, high: 14, medium: 28, low: 12 },
      { date: "2026-01-01", critical: 4, high: 16, medium: 24, low: 10 },
    ],
    [],
  );

  const chartConfig = {
    critical: { label: "Critical", color: "hsl(0 84.2% 60.2%)" },
    high: { label: "High", color: "hsl(24.6 95% 53.1%)" },
    medium: { label: "Medium", color: "hsl(47.9 95.8% 53.1%)" },
    low: { label: "Low", color: "hsl(142.1 76.2% 36.3%)" },
  } satisfies ChartConfig;

  const [activeMetric, setActiveMetric] = React.useState<keyof typeof chartConfig>("critical");

  const total = React.useMemo(
    () => ({
      critical: chartData.reduce((acc, curr) => acc + curr.critical, 0),
      high: chartData.reduce((acc, curr) => acc + curr.high, 0),
      medium: chartData.reduce((acc, curr) => acc + curr.medium, 0),
      low: chartData.reduce((acc, curr) => acc + curr.low, 0),
    }),
    [chartData],
  );

  const heroFeatures = React.useMemo(
    () =>
      [
        {
          icon: BrainCircuitIcon,
          title: "AI summaries",
          description: "Turn updates into a review-ready brief.",
        },
        {
          icon: WorkflowIcon,
          title: "Smart triage",
          description: "Group changes by owner, severity, and status.",
        },
        {
          icon: CalendarClockIcon,
          title: "Cadence tracking",
          description: "Keep review dates visible and actionable.",
        },
        {
          icon: SearchIcon,
          title: "Fast lookup",
          description: "Find risks instantly with keyboard-first search.",
        },
        {
          icon: FingerprintIcon,
          title: "Audit-friendly",
          description: "Owners and decisions stay attributable.",
        },
        {
          icon: BotIcon,
          title: "Assistive prompts",
          description: "Draft mitigations and next steps in context.",
        },
      ] as const,
    [],
  );

  return (
    <div className="flex min-h-svh flex-col bg-background">
      <HeaderNav />
      <main className="flex-1 w-full pb-12 pt-8 sm:pt-12">
        <section className="w-full">
          <div className="mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8">
            <div className="grid gap-16">
              {/* Hero Section */}
              <div className="grid gap-12 lg:grid-cols-[1fr_1.1fr] lg:items-center">
                <div className="grid gap-6 text-center lg:text-left">
                  <div className="inline-flex w-fit items-center gap-2 rounded-full border bg-muted/40 px-3 py-1 text-xs text-muted-foreground mx-auto lg:mx-0">
                    <span className="size-1.5 rounded-full bg-emerald-500" />
                    <span>Live view</span>
                  </div>

                  <h1 className="text-balance text-4xl font-semibold tracking-tight sm:text-5xl lg:text-6xl">
                    Risk register, focused on what changes
                  </h1>
                  <p className="max-w-[600px] text-pretty text-base leading-7 text-muted-foreground sm:text-lg mx-auto lg:mx-0">
                    Capture owners, severity, and review cadence in one view. Track risks, monitor trends, and maintain accountability across your organization.
                  </p>

                  <div className="flex flex-wrap justify-center gap-3 lg:justify-start">
                    <Link to="/login">
                      <Button size="lg" className="rounded-full px-8">Open</Button>
                    </Link>
                    <a
                      href="#analytics"
                      className="inline-flex h-9 items-center justify-center rounded-full border bg-background px-8 text-xs font-medium hover:bg-muted transition-colors"
                    >
                      View Analytics
                    </a>
                  </div>

                  <dl className="grid grid-cols-3 gap-4 pt-4">
                    <div className="rounded-2xl border bg-card/50 p-4">
                      <dt className="text-xs font-medium text-muted-foreground">Open</dt>
                      <dd className="mt-1 text-2xl font-bold tracking-tight">124</dd>
                    </div>
                    <div className="rounded-2xl border bg-card/50 p-4">
                      <dt className="text-xs font-medium text-muted-foreground">Mitigating</dt>
                      <dd className="mt-1 text-2xl font-bold tracking-tight">42</dd>
                    </div>
                    <div className="rounded-2xl border bg-card/50 p-4">
                      <dt className="text-xs font-medium text-muted-foreground">Reviews</dt>
                      <dd className="mt-1 text-2xl font-bold tracking-tight">8</dd>
                    </div>
                  </dl>
                </div>

                <div className="relative">
                  <div className="pointer-events-none absolute -inset-6 -z-10 rounded-3xl bg-gradient-to-br from-primary/20 via-transparent to-transparent blur-3xl" />
                  <Card className="overflow-hidden border-2 bg-card/50 backdrop-blur-sm sm:p-2">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-6">
                      <div className="grid gap-1">
                        <CardTitle className="text-xl font-bold">Advanced + AI</CardTitle>
                        <CardDescription className="text-xs">Power features for faster reviews and tighter ownership.</CardDescription>
                      </div>
                      <div className="inline-flex items-center gap-1.5 rounded-full border bg-primary/10 px-3 py-1 text-xs font-medium text-primary shadow-sm">
                        <SparklesIcon className="size-3.5 fill-current" />
                        <span>AI-assisted</span>
                      </div>
                    </CardHeader>
                    <CardContent className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                      {heroFeatures.map((feature) => {
                        const Icon = feature.icon;
                        return (
                          <div key={feature.title} className="group rounded-xl border bg-background/50 p-4 transition-all hover:bg-background hover:shadow-md">
                            <div className="flex items-start gap-3">
                              <div className="inline-flex size-9 items-center justify-center rounded-lg border bg-muted/50 transition-colors group-hover:border-primary/50 group-hover:bg-primary/5">
                                <Icon className="size-4.5 text-primary" />
                              </div>
                              <div className="grid gap-0.5">
                                <div className="text-sm font-semibold">{feature.title}</div>
                                <div className="text-[12px] leading-snug text-muted-foreground">{feature.description}</div>
                              </div>
                            </div>
                          </div>
                        );
                      })}
                    </CardContent>
                    <div className="mx-6 mb-6 mt-2 flex items-center justify-between rounded-xl border bg-muted/30 px-4 py-3">
                      <div className="flex items-center gap-2 text-xs font-medium text-muted-foreground">
                        <LineChartIcon className="size-4 text-primary" />
                        <span>Change-focused signal, not noise</span>
                      </div>
                      <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">Live view</span>
                    </div>
                  </Card>
                </div>
              </div>

              {/* Analytics Section */}
              <div id="analytics" className="relative">
                <div className="pointer-events-none absolute -inset-24 -z-10 bg-[radial-gradient(circle_at_center,var(--primary-foreground)_0%,transparent_70%)] opacity-20" />
                <Card className="border-none bg-card/50 shadow-xl backdrop-blur-md">
                  <CardHeader className="flex flex-col items-stretch space-y-0 border-b p-0 sm:flex-row">
                    <div className="flex flex-1 flex-col justify-center gap-1 px-6 py-5 sm:py-6">
                      <CardTitle className="text-xl font-bold">Risk Exposure Trends</CardTitle>
                      <CardDescription>
                        Visualizing risk severity distribution over the last 6 months.
                      </CardDescription>
                    </div>
                    <div className="flex border-t sm:border-t-0 sm:border-l">
                      {(["critical", "high", "medium", "low"] as const).map((key) => {
                        return (
                          <button
                            key={key}
                            data-active={activeMetric === key}
                            className="flex flex-1 flex-col justify-center gap-1 px-6 py-4 text-left transition-all hover:bg-muted/50 data-[active=true]:bg-muted/80 sm:px-8 sm:py-6 first:border-l-0 border-l relative"
                            onClick={() => setActiveMetric(key)}
                          >
                            <span 
                              className="absolute top-0 left-0 h-1 w-full transition-opacity" 
                              style={{ 
                                backgroundColor: chartConfig[key].color,
                                opacity: activeMetric === key ? 1 : 0 
                              }} 
                            />
                            <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground flex items-center gap-2">
                              <span 
                                className="size-2 rounded-full" 
                                style={{ backgroundColor: chartConfig[key].color }} 
                              />
                              {chartConfig[key].label}
                            </span>
                            <span className="text-xl font-bold leading-none sm:text-3xl">
                              {total[key]}
                            </span>
                          </button>
                        );
                      })}
                    </div>
                  </CardHeader>
                  <CardContent className="p-4 sm:p-6">
                    <ChartContainer
                      config={chartConfig}
                      className="aspect-auto h-[350px] w-full"
                    >
                      <AreaChart
                        data={chartData}
                        margin={{
                          left: 12,
                          right: 12,
                          top: 20,
                        }}
                      >
                        <defs>
                          {Object.entries(chartConfig).map(([key, config]) => (
                            <linearGradient key={key} id={`fill-${key}`} x1="0" y1="0" x2="0" y2="1">
                              <stop
                                offset="5%"
                                stopColor={config.color}
                                stopOpacity={0.3}
                              />
                              <stop
                                offset="95%"
                                stopColor={config.color}
                                stopOpacity={0}
                              />
                            </linearGradient>
                          ))}
                        </defs>
                        <CartesianGrid vertical={false} strokeDasharray="3 3" className="stroke-muted" />
                        <XAxis
                          dataKey="date"
                          tickLine={false}
                          axisLine={false}
                          tickMargin={12}
                          minTickGap={32}
                          tickFormatter={(value) => {
                            const date = new Date(value);
                            return date.toLocaleDateString("en-US", {
                              month: "short",
                            });
                          }}
                          className="text-xs font-medium text-muted-foreground"
                        />
                        <ChartTooltip
                          cursor={false}
                          content={
                            <ChartTooltipContent
                              labelFormatter={(value) => {
                                return new Date(value).toLocaleDateString("en-US", {
                                  month: "long",
                                  year: "numeric",
                                });
                              }}
                              indicator="dot"
                            />
                          }
                        />
                        <Area
                          dataKey={activeMetric}
                          type="monotone"
                          fill={`url(#fill-${activeMetric})`}
                          stroke={chartConfig[activeMetric].color}
                          strokeWidth={2}
                          stackId="a"
                        />
                      </AreaChart>
                    </ChartContainer>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="w-full border-t bg-background py-6">
        <div className="mx-auto flex max-w-7xl flex-col items-center justify-between gap-4 px-4 sm:flex-row sm:px-6 lg:px-8">
          <div className="flex items-center gap-2">
            <div className="size-6 rounded bg-primary" />
            <span className="text-sm font-bold tracking-tight">RiskRegister</span>
          </div>
          <p className="text-xs text-muted-foreground">
            © {new Date().getFullYear()} Risk Register. All rights reserved.
          </p>
          <div className="flex items-center gap-6">
            <Link to="/login" className="text-xs font-medium text-muted-foreground hover:text-foreground">
              Login
            </Link>
            <a href="#" className="text-xs font-medium text-muted-foreground hover:text-foreground">
              Privacy
            </a>
            <a href="#" className="text-xs font-medium text-muted-foreground hover:text-foreground">
              Terms
            </a>
          </div>
        </div>
      </footer>
    </div>
  );
}
