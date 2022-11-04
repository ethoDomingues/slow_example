package model

import "github.com/ethodomingues/slow"

type Solicitation struct {
	Model
	Req string
	Rec string
}

func (s *Solicitation) Accept() {}
func (s *Solicitation) Delete() {}

type SliceSol []*Solicitation

func (s *SliceSol) ToJSON(cUser string, rq *slow.Request) map[string]any {

	rec := []map[string]any{}
	req := []map[string]any{}
	idsRec := []string{}
	idsReq := []string{}

	for _, sol := range *s {
		if sol.Rec == cUser {
			u_, _ := FindByID(sol.Req)
			u := u_.(*User)
			idsReq = append(idsReq, sol.Req)
			req = append(req, u.ToJSON(rq))
		} else {
			u_, _ := FindByID(sol.Rec)
			u := u_.(*User)
			idsRec = append(idsRec, sol.Req)
			rec = append(rec, u.ToJSON(rq))
		}
	}
	return map[string]any{
		"rec": map[string]any{
			"users": rec,
			"ids":   idsRec,
		},
		"req": map[string]any{
			"users": req,
			"ids":   idsReq,
		},
	}
}
