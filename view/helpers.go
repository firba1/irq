package view

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/firba1/irq/model"
)

func maxPage(totalItems, perPage int) int {
	return (totalItems-1)/perPage + 1
}

func queryString(page, count int, query model.Query) string {
	v := url.Values{}
	v.Add("page", strconv.Itoa(page))
	v.Add("count", strconv.Itoa(count))
	if len(query.Search) > 0 {
		v.Add("query", query.Search)
	}
	if query.MaxLines > 0 {
		v.Add("max-lines", strconv.Itoa(query.MaxLines))
	}

	for _, tag := range query.IncludeTags {
		v.Add("tags", tag)
	}
	for _, tag := range query.ExcludeTags {
		v.Add("exclude-tags", tag)
	}

	return "?" + v.Encode()
}

func QuotesBase(title string, orderBy []string) martini.Handler {
	return func(db model.Model, r render.Render, req *http.Request, isJson IsJson) {
		qs := req.URL.Query()

		search := qs.Get("query")

		page, err := strconv.Atoi(qs.Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		count, err := strconv.Atoi(qs.Get("count"))
		if err != nil || count < 1 {
			count = 20
		}

		maxLines, err := strconv.Atoi(qs.Get("max-lines"))
		if err != nil || maxLines < 1 {
			maxLines = 0
		}

		offset := (page - 1) * count
		query := model.Query{
			Limit:       count,
			Offset:      offset,
			MaxLines:    maxLines,
			Search:      search,
			OrderBy:     orderBy,
			IncludeTags: qs["tags"],
			ExcludeTags: qs["exclude-tags"],
		}
		quotes, err := db.GetQuotes(query)
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)

			RenderError(r, 500, isJson, fmt.Sprint("failed to get quotes", err))
			return
		}

		total, err := db.CountQuotes(query)
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
			RenderError(r, 500, isJson, "failed to get quotes")
			return
		}

		if isJson {
			env := struct {
				Quotes []model.Quote `json:"quotes"`
				Total  int           `json:"total"`
			}{quotes, total}
			r.JSON(200, env)
			return
		}

		maxPage := maxPage(total, count)

		var previousPageURL, nextPageURL string
		if page > 1 {
			previousPageURL = queryString(page-1, count, query)
		}
		if page < maxPage {
			nextPageURL = queryString(page+1, count, query)
		}

		env := quotePageEnv{
			PageEnv: PageEnv{
				Title: title,
				Query: search,
			},
			Quotes:          quotes,
			ShowPagination:  true,
			Count:           count,
			Page:            page,
			PreviousPageURL: previousPageURL,
			NextPageURL:     nextPageURL,
			Total:           total,
			MaxPage:         maxPage,
		}
		r.HTML(200, "quote", env)
	}
}

func RenderError(r render.Render, code int, isJson IsJson, errorMessage string) {
	env := ErrorEnv{ErrorMessage: errorMessage}
	if isJson {
		r.JSON(code, env)
	} else {
		env := ErrorPageEnv{
			PageEnv{Title: "error"},
			env,
		}
		r.HTML(code, "error", env)
	}
}
