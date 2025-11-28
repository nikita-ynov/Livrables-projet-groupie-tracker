package pages

import (
	"html/template"
)

var Temp *template.Template

func Init() {
	funcMap := template.FuncMap{
		"seq": func(start, end int) []int {
			s := make([]int, end-start+1)
			for i := range s {
				s[i] = start + i
			}
			return s
		},
	}
	Temp = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
	template.Must(Temp.ParseGlob("pages/*.html"))
}
