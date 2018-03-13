package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	// TODO: msgという名前のリクエストパラメタを取得
	msg := r.FormValue("msg")
	if msg == "" {
		msg = "NO MESSAGE"
	}
	fmt.Fprintln(w, msg)
}
