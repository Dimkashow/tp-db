DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS users_on_forum;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS users;

CREATE UNLOGGED TABLE users (
                       id          serial  primary key,
                       email       varchar(80) unique not null,
                       nickname    varchar(80) unique not null,
                       fullname    varchar(80) not null,
                       about       text
);
CREATE INDEX idx_user_nikcname ON users (lower(nickname));

CREATE UNLOGGED TABLE forums (
                        id      serial  primary key,
                        slug    varchar(80) unique not null,
                        admin   integer not null,
                        title   varchar(120) not null,
                        threads integer default 0,
                        posts   integer default 0,
                        FOREIGN KEY (admin) REFERENCES "users" (id)
);
CREATE INDEX idx_forums_users ON forums (admin);
CREATE INDEX idx_forums_slug ON forums (lower(slug));

CREATE UNLOGGED TABLE threads (
                         id      serial not null primary key,
                         author  integer not null,
                         created timestamp (6) WITH TIME ZONE not null,
                         forum   integer not null,
                         message text not null,
                         slug    varchar(80) unique,
                         title   varchar(120) not null,
                         votes   integer default 0,
                         FOREIGN KEY (forum)     REFERENCES  "forums"    (id),
                         FOREIGN KEY (author)    REFERENCES  "users"     (id)
);
CREATE INDEX idx_threads_forums ON threads (forum, created);
CREATE INDEX idx_threads_users ON threads (author);
CREATE INDEX idx_threads_slug ON threads (lower(slug));

CREATE  UNLOGGED TABLE posts (
                       id          serial not null primary key,
                       author      integer not null,
                       forum       integer not null,
                       created     timestamp (6) WITH TIME ZONE not null default current_timestamp,
                       message     text not null,
                       isEdited    boolean default false,
                       path        integer[],
                       parent      integer,
                       thread      integer not null,
                       FOREIGN KEY (author)   REFERENCES  "users"      (id),
                       FOREIGN KEY (thread)   REFERENCES  "threads"    (id),
                       FOREIGN KEY (forum)    REFERENCES  "forums"     (id)
);
CREATE INDEX idx_posts_users ON posts (author);
CREATE INDEX idx_posts_threads_created ON posts (thread, created);
CREATE INDEX idx_posts_threads_path ON posts (thread, path);
CREATE INDEX idx_posts_threads_array ON posts (thread, (array_length(path, 1)));
CREATE INDEX idx_posts_forum ON posts (forum);
CREATE INDEX idx_posts_path_1 ON posts ((path[1]));

CREATE UNLOGGED TABLE votes (
                       id      serial  not null primary key,
                       thread  integer not null,
                       author  integer  not null,
                       vote    integer    not null,
                       FOREIGN KEY (thread)    REFERENCES  "threads"   (id),
                       FOREIGN KEY (author)    REFERENCES  "users"     (id)
);
CREATE INDEX idx_votes_uesrs ON votes (author);
CREATE INDEX idx_votes_thread ON votes (thread);
CREATE INDEX idx_votes_thread_username ON votes (thread, author);

CREATE UNLOGGED TABLE users_on_forum (
                                id      serial primary key,
                                user_id   integer not null,
                                forum_id  integer not null,
                                FOREIGN KEY (user_id) REFERENCES  "users" (id),
                                FOREIGN KEY (forum_id) REFERENCES "forums" (id)
);
ALTER table users_on_forum add unique(user_id, forum_id);
CREATE INDEX idx_users_of_forums ON users_on_forum (forum_id, user_id);

CREATE OR REPLACE FUNCTION add_post_path() RETURNS TRIGGER AS
$add_post_path$
BEGIN
    new.path = (SELECT path FROM posts WHERE id = new.parent) || new.id;
    UPDATE forums
    SET posts = posts + 1
    WHERE id = new.forum;
    RETURN new;
END;
$add_post_path$ LANGUAGE plpgsql;

CREATE TRIGGER add_post_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE add_post_path();

CREATE OR REPLACE FUNCTION add_forum_user_thread() RETURNS TRIGGER AS
$add_forum_user_thread$
BEGIN
    INSERT INTO users_on_forum (user_id, forum_id)
    VALUES (new.author, new.forum)
    ON CONFLICT DO NOTHING;
    UPDATE forums
    SET threads = threads + 1
    WHERE id = new.forum;

    RETURN new;
END;
$add_forum_user_thread$ LANGUAGE plpgsql;

CREATE TRIGGER add_forum_thread_user
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_user_thread();
