package driver

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-sqlite3"
)

type ext struct {
	lib   string
	entry string
}

var vecExts = []ext{
	{"vector0", "sqlite3_vector_init"},
}

var vssExts = []ext{
	{"vss0", "sqlite3_vss_init"},
}

func init() {
	if e := os.Getenv("SQLITE_VSS_EXT_PATH"); e != "" {
		vecExts = append([]ext{{filepath.Join(e, "vector0"), "sqlite3_vector_init"}}, vecExts...)
		vssExts = append([]ext{{filepath.Join(e, "vss0"), "sqlite3_vss_init"}}, vssExts...)
	}
	sql.Register("sqlite-vss", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			vecLoaded := false
			var errs []error
			for _, v := range vecExts {
				err := conn.LoadExtension(v.lib, v.entry)
				if err == nil {
					vecLoaded = true
					break
				}
				errs = append(errs, err)
			}
			if !vecLoaded {
				return fmt.Errorf("vector0 extension load error: %w\nhint: the extension must be located in the current directory or in the directory specified by the environment variable SQLITE_VSS_EXT_PATH.", errors.Join(errs...))
			}

			vssLoaded := false
			errs = nil
			for _, v := range vssExts {
				err := conn.LoadExtension(v.lib, v.entry)
				if err == nil {
					vssLoaded = true
					break
				}
				errs = append(errs, err)
			}
			if !vssLoaded {
				return fmt.Errorf("vss0 extension load error: %w\nhint: the extension must be located in the current directory or in the directory specified by the environment variable SQLITE_VSS_EXT_PATH.", errors.Join(errs...))
			}

			return nil
		},
	})
}
