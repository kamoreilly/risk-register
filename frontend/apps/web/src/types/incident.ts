export type IncidentStatus = 'new' | 'acknowledged' | 'in_progress' | 'on_hold' | 'resolved' | 'closed';
export type IncidentPriority = 'p1' | 'p2' | 'p3' | 'p4';

export interface IncidentCategory {
  id: string;
  name: string;
  description?: string;
  created_at: string;
}

export interface Incident {
  id: string;
  title: string;
  description?: string;
  category_id?: string;
  category?: IncidentCategory;
  priority: IncidentPriority;
  status: IncidentStatus;
  assignee_id?: string;
  assignee?: {
    id: string;
    name: string;
    email: string;
  };
  reporter_id: string;
  reporter?: {
    id: string;
    name: string;
    email: string;
  };
  service_affected?: string;
  root_cause?: string;
  resolution_notes?: string;
  occurred_at: string;
  detected_at: string;
  resolved_at?: string;
  created_at: string;
  updated_at: string;
  created_by: string;
  updated_by: string;
}

export interface CreateIncidentInput {
  title: string;
  description?: string;
  category_id?: string;
  priority?: IncidentPriority;
  status?: IncidentStatus;
  assignee_id?: string;
  service_affected?: string;
  occurred_at?: string;
  detected_at?: string;
}

export interface UpdateIncidentInput {
  title?: string;
  description?: string;
  category_id?: string;
  priority?: IncidentPriority;
  status?: IncidentStatus;
  assignee_id?: string;
  service_affected?: string;
  root_cause?: string;
  resolution_notes?: string;
  resolved_at?: string;
}

export interface IncidentListParams {
  status?: IncidentStatus;
  priority?: IncidentPriority;
  category_id?: string;
  assignee_id?: string;
  search?: string;
  sort?: string;
  order?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}

export interface IncidentListResponse {
  data: Incident[];
  meta: {
    page: number;
    limit: number;
    total: number;
  };
}

export interface IncidentRisk {
  id: string;
  incident_id: string;
  risk_id: string;
  risk?: {
    id: string;
    title: string;
    status: string;
    severity: string;
  };
  created_at: string;
  created_by: string;
}

export interface LinkIncidentRiskInput {
  risk_id: string;
}
