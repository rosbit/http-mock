package main

import (
	"net/http"
	"regexp"
	"text/template"
	"github.com/gernest/alien"
	"path"
	"mime"
)

var (
	_DEFAULT_METHODS = []string{http.MethodGet}
	_RE = regexp.MustCompile("[\\s,]+")

	_funcs = template.FuncMap{
		"pathVar": func(r *http.Request, n string)string {
			p := alien.GetParams(r)
			return p.Get(n)
		},
	}
)

type Endpoint struct {
	action *Action
}

func (e *Endpoint) GetMethods() []string {
	a := e.action
	if a.Method == "" {
		return _DEFAULT_METHODS
	}
	return _RE.Split(a.Method, -1)
}

func (e *Endpoint) GetUri() string {
	return e.action.Uri
}

func (e *Endpoint) GetHandler() http.HandlerFunc {
	a := e.action
	return func(w http.ResponseWriter, r *http.Request) {
		// redirect
		if a.RedirectBody != "" {
			http.Redirect(w, r, a.RedirectBody, http.StatusMovedPermanently)
			return
		}

		h := w.Header()
		h.Set("Content-Type", ServiceConf.DefaultContentType)
		setHeaders(a, h)
		setCookies(a, h)

		// file body
		if a.FileBody != "" {
			showFile(w, r, a.FileBody, a.Status)
			return
		}

		if a.Status != 0 {
			w.WriteHeader(a.Status)
		}

		// template
		if a.TmplBody != "" {
			execTemplate(w, r, a.TmplBody)
			return
		}

		// average body
		if a.Body != "" {
			w.Write([]byte(a.Body))
		}
	}
}

func CreateEndpoints() []*Endpoint {
	actions := ServiceConf.Actions
	es := make([]*Endpoint, len(ServiceConf.Actions))
	for i := range actions {
		es[i] = &Endpoint{&actions[i]}
	}
	return es
}

func setHeaders(a *Action, h http.Header) {
	if a.Headers == nil || len(a.Headers) == 0 {
		return
	}
	for k,v := range a.Headers {
		h.Set(k, v)
	}
}

func setCookies(a *Action, h http.Header) {
	if a.Cookies == nil || len(a.Cookies) == 0 {
		return
	}
	for n,v := range a.Cookies {
		c := http.Cookie{Name:n, Value:v}
		h.Add("Set-Cookie", c.String())
	}
}

func showFile(w http.ResponseWriter, r *http.Request, file string, status int) {
	ext := path.Ext(file)
	mimeType := mime.TypeByExtension(ext)
	w.Header().Set("Content-Type", mimeType)
	fn := path.Join(ServiceConf.Root, file)
	if status != 0 {
		w.WriteHeader(status)
	}
	http.ServeFile(w, r, fn)
}

func execTemplate(w http.ResponseWriter, r *http.Request, tmpl string) {
	t, err := template.New("webpage").Funcs(_funcs).Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, map[string]interface{}{"r": r})
}

