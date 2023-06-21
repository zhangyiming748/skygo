package mongo_model

import "github.com/globalsign/mgo/bson"

type FirmWareRtsElf struct {
	ID                bson.ObjectId         `bson:"_id" json:"_id"`
	UploadLogId       string                `bson:"upload_log_id" json:"upload_log_id"`
	ProjectId         string                `bson:"project_id" json:"project_id"`
	TasksId           string                `bson:"tasks_id" json:"tasks_id"`
	FirmwareMd5       string                `bson:"firmware_md5" json:"firmware_md5"`
	Executable        []string              `bson:"executable" json:"executable"`
	FilsharedLibeInfo []string              `bson:"filshared_libe_info" json:"shared_lib"`
	KernelModule      []string              `bson:"kernel_module" json:"kernel_module"`
	IsElf             int                   `bson:"is_elf" json:"is_elf"`
	BinaryHardening   []BinaryHardeningInfo `bson:"binary_hardening" json:"binary_hardening"`
}

type BinaryHardeningInfo struct {
	Result    ResultInfo `bson:"result" json:"result"`
	FileName  string     `bson:"file_name" json:"file_name"`
	FullPath  string     `bson:"full_path" json:"full_path"`
	MagicInfo string     `bson:"magic_info" json:"magic_info"`
}

type ResultInfo struct {
	Nx            string `bson:"nx" json:"nx"`
	Pie           string `bson:"pie" json:"pie"`
	Relro         string `bson:"relro" json:"relro"`
	Rpath         string `bson:"rpath" json:"rpath"`
	Canary        string `bson:"canary" json:"canary"`
	Runpath       string `bson:"runpath" json:"runpath"`
	Stripped      string `bson:"stripped" json:"stripped"`
	Fortified     string `bson:"fortified" json:"fortified"`
	FortifyAble   string `bson:"fortify_able" json:"fortify_able"`
	FortifySource string `bson:"fortify_source" json:"fortify_source"`
}
