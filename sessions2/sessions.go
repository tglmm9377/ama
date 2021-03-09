package sessions2

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

const (
	DefaultKey = "tglmm"
	errorFormat = "[sessions] ERROR %s\n"
)


type Store interface{
	sessions.Store
	Options(Options)
}

type Session interface {
	ID() string
	Get(key interface{}) interface{}
	Set(key interface{} , val interface{})
	Delete(key interface{})
	Clear()
	AddFlash(value interface{},vars ...string)
	Flashes(vars ...string) []interface{}
	Options(Options)
}

type session struct {
	name string
	request *http.Request
	store Store
	session *sessions.Session
	written bool
	writer http.ResponseWriter
}

func Sessions(name string, store Store) gin.HandlerFunc{
	return func(c *gin.Context){
		s := &session{name ,c.Request,store,nil,false,c.Writer}
		c.Set(DefaultKey , s)
		//请求生存期结束后清理变量
		defer context.Clear(c.Request)
		//Next应该只在中间件内部使用。它执行调用处理程序内链中的挂起处理程序
		c.Next()
	}

}

func SessionsMany(names []string, store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions := make(map[string]Session, len(names))
		for _, name := range names {
			sessions[name] = &session{name, c.Request, store, nil, false, c.Writer}
		}
		c.Set(DefaultKey, sessions)
		defer context.Clear(c.Request)
		c.Next()
	}
}

func (s *session) ID() string{
	return s.Session().ID
}


func (s *session) Session() *sessions.Session{
	if s.session == nil{
		var err error
		s.session , err = s.store.Get(s.request , s.name)
		if err != nil{
			log.Printf(errorFormat , err)
		}
	}
	return s.session
}
func (s *session) Get(key interface{})interface{}{
	return s.Session().Values[key]
}

func (s *session) Set(key interface{} , val interface{}){
	s.Session().Values[key] = val

}

func (s *session) Delete(key interface{}){
	delete(s.Session().Values,key)
	s.written = true
}

func (s *session) Clear(){
	for key := range s.Session().Values{
		s.Delete(key)
	}
}

func (s *session) AddFlash(value interface{} , vars ...string){
	s.Session().AddFlash(value , vars...)
	s.written = true
}

func (s *session) Flashes(vars ...string)[]interface{}{
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *session) Options(options Options){
	s.Session().Options = options.ToGorillaOptions()
}

func (s *session) Save() error{
	if s.Written(){
		e := s.Session().Save(s.request , s.writer)
		if e == nil{
			s.written = false
		}
		return e
	}
	return nil
}

func (s *session) Written()bool{
	return s.written
}

func Default(c *gin.Context) Session{
	return c.MustGet(DefaultKey).(Session)
}

func DefaultMany(c *gin.Context , name string)Session{
	//通过name	找到对应的session 断言
	return c.MustGet(DefaultKey).(map[string]Session)[name]
}