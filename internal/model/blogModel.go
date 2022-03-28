package model

type Post struct {
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}

type PostDb struct {
	Id int `json:"id"`
	Post
	DateCreated int      `json:"dateCreated"`
	Tags        []TagsDb `json:"tags"`
}

type TagsDb struct {
	Id  int    `json:"id"`
	Tag string `json:"tag"`
}

type GetPostsResult struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []PostDb `json:"results"`
}
