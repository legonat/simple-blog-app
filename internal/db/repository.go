package db

import (
	"awesomeProjectRentaTeam/internal/model"
	"database/sql"
)

type BlogDb interface {
	InsertPostData(post model.Post) (int, error)
	InsertTagsData(post model.Post) error
	InsertPostTags(post model.PostDb) error
	GetAllPostsDataWithLimit(limit, offset int) (*[]model.PostDb, error)
	GetPostsDataByTagWithLimit(tag string, limit, offset int) (*[]model.PostDb, error)
	GetPostTags(id int) ([]string, error)
	GetTagsIds(tags []string) ([]model.TagsDb, error)
	GetPostsTags(ids []int) (map[int][]model.TagsDb, error)
	GetPostsCount() (int, error)
	GetPostsByTagCount(tag string) (int, error)
}

type Repository struct {
	BlogDb
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		NewDbSqlite(db),
	}
}
