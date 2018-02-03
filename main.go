package main

import (
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/search"
)

type Person struct {
	Name     string
	Age      string
	IntroMsg string
}

type newModel struct {
	CreatedID string
}

var tpl *template.Template

func init() {
	http.HandleFunc("/", new)
	http.HandleFunc("/new", new)
	http.HandleFunc("/put", handlePut)
	http.HandleFunc("/show", show)
	http.HandleFunc("/search", searchin)
	http.HandleFunc("/searchprocessing", searchHandler)
	// http.HandleFunc("/get", handleGet)
	tpl = template.Must(template.ParseGlob("templates/*"))

}

func new(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "new.gohtml", nil)
}
func handlePut(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	model := &newModel{}

	name := req.FormValue("Name")
	age := req.FormValue("Age")
	intro := req.FormValue("intro")

	person := &Person{
		Name:     name,
		Age:      age,
		IntroMsg: intro,
	}

	index, err := search.Open("example")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := index.Put(ctx, "", person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	model.CreatedID = id
	tpl.ExecuteTemplate(w, "put.gohtml", model)

}

func show(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	index, err := search.Open("example")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := "325683c2-aeda-4c93-b959-26f3274b167e"
	var user Person
	if err := index.Get(ctx, id, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Print(w, "Retrieved document: ", user)
	tpl.ExecuteTemplate(w, "profile.gohtml", &user)
}

func searchin(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "search.gohtml", nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	query := r.FormValue("query")
	index, err := search.Open("example")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for t := index.Search(ctx, query, nil); ; {
		var user Person
		id, err := t.Next(&user)
		if err == search.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "Search error: %v\n", err)
			break
		}
		fmt.Fprintf(w, "%s -> %#v\n", id, user)
	}
}

//
//
// func handleGet(w http.ResponseWriter, r *http.Request) {
// 	ctx := appengine.NewContext(r)
// 	index, err := search.Open("example")
// 	if err != nil {
// 		panic(err)
// 	}
// 	var person Person
// 	err = index.Get(ctx, "", &person)
//
// 	tpl.ExecuteTemplate(w, "profile.gohtml", person)
// 	fmt.Fprintf(w, "%v", person)
// }
