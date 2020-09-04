package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Templates
const (
	IndexTemplate     = "template/index.html"
	BlogTemplate      = "template/blog.html"
	NotFoundTemplate  = "template/404.html"
	PostTemplate      = "template/post.html"
	ManifestoTemplate = "template/manifesto.html"
)

// Client mongo Db
var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("posts").Collection("post")

}

type P struct {
	Titulo    string             `bson:"titulo"`
	SubTitulo string             `bson:"sub_titulo"`
	Conteudo  string             `bson:"conteudo"`
	Fotos     string             `bson:"fotos"`
	Autor     string             `bson:"autor"`
	Data      time.Time          `bson:"data"`
	ID        primitive.ObjectID `bson:"_id"`
}

type PP struct {
	Posts []*P
}

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

	// homeDir, _ := os.UserHomeDir()
	// projectDir := homeDir + "/go/brigadacatarina/"
	// index := template.Must(template.ParseFiles("template/index.html"))
	// manifesto := template.Must(template.ParseFiles("template/manifesto.html"))
	// blog := template.Must(template.ParseFiles("template/blog.html"))
	// _404 := template.Must(template.ParseFiles("template/404.html"))
	// postPage := template.Must(template.ParseFiles("template/post.html"))

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", Index)
	http.HandleFunc("/manifesto", Manifesto)
	http.HandleFunc("/blog", Blog)
	http.HandleFunc("/post/{permalink}", GetPost)

	fmt.Println("Listening port 9990")
	http.ListenAndServe(":9990", nil)

}

// Test
func Test(w http.ResponseWriter, r *http.Request) {
	permalink := r.URL.Query().Get("permalink")
	fmt.Println(permalink)
}

// Manifesto page
func Manifesto(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(ManifestoTemplate)).Execute(w, struct{ Success bool }{true})
}

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(IndexTemplate)).Execute(w, struct{ Success bool }{true})
}

// Blog page
func Blog(w http.ResponseWriter, r *http.Request) {
	posts, _ := getAll()
	for _, v := range posts {
		fmt.Print(v.Titulo)
	}

	data := PP{
		Posts: posts,
	}

	blog := template.Must(template.ParseFiles(BlogTemplate))
	blog.Execute(w, data)
}

// Get Post page
func GetPost(w http.ResponseWriter, r *http.Request) {

	permalink := "permalink"

	data, _ := getByID(permalink)

	if len(data) < 1 {
		template.Must(template.ParseFiles(NotFoundTemplate))
	}
	post := template.Must(template.ParseFiles(PostTemplate))
	post.Execute(w, data)
}

func createPost(post *Post) error {
	_, err := collection.InsertOne(ctx, post)
	return err
}

func getAll() ([]*P, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	return filterPosts(filter)
}

func getByID(permalink string) ([]*P, error) {
	filter := bson.D{
		primitive.E{Key: "_id", Value: permalink},
	}

	return filterPosts(filter)
}

func filterPosts(filter interface{}) ([]*P, error) {
	// A slice of tasks for storing the decoded documents
	var posts []*P

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return posts, err
	}

	for cur.Next(ctx) {
		var t P
		err := cur.Decode(&t)
		if err != nil {
			return posts, err
		}

		posts = append(posts, &t)
	}

	if err := cur.Err(); err != nil {
		return posts, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(posts) == 0 {
		return posts, mongo.ErrNoDocuments
	}

	return posts, nil
}
