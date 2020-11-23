package tapfn

import (
	"fmt"

	"github.com/bmheenan/tapdb"
)

// Init initializes and returns a TapController, which the API server can use to make changes and read data.
// It requires the username and password for the database
func Init(user, pass, connName string) (TapController, error) {
	db, errDb := tapdb.Init(user, pass, connName)
	if errDb != nil {
		return &cnTapdb{}, fmt.Errorf("Could not initialize the database connection: %v", errDb)
	}
	cn := &cnTapdb{
		db: db,
	}
	return cn, nil
}
