export type Granularity = 'monthly' | 'weekly';

export interface TimeDataPoint {
  period: string;
  count: number;
}

export interface StatusTimeDataPoint {
  period: string;
  open: number;
  closed: number;
}

export interface CategoryCount {
  category_id: string;
  category_name: string;
  count: number;
}

export interface AnalyticsResponse {
  // Current State
  total_risks: number;
  by_severity: Record<string, number>;
  by_status: Record<string, number>;
  by_category: CategoryCount[];

  // Trends
  created_over_time: TimeDataPoint[];
  status_over_time: StatusTimeDataPoint[];
}
