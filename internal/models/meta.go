package models

type Description struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

type META map[string]Description

func Desc(Type string, Labels map[string]string) Description {
	return Description{
		Type:   Type,
		Labels: Labels,
	}
}
