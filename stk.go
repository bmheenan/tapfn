package tapfn

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bmheenan/tapdb"
	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) Stk(email string) (stk taps.Stakeholder, err error) {
	stkp, err := cn.db.GetStk(email)
	if errors.Is(err, tapdb.ErrNotFound) {
		err = fmt.Errorf("Stakeholder not found: %w", ErrNotFound)
		return
	}
	if err != nil {
		err = fmt.Errorf("Could not get stakeholder: %v", err)
		return
	}
	stk = *stkp
	return
}

func (cn *cnTapdb) StkNew(
	email,
	name,
	abbrev,
	colorf,
	colorb string,
	cadence taps.Cadence,
	parents []string,
) error {
	// TODO check arguments
	// TODO return special error if it already exists
	if email == "" {
		return errors.New("Email cannot be blank")
	}
	ePcs := strings.Split(email, "@")
	if len(ePcs) != 2 {
		return errors.New("Invalid email")
	}
	errNew := cn.db.NewStk(email, ePcs[1], name, abbrev, colorf, colorb, cadence)
	if errNew != nil {
		return fmt.Errorf("Could not insert new Stakeholder into the database: %v", errNew)
	}
	for _, p := range parents {
		errPC := cn.db.NewStkHierLink(p, email, ePcs[1])
		if errPC != nil {
			return fmt.Errorf("Could not make %v a child of %v: %v", email, p, errPC)
		}
	}
	return nil
}

func (cn *cnTapdb) StksByDomain(domain string) (stks []taps.StkInHier) {
	stksp, err := cn.db.GetStksForDomain(domain)
	if err != nil {
		panic(fmt.Sprintf("Could not get stakeholders for domain %v: %v", domain, err))
	}
	for _, s := range stksp {
		stks = append(stks, *s)
	}
	return
}
