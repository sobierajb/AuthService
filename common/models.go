package common

type GraterLess[T interface{}] struct {
	LessThen   T `json:"lessThan"`
	GraterThen T `json:"graterThan"`
}

type Request[B,Q,P interface{}] struct {
	Body B 
	Query Q
	Parameter P   
}

type Response[T interface{}] struct {
	Data  T      `json:"data"`
	Error string `json:"error"`
}

type Id[T any] struct {
	Id T `json:"id"`
}

type Session struct {
	
}