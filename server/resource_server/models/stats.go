package models

type UsageStatisticResponse struct {
	MetricType              string               `json:"metric_type"`
	UnitOfMeasurement       string               `json:"unit_of_measurement"`
	MonthToDate             MonthToDateStatistic `json:"month_to_date"`
	LastThirtyDaysStatistic Statistics           `json:"last_thirty_days_statistic"`
}

type MonthToDateStatistic struct {
	Month                     string  `json:"month"`
	MonthToDateUsage          float64 `json:"month_to_date_usage"`
	ForecastedEndOfMonthUsage float64 `json:"forecasted_end_of_month_usage"`
}
type Statistics struct {
	Total            float64                 `json:"total"`
	Average          float64                 `json:"average"`
	StatisticDetails []UsageStatisticPerDate `json:"details"`
}

type UsageStatisticPerDate struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}
