package expense

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGETExpenseByID(t *testing.T) {

	expected := `{
		"id": 1,
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`

	t.Run("Get expenses by id succecss", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
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

		mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
				AddRow(1, "strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})))

		var w, g Expense
		if err = json.Unmarshal([]byte(expected), &w); err != nil {
			t.Fatal("Error convert 'expected' with json.Unmarshall", err)
		}

		con := conDB{db}
		p, _ := strconv.Atoi(c.Param("id"))
		if assert.NoError(t, con.GetExpenseHandlerByID(c)) {
			if err = json.Unmarshal(rec.Body.Bytes(), &g); err != nil {
				t.Fatal("Error convert 'rec.Body.Bytes' with json.Unmarshall", err)
			}

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, w, g)
			assert.Equal(t, w.ID, p)
		}
	})

	t.Run("Error get expenses by id not found", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
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

		mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}))

		con := conDB{db}
		if assert.NoError(t, con.GetExpenseHandlerByID(c)) {

			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})

	t.Run("Error get expenses by id can't scan expenses", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
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

		mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(1))

		con := conDB{db}
		if assert.NoError(t, con.GetExpenseHandlerByID(c)) {

			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestGETAllExpense(t *testing.T) {

	expected := `[
		{
			"id": 1,
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		},
		{
			"id": 2,
			"title": "iPhone 14 Pro Max 1TB",
			"amount": 66900,
			"note": "birthday gift from my love",
			"tags": ["gadget"]
		}
	]`

	t.Run("Get all expenses success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/expenses")

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
				AddRow(1, "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"})).
				AddRow(2, "iPhone 14 Pro Max 1TB", 66900, "birthday gift from my love", pq.Array([]string{"gadget"})))

		var w, g []Expense
		if err = json.Unmarshal([]byte(expected), &w); err != nil {
			t.Fatal(err)
		}

		con := conDB{db}
		if assert.NoError(t, con.GetAllExpenseHandler(c)) {
			if err = json.Unmarshal(rec.Body.Bytes(), &g); err != nil {
				t.Fatal("Error convert 'rec.Body.Bytes' with json.Unmarshall", err)
			}

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, w, g)
		}
	})

	t.Run("can't query all expense", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/expenses")

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses")

		con := conDB{db}
		if assert.NoError(t, con.GetAllExpenseHandler(c)) {

			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

}
