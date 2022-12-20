package models

type Solicitation struct {
	Model
	Req string
	Rec string
}

func (s *Solicitation) ToMap() map[string]any {
	return map[string]any{
		"receiver":  s.Rec,
		"requester": s.Req,
	}
}
