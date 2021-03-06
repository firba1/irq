package view

import (
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/firba1/irq/model"
)

func Delete(db model.Model, req *http.Request, r render.Render, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		r.JSON(404, ErrorEnv{"invalid quote id"})
		return
	}

	err = db.DeleteQuote(id)
	if err != nil {
		r.JSON(500, ErrorEnv{"unable to delete quote"})
		return
	}

	r.JSON(200, nil)
}
