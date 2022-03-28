package service

import (
	"awesomeProjectRentaTeam/internal/db"
	"awesomeProjectRentaTeam/internal/model"
)

type Blog interface {
	InitDb() error
	ForceInitDb() error
	InsertPostData(post model.Post) error
	GetPostsWithLimit(tag string, limit, offset int) (model.GetPostsResult, error)
	GenerateRandomPosts(count int) error
}

type Service struct {
	Blog
}

func NewService(repo *db.Repository) *Service {
	return &Service{Blog: NewBlogService(repo.BlogDb)}
}
