package mongo_model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EvaluateAbandon struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	EvaluateTargetId primitive.ObjectID `bson:"evaluate_target_id"` // 评估目标id
	EvaluateTime     int                `bson:"evaluate_time"`      // 评估次数
	EvaluateGoal     string             `bson:"evaluate_goal"`      // 评估目标
	EvaluateProcess  string             `bson:"evaluate_process"`   // 评估过程
	EvaluateResult   string             `bson:"evaluate_result"`    // 评估结论
	HasVul           int                `bson:"has_vul"`            // 是否存在漏洞(0:不存在,1:存在)
	StartTime        int                `bson:"start_time"`         // 评估开始时间
	EndTime          int                `bson:"end_time"`           // 评估结束时间
	CreateTime       int                `bson:"create_time"`        // 评估创建时间
}
