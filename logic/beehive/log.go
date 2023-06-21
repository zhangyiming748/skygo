package beehive

import (
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/mysql_model"
)

func GetBeehiveTaskLog(tid int, q string) map[string]interface{} {
	s := mysql.GetSession()
	s.Where("task_id = ?", tid)

	widget := orm.PWidget{}
	widget.SetQueryStr(q)
	widget.AddSorter(*(orm.NewSorter("create_time", orm.DESCENDING)))
	all := widget.PaginatorFind(s, &[]mysql_model.BeehiveLog{})
	return all

}
