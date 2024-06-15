package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := getRouter()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func getRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/test", get)
	// With(auth)で認証ミドルウェアを実行している
	r.With(auth).Post("/test", post)
	r.With(auth).Post("/csv", uploadCSV)
	return r
}

// 認証のミドルウェア
// ヘッダーのAuthorizationに「ok_token」が入ってない場合は401エラーで返す
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "ok_token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GETメソッド処理
// 受け取ったIDを返すだけ
func get(w http.ResponseWriter, r *http.Request) {
	log.Println("Get!!")
	id := r.FormValue("id")
	type response struct {
		ID string `json:"id"`
	}
	res := response{ID: id}
	http.Header.Add(w.Header(), "content-type", "application/json")
	http.Header.Add(w.Header(), "Access-Control-Allow-Origin", "*")
	v, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(v)
}

// POSTメソッド処理
// 受け取ったID、名前を返すだけ
func post(w http.ResponseWriter, r *http.Request) {
	log.Println("Post!!")

	type response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var res response

	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(res)
	http.Header.Add(w.Header(), "content-type", "application/json")
	http.Header.Add(w.Header(), "Access-Control-Allow-Origin", "*")
	v, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(v)
}

// CSVのアップロード
// CSVの中身を出力している
func uploadCSV(w http.ResponseWriter, r *http.Request) {
	log.Println("CSV!!")
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	for i := 0; ; i++ {
		record, err := csvReader.Read()
		// 最後の行なので終了
		if err == io.EOF {
			log.Println("EOF")
			break
		}
		if err != nil {
			log.Println("error", err)
			break
		}
		log.Println(record)
	}
}
