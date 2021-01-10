package main

import (
	"main/app/mod"
	"net/http"
)

func getInfo(w http.ResponseWriter, r *http.Request) {
	var res mod.Status
	res.ForumsCount, res.PostsCount, res.ThreadsCount, res.UsersCount = getInfoDB()
	WriteJson(w, res, http.StatusOK)
}

func clearServer(w http.ResponseWriter, r *http.Request) {
	ClearDB()
	WriteJson(w, "", http.StatusOK)
}
