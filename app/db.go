package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"io/ioutil"
	"main/app/mod"
	"strings"
)

var db *pgx.ConnPool

func AddUserDB(buf mod.User) (err error) {
	_, err = db.Exec(
		"INSERT INTO users (nickname, email, fullname, about) VALUES ($1, $2, $3, $4)",
		buf.Nickname, buf.Email, buf.FullName, buf.About,
	)
	return err
}

func GetUserByNicknameDB(nickname string) (res *mod.User, err error) {
	user := &mod.User{}

	err = db.QueryRow(
		"SELECT * FROM users WHERE lower(nickname) = lower($1)", nickname).Scan(
			&user.ID, &user.Email, &user.Nickname, &user.FullName, &user.About)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmailDB(email string) (res *mod.User, err error) {
	user := &mod.User{}

	err = db.QueryRow(
		"SELECT * FROM users WHERE lower(email) = lower($1)", email).Scan(
		&user.ID, &user.Email, &user.Nickname, &user.FullName, &user.About)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func editUserDB(user mod.User) {
	_, _ = db.Exec("UPDATE users SET email = $2, fullname = $3, about = $4 WHERE lower(nickname) = lower($1)", user.Nickname, user.Email, user.FullName, user.About)
}

func AddForumDB(forum mod.Forum) (err error) {
	_, err = db.Exec(
		"INSERT INTO forums (slug, admin, title) VALUES ($1, $2, $3)",
		forum.Slug, forum.AdminID, forum.Title,
	)
	return err
}

func GetForumDB(slug string) (res *mod.Forum, err error) {
	forum := &mod.Forum{}

	err = db.QueryRow(
		"SELECT forums.id, forums.slug, users.nickname, forums.title, forums.threads, forums.posts FROM forums JOIN users ON (users.id = forums.admin) WHERE lower(slug) = lower($1) ", slug).
		Scan(&forum.ID, &forum.Slug, &forum.AdminNickname,
		&forum.Title, &forum.ThreadsCount, &forum.PostsCount)

	if err != nil {
		return nil, err
	}

	return forum, nil
}

func AddThreadDB(thread mod.Thread) (id uint64, err error) {
	id = 0
	err = db.QueryRow(
		"INSERT INTO threads (slug, author, title, message, forum, created) VALUES (NULLIF ($1, ''), $2, $3, $4, $5, $6) RETURNING id",
		thread.Slug, thread.AuthorID, thread.Title, thread.About, thread.ForumID, thread.CreationDate,
	).Scan(&id)
	return id, err
}

func getThreadDB(slug string, limit int, since string, desc string) (threads []*mod.Thread, err error) {
	queryString := "SELECT t.id, u.nickname, t.created, f.slug, t.message, coalesce (t.slug, ''), t.title, t.votes FROM threads AS t JOIN users AS u ON (t.author = u.id) JOIN forums AS f ON (f.id = t.forum) WHERE lower(f.slug) = lower($1) "

	if since != "" {
		if desc == "true" {
			queryString += "AND t.created <= '" + since + "'"
		} else {
			queryString += "AND t.created >= '" + since + "'"
		}
	}
	queryString += " ORDER BY t.created "

	if desc == "true" {
		queryString += "DESC "
	}

	if limit != 0 {
		queryString += fmt.Sprintf(" LIMIT %d ", limit)
	}

	var rows *pgx.Rows

	rows, err = db.Query(queryString, slug)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		t := &mod.Thread{}
		err = rows.Scan(&t.ID, &t.Author, &t.CreationDate, &t.Forum, &t.About, &t.Slug, &t.Title, &t.Votes)
		if err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}

func getThreadBySlagDB(slug string) (thread *mod.Thread, err error) {
	t := &mod.Thread{}
	err = db.QueryRow("SELECT t.id, u.nickname, t.created, t.forum, f.slug, t.message, coalesce (t.slug, ''), t.title, t.votes FROM threads AS t JOIN users AS u ON (t.author = u.id) JOIN forums AS f ON (f.id = t.forum) WHERE lower(t.slug) = lower($1)", slug).
		Scan(&t.ID, &t.Author, &t.CreationDate, &t.ForumID,
			&t.Forum, &t.About, &t.Slug, &t.Title, &t.Votes)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func getThreadByIDDB(id uint64) (thread *mod.Thread, err error) {
	t := &mod.Thread{}
	err = db.QueryRow("SELECT t.id, u.nickname, t.created, t.forum, f.slug, t.message, coalesce (t.slug, ''), t.title, t.votes FROM threads AS t JOIN users AS u ON (t.author = u.id) JOIN forums AS f ON (f.id = t.forum) WHERE t.id = $1", id).
		Scan(&t.ID, &t.Author, &t.CreationDate, &t.ForumID,
			&t.Forum, &t.About, &t.Slug, &t.Title, &t.Votes)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func AddPostsDB(posts *mod.EzPost) (*mod.EzPost, error) {
	var bilderSqlRow strings.Builder

	bilderSqlRow.WriteString("INSERT INTO posts (author, forum, message, parent, thread) VALUES ")

	last := len(*posts)

	for idx, p := range *posts {
		bilderSqlRow.WriteString(fmt.Sprintf("(%d, %d, '%s', %d, %d)", p.AuthorID, p.ForumID, p.Message, p.ParentID, p.ThreadID))
		if idx != last - 1 {
			bilderSqlRow.WriteString(",")
		}

		_, err := db.Exec("INSERT INTO users_on_forum (user_id, forum_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", p.AuthorID, p.ForumID)
		if err != nil {
			return nil, err
		}
	}

	bilderSqlRow.WriteString(" RETURNING id, created")

	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(bilderSqlRow.String())

	if err != nil {
		return nil, err
	}

	defer func() {
		rows.Close()
	}()

	postIndex := 0
	for rows.Next() {
		if err := rows.Scan(&(*posts)[postIndex].ID, &(*posts)[postIndex].CreationDate); err != nil {
			return nil, err
		}
		postIndex++
	}

	return posts, tx.Commit()
}

func addVoteDB(vote mod.Vote) (err error) {
	_, err = db.Exec(
		"INSERT INTO votes (thread, author, vote) VALUES ($1, $2, $3) ",
		vote.ThreadID, vote.UserID, vote.Voice)
	return err
}

func getVoteDB(thread, author uint64) (v *mod.Vote, err error) {
	vote := &mod.Vote{}

	err = db.QueryRow(
		"SELECT id, thread, author, vote FROM votes WHERE thread = $1 and author = $2", thread, author).Scan(
		&vote.ID, &vote.ThreadID, &vote.UserID, &vote.Voice)

	if err != nil {
		return nil, err
	}

	return vote, nil
}

func editThreadVoteDB(thread uint64, votes int64) {
	_, _ = db.Exec("UPDATE threads SET votes = $1 WHERE id = $2", votes, thread)
}

func editVoteDB(id uint64, vote int64) {
	_, _ = db.Exec("UPDATE votes SET vote = $1 WHERE id = $2", vote, id)
}

func getPostFlatDB(id uint64, limit int, since int, desc string) (posts []*mod.Post, err error) {
	queryString := "SELECT p.id, p.parent, u.nickname, p.message, p.isEdited, f.slug, p.thread, p.created FROM posts AS p JOIN users AS u ON (p.author = u.id) JOIN forums AS f ON (f.id = p.forum) WHERE p.thread = $1 "
	if since != 0 {
		if desc == "true" {
			queryString += fmt.Sprintf(" AND p.id < %d ", since)
		} else {
			queryString += fmt.Sprintf(" AND p.id > %d ", since)
		}
	}

	queryString += "ORDER BY p.created "

	if desc == "true" {
		queryString += " DESC"
	}
	queryString += ", p.id"
	if desc == "true" {
		queryString += " DESC"
	}

	if limit != 0 {
		queryString += fmt.Sprintf(" LIMIT %d ", limit)
	}

	var rows *pgx.Rows

	rows, err = db.Query(queryString, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		p := &mod.Post{}
		err = rows.Scan(&p.ID, &p.ParentID, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.ThreadID, &p.CreationDate)

		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}


func getPostTreeDB(id uint64, limit int, since int, desc string) (posts []*mod.Post, err error) {
	queryString := "SELECT p.id, p.parent, u.nickname, p.message, p.isEdited, f.slug, p.thread, p.created FROM posts AS p JOIN users AS u ON (p.author = u.id) JOIN forums AS f ON (f.id = p.forum) WHERE p.thread = $1 "

	if since != 0 {
		if desc == "true" {
			queryString += fmt.Sprintf(" AND coalesce(path < (select path FROM posts where id = %d), true) ", since)
		} else {
			queryString += fmt.Sprintf(" AND coalesce(path > (select path FROM posts where id = %d), true) ", since)
		}
	}

	queryString += "ORDER BY p.path[1]"
	if desc == "true" {
		queryString += " DESC"
	}
	queryString += ", p.path[2:]"
	if desc == "true" {
		queryString += " DESC"
	}
	queryString += " NULLS FIRST"
	if limit != 0 {
		queryString += fmt.Sprintf(" LIMIT %d ", limit)
	}

	var rows *pgx.Rows
	rows, err = db.Query(queryString, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		p := &mod.Post{}
		err = rows.Scan(&p.ID, &p.ParentID, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.ThreadID, &p.CreationDate)

		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func getPostDB(id uint64) (p *mod.Post, err error) {
	post := &mod.Post{}

	err = db.QueryRow(
		"SELECT p.id, u.nickname, f.slug, p.thread, p.message, p.created, p.isEdited, coalesce(path[array_length(path, 1) - 1], 0) FROM posts AS p JOIN users AS u ON (u.id = p.author) JOIN forums AS f ON (f.id = p.forum) WHERE p.id = $1", id).
		Scan(&post.ID, &post.Author, &post.Forum, &post.ThreadID, &post.Message,
			&post.CreationDate, &post.IsEdited, &post.ParentID)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func getPostParentDB(id uint64) (pId uint64, err error) {
	err = db.QueryRow(
		"SELECT thread FROM posts AS p WHERE p.id = $1", id).
		Scan(&pId)

	if err != nil {
		return 0, err
	}

	return pId, nil
}

func getPostParentTreeDB(id uint64, limit int, since int, desc string) (posts []*mod.Post, err error) {
	queryString := "SELECT p.id, p.parent, u.nickname, p.message, p.isEdited, f.slug, p.thread, p.created FROM posts AS p JOIN users AS u ON (p.author = u.id) JOIN forums AS f ON (f.id = p.forum) WHERE p.path[1] IN (SELECT path[1] FROM posts WHERE thread = $1 AND array_length(path, 1) = 1 "

	if since != 0 {
		if desc == "true" {
			queryString += fmt.Sprintf(" AND id < (SELECT path[1] FROM posts WHERE id = %d) ", since)
		} else {
			queryString += fmt.Sprintf(" AND id > (SELECT path[1] FROM posts WHERE id = %d) ", since)
		}
	}

	queryString += " ORDER BY id"
	if desc == "true" {
		queryString += " DESC"
	}

	if limit != 0 {
		queryString += fmt.Sprintf(" LIMIT %d ", limit)
	}

	queryString += ") ORDER BY p.path[1]"
	if desc == "true" {
		queryString += " DESC"
	}
	queryString += ", p.path[2:] NULLS FIRST"

	var rows *pgx.Rows

	rows, err = db.Query(queryString, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		p := &mod.Post{}
		err = rows.Scan(&p.ID, &p.ParentID, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.ThreadID, &p.CreationDate)

		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func editThreadDB(t mod.Thread) {
	_, _ = db.Exec("UPDATE threads SET title = $2, message = $3 WHERE ID = $1", t.ID, t.Title, t.About)
}

func GetForumUsersDB(id uint64, limit uint64, since string, desc string) (users []*mod.User, err error) {
	returnUsers := []*mod.User{}

	queryString := "SELECT u.nickname, u.email, u.fullname, u.about FROM users_on_forum fu JOIN users u ON (fu.user_id = u.id) WHERE fu.forum_id = $1"

	if since != "" {
		if desc == "true" {
			queryString += fmt.Sprintf(" AND lower(u.nickname) COLLATE \"POSIX\" < lower('" + since + "') COLLATE \"POSIX\" ")
		} else {
			queryString += fmt.Sprintf(" AND lower(u.nickname) COLLATE \"POSIX\" > lower('" + since + "') COLLATE \"POSIX\" ")
		}
	}

	queryString += " ORDER BY lower(u.nickname) COLLATE \"POSIX\""

	if desc == "true" {
		queryString += " DESC"
	}

	if limit != 0 {
		queryString += fmt.Sprintf(" LIMIT %d", limit)
	}

	var rows *pgx.Rows
	rows, err = db.Query(queryString, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := &mod.User{}
		err = rows.Scan(&u.Nickname, &u.Email, &u.FullName, &u.About)
		if err != nil {
			return nil, err
		}

		returnUsers = append(returnUsers, u)
	}

	return returnUsers, nil
}

func editPostDB(post mod.Post) {
	_, _ = db.Exec("UPDATE posts SET message = $2, isEdited = $3 WHERE id = $1", post.ID, post.Message, post.IsEdited)
}

func getInfoDB() (Fcnt, Pcnt, Tcnt, Ucnt uint64) {
	_ = db.QueryRow("SELECT count(*) from forums").Scan(&Fcnt)
	_ = db.QueryRow("SELECT count(*) from posts").Scan(&Pcnt)
	_ = db.QueryRow("SELECT count(*) from threads").Scan(&Tcnt)
	_ = db.QueryRow("SELECT count(*) from users").Scan(&Ucnt)
	return Fcnt, Pcnt, Tcnt, Ucnt
}

func ClearDB() {
	_, _ = db.Exec("TRUNCATE TABLE forums CASCADE")
	_, _ = db.Exec("TRUNCATE TABLE posts CASCADE")
	_, _ = db.Exec("TRUNCATE TABLE threads CASCADE")
	_, _ = db.Exec("TRUNCATE TABLE users CASCADE")
	_, _ = db.Exec("TRUNCATE TABLE votes CASCADE")
	_, _ = db.Exec("TRUNCATE TABLE users_on_forum CASCADE")
}

func InitDB(db *pgx.ConnPool) error {
	file, err := ioutil.ReadFile("database/init.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(file))
	if err != nil {
		return err
	}

	return nil
}


