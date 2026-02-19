export type RiskStatus = 'open' | 'mitigating' | 'resolved' | 'accepted';
export type RiskSeverity = 'low' | 'medium' | 'high' | 'critical';

export interface Category {
  id: string;
  name: string;
  description?: string;
  created_at: string;
}

export interface Risk {
  id: string;
  title: string;
  description?: string;
  owner_id: string;
  owner?: {
    id: string;
    name: string;
    email: string;
  };
  status: RiskStatus;
  severity: RiskSeverity;
  category_id?: string;
  category?: Category;
  review_date?: string;
  created_at: string;
  updated_at: string;
  created_by: string;
  updated_by: string;
}

export interface CreateRiskInput {
  title: string;
  description?: string;
  owner_id: string;
  status?: RiskStatus;
  severity?: RiskSeverity;
  category_id?: string;
  review_date?: string;
}

export interface UpdateRiskInput {
  title?: string;
  description?: string;
  owner_id?: string;
  status?: RiskStatus;
  severity?: RiskSeverity;
  category_id?: string;
  review_date?: string;
}

export interface RiskListParams {
  status?: RiskStatus;
  severity?: RiskSeverity;
  category_id?: string;
  owner_id?: string;
  search?: string;
  sort?: string;
  order?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}

export interface RiskListResponse {
  data: Risk[];
  meta: {
    page: number;
    limit: number;
    total: number;
  };
}
