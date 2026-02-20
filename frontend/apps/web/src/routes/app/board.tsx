import { createFileRoute, Link } from "@tanstack/react-router";
import * as React from "react";
import {
  DndContext,
  DragOverlay,
  closestCenter,
  PointerSensor,
  useSensor,
  useSensors,
  useDroppable,
  useDraggable,
} from "@dnd-kit/core";
import type { DragEndEvent, DragStartEvent } from "@dnd-kit/core";

import { useRisks, useUpdateRiskStatus } from "@/hooks/useRisks";
import { cn } from "@/lib/utils";
import type { Risk, RiskStatus, RiskSeverity } from "@/types/risk";

export const Route = createFileRoute("/app/board")({
  component: Board,
});

const STATUS_COLUMNS: { id: RiskStatus; label: string }[] = [
  { id: "open", label: "Open" },
  { id: "mitigating", label: "Mitigating" },
  { id: "resolved", label: "Resolved" },
  { id: "accepted", label: "Accepted" },
];

const SEVERITY_COLORS: Record<RiskSeverity, string> = {
  low: "bg-gray-100 text-gray-800 border-gray-200",
  medium: "bg-yellow-100 text-yellow-800 border-yellow-200",
  high: "bg-orange-100 text-orange-800 border-orange-200",
  critical: "bg-red-100 text-red-800 border-red-200",
};

const COLUMN_COLORS: Record<RiskStatus, string> = {
  open: "border-t-yellow-400",
  mitigating: "border-t-blue-400",
  resolved: "border-t-green-400",
  accepted: "border-t-gray-400",
};

function Board() {
  const { data, isLoading } = useRisks({ limit: 100 });
  const updateRiskStatus = useUpdateRiskStatus();
  const [activeRisk, setActiveRisk] = React.useState<Risk | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 5,
      },
    })
  );

  const risks = data?.data ?? [];

  const risksByStatus = React.useMemo(() => {
    const grouped: Record<RiskStatus, Risk[]> = {
      open: [],
      mitigating: [],
      resolved: [],
      accepted: [],
    };

    for (const risk of risks) {
      grouped[risk.status].push(risk);
    }

    return grouped;
  }, [risks]);

  const handleDragStart = (event: DragStartEvent) => {
    const { active } = event;
    const risk = risks.find((r) => r.id === active.id);
    if (risk) {
      setActiveRisk(risk);
    }
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    setActiveRisk(null);

    if (!over) return;

    const riskId = active.id as string;
    const newStatus = over.id as RiskStatus;

    const currentRisk = risks.find((r) => r.id === riskId);
    if (currentRisk && currentRisk.status !== newStatus) {
      updateRiskStatus.mutate({ id: riskId, status: newStatus });
    }
  };

  if (isLoading) {
    return (
      <div className="p-8">
        <div className="text-center text-muted-foreground">Loading...</div>
      </div>
    );
  }

  return (
    <div className="p-8 h-full">
      <div className="mb-6">
        <h1 className="text-2xl font-bold">Risk Board</h1>
        <p className="text-muted-foreground">
          Drag and drop risks to update their status
        </p>
      </div>

      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
      >
        <div className="grid grid-cols-4 gap-4 h-[calc(100vh-200px)]">
          {STATUS_COLUMNS.map((column) => (
            <Column
              key={column.id}
              id={column.id}
              label={column.label}
              risks={risksByStatus[column.id]}
              activeRiskId={activeRisk?.id}
            />
          ))}
        </div>

        <DragOverlay>
          {activeRisk ? <DragOverlayCard risk={activeRisk} /> : null}
        </DragOverlay>
      </DndContext>
    </div>
  );
}

interface ColumnProps {
  id: RiskStatus;
  label: string;
  risks: Risk[];
  activeRiskId?: string;
}

function Column({ id, label, risks, activeRiskId }: ColumnProps) {
  const { setNodeRef, isOver } = useDroppable({
    id,
  });

  return (
    <div
      ref={setNodeRef}
      className={cn(
        "flex flex-col bg-muted/30 rounded-lg border-t-4 overflow-hidden transition-colors",
        COLUMN_COLORS[id],
        isOver && "bg-muted/50"
      )}
    >
      <div className="p-4 border-b bg-muted/50">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold">{label}</h3>
          <span className="text-sm text-muted-foreground bg-background px-2 py-0.5 rounded-full">
            {risks.length}
          </span>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto p-2 space-y-2">
        {risks.length === 0 ? (
          <div className="text-center text-muted-foreground text-sm py-8">
            No risks
          </div>
        ) : (
          risks.map(
            (risk) =>
              risk.id !== activeRiskId && (
                <RiskCard key={risk.id} risk={risk} />
              )
          )
        )}
      </div>
    </div>
  );
}

interface RiskCardProps {
  risk: Risk;
}

function RiskCard({ risk }: RiskCardProps) {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: risk.id,
  });

  const style = transform
    ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
      }
    : undefined;

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      className={cn(isDragging && "opacity-50")}
    >
      <Link
        to="/app/risks/$id"
        params={{ id: risk.id }}
        className="block"
        onClick={(e) => e.preventDefault()}
      >
        <div className="bg-background border rounded-lg p-3 shadow-sm hover:shadow-md transition-shadow cursor-grab active:cursor-grabbing">
          <div className="flex items-start justify-between gap-2 mb-2">
            <h4 className="font-medium text-sm line-clamp-2">{risk.title}</h4>
          </div>
          <div className="flex items-center justify-between gap-2">
            <span
              className={cn(
                "px-2 py-0.5 rounded text-xs font-medium border",
                SEVERITY_COLORS[risk.severity]
              )}
            >
              {risk.severity}
            </span>
            {risk.owner && (
              <span className="text-xs text-muted-foreground truncate max-w-[100px]">
                {risk.owner.name}
              </span>
            )}
          </div>
        </div>
      </Link>
    </div>
  );
}

interface DragOverlayCardProps {
  risk: Risk;
}

function DragOverlayCard({ risk }: DragOverlayCardProps) {
  return (
    <div className="bg-background border rounded-lg p-3 shadow-lg cursor-grabbing rotate-2">
      <div className="flex items-start justify-between gap-2 mb-2">
        <h4 className="font-medium text-sm line-clamp-2">{risk.title}</h4>
      </div>
      <div className="flex items-center justify-between gap-2">
        <span
          className={cn(
            "px-2 py-0.5 rounded text-xs font-medium border",
            SEVERITY_COLORS[risk.severity]
          )}
        >
          {risk.severity}
        </span>
        {risk.owner && (
          <span className="text-xs text-muted-foreground truncate max-w-[100px]">
            {risk.owner.name}
          </span>
        )}
      </div>
    </div>
  );
}
