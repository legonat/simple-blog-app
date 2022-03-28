package service

import (
	"awesomeProjectRentaTeam/internal/config"
	"awesomeProjectRentaTeam/internal/db"
	"awesomeProjectRentaTeam/internal/model"
	"awesomeProjectRentaTeam/pkg"
	"awesomeProjectRentaTeam/pkg/erx"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type BlogService struct {
	repo db.BlogDb
}

func NewBlogService(repo db.BlogDb) *BlogService {
	return &BlogService{repo: repo}
}

func (s *BlogService) InitDb() error {
	err := db.InitBlogDb("./data")
	if err != nil {
		return erx.New(err)
	}

	return nil
}

func (s *BlogService) ForceInitDb() error {
	err := db.ForceInitBlogDb("./data")
	if err != nil {
		return erx.New(err)
	}

	return nil
}

func (s *BlogService) InsertPostData(post model.Post) error {
	postId, err := s.repo.InsertPostData(post)
	if err != nil {
		return erx.New(err)
	}

	if len(post.Tags) == 0 {
		return nil
	}

	err = s.repo.InsertTagsData(post)
	if err != nil {
		return erx.New(err)
	}

	tags, err := s.repo.GetTagsIds(post.Tags)
	if err != nil {
		return erx.New(err)
	}

	err = s.repo.InsertPostTags(model.PostDb{
		Post:        post,
		Id:          postId,
		DateCreated: 0,
		Tags:        tags,
	})

	if err != nil {
		return erx.New(err)
	}

	return nil
}

func (s *BlogService) GetPostsWithLimit(tag string, limit, offset int) (model.GetPostsResult, error) {
	var res model.GetPostsResult
	var posts *[]model.PostDb
	var err error
	var count int

	path := "/get_posts/"
	queryParams := map[string]string{
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(offset),
		"tag":    tag,
	}

	if tag != "" {
		posts, err = s.repo.GetPostsDataByTagWithLimit(tag, limit, offset)
		if err != nil {
			return res, erx.New(err)
		}
		count, err = s.repo.GetPostsByTagCount(tag)
		if err != nil {
			return res, erx.New(err)
		}
	}
	if tag == "" {
		posts, err = s.repo.GetAllPostsDataWithLimit(limit, offset)
		if err != nil {
			return res, erx.New(err)
		}
		count, err = s.repo.GetPostsCount()
		if err != nil {
			return res, erx.New(err)
		}
	}

	ids := make([]int, 0, len(*posts))
	for _, v := range *posts {
		ids = append(ids, v.Id)
	}

	tags, err := s.repo.GetPostsTags(ids)
	if err != nil {
		return res, erx.New(err)
	}

	for i := range *posts {
		(*posts)[i].Tags = append((*posts)[i].Tags, tags[(*posts)[i].Id]...)
	}

	res.Count = count
	res.Results = *posts

	cfg := config.GetConfigInstance()
	tlsEnabled := cfg.Server.TlsEnabled
	host := cfg.Server.Domain
	if cfg.Server.Port != 0 {
		host = fmt.Sprintf("%v:%d", host, cfg.Server.Port)
	}

	switch {
	case offset == 0 && count > limit:
		queryParams["offset"] = strconv.Itoa(offset + limit)
		res.Previous = ""
		res.Next = NewURLString(host, path, tlsEnabled, queryParams)
	case offset > 0 && offset+limit < res.Count:
		queryParams["offset"] = strconv.Itoa(offset - limit)
		res.Previous = NewURLString(host, path, tlsEnabled, queryParams)
		queryParams["offset"] = strconv.Itoa(offset + limit)
		res.Next = NewURLString(host, path, tlsEnabled, queryParams)
	case offset+limit > res.Count && offset != 0:
		queryParams["offset"] = strconv.Itoa(offset - limit)
		res.Previous = NewURLString(host, path, tlsEnabled, queryParams)
		res.Next = ""
	}

	return res, nil

}

func (s *BlogService) GenerateRandomPosts(count int) error {
	post := model.Post{}
	path := config.GetConfigInstance().DummyText.Path
	_, err := os.Stat(path)
	if err != nil {
		return erx.New(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return erx.New(err)
	}

	words := bytes.SplitN(b, []byte(" "), 27)

	rnd := func(min, max int) int {
		return pkg.RandomIntFromRange(min, max)
	}

	for i := 0; i <= count; i++ {
		post = model.Post{
			Title: string(byte(rnd(65, 90))) + string(b[rnd(1, 5):rnd(5, 70)]),
			Text:  string(b[0:rnd(250, 2800)]),
			Tags:  []string{string(words[rnd(0, 5)])},
		}

		err = s.InsertPostData(post)
		if err != nil {
			return erx.New(err)
		}
	}

	return nil
}

func NewURLString(host, path string, tlsEnabled bool, params map[string]string) string {
	proto := "http"
	if tlsEnabled {
		proto = "https"
	}
	var v = make(url.Values)
	for i, val := range params {
		if val != "" {
			v.Set(i, val)
		}
	}
	var u = url.URL{
		Scheme:   proto,
		Host:     host,
		Path:     path,
		RawQuery: v.Encode(),
	}
	return u.String()
}
