package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func userVote(w http.ResponseWriter, r *http.Request) {
	slug_or_id, _ := mux.Vars(r)["slug_or_id"]

	buf := &Vote{}

	err := json.NewDecoder(r.Body).Decode(&buf)
	if err != nil {}

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

	user, err := GetUserByNicknameDB((*buf).Nickname)
	if err != nil {
		var message = Mes{Message: "Can't find user with id #" + (*buf).Nickname + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}

	vote, err := getVoteDB(thread.ID, user.ID)
	(*buf).UserID = user.ID
	(*buf).ThreadID = thread.ID
	if err != nil {
		_ = addVoteDB(*buf)
		(*thread).Votes = (*thread).Votes + (*buf).Voice
		editThreadVoteDB((*thread).ID, (*thread).Votes)

		WriteJson(w, thread, http.StatusOK)
		return
	}

	(*thread).Votes = (*thread).Votes + (*buf).Voice - (*vote).Voice
	editThreadVoteDB((*thread).ID, (*thread).Votes)
	editVoteDB((*vote).ID, (*buf).Voice)

	WriteJson(w, thread, http.StatusOK)
	return
}

