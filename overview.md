# Architectural overview

## Pseudocode to describe high-level logic

**NewThread**(name, owner, iteration, cost, parents, children):
    thread = db.NewThread(name, owner, iteration, cost)
    **AddStakeholderToThread**(thread, owner)
    for parent in parents:
        **LinkThreads**(parent, thread)
    for child in children:
        **LinkThreads**(thread, child) 

**LinkThreads**(parent, child):
    if **wouldMakeLoop**(parent, child):
        error
    iter = **iterResulting**(parent.Owner.Cadence, child.Iter)
    ord = db.GetOrderBeforeForParent(parent, iter, MAX)
    db.LinkThreads(parent, child, iter, ord)
    for ancestor in db.GetAncestors(parent):
        **recalcTotalCost**(ancestor)
    **balanceParent**(parent, iter)
    **recalcAllStakeholderCosts**(parent)

**UnlinkThreads**(parent, child):
    db.UnlinkThreads(parent, child)
    for ancestor in db.GetAncestors(parent):
        **recalcTotalCost**(ancestor)
    **recalcAllStakeholderCosts**(parent)

**AddStakeholderToThread**(thread, stakeholder):
    iter = **iterResulting**(stakeholder.Cadence, thread.Iter)
    db.AddStakeholderToThread(thread, stakeholder, iter, MAX)
    **balanceStakeholder**(stakeholder, iter)
    **recalcAllStakeholderCosts**(thread)

**RemoveStakeholderFromThread**(thread, stakeholder):
    if thread.Owner = stakeholder:
        error
    db.RemoveStakeholderFromThread(thread, stakeholder)
    **recalcAllStakeholderCosts**(thread)

**SetThreadIter**(thread, iter):
    for descendant in db.GetDescendants(thread):
        if **iterResulting**(thread.Owner.Cadence, descendant.Iter) != thread.Iter:
            break
        iter = **iterResulting**(descendant.Owner.Cadence, iter)
        db.SetIter(descendant, iter)
        for parent in descendant.Parents:
            iter = **iterResulting**(parent.Owner.Cadence, iter)
            if iter = descendants.Parents[parent].Iter:
                break
            else if iter < descendants.Parents[parent].Iter:
                place = MoveToEnd
            else:
                place = MoveToStart
            db.SetIterForParent(descendant, parent, iter)
            db.MoveThreadForParent(descendant, nil, parent, place)
        for stakeholder in descendant.Stakeholders:
            iter = **iterResulting**(stakeholder.Cadence, iter)
            if iter = descendant.Stakeholders[stakeholder].Iter:
                break
            if iter < descendant.Stakeholders[stakeholder].Iter:
                place = MoveToEnd
            else:
                place = MoveToStart
            db.SetIterForStakeholder(descendant, stakeholder, iter)
            db.MoveThreadForStakeholder(descendant, nil, stakeholder, place)
        **recalcAllStakeholderCosts**(descendant)

**SetThreadOwner**(thread, owner):
    **AddStakeholderToThread(thread, owner)
    db.SetThreadOwner(thread, owner)
    **recalcAllStakeholderCosts**(thread)

**wouldMakeLoop**(parent, child) bool:
    ancestors = db.GetAncestors(parent)
    if child exists in ancestors:
        return true
    else:
        return false

**recalcTotalCost**(thread):
    for descendant in db.GetDescendants(thread):
        += sum
    db.SetTotalCost(thread, sum)

**recalcAllStakeholderCosts**(thread):
    for ancestor in db.GetAncestors:
        for stakeholder in ancestor.Stakeholders:
            **recalcStakeholderCost**(ancestor, stakeholder)

**recalcStakeholderCost**(thread, stakeholder):
    members = db.GetStakeholderDescendants(stakeholder)
    for descendant in db.GetDescendants(thread):
        if descendant.Owner exists in members and **iterResulting**(stakeholder.Cadence, descendant.Iter) = thread.Iter:
            += sum
    db.SetCostForStakeholder(thread, stakeholder, sum)

**iterResulting**(cadence, iter) iter:
    ...

**balanceParent**(parent, iter):
    ...

**balanceStakeholder**(stakeholder, iter):
    ...

## Things to rationalize still

* Where to balance stk and parent