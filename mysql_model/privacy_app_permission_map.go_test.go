package mysql_model

import "testing"

func TestFindNameByPermissionId(t *testing.T) {
	TestingInit()
	pid := 19
	pname, _ := FindNameByPermissionId(pid)
	t.Log(pname)
}
