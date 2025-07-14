package pagination

type PageRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
