package driver_test

import (
	"database/sql"
	"fmt"
	"log"
)

func ExampleOpen() {
	db, err := sql.Open("sqlite-vss", "tmp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := db.QueryRow("select vss_version();")
	if err := r.Err(); err != nil {
		log.Fatal(err)
	}

	var version string
	if err := r.Scan(&version); err != nil {
		log.Fatal(err)
	}

	fmt.Println(version)
}
