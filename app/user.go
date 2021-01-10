package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"main/app/mod"
	"net/http"
	"strings"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	buf := &mod.User{}
	err := json.NewDecoder(r.Body).Decode(&buf)
	if err != nil {}
	nickname, er := mux.Vars(r)["nickname"]
	if er {}
	(*buf).Nickname = nickname

	old0 := &mod.User{}
	old1 := &mod.User{}
	var UserOld []*mod.User

	old0, err = GetUserByNicknameDB((*buf).Nickname)
	if err == nil {
		UserOld = append(UserOld, old0)
	}

	old1, err = GetUserByEmailDB((*buf).Email)
	if err == nil {
		if old0 != nil {
			if old1.Nickname == old0.Nickname {
				UserOld = UserOld[:0]
			}
		}
		UserOld = append(UserOld, old1)
	}

	if len(UserOld) != 0 {
		WriteJson(w, UserOld, http.StatusConflict)
		return
	}

	err = AddUserDB(*buf)
	if err != nil {}
	WriteJson(w, buf, http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	nickname, er := mux.Vars(r)["nickname"]
	if er {}

	user, err := GetUserByNicknameDB(nickname)
	if err != nil {
		var message = mod.Mes{Message: "Can't find user with id #" + nickname + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}
	WriteJson(w, user, http.StatusOK)
}

func editUser(w http.ResponseWriter, r *http.Request) {
	buf := &mod.User{}
	err := json.NewDecoder(r.Body).Decode(&buf)
	if err != nil {}
	nickname, er := mux.Vars(r)["nickname"]
	if er {}
	(*buf).Nickname = nickname

	user := &mod.User{}
	old1 := &mod.User{}

	user, err = GetUserByNicknameDB((*buf).Nickname)
	if err != nil {
		var message = mod.Mes{Message: "Can't find user with id #" + nickname + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}

	old1, err = GetUserByEmailDB((*buf).Email)
	if err == nil {
		if strings.ToLower((*old1).Nickname) != strings.ToLower((*buf).Nickname) {
			var message = mod.Mes{Message: "Can't find user with id #" + nickname + "\n"}
			WriteJson(w, message, http.StatusConflict)
			return
		}
	}

	if (*buf).Email != "" {
		(*user).Email = (*buf).Email
	}
	if (*buf).About != "" {
		(*user).About = (*buf).About
	}
	if (*buf).FullName != "" {
		(*user).FullName = (*buf).FullName
	}

	editUserDB(*user)
	WriteJson(w, user, http.StatusOK)
}
