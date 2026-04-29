package handler

import (
	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"net/http"
	"service/internal/activity/repository"
	parser2 "service/internal/pkg/parser"
)

type ActivityHandler struct{}

func (ctr ActivityHandler) Get(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewActivityRepository()
	activities, pagination, _ := repo.Find(r.URL.Query())

	parser := parser2.ActivityParser{Array: activities}

	res := xtremeres.Response{Array: parser.Get(), Pagination: &pagination}
	res.Success(w)
}
