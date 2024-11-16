package entity

type RedisCounter struct {
	Name       string `json:"name"`
	Date       string `json:"date"`
	HttpStatus int    `json:"http_status"`
}
