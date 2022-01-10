package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func compareBytesAndInterface(b []byte, i interface{}, t *testing.T) {

	t.Helper()

	interfaceB, err := json.Marshal(i)
	if err != nil {
		t.Fatalf("unable to marshal expected struct to JSON; error: %v ", err.Error())
	}

	if r := bytes.Compare(interfaceB, b); r != 0 {
		t.Fatalf("bytes are different from each other")
	}
}

// positive test case
func TestAppPost(t *testing.T) {

	//initialising new mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error in initialising sql mock: %s", err.Error())
	}
	defer db.Close()

	apiObj := api{db: db}

	r, err := http.NewRequest("GET", "http://localhost/posts", nil)
	if err != nil {
		t.Fatalf("error in initialising new request: %s", err.Error())
	}
	w := httptest.NewRecorder()

	expectedPosts := []*post{
		{ID: 1, Title: "Title for 1", Body: "Body for 1"},
		{ID: 2, Title: "Title for 2", Body: "Body for 2"},
		{ID: 3, Title: "Title for 3", Body: "Body for 3"},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "post"})

	for _, post := range expectedPosts {
		rows.AddRow(post.ID, post.Title, post.Body)
	}

	mock.ExpectQuery("^SELECT (.+) FROM posts$").WillReturnRows(rows)

	apiObj.posts(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("not OK status received from httpWriter")
	}

	resp := w.Body.Bytes()

	compareBytesAndInterface(resp, struct{ Posts []*post }{expectedPosts}, t)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("exepctations were not met from the mockings; %v", err.Error())
	}
}

// negative test case
func TestFailedAppPost(t *testing.T) {

	//initialising new mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error in initialising sql mock: %s", err.Error())
	}
	defer db.Close()

	apiObj := api{db: db}

	r, err := http.NewRequest("GET", "http://localhost/posts", nil)
	if err != nil {
		t.Fatalf("error in initialising new request: %s", err.Error())
	}
	w := httptest.NewRecorder()

	errorToThrow := "failed to execute get posts query"
	mock.ExpectQuery("^SELECT (.+) FROM posts$").WillReturnError(fmt.Errorf(errorToThrow))

	apiObj.posts(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("failed status not received from httpWriter")
	}

	resp := w.Body.Bytes()

	compareBytesAndInterface(resp, struct {
		Error string
	}{Error: "failed to fetch posts: " + errorToThrow}, t)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("exepctations were not met from the mockings; %v", err.Error())
	}
}
