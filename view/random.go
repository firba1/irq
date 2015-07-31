package view

import (
	"fmt"
	"github.com/firba1/irq/model"
	"github.com/martini-contrib/render"
	"net/http"
)

func Random(r render.Render) {
	db, err := model.NewModel("quotes.db")
	if err != nil {
		env := map[string]interface{}{
			"title": "error",
			"error": "db connection failed",
		}
		r.HTML(500, "error", env)
		return
	}

	quotes, err := db.GetQuotes(model.Query{
		Limit:   1,
		OrderBy: []string{"random()"},
	})
	if err != nil || len(quotes) == 0 {
		env := map[string]interface{}{
			"title": "error",
			"error": "quote not found",
		}
		r.HTML(500, "error", env)
		return
	}

	r.Redirect(fmt.Sprintf("/quote/%d", quotes[0].ID))
}

func RandomJson(r render.Render, req *http.Request) {
	qs := req.URL.Query()

	query := qs.Get("query")

	db, err := model.NewModel("quotes.db")

	if err != nil {
        r.JSON(500, map[string]interface{}{
            "error": "db connection failed",
        })
		return
	}

	quotes, err := db.GetQuotes(model.Query{
        Limit: 1,
        Search: query,
		OrderBy: []string{"random()"},
    })

	if err != nil || len(quotes) == 0 {
        r.JSON(500, map[string]interface{}{
            "error": "quote not found",
        })
		return
	}

    quote := quotes[0]

    r.JSON(200, map[string]interface{}{
        "id": quote.ID,
        "text": quote.Text,
        "score": quote.Score,
    })
}
