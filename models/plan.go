package models

type Plan struct {
	ID          int `query:"id" param:"id" json:"id"`
	Name        string
	Description string
	Steps       []Step
}

func (p Plan) CompletionStatus() int {
	result := 0
	incr := 100 / len(p.Steps)
	for _, s := range p.Steps {
		if s.Done {
			result += incr
		}
	}
	return result
}
