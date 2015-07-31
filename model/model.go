package model

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	db *sql.DB
}

/*
NewModel creates a new model with a DB connection to the give dbPath (currently sqlite)
*/
func NewModel(dbPath string) (m Model, err error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}
	m.db = db
	return
}

func (m Model) GetQuote(id int) (quote Quote, err error) {
	rawQ := rawQuote{}

	err = m.db.QueryRow(
		"SELECT id, text, score, time_created, is_offensive, is_nishbot from quote where id = ?",
		id,
	).Scan(
		&rawQ.ID,
		&rawQ.Text,
		&rawQ.Score,
		&rawQ.TimeCreated,
		&rawQ.IsOffensive,
		&rawQ.IsNishbot,
	)
	if err != nil {
		return
	}

	return toQuote(rawQ), nil
}

func (m Model) GetQuotes(q Query) (quotes []Quote, err error) {
	query := strings.Join(
		[]string{
			"SELECT id, text, score, time_created, is_offensive, is_nishbot",
			"FROM quote",
			q.toSQL(),
		},
		"\n",
	)
	rows, err := m.db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		rawQ := rawQuote{}
		err = rows.Scan(
			&rawQ.ID,
			&rawQ.Text,
			&rawQ.Score,
			&rawQ.TimeCreated,
			&rawQ.IsOffensive,
			&rawQ.IsNishbot,
		)
		if err != nil {
			return
		}
		quotes = append(quotes, toQuote(rawQ))
	}
	return
}

func (m Model) CountQuotes(search string) (count int, err error) {
	err = m.db.QueryRow(
		strings.Join(
			[]string{
				"SELECT count(*)",
				"FROM quote",
				searchWhereClause(search),
			},
			"\n",
		),
	).Scan(&count)
	return
}

func (m Model) CountAllQuotes() (count int, err error) {
	return m.CountQuotes("")
}

func (m Model) AddQuote(q Quote) (err error) {
	rawQ := fromQuote(q)
	_, err = m.db.Exec(
		"INSERT INTO quote(text, score, time_created, is_offensive, is_nishbot) values(?, ?, ?, ?, ?)",
		rawQ.Text,
		rawQ.Score,
		time.Now().Unix(),
		rawQ.IsOffensive,
		rawQ.IsNishbot,
	)
	return
}

func (m Model) DeleteQuote(id int) (err error) {
	quote, err := m.GetQuote(id)
	if err != nil {
		return
	}
	_, err = m.db.Exec("DELETE FROM quote WHERE id = ?", id)
	if err != nil {
		return
	}

	rawQ := fromQuote(quote)
	_, err = m.db.Exec(
		"INSERT INTO deleted_quote(id, text, score, time_created, is_offensive, is_nishbot) values(?, ?, ?, ?, ?, ?)",
		rawQ.ID,
		rawQ.Text,
		rawQ.Score,
		rawQ.TimeCreated,
		rawQ.IsOffensive,
		rawQ.IsNishbot,
	)
	return
}

func (m Model) VoteQuote(id int, vote int) (err error) {
	q, err := m.GetQuote(id)
	if err != nil {
		return
	}

	newScore := q.Score + vote

	_, err = m.db.Exec("UPDATE quote SET score = ? where id = ?", newScore, id)
	return
}

func (m Model) Close() error {
	return m.db.Close()
}
