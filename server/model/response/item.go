package response

type Item struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Colloq      string   `json:"colloq"`
	From        []string `json:"from"`
	Image       string   `json:"image"`
	Gold        int      `json:"gold"`
	Depth       int      `json:"depth"`
}
