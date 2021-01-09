package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	_ "io/ioutil"
	"net/http"
	"strconv"
)

func createForum(w http.ResponseWriter, r *http.Request) {
	buf := &Forum{}
	err := json.NewDecoder(r.Body).Decode(&buf)
	if err != nil {}

	user, err := GetUserByNicknameDB((*buf).AdminNickname)
	if err != nil {
		var message = Mes{Message: "Can't find user with id #" + (*buf).AdminNickname + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}
	(*buf).AdminNickname = (*user).Nickname

	old, err := GetForumDB((*buf).Slug)
	if err == nil {
		WriteJson(w, old, http.StatusConflict)
		return
	}

	(*buf).AdminID = user.ID
	err = AddForumDB(*buf)
	WriteJson(w, buf, http.StatusCreated)
	return
}

func getForum(w http.ResponseWriter, r *http.Request) {
	slug, er := mux.Vars(r)["slug"]
	if er {}

	forum, err := GetForumDB(slug)
	if err != nil {
		var message = Mes{Message: "Can't find user with id #" + slug + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}
	WriteJson(w, forum, http.StatusOK)
}

func getForumUser(w http.ResponseWriter, r *http.Request) {
	slug, er := mux.Vars(r)["slug"]
	if er {}

	forum, err := GetForumDB(slug)
	if err != nil {
		var message = Mes{Message: "Can't find user with id #" + slug + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	desc := r.URL.Query().Get("desc")
	since := r.URL.Query().Get("since")

	users, err := GetForumUsersDB((*forum).ID, uint64(limit), since, desc)
	WriteJson(w, users, http.StatusOK)
}
