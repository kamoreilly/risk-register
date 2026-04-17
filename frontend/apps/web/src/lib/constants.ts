import type { RiskStatus, RiskSeverity } from '@/types/risk';
import type { IncidentStatus, IncidentPriority } from '@/types/incident';

export const RISK_STATUS_COLORS: Record<RiskStatus, string> = {
  open: 'bg-yellow-100 text-yellow-800',
  mitigating: 'bg-blue-100 text-blue-800',
  resolved: 'bg-green-100 text-green-800',
  accepted: 'bg-gray-100 text-gray-800',
};

export const RISK_SEVERITY_COLORS: Record<RiskSeverity, string> = {
  low: 'bg-gray-100 text-gray-800',
  medium: 'bg-yellow-100 text-yellow-800',
  high: 'bg-orange-100 text-orange-800',
  critical: 'bg-red-100 text-red-800',
};

export const INCIDENT_STATUS_COLORS: Record<IncidentStatus, string> = {
  new: 'bg-blue-100 text-blue-800',
  acknowledged: 'bg-indigo-100 text-indigo-800',
  in_progress: 'bg-yellow-100 text-yellow-800',
  on_hold: 'bg-orange-100 text-orange-800',
  resolved: 'bg-green-100 text-green-800',
  closed: 'bg-gray-100 text-gray-800',
};

export const INCIDENT_PRIORITY_COLORS: Record<IncidentPriority, string> = {
  p1: 'bg-red-100 text-red-800',
  p2: 'bg-orange-100 text-orange-800',
  p3: 'bg-yellow-100 text-yellow-800',
  p4: 'bg-blue-100 text-blue-800',
};