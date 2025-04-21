package objects

type PvzQuery struct {
	StartDate string `schema:"startDate"`
	EndDate   string `schema:"endDate"`
	Page      int    `schema:"page"`
	Limit     int    `schema:"limit"`
}
