export interface SummarizeRequest {
  title: string;
  description?: string;
  severity: string;
  status: string;
}

export interface SummarizeResponse {
  summary: string;
}

export interface DraftMitigationRequest {
  risk_title: string;
  risk_description?: string;
  severity: string;
}

export interface DraftMitigationResponse {
  draft: string;
}
