export type MitigationStatus = 'planned' | 'in_progress' | 'completed' | 'cancelled';

export interface Mitigation {
  id: string;
  risk_id: string;
  description: string;
  owner?: string;
  status: MitigationStatus;
  due_date?: string;
  created_at: string;
  updated_at: string;
  created_by: string;
  updated_by: string;
}

export interface CreateMitigationInput {
  risk_id: string;
  description: string;
  owner?: string;
  status?: MitigationStatus;
  due_date?: string;
}

export interface UpdateMitigationInput {
  description?: string;
  owner?: string;
  status?: MitigationStatus;
  due_date?: string;
}
