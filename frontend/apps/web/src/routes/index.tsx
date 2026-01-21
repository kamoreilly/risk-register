import { createFileRoute, Link } from "@tanstack/react-router";
import * as React from "react";
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

export const Route = createFileRoute("/")({
  component: HomeComponent,
});

function HomeComponent() {
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
