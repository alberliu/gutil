package httpcli

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestHttpServer(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		fmt.Println(request.PostFormValue("name"))
		fmt.Fprintf(writer, `{"id":1,"name":"alber"}`)
	})
	http.ListenAndServe(":8000", nil)
}

func TestGet(t *testing.T) {
	buf, err := Get("http://localhost:8000").SetTimeOut(1 * time.Second).ToBytes()
	fmt.Println(string(buf), err)
}

func TestPost(t *testing.T) {
	var s struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	s.Id = 1
	s.Name = "alber"

	err := Post("http://localhost:8000").BodyWithJson(s).SetTimeOut(3 * time.Second).ToJson(&s)
	fmt.Println(s, err)
}

func TestBodyWithParam(t *testing.T) {
	var s struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	s.Id = 1
	s.Name = "alber"

	err := Post("http://localhost:8000").
		BodyWithForm(map[string]string{"name": "world"}).
		SetTimeOut(3 * time.Second).ToJson(&s)
	fmt.Println(s, err)
}
