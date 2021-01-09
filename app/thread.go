package main

import (
	"encoding/json"
	_ "fmt"
	"github.com/gorilla/mux"
	_ "io/ioutil"
	"net/http"
	"strconv"
)

func createThread(w http.ResponseWriter, r *http.Request) {
	slug, er := mux.Vars(r)["slug"]
	if er {}
	buf := &Thread{}
	err := json.NewDecoder(r.Body).Decode(&buf)
	if err != nil {}

	user, err := GetUserByNicknameDB((*buf).Author)
	if err != nil {
		var message = Mes{Message: "Can't find user with id #" + (*buf).Author + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}

	forum, err := GetForumDB(slug)
	if err != nil {
		var message = Mes{Message: "Can't find user with id #" + slug + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}
	(*buf).Forum = (*forum).Slug
	(*buf).ForumID = (*forum).ID
	(*buf).AuthorID = (*user).ID

	old, err := getThreadBySlagDB((*buf).Slug)
	if err == nil {
		WriteJson(w, old, http.StatusConflict)
		return
	}

	id, _ := AddThreadDB(*buf)
	(*buf).ID = id

	//AddThreadToForum((*buf).ForumID, (*forum).ThreadsCount + 1)
	//NewUserOnForum((*user).ID, (*buf).ForumID)

	WriteJson(w, buf, http.StatusCreated)
}

func getThread(w http.ResponseWriter, r *http.Request) {
	slug, _ := mux.Vars(r)["slug"]
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	desc := r.URL.Query().Get("desc")
	since := r.URL.Query().Get("since")

	forum, err := GetForumDB(slug)
	if err != nil {
		var message = Mes{Message: "Can't find user with id #" + slug + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}

	result, err := getThreadDB(forum.Slug, limit, since, desc)
	if result == nil {
		var a []int
		a = append(a, 1)
		WriteJson(w, a[:0], http.StatusOK)
		return
	}
	WriteJson(w, result, http.StatusOK)
}

func getThreadSimple(w http.ResponseWriter, r *http.Request) {
	slug_or_id, _ := mux.Vars(r)["slug_or_id"]

	thread, err := getThreadBySlagDB(slug_or_id)
	if err != nil {
		or_id, _ := strconv.Atoi(slug_or_id)
		thread, err = getThreadByIDDB(uint64(or_id))
		if err != nil {
			var message = Mes{Message: "Can't find user with id #" + slug_or_id + "\n"}
			WriteJson(w, message, http.StatusNotFound)
			return
		}
	}
	WriteJson(w, thread, http.StatusOK)
}

func editThread(w http.ResponseWriter, r *http.Request) {
	slug_or_id, _ := mux.Vars(r)["slug_or_id"]

	thread, err := getThreadBySlagDB(slug_or_id)
	if err != nil {
		or_id, _ := strconv.Atoi(slug_or_id)
		thread, err = getThreadByIDDB(uint64(or_id))
		if err != nil {
			var message = Mes{Message: "Can't find user with id #" + slug_or_id + "\n"}
			WriteJson(w, message, http.StatusNotFound)
			return
		}
	}

	buf := &Thread{}
	err = json.NewDecoder(r.Body).Decode(&buf)
	if err != nil {}

	(*buf).ID = thread.ID

	if (*buf).Title != "" {
		(*thread).Title = (*buf).Title
	}

	if (*buf).About != "" {
		(*thread).About = (*buf).About
	}

	editThreadDB(*thread)
	WriteJson(w, thread, http.StatusOK)
}