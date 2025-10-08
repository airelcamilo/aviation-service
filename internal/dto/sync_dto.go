package dto

type SyncResponse struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Failed  int `json:"failed"`
	Error   int `json:"error"`
}
