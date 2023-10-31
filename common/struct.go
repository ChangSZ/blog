package common

import (
	"html/template"
	"time"

	"github.com/ChangSZ/blog/model"
)

type PostStore struct {
	Title    string `json:"title"`
	Category int    `json:"category"`
	Tags     []int  `json:"tags"`
	Summary  string `json:"summary"`
	Content  string `json:"content"`
}

type CateStore struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	ParentId    int    `json:"parentId"`
	SeoDesc     string `json:"seoDesc"`
}

type TagStore struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	SeoDesc     string `json:"seoDesc"`
}

type LinkStore struct {
	Name  string `json:"name"`
	Link  string `json:"link"`
	Order int    `json:"order"`
}

type AuthLogin struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Captcha    string `json:"captcha"`
	CaptchaKey string `json:"captchaKey"`
}

type AuthRegister struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConsolePostList struct {
	Post     ConsolePost  `json:"post,omitempty"`
	Tags     []ConsoleTag `json:"tags,omitempty"`
	Category ConsoleCate  `json:"category,omitempty"`
	View     ConsoleView  `json:"view,omitempty"`
	Author   ConsoleUser  `json:"author,omitempty"`
}

type ConsolePost struct {
	Id        int       `json:"id,omitempty"`
	Uid       string    `json:"uid,omitempty"`
	UserId    int       `json:"userId,omitempty"`
	Title     string    `json:"title,omitempty"`
	Summary   string    `json:"summary,omitempty"`
	Original  string    `json:"original,omitempty"`
	Content   string    `json:"content,omitempty"`
	Password  string    `json:"password,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type ConsoleTag struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	SeoDesc     string `json:"seoDesc,omitempty"`
	Num         int    `json:"num,omitempty"`
}

type ConsoleCate struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	SeoDesc     string `json:"seoDesc,omitempty"`
}

type ConsoleUser struct {
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Status int    `json:"status,omitempty"`
}

type ConsoleSystem struct {
	Title        string `json:"title;omitempty"`
	Keywords     string `json:"keywords;omitempty"`
	Theme        int    `json:"theme;omitempty"`
	Description  string `json:"description;omitempty"`
	RecordNumber string `json:"recordNumber;omitempty"`
}

type ConsoleView struct {
	Num int `json:"num,omitempty"`
}

type IndexPostList struct {
	PostListArr []*ConsolePostList
	Paginate    Paginate
}

type Paginate struct {
	Limit   int `json:"limit"`
	Count   int `json:"count"`
	Total   int `json:"total"`
	Last    int `json:"last"`
	Current int `json:"current"`
	Next    int `json:"next"`
}

type IndexPost struct {
	Id        int           `json:"id,omitempty"`
	Uid       string        `json:"uid,omitempty"`
	UserId    int           `json:"userId,omitempty"`
	Title     string        `json:"title,omitempty"`
	Summary   string        `json:"summary,omitempty"`
	Original  string        `json:"original,omitempty"`
	Content   template.HTML `json:"content,omitempty"`
	Password  string        `json:"password,omitempty"`
	DeletedAt time.Time     `json:"deletedAt,omitempty"`
	CreatedAt time.Time     `json:"createdAt,omitempty"`
	UpdatedAt time.Time     `json:"updatedAt,omitempty"`
}

type IndexRss struct {
	Id        int       `json:"id,omitempty"`
	Uid       string    `json:"uid,omitempty"`
	Author    string    `json:"author,omitempty"`
	Title     string    `json:"title,omitempty"`
	Summary   string    `json:"summary,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type IndexPostDetail struct {
	Post     IndexPost    `json:"post,omitempty"`
	Tags     []ConsoleTag `json:"tags,omitempty"`
	Category ConsoleCate  `json:"category,omitempty"`
	View     ConsoleView  `json:"view,omitempty"`
	Author   ConsoleUser  `json:"author,omitempty"`
	LastPost *model.Posts `json:"lastPost,omitempty"`
	NextPost *model.Posts `json:"nextPost,omitempty"`
}

type IndexGithubParam struct {
	GithubName         string
	GithubRepo         string
	GithubClientId     string
	GithubClientSecret string
	GithubLabels       string
}

type Category struct {
	Cates model.Categories `json:"cates"`
	Html  string           `json:"html"`
}

type IndexCategory struct {
	Cates model.Categories `json:"cates"`
	Html  template.HTML    `json:"html"`
}
