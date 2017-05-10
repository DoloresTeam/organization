package godn

import "fmt"

// DoloresType base DN
func DoloresType(subffix string, isUnit bool) string {
	if isUnit {
		return fmt.Sprintf(`ou=unit,ou=type,%s`, subffix)
	}
	return fmt.Sprintf(`ou=member,ou=type,%s`, subffix)
}

// Permission base DN
func Permission(subffix string, isUnit bool) string {
	if isUnit {
		return fmt.Sprintf(`ou=unit,ou=permission,%s`, subffix)
	}
	return fmt.Sprintf(`ou=member,ou=permission,%s`, subffix)
}

// Role base DN
func Role(subffix string) string {
	return fmt.Sprintf(`ou=role,%s`, subffix)
}

// Member base DN
func Member(subffix string) string {
	return fmt.Sprintf(`ou=member,%s`, subffix)
}

// Unit base DN
func Unit(subffix string) string {
	return fmt.Sprintf(`o=1,ou=unit, %s`, subffix) // TODO 后续支持多公司
}
