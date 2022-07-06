package tree

func isNum(v valueKind) bool {
	if v == Integer || v == Decimal {
		return true
	}
	return false
}
