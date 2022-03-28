package db

import (
	"awesomeProjectRentaTeam/internal/model"
	"awesomeProjectRentaTeam/pkg/erx"
	"database/sql"
	"strings"
	"time"
)

const (
	INSERT_POST            = `INSERT INTO posts (title, text, date_created) VALUES ($1, $2, $3) RETURNING id`
	INSERT_TAG             = `INSERT OR IGNORE INTO tags (tag) VALUES (?);`
	INSERT_POST_TAG        = `INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?);`
	GET_TAG                = `SELECT * FROM tags WHERE tag LIKE ?`
	GET_POSTS              = `SELECT * FROM posts LIMIT $1 OFFSET $2;`
	GET_POSTS_COUNT        = `SELECT COUNT(*) FROM posts;`
	GET_POST_TAGS          = `SELECT t.tag FROM tags t JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = ?;`
	GET_POSTS_TAGS         = `SELECT pt.post_id, t.id, t.tag FROM tags t JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = ?`
	GET_POSTS_BY_TAG_COUNT = `SELECT COUNT(*) FROM posts p JOIN (SELECT DISTINCT post_id FROM post_tags pt JOIN tags t ON t.id = pt.tag_id WHERE t.tag LIKE $1) pi on p.id = pi.post_id;`
	GET_POSTS_BY_TAG       = `SELECT id, title, text, date_created FROM posts p JOIN (SELECT DISTINCT post_id FROM post_tags pt JOIN tags t ON t.id = pt.tag_id WHERE t.tag LIKE $1) pi on p.id = pi.post_id LIMIT $2 OFFSET $3;`
)

type DbSqlite struct {
	db *sql.DB
}

func NewDbSqlite(db *sql.DB) *DbSqlite {
	return &DbSqlite{db: db}
}

func (r *DbSqlite) InsertPostData(post model.Post) (int, error) {
	var id int
	row := r.db.QueryRow(INSERT_POST, post.Title, post.Text, time.Now().Unix())

	err := row.Scan(&id)
	if err != nil {
		return 0, erx.New(err)
	}

	return id, nil
}

func (r *DbSqlite) InsertTagsData(post model.Post) error {
	req, args, err := prepareInsertTagsRequest(post.Tags)
	if err != nil {
		return erx.New(err)
	}

	_, err = r.db.Exec(req, args...)
	if err != nil {
		return erx.New(err)
	}

	return nil
}

func (r *DbSqlite) InsertPostTags(post model.PostDb) error {
	req, args, err := prepareInsertPostTagsRequest(&post)
	if err != nil {
		return erx.New(err)
	}

	_, err = r.db.Exec(req, args...)
	if err != nil {
		return erx.New(err)
	}

	return nil
}

func (r *DbSqlite) GetAllPostsDataWithLimit(limit, offset int) (*[]model.PostDb, error) {
	var postDb model.PostDb
	postDbs := make([]model.PostDb, 0, limit)

	rows, err := r.db.Query(GET_POSTS, limit, offset)
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)
	if err != nil {
		return nil, erx.New(err)
	}
	for rows.Next() {
		err = rows.Scan(&postDb.Id, &postDb.Title, &postDb.Text, &postDb.DateCreated)
		if err != nil {
			return nil, erx.New(err)
		}
		postDbs = append(postDbs, postDb)
	}

	return &postDbs, nil
}

func (r *DbSqlite) GetPostsDataByTagWithLimit(tag string, limit, offset int) (*[]model.PostDb, error) {
	var postDb model.PostDb
	postDbs := make([]model.PostDb, 0, limit)

	rows, err := r.db.Query(GET_POSTS_BY_TAG, tag, limit, offset)
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)
	if err != nil {
		return nil, erx.New(err)
	}
	for rows.Next() {
		err = rows.Scan(&postDb.Id, &postDb.Title, &postDb.Text, &postDb.DateCreated)
		if err != nil {
			return nil, erx.New(err)
		}
		postDbs = append(postDbs, postDb)
	}

	return &postDbs, nil
}

