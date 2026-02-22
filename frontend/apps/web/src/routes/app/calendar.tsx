import { createFileRoute, Link } from "@tanstack/react-router";
import * as React from "react";

import { useUpcomingReviews, useOverdueReviews } from "@/hooks/useDashboard";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import type { ReviewRisk } from "@/types/dashboard";

export const Route = createFileRoute("/app/calendar")({
  component: ReviewCalendar,
});

// Helper functions for calendar
function daysInMonth(year: number, month: number): number {
  return new Date(year, month + 1, 0).getDate();
}

function firstDayOfMonth(year: number, month: number): number {
  return new Date(year, month, 1).getDay();
}

function formatDateKey(date: Date): string {
  return date.toISOString().split("T")[0];
}

function parseDate(dateString: string): Date {
  return new Date(dateString + "T00:00:00");
}

function getMarkerColor(reviewDate: Date, today: Date): string {
  const diffDays = Math.floor(
    (reviewDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24)
  );

  if (diffDays < 0) {
    return "bg-red-500"; // Overdue
  } else if (diffDays <= 7) {
    return "bg-yellow-500"; // Within 7 days
  }
  return "bg-blue-500"; // Later
}

const MONTHS = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

const DAYS = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

function ReviewCalendar() {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const [currentMonth, setCurrentMonth] = React.useState(today.getMonth());
  const [currentYear, setCurrentYear] = React.useState(today.getFullYear());

  const { data: upcomingData, isLoading: upcomingLoading } = useUpcomingReviews(30);
  const { data: overdueData, isLoading: overdueLoading } = useOverdueReviews();

  const isLoading = upcomingLoading || overdueLoading;

  // Combine all reviews
  const allReviews: ReviewRisk[] = [
    ...(upcomingData?.risks ?? []),
    ...(overdueData?.risks ?? []),
  ];

  // Group reviews by date
  const reviewsByDate = React.useMemo(() => {
    const grouped: Record<string, ReviewRisk[]> = {};
    for (const review of allReviews) {
      const dateKey = review.review_date.split("T")[0];
      if (!grouped[dateKey]) {
        grouped[dateKey] = [];
      }
      grouped[dateKey].push(review);
    }
    return grouped;
  }, [allReviews]);

  const goToPreviousMonth = () => {
    if (currentMonth === 0) {
      setCurrentMonth(11);
      setCurrentYear(currentYear - 1);
    } else {
      setCurrentMonth(currentMonth - 1);
    }
  };

  const goToNextMonth = () => {
    if (currentMonth === 11) {
      setCurrentMonth(0);
      setCurrentYear(currentYear + 1);
    } else {
      setCurrentMonth(currentMonth + 1);
    }
  };

  // Generate calendar days
  const totalDays = daysInMonth(currentYear, currentMonth);
  const startDay = firstDayOfMonth(currentYear, currentMonth);
  const calendarDays: (number | null)[] = [];

  // Add empty cells for days before the first day of the month
  for (let i = 0; i < startDay; i++) {
    calendarDays.push(null);
  }

  // Add the days of the month
  for (let day = 1; day <= totalDays; day++) {
    calendarDays.push(day);
  }

  // Get overdue reviews for sidebar
  const overdueReviews = overdueData?.risks ?? [];
  const upcomingReviewsList = upcomingData?.risks ?? [];

  return (
    <div className="p-8">
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Calendar Section */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
              <CardTitle className="text-lg">
                {MONTHS[currentMonth]} {currentYear}
              </CardTitle>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={goToPreviousMonth}>
                  &lt;
                </Button>
                <Button variant="outline" size="sm" onClick={goToNextMonth}>
                  &gt;
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <p className="text-muted-foreground text-center py-8">
                  Loading...
                </p>
              ) : (
                <div className="grid grid-cols-7 gap-1">
                  {/* Day headers */}
                  {DAYS.map((day) => (
                    <div
                      key={day}
                      className="text-center text-sm font-medium text-muted-foreground p-2"
                    >
                      {day}
                    </div>
                  ))}

                  {/* Calendar days */}
                  {calendarDays.map((day, index) => {
                    if (day === null) {
                      return <div key={`empty-${index}`} className="p-2" />;
                    }

                    const dateKey = formatDateKey(
                      new Date(currentYear, currentMonth, day)
                    );
                    const reviewsOnDay = reviewsByDate[dateKey] ?? [];
                    const isToday =
                      day === today.getDate() &&
                      currentMonth === today.getMonth() &&
                      currentYear === today.getFullYear();

                    return (
                      <div
                        key={day}
                        className={`min-h-[80px] border p-1 ${
                          isToday ? "bg-muted/50 border-primary" : "border-border"
                        }`}
                      >
                        <div
                          className={`text-sm ${
                            isToday ? "font-bold text-primary" : ""
                          }`}
                        >
                          {day}
                        </div>
                        <div className="mt-1 space-y-1">
                          {reviewsOnDay.slice(0, 3).map((review) => {
                            const reviewDate = parseDate(review.review_date);
                            const markerColor = getMarkerColor(reviewDate, today);

                            return (
                              <Link
                                key={review.id}
                                to="/app/risks/$id"
                                params={{ id: review.id }}
                                className="block"
                              >
                                <div
                                  className={`text-xs truncate px-1 py-0.5 rounded text-white ${markerColor} hover:opacity-80 transition-opacity`}
                                  title={review.title}
                                >
                                  {review.title}
                                </div>
                              </Link>
                            );
                          })}
                          {reviewsOnDay.length > 3 && (
                            <div className="text-xs text-muted-foreground px-1">
                              +{reviewsOnDay.length - 3} more
                            </div>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}

              {/* Legend */}
              <div className="flex gap-4 mt-4 pt-4 border-t">
                <div className="flex items-center gap-2">
                  <div className="size-3 rounded bg-red-500" />
                  <span className="text-xs text-muted-foreground">Overdue</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="size-3 rounded bg-yellow-500" />
                  <span className="text-xs text-muted-foreground">Within 7 days</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="size-3 rounded bg-blue-500" />
                  <span className="text-xs text-muted-foreground">Later</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Overdue Reviews */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg text-red-600">
                Overdue Reviews ({overdueReviews.length})
              </CardTitle>
            </CardHeader>
            <CardContent>
              {overdueLoading ? (
                <p className="text-muted-foreground text-sm">Loading...</p>
              ) : overdueReviews.length === 0 ? (
                <p className="text-muted-foreground text-sm">
                  No overdue reviews
                </p>
              ) : (
                <div className="space-y-2">
                  {overdueReviews.map((review) => (
                    <Link
                      key={review.id}
                      to="/app/risks/$id"
                      params={{ id: review.id }}
                      className="block p-2 rounded border border-red-200 bg-red-50 hover:bg-red-100 transition-colors"
                    >
                      <div className="font-medium text-sm truncate">
                        {review.title}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Due: {new Date(review.review_date).toLocaleDateString()}
                      </div>
                    </Link>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Upcoming Reviews */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">
                Upcoming Reviews ({upcomingReviewsList.length})
              </CardTitle>
            </CardHeader>
            <CardContent>
              {upcomingLoading ? (
                <p className="text-muted-foreground text-sm">Loading...</p>
              ) : upcomingReviewsList.length === 0 ? (
                <p className="text-muted-foreground text-sm">
                  No upcoming reviews
                </p>
              ) : (
                <div className="space-y-2 max-h-[300px] overflow-y-auto">
                  {upcomingReviewsList.map((review) => {
                    const reviewDate = parseDate(review.review_date);
                    const markerColor = getMarkerColor(reviewDate, today);
                    const borderColor =
                      markerColor === "bg-yellow-500"
                        ? "border-yellow-200 bg-yellow-50"
                        : "border-blue-200 bg-blue-50";

                    return (
                      <Link
                        key={review.id}
                        to="/app/risks/$id"
                        params={{ id: review.id }}
                        className={`block p-2 rounded border ${borderColor} hover:opacity-80 transition-opacity`}
                      >
                        <div className="font-medium text-sm truncate">
                          {review.title}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          Due: {new Date(review.review_date).toLocaleDateString()}
                        </div>
                      </Link>
                    );
                  })}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
