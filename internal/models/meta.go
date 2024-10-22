package models

type Description struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
	Order  int               `json:"order"`
}

type META map[string]Description

func Desc(Type string, Labels map[string]string, order int) Description {
	return Description{
		Type:   Type,
		Labels: Labels,
		Order:  order,
	}
}
