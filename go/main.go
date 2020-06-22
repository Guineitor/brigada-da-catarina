package main

import (
    "html/template"
    "net/http"
)

func main() {
	
	index := template.Must(template.ParseFiles("template/index.html"))
	manifesto := template.Must(template.ParseFiles("template/manifesto.html"))
	blog := template.Must(template.ParseFiles("template/blog.html"))
	_404 := template.Must(template.ParseFiles("template/404.html"))


	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))



	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
            index.Execute(w, nil)
            return
        }
		index.Execute(w, struct{ Success bool }{true})
	})

	http.HandleFunc("/manifesto", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
            _404.Execute(w, nil)
            return
        }
		manifesto.Execute(w, struct{ Success bool }{true})
	})

	http.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
            _404.Execute(w, nil)
            return
        }
		blog.Execute(w, struct{ Success bool }{true})
	})

	http.HandleFunc("/not-found", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
            _404.Execute(w, nil)
            return
        }
		_404.Execute(w, struct{ Success bool }{true})
	})

    http.ListenAndServe(":80", nil)

}