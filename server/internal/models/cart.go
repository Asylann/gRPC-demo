package models

type Cart struct {
	Id      int `json:"id" sqlx:"id"`
	User_id int `json:"userId" sqlx:"user_id"`
}
