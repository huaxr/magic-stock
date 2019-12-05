package models

type TicketGroupResult struct {
	Type  string `json:"type"`
	Total int    `json:"total"`
}

type EventGroupResult struct {
	Type  string `json:"type"`
	Total int    `json:"total"`
}
