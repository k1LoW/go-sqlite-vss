# go-sqlite-vss

`go-sqlite-vss` is a ["SQLite + SQLite Vector Similarity Search extension"](https://github.com/asg017/sqlite-vss) driver for database/sql package.

## Usage

Install [vector0 and vss0](https://github.com/asg017/sqlite-vss/releases).

``` console
$ # An example installation
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

## Test data

- [testdata/test.vec](testdata/test.vec) is created using data https://fasttext.cc/docs/en/crawl-vectors.html .
    - License
        - The word vectors are distributed under the [Creative Commons Attribution-Share-Alike License 3.0](https://creativecommons.org/licenses/by-sa/3.0/).
    - References
        - [E. Grave*, P. Bojanowski*, P. Gupta, A. Joulin, T. Mikolov, Learning Word Vectors for 157 Languages](https://arxiv.org/abs/1802.06893)
