package mongo_model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/orm_mongo"
)

type FirmWareCategory struct {
	Id                  primitive.ObjectID `bson:"_id" json:"_id"`
	MasterId            string             `bson:"master_id" json:"master_id"`                       //source表Id
	ProjectId           int                `bson:"project_id" json:"project_id"`                     //工程的ID
	TaskId              int                `bson:"task_id" json:"task_id"`                           //任务ID
	TemplateId          int                `bson:"template_id" json:"template_id"`                   //模板ID
	InitFiles           int                `bson:"init_files" json:"init_files"`                     //是或否存在 初始化的检测（1：存在，0：不存在）
	ElfScanner          int                `bson:"elf_scanner" json:"elf_scanner"`                   //是或否存在 elf的检测（1：存在，0：不存在）
	BinaryHardening     int                `bson:"binary_hardening" json:"binary_hardening"`         //是或否存在 二进制的检测（1：存在，0：不存在）
	SymbolsXrefs        int                `bson:"symbols_xrefs" json:"symbols_xrefs"`               //是或否存在 符号交叉的检测（1：存在，0：不存在）
	VersionScanner      int                `bson:"version_scanner" json:"version_scanner"`           //是或否存在 CVE的检测（1：存在，0：不存在）
	CertificatesScanner int                `bson:"certificates_scanner" json:"certificates_scanner"` //是或否存在 签名的检测（1：存在，0：不存在）
	LeaksScanner        int                `bson:"leaks_scanner" json:"json"`                        //是或否存在 敏感信息的检测（1：存在，0：不存在）
	PasswordScanner     int                `bson:"password_scanner" json:"password_scanner"`         //是或否存在 密码信息的检测（1：存在，0：不存在）
	ApkInfo             int                `bson:"apk_info" json:"apk_info"`                         //是或否存在 密码信息的检测（1：存在，0：不存在）
	ApkCommonVul        int                `bson:"apk_common_vul" json:"apk_common_vul"`             //是或否存在 密码信息的检测（1：存在，0：不存在）
	ApkSensitiveInfo    int                `bson:"apk_sensitive_info" json:"apk_sensitive_info"`     //是或否存在 密码信息的检测（1：存在，0：不存在）
	LinuxBasicAudit     int                `bson:"linux_basic_audit" json:"linux_basic_audit"`       //是或否存在 Linux固件基线检测（1：存在，0：不存在）
	IsElf               int                `bson:"is_elf" json:"is_elf"`                             //是或否存在 二进制文件增强的检测（1：存在，0：不存在）
}

func (this *FirmWareCategory) Create(rawInfo qmap.QM) (*FirmWareCategory, error) {
	if val, has := rawInfo.TryString("master_id"); has {
		this.MasterId = val
	}
	if val, has := rawInfo.TryInt("project_id"); has {
		this.ProjectId = val
	}
	if val, has := rawInfo.TryInt("task_id"); has {
		this.TaskId = val
	}
	if val, has := rawInfo.TryInt("template_id"); has {
		this.TemplateId = val
	}
	if _, has := rawInfo.TryString("init_files"); has {
		this.InitFiles = 1
	}
	if _, has := rawInfo.TryString("elf_scanner"); has {
		this.ElfScanner = 1
		this.IsElf = 1
	}
	if _, has := rawInfo.TryString("binary_hardening"); has {
		this.BinaryHardening = 1
	}
	if _, has := rawInfo.TryString("symbols_xrefs"); has {
		this.SymbolsXrefs = 1
	}
	if _, has := rawInfo.TryString("version_scanner"); has {
		this.VersionScanner = 1
	}
	if _, has := rawInfo.TryString("certificates_scanner"); has {
		this.CertificatesScanner = 1
	}
	if _, has := rawInfo.TryString("leaks_scanner"); has {
		this.LeaksScanner = 1
	}
	if _, has := rawInfo.TryString("password_scanner"); has {
		this.PasswordScanner = 1
	}
	if _, has := rawInfo.TryString("linux_basic_audit"); has {
		this.LinuxBasicAudit = 1
	}
	if _, has := rawInfo.TryString("apk_info"); has {
		this.ApkInfo = 1
	}
	if _, has := rawInfo.TryString("apk_sensitive_info"); has {
		this.ApkSensitiveInfo = 1
	}
	if _, has := rawInfo.TryString("apk_common_vul"); has {
		this.ApkCommonVul = 1
	}
	this.Id = primitive.NewObjectID()

	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_FIRMWARE_RTS_CATEGORY)
	_, err := coll.InsertOne(context.Background(), this)
	return this, err
}
