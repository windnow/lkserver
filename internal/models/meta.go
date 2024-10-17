package models

type Description struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

type META map[string]Description
