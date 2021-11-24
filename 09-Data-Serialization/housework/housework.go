package housework

type Chore struct {
	Complete    bool   `json:"complete"`
	Description string `json:"description"`
}
