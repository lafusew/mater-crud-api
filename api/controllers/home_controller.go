package controllers

import (
	"net/http"

	responses "github.com/gob-mater/app/api/res"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Ma'ter: GOB_2022")
}