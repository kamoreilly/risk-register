export interface Framework {
  id: string;
  name: string;
  description?: string;
  created_at: string;
}

export interface RiskFrameworkControl {
  id: string;
  risk_id: string;
  framework_id: string;
  framework_name: string;
  control_ref: string;
  notes?: string;
  created_at: string;
  created_by: string;
}

export interface LinkControlInput {
  framework_id: string;
  control_ref: string;
  notes?: string;
}

export interface CreateFrameworkInput {
  name: string;
  description?: string;
}

export interface UpdateFrameworkInput {
  name?: string;
  description?: string;
}
