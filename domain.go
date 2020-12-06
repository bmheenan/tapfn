package tapfn

func (cn *cnTapdb) DomainClear(domain string) {
	if domain == "" {
		panic("Domain cannot be blank")
	}
	cn.db.ClearThreadStkLinks(domain)
	cn.db.ClearThreadHierLinks(domain)
	cn.db.ClearThreads(domain)
	cn.db.ClearStkHierLinks(domain)
	cn.db.ClearStks(domain)
}
