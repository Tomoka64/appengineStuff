package main

import "net/http"

func init() {
	// TODO: ハンドラの登録
	http.HandleFunc("/", index)
}
