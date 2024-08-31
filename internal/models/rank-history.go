package models

type RankHistory struct {
	Date       JSONTime    `json:"date"`
	Individual Individuals `json:"individual"`
	Rank       Rank        `json:"rank"`
}
