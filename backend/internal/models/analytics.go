package models

// TimeDataPoint represents a single data point for time-series charts
type TimeDataPoint struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

// StatusTimeDataPoint represents opened/closed counts for a period
type StatusTimeDataPoint struct {
	Period string `json:"period"`
	Open   int    `json:"open"`
	Closed int    `json:"closed"`
}

// AnalyticsResponse contains all analytics data for the frontend
type AnalyticsResponse struct {
	// Current State
	TotalRisks int               `json:"total_risks"`
	BySeverity map[string]int    `json:"by_severity"`
	ByStatus   map[string]int    `json:"by_status"`
	ByCategory []CategoryCount   `json:"by_category"`

	// Trends
	CreatedOverTime []TimeDataPoint       `json:"created_over_time"`
	StatusOverTime  []StatusTimeDataPoint `json:"status_over_time"`
}

// AnalyticsGranularity defines the time grouping for trend data
type AnalyticsGranularity string

const (
	GranularityMonthly AnalyticsGranularity = "monthly"
	GranularityWeekly  AnalyticsGranularity = "weekly"
)
