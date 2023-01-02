package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (con *conDB) PutExpenseHandlerByID(c echo.Context) error {
	var ex Expenses
	if err := c.Bind(&ex); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if ex.Title == "" || ex.Amount == 0 || ex.Note == "" || ex.Tags == nil || len(ex.Tags) == 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Missing values:"})
	}

	paramID := c.Param("id")
	row := con.DB.QueryRow("SELECT id FROM expenses WHERE id=$1", paramID)
	err := row.Scan(&ex.ID)

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expenses not found"})

	case nil:
		row := con.DB.QueryRow("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id", ex.Title, ex.Amount, ex.Note, pq.Array(&ex.Tags), paramID)

		if err := row.Scan(&ex.ID); err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, ex)

	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses:" + err.Error()})
	}

}
