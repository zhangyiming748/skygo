package project_file_logic

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"go.uber.org/zap"

	"io/ioutil"
	"skygo_detection/common"
	"skygo_detection/custom_util/clog"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/mysql_model"
)

// 在mongo中检出文档
func FindFileByFileID(fid string) (*mgo.GridFile, []byte, error) {
	fi, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(fid))
	if err != nil {
		clog.Info("field not found")
		return nil, nil, err
	}
	defer fi.Close()

	fileContent, err := ioutil.ReadAll(fi)
	if err != nil {
		return nil, nil, err
	}
	return fi, fileContent, nil
}

// 累加下载次数
func DownloadTimes(fid string) {
	model, err := mysql_model.FindByFileID(fid)
	if err != nil {
		clog.Info("field not found", zap.Any("Info", err))
	}
	err = model.AddDownloadCount(1)
	if err != nil {
		clog.Info("update failure", zap.Any("Info", err))
	}

}
