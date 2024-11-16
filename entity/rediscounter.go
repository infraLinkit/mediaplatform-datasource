package entity

type RedisCounter struct {
	Name       string `json:"name"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	HttpStatus int    `json:"http_status"`
}
