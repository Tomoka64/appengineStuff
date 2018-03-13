package main

import (
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var indexTmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
	<head>
		<title>ゲストブック</title>
	</head>
	<body>
	<form action="/post">
		<input type="text" name="name" placeholder="お名前">
		<input type="text" name="message" placeholder="メッセージ">
		<input type="submit">
	</form>
	<div class="messages">{{range .}}
		<div class="message">
			<h2 class="message-name">{{.Name}}</h2>
			<p class="message-text">{{.Text}}</p>
		</div>
	{{end}}</div>
	</body>
</html>`))

func index(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	msgs := make([]*Message, 0, 10)
	// TODO: Message Kindから作成時間が新しい順に10件取得する
	q := datastore.NewQuery("Message").Order("-createdAt").Limit(cap(msgs))
	for it := q.Run(ctx); ; {
		var msg Message
		// TODO: msgにDatastoreから読み込んだ値を設定する
		_, err := it.Next(&msg)
		if err == datastore.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		msgs = append(msgs, &msg)
	}
	fmt.Println(ctx, msgs)

	if err := indexTmpl.Execute(w, msgs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
