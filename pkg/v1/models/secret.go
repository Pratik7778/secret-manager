package models

type Secret struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateSecret struct {
	Value string `json:"value"`
}
