package tapfn

import "fmt"

func (cn *cnTapdb) ClearDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("Domain cannot be blank: %w", ErrBadArgs)
	}
	errSk := cn.db.ClearStakeholders(domain)
	if errSk != nil {
		return fmt.Errorf("Could not clear stakeholders: %v", errSk)
	}
	errTPC := cn.db.ClearThreadsPC(domain)
	if errTPC != nil {
		return fmt.Errorf("Could not clear thread parent child relationships: %v", errTPC)
	}
	errTh := cn.db.ClearThreads(domain)
	if errTh != nil {
		return fmt.Errorf("Could not clear threads: %v", errTh)
	}
	errPPC := cn.db.ClearPersonteamsPC(domain)
	if errPPC != nil {
		return fmt.Errorf("Could not clear personteam parent child relationship: %v", errPPC)
	}
	errPT := cn.db.ClearPersonteams(domain)
	if errPT != nil {
		return fmt.Errorf("Could not clear personteams: %v", errPT)
	}
	return nil
}
