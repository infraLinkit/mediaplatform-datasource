package entity

type RedisCounter struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	HttpStatus int    `json:"http_status"`
}
