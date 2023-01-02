package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (con *conDB) GetExpenseHandlerByID(c echo.Context) error {
	row := con.DB.QueryRow("SELECT id, title, amount, note, tags FROM expenses WHERE id=$1", c.Param("id"))
	var ex Expenses
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
