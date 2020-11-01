package tapfn

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bmheenan/taps"
)

// NewPersonteam creates a new Personteam with the given data:
//     email: the email of the team or person
//     name: the name of the team or person
//     abbrev: a short (max 3 character) abbrevation of the name. Used for the personteam's icon
//     colorf: the forground color of the personteam's icon
//     colorb: the background colof of the personteam's icon
//     itertiming: what cadence the person or team will plan on
//     parents: the emails of any existing personteams who this personteam is a part of. You cannot create a loop. Leave
//              empty to insert this personteam at the domain's root
func (cn *cnTapdb) NewPersonteam(
	email,
	name,
	abbrev,
	colorf,
	colorb string,
	itertiming taps.IterTiming,
	parents []string,
) error {
	// TODO: Check arguments
	if email == "" {
		return errors.New("Email cannot be blank")
	}
	emailPieces := strings.Split(email, "@")
	if len(emailPieces) != 2 {
		return errors.New("Invalid email")
	}
	errNew := cn.db.NewPersonteam(email, emailPieces[1], name, abbrev, colorf, colorb, itertiming)
	if errNew != nil {
		return fmt.Errorf("Could not insert new Personteam into the database: %v", errNew)
	}
	for _, p := range parents {
		errPC := cn.db.NewPersonteamPC(p, email, emailPieces[1])
		if errPC != nil {
			return fmt.Errorf("Could not make %v a child of %v: %v", email, p, errPC)
		}
	}
	return nil
}
