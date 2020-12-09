package tapfn

func (cn *cnTapdb) ThreadSetName(thread int64, name string) {
	cn.db.SetName(thread, name)
}

func (cn *cnTapdb) ThreadSetDesc(thread int64, desc string) {
	cn.db.SetDesc(thread, desc)
}
