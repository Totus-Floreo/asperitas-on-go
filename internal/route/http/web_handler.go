package route

import (
	"html/template"
	"net/http"
)

func WebHandler(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseGlob("./web/template/*.html"))
	err := template.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}