func (r *DbSqlite) GetPostTags(id int) ([]string, error) {
	var tag string
	var tags []string

	rows, err := r.db.Query(GET_POST_TAGS)
	if err != nil {
		return nil, erx.New(err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	for rows.Next() {
		err = rows.Scan(&tag)
		if err != nil {
			return nil, erx.New(err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *DbSqlite) GetPostsTags(ids []int) (map[int][]model.TagsDb, error) {
	var tag model.TagsDb
	tags := make(map[int][]model.TagsDb)

	req, vals, err := prepareGetPostsTagsRequest(ids)
	rows, err := r.db.Query(req, vals...)
	if err != nil {
		return nil, erx.New(err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)
	var id int
	for rows.Next() {
		err = rows.Scan(&id, &tag.Id, &tag.Tag)
		if err != nil {
			return nil, erx.New(err)
		}
		tags[id] = append(tags[id], tag)
	}

	return tags, nil
}

func (r *DbSqlite) GetTagsIds(tags []string) ([]model.TagsDb, error) {
	req, vals, err := prepareGetTagsIdsRequest(tags)
	if err != nil {
		return nil, erx.New(err)
	}

	rows, err := r.db.Query(req, vals...)
	if err != nil {
		return nil, erx.New(err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	var tag model.TagsDb
	tagsDb := make([]model.TagsDb, 0, len(tags))
	for rows.Next() {
		err = rows.Scan(&tag.Id, &tag.Tag)
		if err != nil {
			return nil, erx.New(err)
		}
		tagsDb = append(tagsDb, tag)
	}

	return tagsDb, nil
}

func (r DbSqlite) GetPostsCount() (int, error) {
	var count int
	row := r.db.QueryRow(GET_POSTS_COUNT)

	err := row.Scan(&count)
	if err != nil {
		return 0, erx.New(err)
	}

	return count, nil
}

func (r DbSqlite) GetPostsByTagCount(tag string) (int, error) {
	var count int
	row := r.db.QueryRow(GET_POSTS_BY_TAG_COUNT, tag)

	err := row.Scan(&count)
	if err != nil {
		return 0, erx.New(err)
	}

	return count, nil
}

func prepareInsertTagsRequest(tags []string) (req string, valArgs []interface{}, err error) {
	rowCount := len(tags)
	valArgs = make([]interface{}, 0, rowCount)
	valStrings := make([]string, 0, rowCount)
	for _, v := range tags {
		valArgs = append(valArgs, v)
		valStrings = append(valStrings, INSERT_TAG)
	}
	req = strings.Join(valStrings, "")

	return
}

func prepareGetTagsIdsRequest(tags []string) (req string, valArgs []interface{}, err error) {
	rowCount := len(tags)
	valArgs = make([]interface{}, 0, rowCount)
	valStrings := make([]string, 0, rowCount)
	for _, v := range tags {
		valArgs = append(valArgs, v)
		valStrings = append(valStrings, GET_TAG)
	}
	req = strings.Join(valStrings, " UNION ALL ")

	return
}

func prepareGetPostsTagsRequest(ids []int) (req string, valArgs []interface{}, err error) {
	rowCount := len(ids)
	valArgs = make([]interface{}, 0, rowCount)
	valStrings := make([]string, 0, rowCount)
	for _, v := range ids {
		valArgs = append(valArgs, v)
		valStrings = append(valStrings, GET_POSTS_TAGS)
	}
	req = strings.Join(valStrings, " UNION ALL ")

	return
}

func prepareInsertPostTagsRequest(post *model.PostDb) (req string, valArgs []interface{}, err error) {
	rowCount := len(post.Tags)
	valArgs = make([]interface{}, 0, rowCount)
	valStrings := make([]string, 0, rowCount)
	for _, v := range post.Tags {
		valArgs = append(valArgs, post.Id)
		valArgs = append(valArgs, v.Id)
		valStrings = append(valStrings, INSERT_POST_TAG)
	}
	req = strings.Join(valStrings, "")
	return
}
