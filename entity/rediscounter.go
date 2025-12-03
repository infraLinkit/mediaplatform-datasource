package entity

type RedisCounter struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	HttpStatus int    `json:"http_status"`
	OS         string `json:"os"`
	IsOperator bool   `json:"is_operator"`
}

type DomainServices struct {
	Domain string `json:"domain"`
	Render string `json:"render"`
}