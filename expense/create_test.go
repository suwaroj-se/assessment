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

	t.Run("Create expenses succecss", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO expenses").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		con := Conn{db}
		var w, g Expenses
		// want := `{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}` + "\n"
		want := `{
			"id": 1,
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath", 
			"tags": ["food", "beverage"]
		}`

		err := json.Unmarshal([]byte(want), &w)
		if err != nil {
			t.Fatal("Error convert 'want' with json.Unmarshall", err)
		}

		if assert.NoError(t, con.CreateExpenseHadler(c)) {
			err = json.Unmarshal(rec.Body.Bytes(), &g)
			if err != nil {
				t.Fatal("Error convert 'rec.Body.Bytes' with json.Unmarshall", err)
			}

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, w, g)
		}
	})
}
