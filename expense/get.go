package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (con *conDB) GetExpenseHandlerByID(c echo.Context) error {
	row := con.DB.QueryRow("SELECT id, title, amount, note, tags FROM expenses WHERE id=$1", c.Param("id"))
	var ex Expense
	err := row.Scan(&ex.ID, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expenses not found"})
	case nil:
		return c.JSON(http.StatusOK, ex)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses:" + err.Error()})
	}

}

func (con *conDB) GetAllExpenseHandler(c echo.Context) error {
	row, err := con.DB.Query("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query all expenses:" + err.Error()})
	}

	expenses := []Expense{}
	for row.Next() {
		var ex Expense
		if err := row.Scan(&ex.ID, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags)); err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
		}
		expenses = append(expenses, ex)
	}

	return c.JSON(http.StatusOK, expenses)
}
