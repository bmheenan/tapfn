package tapfn

import (
	"errors"
	"fmt"
)

func (cn *cnTapdb) ClearDomain(domain string) error {
	if domain == "" {
		return errors.New("Domain cannot be blank")
	}
	errTSPC := cn.db.ClearThreadStkHierLinks(domain)
	if errTSPC != nil {
		return fmt.Errorf("Could not clear thread stakeholder parent child hierarchy links: %v", errTSPC)
	}
	errSk := cn.db.ClearThreadStkLinks(domain)
	if errSk != nil {
		return fmt.Errorf("Could not clear stakeholders: %v", errSk)
	}
	errTPC := cn.db.ClearThreadHierLinks(domain)
	if errTPC != nil {
		return fmt.Errorf("Could not clear thread parent child relationships: %v", errTPC)
	}
	errTh := cn.db.ClearThreads(domain)
	if errTh != nil {
		return fmt.Errorf("Could not clear threads: %v", errTh)
	}
	errPPC := cn.db.ClearStkHierLinks(domain)
	if errPPC != nil {
		return fmt.Errorf("Could not clear stakeholder parent child relationship: %v", errPPC)
	}
	errPT := cn.db.ClearStks(domain)
	if errPT != nil {
		return fmt.Errorf("Could not clear stakeholders: %v", errPT)
	}
	return nil
}
