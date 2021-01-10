package main

import (
	"encoding/json"
	_ "fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	_ "io/ioutil"
	"main/app/mod"
	"net/http"
	"strconv"
	_ "strconv"
	"strings"
	"time"
)

func createPost(w http.ResponseWriter, r *http.Request) {
	buf := mod.EzPost{}
	err := easyjson.UnmarshalFromReader(r.Body, &buf)
	if err != nil {}

	slug_or_id, _ := mux.Vars(r)["slug_or_id"]
	thread, err := getThreadBySlagDB(slug_or_id)
	if err != nil {
		or_id, _ := strconv.Atoi(slug_or_id)
		thread, err = getThreadByIDDB(uint64(or_id))
		if err != nil {
			var message = mod.Mes{Message: "err 1 #" + slug_or_id + "\n"}
			WriteJson(w, message, http.StatusNotFound)
			return
		}
	}

	if len(buf) == 0 {
		var a []int
		a = append(a, 1)
		WriteJson(w, a[:0], http.StatusCreated)
		return
	}

	now := time.Now()
	for id := range buf {
		if (*(buf)[id]).ParentID != 0 {
			//CREATE NEW FUNC WHERE BOOL, ID
			parentThread, err := getPostParentDB((*(buf)[id]).ParentID)
			if err != nil {
				var message = mod.Mes{Message: "err 2 #" + (*(buf)[id]).Author + "\n"}
				WriteJson(w, message, http.StatusConflict)
				return
			}

			if parentThread != thread.ID {
				var message = mod.Mes{Message: "err 3 #" + (*(buf)[id]).Author + "\n"}
				WriteJson(w, message, http.StatusConflict)
				return
			}
		}

		user, err := GetUserByNicknameDB((*(buf)[id]).Author)
		if err != nil {
			var message = mod.Mes{Message: "err 4 #" + (*(buf)[id]).Author + "\n"}
			WriteJson(w, message, http.StatusNotFound)
			return
		}
		(*(buf)[id]).CreationDate = now
		(*(buf)[id]).Forum = thread.Forum
		(*(buf)[id]).ForumID = thread.ForumID

		(*(buf)[id]).ThreadID = thread.ID
		(*(buf)[id]).AuthorID = user.ID
	}

/*	for id := range buf {
		NewUserOnForum((*(buf)[id]).AuthorID, (*(buf)[id]).ForumID)
	}*/

	new_buf, err := AddPostsDB(&buf)

	//forum, err := GetForumDB(thread.Forum)
	//AddPostToForum(thread.ForumID, (*forum).PostsCount + uint64(len(buf)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	easyjson.MarshalToHTTPResponseWriter(*new_buf, w)

	//WriteJson(w, *new_buf, http.StatusCreated)
}

func getPostFlat(w http.ResponseWriter, r *http.Request) {
	slug_or_id, _ := mux.Vars(r)["slug_or_id"]

	thread, err := getThreadBySlagDB(slug_or_id)
	if err != nil {
		or_id, _ := strconv.Atoi(slug_or_id)
		thread, err = getThreadByIDDB(uint64(or_id))
		if err != nil {
			var message = mod.Mes{Message: "Can't find user with id #" + slug_or_id + "\n"}
			WriteJson(w, message, http.StatusNotFound)
			return
		}
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	desc := r.URL.Query().Get("desc")
	since, _ := strconv.Atoi(r.URL.Query().Get("since"))

	result, err := getPostFlatDB((*thread).ID, limit, since, desc)
	if result == nil {
		var a []int
		a = append(a, 1)
		WriteJson(w, a[:0], http.StatusOK)
		return
	}
	WriteJson(w, result, http.StatusOK)
}

func getPostTree(w http.ResponseWriter, r *http.Request) {
	slug_or_id, _ := mux.Vars(r)["slug_or_id"]

	thread, err := getThreadBySlagDB(slug_or_id)
	if err != nil {
		or_id, _ := strconv.Atoi(slug_or_id)
		thread, err = getThreadByIDDB(uint64(or_id))
		if err != nil {
			var message = mod.Mes{Message: "Can't find user with id #" + slug_or_id + "\n"}
			WriteJson(w, message, http.StatusNotFound)
			return
		}
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	desc := r.URL.Query().Get("desc")
	since, _ := strconv.Atoi(r.URL.Query().Get("since"))

	result, err := getPostTreeDB((*thread).ID, limit, since, desc)
	if result == nil {
		var a []int
		a = append(a, 1)
		WriteJson(w, a[:0], http.StatusOK)
		return
	}
	WriteJson(w, result, http.StatusOK)
}

func getPostParentTree(w http.ResponseWriter, r *http.Request) {
	slug_or_id, _ := mux.Vars(r)["slug_or_id"]

	thread, err := getThreadBySlagDB(slug_or_id)
	if err != nil {
		or_id, _ := strconv.Atoi(slug_or_id)
		thread, err = getThreadByIDDB(uint64(or_id))
		if err != nil {
			var message = mod.Mes{Message: "Can't find user with id #" + slug_or_id + "\n"}
			WriteJson(w, message, http.StatusNotFound)
			return
		}
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	desc := r.URL.Query().Get("desc")
	since, _ := strconv.Atoi(r.URL.Query().Get("since"))

	result, err := getPostParentTreeDB((*thread).ID, limit, since, desc)
	if result == nil {
		var a []int
		a = append(a, 1)
		WriteJson(w, a[:0], http.StatusOK)
		return
	}
	WriteJson(w, result, http.StatusOK)
}

// В ОДНУ ФУНКЦИЮ
func getPosts(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	if sort == "tree" {
		getPostTree(w, r)
	} else if sort == "parent_tree" {
		getPostParentTree(w, r)
	} else {
		getPostFlat(w, r)
	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	post, err := getPostDB(uint64(id))
	if err != nil {
		var message = mod.Mes{Message: "Can't find user with id #" + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}
	var result mod.PostFull
	result.PostData = post

	buf := strings.Split(r.URL.Query().Get("related"), ",")
	r.URL.Query().Get("limit")
	for _, val := range buf {
		if val == "user" {
			user, err := GetUserByNicknameDB((*post).Author)
			if err != nil {
				var message = mod.Mes{Message: "Can't find user with id #" + "\n"}
				WriteJson(w, message, http.StatusNotFound)
				return
			}
			result.Author = user
		} else if val == "forum" {
			forum, err := GetForumDB((*post).Forum)
			if err != nil {
				var message = mod.Mes{Message: "Can't find user with id #" + "\n"}
				WriteJson(w, message, http.StatusNotFound)
				return
			}
			result.Forum = forum
		} else if val == "thread" {
			thread, err := getThreadByIDDB((*post).ThreadID)
			if err != nil {
				var message = mod.Mes{Message: "Can't find user with id #" + "\n"}
				WriteJson(w, message, http.StatusNotFound)
				return
			}
			result.Thread = thread
		}
	}

	WriteJson(w, result, http.StatusOK)
}

func editPost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	post, err := getPostDB(uint64(id))
	if err != nil {
		var message = mod.Mes{Message: "Can't find user with id #" + "\n"}
		WriteJson(w, message, http.StatusNotFound)
		return
	}

	buf := &mod.Mes{}
	err = json.NewDecoder(r.Body).Decode(&buf)

	if (*buf).Message != "" && (*buf).Message != (*post).Message {
		(*post).Message = (*buf).Message
		(*post).IsEdited = true
		editPostDB(*post)
	}

	WriteJson(w, post, http.StatusOK)
}
