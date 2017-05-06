package godn

import "fmt"

func DoloresType(subffix string, isUnit bool) string {
	if isUnit {
		return fmt.Sprintf(`ou=unit,ou=type,%s`, subffix)
	}
	return fmt.Sprintf(`ou=person,ou=type,%s`, subffix)
}

func Permission(subffix string, isUnit bool) string {
	if isUnit {
		return fmt.Sprintf(`ou=unit,ou=permission,%s`, subffix)
	}
	return fmt.Sprintf(`ou=person,ou=permission,%s`, subffix)
}

func Role(subffix string) string {
	return fmt.Sprintf(`ou=role,%s`, subffix)
}

func Person(subffix string) string {
	return fmt.Sprintf(`ou=person,%s`, subffix)
}

func Unit(subffix string) string {
	return fmt.Sprintf(`oid=1,ou=unit,%s`, subffix)
}
