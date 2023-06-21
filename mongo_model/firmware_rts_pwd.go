package mongo_model

import (
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

type FirmWareRtsPwd struct {
	FullPath string    `bson:"fullpath" json:"fullpath"`
	List     []PwdInfo `bson:"list" json:"list"`
}

type PwdInfo struct {
	Gid   string `bson:"gid" json:"gid"`
	Uid   string `bson:"uid" json:"uid"`
	Home  string `bson:"home" json:"home"`
	Name  string `bson:"name" json:"name"`
	Pass  string `bson:"pass" json:"pass"`
	User  string `bson:"user" json:"user"`
	Shell string `bson:"shell" json:"shell"`
}

func (this *FirmWareRtsPwd) Create(rawinfo qmap.QM) (*FirmWareRtsPwd, error) {
	if err := mongo.NewMgoSession(common.MC_EVALUATE_VULNERABILITY).Insert(rawinfo); err == nil {
		return this, nil
	} else {
		return nil, err
	}

}
