package driver

import (
	"database/sql"
	"errors"
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
		vecExts = append(vecExts, ext{filepath.Join(e, "vector0"), "sqlite3_vector_init"})
		vssExts = append(vssExts, ext{filepath.Join(e, "vss0"), "sqlite3_vss_init"})
	}
	sql.Register("sqlite-vss", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			vecLoaded := false
			for _, v := range vecExts {
				if err := conn.LoadExtension(v.lib, v.entry); err == nil {
					vecLoaded = true
					break
				}
			}
			if !vecLoaded {
				return errors.New("vector extension not found. the extension must be located in the current directory or in the directory specified by the environment variable SQLITE_VSS_EXT_PATH. ref: https://github.com/asg017/sqlite-vss/releases")
			}

			vssLoaded := false
			for _, v := range vssExts {
				if err := conn.LoadExtension(v.lib, v.entry); err == nil {
					vssLoaded = true
					break
				}
			}
			if !vssLoaded {
				return errors.New("vss extension not found. the extension must be located in the current directory or in the directory specified by the environment variable SQLITE_VSS_EXT_PATH. ref: https://github.com/asg017/sqlite-vss/releases")
			}

			return nil
		},
	})
}
