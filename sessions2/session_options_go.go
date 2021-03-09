package sessions2

import (
	gsessions "github.com/gorilla/sessions"
	"net/http"
)

type Options struct{
	Path string
	Domain string
	MaxAge int
	Secure bool
	HttpOnly bool

	SameSite http.SameSite
}

//完成options 转换
func (options Options) ToGorillaOptions() *gsessions.Options{
	return &gsessions.Options{
		Path: options.Path,
		Domain: options.Domain,
		MaxAge: options.MaxAge,
		Secure: options.Secure,
		HttpOnly: options.HttpOnly,
		//SameSite: options.SameSite,


	}
}