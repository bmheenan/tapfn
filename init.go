package tapfn

import (
	"errors"
	"fmt"

	"github.com/bmheenan/tapdb"
)

// Init initializes and returns a TapController, which the API server can use to make changes and read data
func Init() (TapController, error) {
	db, errDb := tapdb.Init("user", "pass")
	if errDb != nil {
		return &cnTapdb{}, fmt.Errorf("Could not initialize the database connection: %v", errDb)
	}
	cn := &cnTapdb{
		db: db,
	}
	return cn, errors.New("Not implemented")
}
