package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

// NASAの画像データ構造
type NasaImage struct {
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Url         string `json:"url"`
}


func main() {
	// NASA Open API のエンドポイント
	url := "https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY" 

	// APIからデータを取得
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// JSONデータを構造体に変換
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var imageData NasaImage
	json.Unmarshal(body, &imageData)

	// HTMLテンプレートを定義

	tmpl := template.Must(template.New("image").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Title}}</title>
		</head>
		<body>
			<h1>{{.Title}}</h1>
			<p>{{.Explanation}}</p>
			<img src="{{.Url}}" alt="{{.Title}}">
		</body>
		</html>
	`))

	// HTTPハンドラを定義
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, imageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// サーバーを起動
	fmt.Println("サーバーを起動しました。 http://localhost:8080 にアクセスしてください。")
	log.Fatal(http.ListenAndServe(":8080", nil))
}