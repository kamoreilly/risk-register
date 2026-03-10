export interface Framework {
  id: string;
  name: string;
  description?: string;
  created_at: string;
}

export interface FrameworkControl {
  id: string;
  framework_id: string;
  framework_name: string;
  control_ref: string;
  title: string;
  description?: string;
  created_at: string;
  updated_at: string;
  linked_risk_count: number;
}

export interface ControlLinkedRisk {
  id: string;
  title: string;
  status: string;
  severity: string;
  category_name?: string;
  owner_name?: string;
  updated_at: string;
}

export interface RiskFrameworkControl {
  id: string;
  risk_id: string;
  framework_control_id: string;
  framework_id: string;
  framework_name: string;
  control_ref: string;
  control_title: string;
  control_description?: string;
  notes?: string;
  created_at: string;
  created_by: string;
}

export interface LinkControlInput {
  framework_control_id: string;
  notes?: string;
}

export interface CreateFrameworkControlInput {
  framework_id: string;
  control_ref: string;
  title: string;
  description?: string;
}

export interface UpdateFrameworkControlInput {
  control_ref?: string;
  title?: string;
  description?: string;
}

export interface CreateFrameworkInput {
  name: string;
  description?: string;
}

export interface UpdateFrameworkInput {
  name?: string;
  description?: string;
}
