package expense

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateExpense(t *testing.T) {

	reqInput := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`

	t.Run("Create expenses succecss", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(reqInput))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("INSERT INTO expenses").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		var w, g Expense
		want := `{
			"id": 1,
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath", 
			"tags": ["food", "beverage"]
		}`

		if err = json.Unmarshal([]byte(want), &w); err != nil {
			t.Fatal("Error convert 'want' with json.Unmarshall", err)
		}

		con := conDB{db}
		if assert.NoError(t, con.CreateExpenseHandler(c)) {
			if err = json.Unmarshal(rec.Body.Bytes(), &g); err != nil {
				t.Fatal("Error convert 'rec.Body.Bytes' with json.Unmarshall", err)
			}

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, w, g)
		}
	})

	t.Run("Create expenses unsuccecss with error row scan", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(reqInput))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		con := conDB{db}
		if assert.NoError(t, con.CreateExpenseHandler(c)) {

			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Bad request post with wrong input", func(t *testing.T) {
		wrongInput := `[]`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(wrongInput))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("INSERT INTO expenses").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		con := conDB{db}
		if assert.NoError(t, con.CreateExpenseHandler(c)) {

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

}
