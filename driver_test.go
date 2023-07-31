package driver_test

import (
	"bufio"
	"context"
	"database/sql"
	"math"
	"os"
	"strings"
	"testing"

	_ "github.com/k1LoW/go-sqlite-vss"
)

func TestOpen(t *testing.T) {
	db, err := sql.Open("sqlite-vss", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	t.Cleanup(func() {
		_ = os.Remove("test.db")
	})

	r := db.QueryRow("select vss_version();")
	if err := r.Err(); err != nil {
		t.Fatal(err)
	}

	var version string
	if err := r.Scan(&version); err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(version, "v") {
		t.Errorf("version should be start with 'v', but got %s", version)
	}
}

// ref: https://github.com/koron/techdocs/blob/main/sqlite-vss-getting-started/doc.md
func TestVectorSearch(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite-vss", "vec.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
		_ = os.Remove("vec.db")
	})

	// Create table
	if _, err := db.ExecContext(ctx, `CREATE TABLE words (
    label  TEXT,
	  vector BLOB
);`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `CREATE VIRTUAL TABLE vss_words USING vss0(
  vector(300)
);`); err != nil {
		t.Fatal(err)
	}

	// read test.vec line by line.
	f, err := os.Open("testdata/test.vec")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = f.Close()
	})
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		row := scanner.Text()
		splitted := strings.Split(row, " ")
		word := splitted[0]
		vec := "[" + strings.Join(splitted[1:], ",") + "]"
		if _, err := db.ExecContext(ctx, `INSERT INTO words (label, vector) VALUES (?, ?);`, word, vec); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := db.ExecContext(ctx, `UPDATE words SET vector = vector_to_blob(vector_from_json(vector));`); err != nil {
		t.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, `INSERT INTO vss_words(rowid, vector) SELECT rowid, vector FROM words;`); err != nil {
		t.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, `VACUUM;`); err != nil {
		t.Fatal(err)
	}

	rows, err := db.QueryContext(ctx, `SELECT w.label, v.distance FROM vss_words AS v
  JOIN words AS w ON w.rowid = v.rowid
  WHERE vss_search(
    v.vector,
    vss_search_params(
      (select vector from words where label = 'food'),
      10
    )
  );`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	got := map[string]float64{}
	for rows.Next() {
		var (
			label string
			dist  float64
		)
		if err := rows.Scan(&label, &dist); err != nil {
			t.Fatal(err)
		}
		got[label] = dist
	}

	if want := 10; len(got) != want {
		t.Errorf("len(got) = %d, want %d", len(got), want)
	}

	if want := 0.0; got["food"] != want {
		t.Errorf("got[food] = %f, want %f", got["food"], want)
	}

	if want := 1.828758; math.Round(got["health"]*10000)/10000 != math.Round(want*10000)/10000 {
		t.Errorf("got[health] = %f, want %f", got["health"], want)
	}
}
