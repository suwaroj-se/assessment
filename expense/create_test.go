package expense

import (
	"testing"
)

func TestCreate(t *testing.T) {
	jsonInput := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`
	want := `{
		"id": "1",
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`

	got := Create(jsonInput)

	if got != want {
		t.Errorf("Missmatch")
	}
}
