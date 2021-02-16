package handler

type Response interface {
	GetCode() int
	ToJson() string
}
