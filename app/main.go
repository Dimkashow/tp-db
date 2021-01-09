package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	_ "io/ioutil"
	"log"
	"net/http"
)

func WriteJson(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	connectConfig := pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "forums",
		User:     "newuser",
		Password: "password",
	}

	DBConfig := pgx.ConnPoolConfig{
		ConnConfig: connectConfig,
		MaxConnections: 1000,
	}

	DBConn, err := pgx.NewConnPool(DBConfig)
	if err != nil {
		logrus.Fatal(err)
		return
	}
	err = InitDB(DBConn)
	db = DBConn

	if err != nil {
		logrus.Fatal(err)
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/forum/create", createForum).Methods("post")
	r.HandleFunc("/api/forum/{slug}/details", getForum).Methods("get")
	r.HandleFunc("/api/forum/{slug}/create", createThread).Methods("post")
	r.HandleFunc("/api/forum/{slug}/users", getForumUser).Methods("get")
	r.HandleFunc("/api/forum/{slug}/threads", getThread).Methods("get")

	r.HandleFunc("/api/post/{id}/details", getPost).Methods("get")
	r.HandleFunc("/api/post/{id}/details", editPost).Methods("post")

	r.HandleFunc("/api/service/clear", clearServer).Methods("post")
	r.HandleFunc("/api/service/status", getInfo).Methods("get")

	r.HandleFunc("/api/thread/{slug_or_id}/create", createPost).Methods("post")
	r.HandleFunc("/api/thread/{slug_or_id}/details", getThreadSimple).Methods("get")
	r.HandleFunc("/api/thread/{slug_or_id}/details", editThread).Methods("post")
	r.HandleFunc("/api/thread/{slug_or_id}/posts", getPosts).Methods("get")
	r.HandleFunc("/api/thread/{slug_or_id}/vote", userVote).Methods("post")

	r.HandleFunc("/api/user/{nickname}/create", createUser).Methods("post")
	r.HandleFunc("/api/user/{nickname}/profile", getUser).Methods("get")
	r.HandleFunc("/api/user/{nickname}/profile", editUser).Methods("post")

	log.Println("Launching at HTTP port 5000")
	err = http.ListenAndServe(":5000", r)
	if err != nil {
		logrus.Fatal(err)
		return
	}
}
