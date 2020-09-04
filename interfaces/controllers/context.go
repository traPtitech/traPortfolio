package controllers

type Context interface {
	Param(string) string
	Bind(interface{}) error
	JSON(int, interface{}) error
	String(code int, s string) error
	NoContent(code int) error
}
