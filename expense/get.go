package expense

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/lib/pq"

	"github.com/labstack/echo/v4"
)

func (con *Conn) GetExpenseHadlerByID(c echo.Context) error {
	stmt, err := con.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id=$1;")
	if err != nil {
		log.Fatal("can't prepare query one users statment", err)
	}

	rowID := c.Param("id")
	row := stmt.QueryRow(rowID)
	var ex Expenses
	err = row.Scan(&ex.ID, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "user not found"})
	case nil:
		return c.JSON(http.StatusOK, ex)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
	}

}
