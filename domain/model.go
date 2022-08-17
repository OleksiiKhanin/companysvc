package domain

type Company struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
}

type FilterOptions struct {
	Limit  *int
	Params map[string]string
}

type EventType string

type Event struct {
	Type    EventType `json:"type"`
	Subject Company   `json:"subject"`
	OldName string    `json:"oldName"`
	OldCode string    `json:"oldCode"`
}
