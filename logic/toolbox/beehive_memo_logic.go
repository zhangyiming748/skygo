package toolbox

import (
	"fmt"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/mysql_model"
	"time"

	"github.com/gin-gonic/gin"
)

type Memo struct{}

func (m Memo) Save(ctx *gin.Context, taskId int, Content string) (int64, error) {
	memo := mysql_model.BeehiveMemo{}
	s := mysql.GetSession()
	b, err := s.Where("task_id=?", taskId).Get(&memo)
	if err != nil {
		return 0, err
	}
	memo.Content = Content
	if !b {
		memo.TaskId = taskId
		memo.CreateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
		i, err := memo.Create()
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return memo.Update()
}
