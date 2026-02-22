package models

import "time"

// CategoryCount represents the count of risks per category
type CategoryCount struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Count        int    `json:"count"`
}

// DashboardSummaryResponse represents the dashboard summary data
type DashboardSummaryResponse struct {
	TotalRisks     int             `json:"total_risks"`
	ByStatus       map[string]int  `json:"by_status"`
	BySeverity     map[string]int  `json:"by_severity"`
	ByCategory     []CategoryCount `json:"by_category"`
	OverdueReviews int             `json:"overdue_reviews"`
}

// ReviewRisk represents a risk with review date information
type ReviewRisk struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	ReviewDate time.Time `json:"review_date"`
	Severity   string    `json:"severity"`
	Status     string    `json:"status"`
}

// ReviewListResponse represents the response for review endpoints
type ReviewListResponse struct {
	Risks []ReviewRisk `json:"risks"`
}
