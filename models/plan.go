package models

type Plan struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int    `json:"userid"`
	Steps       []Step `json:"steps"`
}

func (p Plan) CompletionStatus() int {
	result := 0
	if len(p.Steps) == 0 {
		return 0
	}
	incr := 100 / len(p.Steps)
	for _, s := range p.Steps {
		if s.Done {
			result += incr
		}
	}
	return result
}
