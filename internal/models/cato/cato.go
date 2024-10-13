package cato

import m "lkserver/internal/models"

type Cato struct {
	Ref         m.JSONByte `json:"ref"`
	ParentRef   m.JSONByte `json:"parentRef"`
	Code        string     `json:"code"`
	K1          string     `json:"k1"`
	K2          string     `json:"k2"`
	K3          string     `json:"k3"`
	K4          string     `json:"k4"`
	K5          string     `json:"k5"`
	Description string     `json:"description"`
	Title       string     `json:"title"`
}
