# go-sqlite-vss

`go-sqlite-vss` is a ["SQLite + SQLite Vector Similarity Search extension"](https://github.com/asg017/sqlite-vss) driver for database/sql package.

## Usage

Install [vector0 and vss0](https://github.com/asg017/sqlite-vss/releases).

``` console
# An example installation
$ gh ext install k1LoW/gh-setup
$ gh setup --repo github.com/asg017/sqlite-vss --bin-dir ${SQLITE_VSS_EXT_PATH} --match sqlite-vss-v.*-loadable --skip-content-type-check
```

And then, use `sqlite-vss` as the driver name.

``` go
package main

import (
	"database/sql"
	"fmt"
	"log"
)

func main() {
	db, err := sql.Open("sqlite-vss", "test.db")
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
	// Output: v0.1.1
}
```
