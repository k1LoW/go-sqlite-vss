package driver_test

import (
	"database/sql"
	"testing"

	_ "github.com/k1LoW/go-sqlite-vss"
)

func TestOpen(t *testing.T) {
	db, err := sql.Open("sqlite-vss", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := db.QueryRow("select vss_version();")
	if err := r.Err(); err != nil {
		t.Fatal(err)
	}

	var version string
	if err := r.Scan(&version); err != nil {
		t.Fatal(err)
	}

	want := "v0.1.1"
	if version != want {
		t.Errorf("got %q, want %q", version, want)
	}
}
