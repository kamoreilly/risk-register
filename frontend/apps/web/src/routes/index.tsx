import { createFileRoute, Link } from "@tanstack/react-router";
import * as React from "react";

import { HeaderNav } from "@/components/header-nav";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/")({
  component: HomeComponent,
});

function HomeComponent() {
  const risks = React.useMemo(
    () =>
      [
        {
          id: "R-014",
          title: "Vendor SLA breach impacts incident response",
          status: "Open" as const,
          severity: "High" as const,
          owner: "Ops",
          nextReview: "Feb 02",
          updated: "2h ago",
          mitigation: "Add escalation path + on-call runbook.",
          notes: "Awaiting revised SLA; interim monitoring in place.",
        },
        {
          id: "R-021",
          title: "Unreviewed access changes in non-prod",
          status: "Mitigating" as const,
          severity: "Medium" as const,
          owner: "Security",
          nextReview: "Jan 29",
          updated: "1d ago",
          mitigation: "Weekly access diff + approvals.",
          notes: "Backfill reviews for last 30 days.",
        },
        {
          id: "R-005",
          title: "Missing dependency update window",
          status: "Accepted" as const,
          severity: "Low" as const,
          owner: "Eng",
          nextReview: "Feb 15",
          updated: "3d ago",
          mitigation: "Track upgrade cadence per repo.",
          notes: "Accepted short-term; revisit after Q1 release.",
        },
        {
          id: "R-032",
          title: "Backups not tested against point-in-time restore",
          status: "Open" as const,
          severity: "Critical" as const,
          owner: "Platform",
          nextReview: "Jan 24",
          updated: "6h ago",
          mitigation: "Run restore drill; document RTO/RPO.",
          notes: "Restore run scheduled; add automated verification.",
        },
        {
          id: "R-018",
          title: "Legacy service lacks alert ownership",
          status: "Closed" as const,
          severity: "Medium" as const,
          owner: "SRE",
          nextReview: "—",
          updated: "2w ago",
          mitigation: "Assign owner; add paging policy.",
          notes: "Owner assigned; alerts routed via on-call rotation.",
        },
      ] as const,
    [],
  );

  const filters = React.useMemo(
    () =>
      [
        { key: "All", label: "All" },
        { key: "Open", label: "Open" },
        { key: "Mitigating", label: "Mitigating" },
        { key: "Accepted", label: "Accepted" },
        { key: "Closed", label: "Closed" },
      ] as const,
    [],
  );

  const [filter, setFilter] = React.useState<(typeof filters)[number]["key"]>("All");
  const visibleRisks = React.useMemo(() => {
    if (filter === "All") return risks;
    return risks.filter((r) => r.status === filter);
  }, [filter, risks]);

  const [selectedId, setSelectedId] = React.useState<string>(() => risks[0]?.id ?? "");
  const selected = React.useMemo(
    () => visibleRisks.find((r) => r.id === selectedId) ?? visibleRisks[0] ?? risks[0],
    [risks, selectedId, visibleRisks],
  );

  const [autoRotate, setAutoRotate] = React.useState(true);
  const resumeTimerRef = React.useRef<number | null>(null);

  const statusCounts = React.useMemo(() => {
    const counts: Record<(typeof risks)[number]["status"], number> = {
      Open: 0,
      Mitigating: 0,
      Accepted: 0,
      Closed: 0,
    };
    for (const r of visibleRisks) counts[r.status] += 1;
    return counts;
  }, [visibleRisks]);

  const severityCounts = React.useMemo(() => {
    const counts: Record<(typeof risks)[number]["severity"], number> = {
      Critical: 0,
      High: 0,
      Medium: 0,
      Low: 0,
    };
    for (const r of visibleRisks) counts[r.severity] += 1;
    return counts;
  }, [visibleRisks]);

  const miniTrend = React.useMemo(() => {
    const seed = (selected?.id ?? "R").split("").reduce((acc, ch) => acc + ch.charCodeAt(0), 0);
    return Array.from({ length: 10 }, (_, i) => ((seed * (i + 3) + i * 7) % 10) + 1);
  }, [selected?.id]);

  const upcoming = React.useMemo(() => {
    const list = visibleRisks.filter((r) => r.nextReview !== "—").slice(0, 3);
    return list;
  }, [visibleRisks]);

  function statusPillClasses(status: (typeof risks)[number]["status"]) {
    return cn(
      "inline-flex h-6 items-center rounded-md border px-2 text-[11px]",
      status === "Open" && "border-rose-200 bg-rose-50 text-rose-700 dark:border-rose-950 dark:bg-rose-950/30 dark:text-rose-200",
      status === "Mitigating" && "border-amber-200 bg-amber-50 text-amber-800 dark:border-amber-950 dark:bg-amber-950/30 dark:text-amber-200",
      status === "Accepted" && "border-slate-200 bg-slate-50 text-slate-800 dark:border-slate-800 dark:bg-slate-900/40 dark:text-slate-200",
      status === "Closed" && "border-emerald-200 bg-emerald-50 text-emerald-800 dark:border-emerald-950 dark:bg-emerald-950/30 dark:text-emerald-200",
    );
  }

  React.useEffect(() => {
    if (!autoRotate) return;
    if (!visibleRisks.length) return;

    const interval = window.setInterval(() => {
      const currentIndex = Math.max(
        0,
        visibleRisks.findIndex((r) => r.id === selected?.id),
      );
      const nextIndex = (currentIndex + 1) % visibleRisks.length;
      setSelectedId(visibleRisks[nextIndex].id);
    }, 3500);

    return () => window.clearInterval(interval);
  }, [autoRotate, selected?.id, visibleRisks]);

  React.useEffect(() => {
    return () => {
      if (resumeTimerRef.current !== null) {
        window.clearTimeout(resumeTimerRef.current);
        resumeTimerRef.current = null;
      }
    };
  }, []);

  React.useEffect(() => {
    // Keep selection valid when filtering.
    if (!visibleRisks.length) return;
    if (!visibleRisks.some((r) => r.id === selectedId)) {
      setSelectedId(visibleRisks[0].id);
    }
  }, [selectedId, visibleRisks]);

  function onPickRisk(id: string) {
    setSelectedId(id);
    setAutoRotate(false);

    if (resumeTimerRef.current !== null) {
      window.clearTimeout(resumeTimerRef.current);
    }
    resumeTimerRef.current = window.setTimeout(() => {
      setAutoRotate(true);
      resumeTimerRef.current = null;
    }, 10_000);
  }

  return (
    <div className="min-h-svh bg-background">
      <HeaderNav />
      <main className="mx-auto w-full max-w-5xl px-4 pb-12 pt-8 sm:pt-12">
        <section className="grid gap-8">
          <div className="mx-auto grid w-full max-w-2xl justify-items-center gap-4 text-center sm:justify-items-start sm:text-left">
            <div className="inline-flex w-fit items-center gap-2 rounded-full border bg-muted/40 px-3 py-1 text-xs text-muted-foreground">
              <span className={cn("size-1.5 rounded-full", autoRotate ? "bg-emerald-500" : "bg-amber-500")} />
              <span>{autoRotate ? "Live view (auto)" : "Pinned (tap a row)"}</span>
            </div>

            <h1 className="text-balance text-3xl font-semibold tracking-tight sm:text-4xl">
              Risk register, focused on what changes
            </h1>
            <p className="text-pretty text-sm leading-6 text-muted-foreground sm:text-base">
              Capture owners, severity, and review cadence in one view.
            </p>

            <div className="flex flex-wrap justify-center gap-2 md:justify-start">
              <Link to="/login">
                <Button size="sm">Open the app</Button>
              </Link>
              <a
                href="#screen"
                className="hover:bg-muted inline-flex h-9 items-center rounded-md border px-3 text-sm font-medium"
              >
                Jump to preview
              </a>
            </div>

            <dl className="grid w-full max-w-sm grid-cols-3 gap-3 pt-2">
              <div className="rounded-lg border bg-card px-3 py-2">
                <dt className="text-xs text-muted-foreground">Open</dt>
                <dd className="text-sm font-semibold">
                  {risks.filter((r) => r.status === "Open").length}
                </dd>
              </div>
              <div className="rounded-lg border bg-card px-3 py-2">
                <dt className="text-xs text-muted-foreground">Mitigating</dt>
                <dd className="text-sm font-semibold">
                  {risks.filter((r) => r.status === "Mitigating").length}
                </dd>
              </div>
              <div className="rounded-lg border bg-card px-3 py-2">
                <dt className="text-xs text-muted-foreground">Next review</dt>
                <dd className="text-sm font-semibold">{selected?.nextReview ?? "—"}</dd>
              </div>
            </dl>
          </div>

          <div id="screen" className="relative w-full">
            <div className="pointer-events-none absolute -inset-8 -z-10 bg-gradient-to-b from-primary/10 via-transparent to-transparent blur-2xl" />
            <div className="rounded-2xl border bg-card shadow-sm">
              <div className="flex items-center justify-between gap-3 border-b px-4 py-3">
                <div className="min-w-0">
                  <div className="truncate text-sm font-semibold">Workspace</div>
                  <div className="truncate text-xs text-muted-foreground">Risk register · {filter}</div>
                </div>
                <div className="flex items-center gap-2">
                  <div className="hidden sm:flex items-center gap-1 rounded-md border bg-background px-2 py-1 text-xs text-muted-foreground">
                    <span className="text-foreground">⌘K</span>
                    <span>Search</span>
                  </div>
                  <span className="inline-flex items-center rounded-md border bg-muted px-2 py-1 text-xs">
                    {visibleRisks.length} items
                  </span>
                </div>
              </div>

              <div className="grid gap-4 p-4">
                <div className="flex flex-wrap gap-2">
                  {filters.map((f) => {
                    const active = f.key === filter;
                    return (
                      <button
                        key={f.key}
                        type="button"
                        aria-pressed={active}
                        onClick={() => setFilter(f.key)}
                        className={cn(
                          "h-8 rounded-md border px-2.5 text-xs font-medium",
                          active ? "bg-primary text-primary-foreground border-primary" : "hover:bg-muted",
                        )}
                      >
                        {f.label}
                      </button>
                    );
                  })}
                </div>

                <div className="grid gap-3">
                  <div className="rounded-xl border bg-background">
                    <div className="grid grid-cols-[0.9fr_2.1fr] gap-2 border-b px-3 py-2 text-[11px] font-medium text-muted-foreground sm:grid-cols-[0.8fr_2.2fr_0.9fr_0.9fr]">
                      <div>ID</div>
                      <div>Risk</div>
                      <div className="hidden sm:block">Status</div>
                      <div className="hidden text-right sm:block">Review</div>
                    </div>

                    <div className="divide-y">
                      {visibleRisks.slice(0, 4).map((r) => {
                        const active = r.id === selected?.id;
                        return (
                          <button
                            key={r.id}
                            type="button"
                            aria-current={active}
                            onClick={() => onPickRisk(r.id)}
                            className={cn(
                              "grid w-full grid-cols-[0.9fr_2.1fr] gap-2 px-3 py-2 text-left sm:grid-cols-[0.8fr_2.2fr_0.9fr_0.9fr]",
                              "hover:bg-muted/50",
                              active && "bg-muted/60",
                            )}
                          >
                            <div className="text-xs font-mono text-muted-foreground">{r.id}</div>
                            <div className="min-w-0">
                              <div className="truncate text-xs font-medium">{r.title}</div>
                              <div className="truncate text-[11px] text-muted-foreground">
                                {r.owner} · {r.severity} · updated {r.updated}
                                <span className="sm:hidden"> · {r.status}</span>
                              </div>
                            </div>
                            <div className="hidden items-center sm:flex">
                              <span
                                className={statusPillClasses(r.status)}
                              >
                                {r.status}
                              </span>
                            </div>
                            <div className="hidden text-right text-xs text-muted-foreground sm:block">{r.nextReview}</div>
                          </button>
                        );
                      })}
                    </div>
                  </div>

                  <div className="overflow-hidden rounded-xl border bg-background">
                    <div className="border-b px-3 py-2">
                      <div className="flex min-w-0 items-center gap-2">
                        <div className="min-w-0 flex-1">
                          <div className="truncate text-xs font-semibold">{selected?.title ?? "—"}</div>
                          <div className="truncate text-[11px] text-muted-foreground">
                            {selected?.id ?? ""} · {selected?.owner ?? ""}
                          </div>
                        </div>
                        <div className="hidden shrink-0 sm:block text-[11px] text-muted-foreground">Dashboard</div>
                      </div>
                    </div>
                    <div className="grid gap-3 p-3">
                      <div className="grid grid-cols-2 gap-2">
                        <div className="rounded-lg border bg-card px-2.5 py-2">
                          <div className="text-[11px] text-muted-foreground">Selected</div>
                          <div className="flex flex-wrap items-center gap-2 pt-1">
                            <span className={statusPillClasses(selected?.status ?? "Open")}>{selected?.status ?? "—"}</span>
                            <span className="inline-flex h-6 items-center rounded-md border bg-background px-2 text-[11px]">
                              {selected?.severity ?? "—"}
                            </span>
                            <span className="inline-flex h-6 items-center rounded-md border bg-background px-2 text-[11px] text-muted-foreground">
                              review {selected?.nextReview ?? "—"}
                            </span>
                          </div>
                        </div>

                        <div className="rounded-lg border bg-card px-2.5 py-2">
                          <div className="flex items-center justify-between">
                            <div className="text-[11px] text-muted-foreground">Activity</div>
                            <div className="text-[11px] text-muted-foreground">7d</div>
                          </div>
                          <div className="mt-2 flex h-10 items-end gap-1">
                            {miniTrend.map((v, i) => (
                              <div
                                // eslint-disable-next-line react/no-array-index-key
                                key={i}
                                className={cn(
                                  "w-full rounded-sm",
                                  i === miniTrend.length - 1 ? "bg-primary" : "bg-muted-foreground/30",
                                )}
                                style={{ height: `${v * 8}%` }}
                              />
                            ))}
                          </div>
                        </div>
                      </div>

                      <div className="grid grid-cols-2 gap-2">
                        <div className="rounded-lg border bg-card px-2.5 py-2">
                          <div className="text-[11px] text-muted-foreground">Status</div>
                          <div className="mt-2 grid gap-1">
                            {(
                              [
                                { k: "Open" as const, color: "bg-rose-500" },
                                { k: "Mitigating" as const, color: "bg-amber-500" },
                                { k: "Accepted" as const, color: "bg-slate-500" },
                                { k: "Closed" as const, color: "bg-emerald-500" },
                              ] as const
                            ).map((row) => {
                              const total = Math.max(1, visibleRisks.length);
                              const pct = Math.round((statusCounts[row.k] / total) * 100);
                              return (
                                <div key={row.k} className="grid gap-1">
                                  <div className="flex items-center justify-between text-[11px]">
                                    <span className="text-muted-foreground">{row.k}</span>
                                    <span className="text-muted-foreground">{statusCounts[row.k]}</span>
                                  </div>
                                  <div className="h-1.5 overflow-hidden rounded-full bg-muted">
                                    <div className={cn("h-full", row.color)} style={{ width: `${pct}%` }} />
                                  </div>
                                </div>
                              );
                            })}
                          </div>
                        </div>

                        <div className="rounded-lg border bg-card px-2.5 py-2">
                          <div className="text-[11px] text-muted-foreground">Severity</div>
                          <div className="mt-2 grid gap-2">
                            <div className="flex h-2 w-full overflow-hidden rounded-full bg-muted">
                              {(
                                [
                                  { k: "Critical" as const, color: "bg-rose-600" },
                                  { k: "High" as const, color: "bg-amber-600" },
                                  { k: "Medium" as const, color: "bg-sky-600" },
                                  { k: "Low" as const, color: "bg-emerald-600" },
                                ] as const
                              ).map((seg) => {
                                const total = Math.max(1, visibleRisks.length);
                                const pct = (severityCounts[seg.k] / total) * 100;
                                return (
                                  <div
                                    key={seg.k}
                                    className={seg.color}
                                    style={{ width: `${pct}%` }}
                                    aria-hidden
                                  />
                                );
                              })}
                            </div>

                            <div className="grid gap-1">
                              {(
                                [
                                  { k: "Critical" as const },
                                  { k: "High" as const },
                                  { k: "Medium" as const },
                                  { k: "Low" as const },
                                ] as const
                              ).map((row) => (
                                <div key={row.k} className="flex items-center justify-between text-[11px]">
                                  <span className="text-muted-foreground">{row.k}</span>
                                  <span className="text-muted-foreground">{severityCounts[row.k]}</span>
                                </div>
                              ))}
                            </div>
                          </div>
                        </div>
                      </div>

                      <div className="rounded-lg border bg-card px-2.5 py-2">
                        <div className="flex items-center justify-between">
                          <div className="text-[11px] text-muted-foreground">Upcoming reviews</div>
                          <div className="text-[11px] text-muted-foreground">next</div>
                        </div>
                        <div className="mt-2 grid gap-1.5">
                          {upcoming.length ? (
                            upcoming.map((r) => (
                              <div key={r.id} className="flex min-w-0 items-center justify-between gap-2">
                                <div className="min-w-0 truncate text-xs">
                                  <span className="font-mono text-muted-foreground">{r.id}</span>
                                  <span className="text-muted-foreground"> · </span>
                                  <span className="text-muted-foreground">{r.owner}</span>
                                </div>
                                <div className="shrink-0 text-xs text-muted-foreground">{r.nextReview}</div>
                              </div>
                            ))
                          ) : (
                            <div className="text-xs text-muted-foreground">No upcoming items in this view.</div>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
