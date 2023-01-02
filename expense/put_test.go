package expense

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	// "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPutExpenseByID(t *testing.T) {

	reqInput := `{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`

	t.Run("Put method to update expenses succecss", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(reqInput))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/expenses/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM expenses").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("UPDATE expenses").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		var w, g Expenses
		want := `{
			"id": 1,
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`

		if err = json.Unmarshal([]byte(want), &w); err != nil {
			t.Fatal("Error convert 'want' with json.Unmarshall", err)
		}

		con := conDB{db}
		if assert.NoError(t, con.PutExpenseHandlerByID(c)) {
			if err = json.Unmarshal(rec.Body.Bytes(), &g); err != nil {
				t.Fatal("Error convert 'rec.Body.Bytes' with json.Unmarshall", err)
			}

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, w, g)
		}
	})

	t.Run("Wrong input expenses struct request to Put method", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(`-`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/expenses/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM expenses").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("UPDATE expenses").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		con := conDB{db}
		if assert.NoError(t, con.PutExpenseHandlerByID(c)) {

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Missing values to Put method", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(``))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/expenses/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM expenses").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("UPDATE expenses").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		want := `{"message":"Missing values:"}` + "\n"

		con := conDB{db}
		if assert.NoError(t, con.PutExpenseHandlerByID(c)) {

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, want, rec.Body.String())
		}
	})

	t.Run("Expenses not found", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(reqInput))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/expenses/:id")
		c.SetParamNames("id")
		c.SetParamValues("")

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM expenses").WillReturnRows(sqlmock.NewRows([]string{"id"}))

		want := `{"message":"Expenses not found"}` + "\n"

		con := conDB{db}
		if assert.NoError(t, con.PutExpenseHandlerByID(c)) {

			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.Equal(t, want, rec.Body.String())
		}
	})

}
