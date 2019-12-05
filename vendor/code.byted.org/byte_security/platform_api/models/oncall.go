package models

type EnabledOCOpt struct {
	Ids     []int `json:"ids"`
	Enabled int   `json:"enabled"`
}

type DeleteOpt struct {
	Ids []int `json:"ids"`
}
