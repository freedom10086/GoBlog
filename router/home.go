package router

import "net/http"

type HomeHandler struct {
	BaseHandler
}

func (h *HomeHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

//w http.ResponseWriter, data interface{}, res []string, tmpls ...string
func (*HomeHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	Template(w, &TemplateData{
		Css:  []string{"style.css","nav.css","footer.css"},
		Js:   nil,
		Data: nil, },
		"home", "nav", "footer")
}
