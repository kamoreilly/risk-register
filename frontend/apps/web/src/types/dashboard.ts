export interface CategoryCount {
  category_id: string;
  category_name: string;
  count: number;
}

export interface DashboardSummary {
  total_risks: number;
  by_status: Record<string, number>;
  by_severity: Record<string, number>;
  by_category: CategoryCount[];
  overdue_reviews: number;
}

export interface ReviewRisk {
  id: string;
  title: string;
  review_date: string;
  severity: string;
  status: string;
}

export interface ReviewListResponse {
  risks: ReviewRisk[];
}
