package post

type ID string

type Text string

type UserID string

type Post struct {
	PostID           ID     `json:"postId"`
	PostText         Text   `json:"postText"`
	PostAuthorUserID UserID `json:"author_user_id"`
}
