package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Post struct {
	Titulo    string `json:"titulo"`
	SubTitulo string `json:"sub_titulo"`
	Conteudo  string `json:"conteudo"`
	Fotos     string `json:"fotos"`
	Autor     string `json:"autor"`
	Data      string `json:"data"`
	Permalink string `json:"permalink"`
}

type Posts struct {
	Posts []Post
}

func FindPost() []Post {

	jsonFile, err := os.Open("../data/posts.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened posts.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var p []Post

	json.Unmarshal(byteValue, &p)

	return p
}

func GetPosts() Posts {
	data := Posts{
		Posts: FindPost(),
	}
	return data
}

func main() {

	index := template.Must(template.ParseFiles("template/index.html"))
	manifesto := template.Must(template.ParseFiles("template/manifesto.html"))
	blog := template.Must(template.ParseFiles("template/blog.html"))
	_404 := template.Must(template.ParseFiles("template/404.html"))

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("index")
		if r.Method != http.MethodGet {
			index.Execute(w, nil)
			return
		}
		index.Execute(w, struct{ Success bool }{true})
	})

	http.HandleFunc("/manifesto", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("manifesto")
		if r.Method != http.MethodGet {
			_404.Execute(w, nil)
			return
		}
		manifesto.Execute(w, struct{ Success bool }{true})
	})

	http.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("blog")
		if r.Method != http.MethodGet {
			_404.Execute(w, nil)
			return
		}

		data := GetPosts()

		blog.Execute(w, data)
	})

	http.HandleFunc("/not-found", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("not-found")
		if r.Method != http.MethodGet {
			_404.Execute(w, nil)
			return
		}
		_404.Execute(w, struct{ Success bool }{true})
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("not-found")
		if r.Method != http.MethodGet {
			_404.Execute(w, nil)
			return
		}
		_404.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(":9990", nil)

}
