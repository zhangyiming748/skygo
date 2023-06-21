package mongo_model

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
)

type Project struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty"`
	Name                string             `bson:"name"`                 // 项目名称
	StartTime           int                `bson:"start_time"`           // 项目开始时间
	EndTime             int                `bson:"end_time"`             // 项目结束时间
	Company             string             `bson:"company"`              // 项目所属车厂
	Brand               string             `bson:"brand"`                // 车型品牌
	CodeName            string             `bson:"code_name"`            // 车型代号
	ManagerId           int                `bson:"manager_id"`           // 项目经理id
	MemberIds           []int              `bson:"member_ids"`           // 项目成员id
	AllUsers            []int              `bson:"all_users"`            // 所有项目人员(用于方便项目人员查询)
	Amount              int                `bson:"amount"`               // 项目金额(单位:万元)
	Description         string             `bson:"description"`          // 项目描述
	EvaluateRequirement string             `bson:"evaluate_requirement"` // 评估要求
	Status              int                `bson:"status"`               // 项目状态
	IsDeleted           int                `bson:"is_deleted"`           // 软删除  0未删除，1已删除
	UpdateTime          int64              `bson:"update_time"`
	CreateTime          int64              `bson:"create_time"`
}

func (this *Project) Create(ctx context.Context, rawInfo *qmap.QM, userId int64) (*Project, error) {
	// 检查项目名称是否重复
	this.Name = rawInfo.MustString("name")

	// 使用Widget.Get查询
	info, err := orm_mongo.NewWidgetWithCollectionName(common.MC_PROJECT).SetParams(qmap.QM{"e_name": this.Name}).Get()
	if err == nil && info["_id"] != nil {
		return nil, errors.New("项目名称已存在")
	}

	this.Id = primitive.NewObjectID()
	this.StartTime = rawInfo.MustInt("start_time")
	this.EndTime = rawInfo.MustInt("end_time")
	this.Company = rawInfo.MustString("company")
	this.Brand = rawInfo.MustString("brand")
	this.CodeName = rawInfo.String("code_name")
	this.Description = rawInfo.MustString("description")
	this.EvaluateRequirement = rawInfo.MustString("evaluate_requirement")
	this.ManagerId = int(userId)
	this.Amount = rawInfo.MustInt("amount")
	this.Status = common.PS_NEW
	this.IsDeleted = common.PSD_DEFAULT
	this.UpdateTime = custom_util.GetCurrentMilliSecond()
	this.CreateTime = custom_util.GetCurrentMilliSecond()
	this.MemberIds = []int{}
	mbIds := rawInfo.MustSlice("member_ids")
	for _, val := range mbIds {
		var t int
		switch val.(type) {
		case float64:
			t = int(val.(float64))
		case int64:
			t = int(val.(int64))
		case int32:
			t = int(val.(int32))
		case string:
			if temp, err := strconv.Atoi(val.(string)); err == nil {
				t = temp
			} else {
				panic(err)
			}
		default:
			t = val.(int)
		}
		this.MemberIds = append(this.MemberIds, t)
	}
	this.AllUsers = []int{this.ManagerId}
	this.AllUsers = append(this.AllUsers, this.MemberIds...)

	if _, err := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_PROJECT).InsertOne(ctx, this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *Project) Update(projectId string, rawInfo qmap.QM) (*Project, error) {
	_id, _ := primitive.ObjectIDFromHex(projectId)
	params := qmap.QM{
		"e__id": _id,
	}

	widget := orm_mongo.NewWidgetWithCollectionName(common.MC_PROJECT).SetParams(params)
	if err := widget.One(&this); err == nil {
		if status, has := rawInfo.TryInt("status"); has {
			if _, err := this.ChangeStatus(this.Id.Hex(), status); err != nil {
				return nil, err
			}
		}

		if val, has := rawInfo.TryString("name"); has {
			this.Name = val
			info, err := orm_mongo.NewWidgetWithCollectionName(common.MC_PROJECT).SetParams(qmap.QM{"e_name": this.Name}).Get()
			if err == nil && info["_id"] != this.Id {
				return nil, errors.New("项目名称已存在")
			}
		}
		if val, has := rawInfo.TryInt("start_time"); has {
			this.StartTime = val
		}
		if val, has := rawInfo.TryInt("end_time"); has {
			this.EndTime = val
		}
		if val, has := rawInfo.TryString("company"); has {
			this.Company = val
		}
		if val, has := rawInfo.TryString("brand"); has {
			this.Brand = val
		}
		if val, has := rawInfo.TryString("code_name"); has {
			this.CodeName = val
		}
		if val, has := rawInfo.TryString("description"); has {
			this.Description = val
		}
		if val, has := rawInfo.TryString("evaluate_requirement"); has {
			this.EvaluateRequirement = val
		}
		if val, has := rawInfo.TryInt("amount"); has {
			this.Amount = val
		}
		if val, has := rawInfo.TryInt("is_deleted"); has {
			this.IsDeleted = val
		}
		if members, has := rawInfo.TrySlice("member_ids"); has {
			tempMember := []int{this.ManagerId}
			for _, val := range members {
				var t int
				switch val.(type) {
				case float64:
					t = int(val.(float64))
				case int64:
					t = int(val.(int64))
				case int32:
					t = int(val.(int32))
				case string:
					if temp, err := strconv.Atoi(val.(string)); err == nil {
						t = temp
					} else {
						panic(err)
					}
				default:
					t = val.(int)
				}
				if custom_util.InIntSlice(t, tempMember) == false {
					tempMember = append(tempMember, t)
				}
			}
			this.MemberIds = tempMember
			this.AllUsers = []int{this.ManagerId}
			this.AllUsers = append(this.AllUsers, this.MemberIds...)
		}

		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_PROJECT)
		if _, err := coll.UpdateOne(context.Background(), bson.M{"_id": this.Id}, bson.M{"$set": this}); err != nil {
			return nil, err
		}
	} else {
		panic(err)
	}
	return this, nil
}

func (this *Project) ChangeStatus(projectId string, status int) (bool, error) {
	_id, _ := primitive.ObjectIDFromHex(projectId)
	params := qmap.QM{
		"e__id": _id,
	}
	widget := orm_mongo.NewWidgetWithParams(common.MC_PROJECT, params)
	err := widget.One(&this)
	if err == nil {
		if status == common.PS_TEST {
			// 当有任务进入测试阶段且当前为创建阶段，才触发项目状态为测试阶段
			if this.Status == common.PS_NEW {
				this.Status = common.PS_TEST
			}
		} else if status == common.PS_COMPLETE {
			// 检查项目状态是否为测试阶段
			if this.Status != common.PS_TEST {
				return false, errors.New("项目状态切换失败，项目当前阶段不能切换到已完成")
			}

			// 检查是否所有任务都已完成
			params := qmap.QM{
				"ne_status":    common.PTS_FINISH,
				"e_project_id": projectId,
			}
			if tasks, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK, params).Find(); err == nil && len(tasks) > 0 {
				return false, errors.New("项目状态切换失败，存在未完成的任务或未发布的报告")
			}
			if this.Status == common.PS_TEST {
				this.Status = common.PS_COMPLETE
				this.UpdateTime = custom_util.GetCurrentMilliSecond()
			}
		}
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_PROJECT)
		if _, err := coll.UpdateOne(context.Background(), bson.M{"_id": this.Id}, bson.M{"$set": this}); err == nil {
			return true, err
		}
	}
	return false, err
}

// 查询某一个人的授权项目列表
func GetAuthorityProjects(userId int64) []string {
	return []string{}
}

// 判断某一个人是否是项目管理员
func CheckIsProjectManager(projectId string, userId int64) error {
	_id, _ := primitive.ObjectIDFromHex(projectId)
	params := qmap.QM{
		"e__id": _id,
	}
	if project, err := orm_mongo.NewWidgetWithParams(common.MC_PROJECT, params).Get(); err != nil {
		return errors.New(fmt.Sprintf("项目: %s 不存在！", projectId))
	} else {
		if manageId := project.Int("manager_id"); manageId == 0 || manageId != int(userId) {
			return errors.New("您没有操作该项目的权限")
		}
	}
	return nil
}

// 判断某一个人是否是项目成员
func CheckIsProjectUser(projectId string, userId int64) error {
	_id, _ := primitive.ObjectIDFromHex(projectId)
	params := qmap.QM{
		"e__id": _id,
	}
	if project, err := orm_mongo.NewWidgetWithParams(common.MC_PROJECT, params).Get(); err != nil {
		return errors.New(fmt.Sprintf("项目: %s 不存在！", projectId))
	} else {
		if allUsers := project.SliceInt("all_users"); !custom_util.InIntSlice(int(userId), allUsers) {
			return errors.New("您没有操作该项目的权限")
		}
	}
	return nil
}

// 判断项目是否存在
func CheckProjectExist(projectId string) error {
	_id, _ := primitive.ObjectIDFromHex(projectId)
	params := qmap.QM{
		"e__id": _id,
	}
	if _, err := orm_mongo.NewWidgetWithParams(common.MC_PROJECT, params).Get(); err != nil {
		return errors.New(fmt.Sprintf("项目: %s 不存在！", projectId))
	}
	return nil
}

func GetProjectSlice() qmap.QM {
	result := qmap.QM{}
	w := orm_mongo.NewWidgetWithCollectionName(common.MC_PROJECT)
	w.SetLimit(10000)
	if list, err := w.Find(); err == nil {
		for _, item := range list {
			result[item["_id"].(primitive.ObjectID).Hex()] = item["name"]
		}
	}
	return result
}

func (this *Project) One(id string) (*Project, error) {
	_id, _ := primitive.ObjectIDFromHex(id)
	params := qmap.QM{
		"e__id": _id,
	}
	err := orm_mongo.NewWidgetWithParams(common.MC_PROJECT, params).One(this)
	return this, err
}

// --------------------------------------------------
// 根据用户uid，得到下拉列表，当前用户所属项目列表以及每个项目关联的资产
func (this *Project) SelectListProjectAsset(uid int) []SelectListProjectAsset {
	result := make([]SelectListProjectAsset, 0)
	projectListModel := make([]ProjectObjectIdList, 0)

	{
		filter := bson.M{
			"all_users": uid, // all_user是[]int，直接用等于匹配相当于mysql的in查询
		}

		// mgoSession := mongo.NewMgoSession(common.MC_PROJECT)
		// mgoSession.Session.Find(filter).Sort("-create_time").Limit(100).All(&projectListModel)
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_PROJECT)
		opts := options.Find().SetLimit(100).SetSort(bson.D{{"create_time", -1}})
		cursor, err := coll.Find(context.Background(), filter, opts)
		if err != nil {
			log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
			return result
		}
		cursor.All(context.Background(), &projectListModel)
	}

	{
		for _, pl := range projectListModel {
			filter := bson.M{
				"project_id": pl.ProjectObjectId.Hex(), // 集合里面projectId用的hex
			}
			assetObjectIdList := make([]AssetObjectIdList, 0)
			mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_ASSET)
			mgoSession.Session.Find(filter).All(&assetObjectIdList)

			if len(assetObjectIdList) > 0 {
				assetList := make([]AssetList, 0)
				for _, v := range assetObjectIdList {
					assetList = append(assetList, AssetList{
						v.AssetObjectId,
						v.Name,
					})
				}

				s := SelectListProjectAsset{
					ProjectId: pl.ProjectObjectId.Hex(),
					Name:      pl.ProjectName,
					AssetList: assetList,
				}

				result = append(result, s)
			}
		}
	}

	return result
}

type SelectListProjectAsset struct {
	ProjectId string      `json:"project_id"`
	Name      string      `json:"name"`
	AssetList []AssetList `json:"asset_list"`
}

type AssetList struct {
	AssetId string `json:"asset_id"`
	Name    string `json:"name"`
}

type ProjectObjectIdList struct {
	ProjectObjectId primitive.ObjectID `bson:"_id,omitempty"`
	ProjectName     string             `bson:"name"`
}

type AssetObjectIdList struct {
	AssetObjectId string `bson:"_id,omitempty"` // asset表里面存的本身就是string
	Name          string `bson:"name"`
}

/**
 * @Description: 查询指定用户管理的项目列表
 * @param userId 用户id
 * @param status 项目状态
 * @return []string
 */
func (this *Project) GetManageProjects(userId int64, status int, startTime int64) []string {
	params := qmap.QM{
		"e_manager_id": userId,
		"e_status":     status,
	}
	if startTime > 0 {
		params["gt_update_time"] = startTime
	}
	projectIds := []string{}
	if list, err := orm_mongo.NewWidgetWithParams(common.MC_PROJECT, params).SetLimit(10000).Find(); err == nil {
		for _, item := range list {
			projectIds = append(projectIds, item["_id"].(primitive.ObjectID).Hex())
		}
	}
	return projectIds
}

/**
 * @Description:统计指定用户管理的项目包含的项目成员数量
 * @param userId
 * @param status
 * @return int
 */
func (this *Project) CountProjectsMembers(userId int64, status []int) int {
	match := bson.M{
		"manager_id": bson.M{"$eq": userId},
		"status":     bson.M{"$in": status},
	}
	group := bson.M{"_id": bson.M{"umember_id": "$member_ids"}}

	operations := []bson.M{
		{"$match": match},
		{"$unwind": "$member_ids"},
		{"$group": group},
		{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
	}
	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_PROJECT)
	cursor, err := coll.Aggregate(context.Background(), operations)
	if err != nil {
		panic(err)
	}
	resp := []bson.M{}
	if err := cursor.All(context.Background(), &resp); err == nil {
		if len(resp) > 0 {
			data := resp[0]
			return int(data["count"].(int32))
		}
	} else {
		panic(err)
	}
	return 0
}

/**
 * @Description:统计指定用户管理的项目包含的车型数量
 * @param userId
 * @param status
 * @return int
 */
func (this *Project) CountProjectsVehicleCode(userId int64, status []int) int {
	match := bson.M{
		"manager_id": bson.M{"$eq": userId},
		"status":     bson.M{"$in": status},
	}
	group := bson.M{"_id": bson.M{"ucode_name": "$code_name"}}

	operations := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
	}

	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_PROJECT)
	cursor, err := coll.Aggregate(context.Background(), operations)
	if err != nil {
		panic(err)
	}
	resp := []bson.M{}

	if err := cursor.All(context.Background(), &resp); err == nil {
		if len(resp) > 0 {
			data := resp[0]
			return int(data["count"].(int32))
		}
	} else {
		panic(err)
	}
	return 0
}

/**
 * @Description: 查询指定用户管理的项目概要信息
 * @param userId 用户id
 * @param status 项目状态
 * @return []string
 */
func (this *Project) GetManageProjectInfo(userId int64, status []int) []qmap.QM {
	params := qmap.QM{
		"e_manager_id": userId,
		"in_status":    status,
	}
	projectIds := []qmap.QM{}
	if list, err := orm_mongo.NewWidgetWithParams(common.MC_PROJECT, params).SetLimit(10000).Find(); err == nil {
		for _, item := range list {
			var itemQM qmap.QM = item
			project := qmap.QM{
				"id":   itemQM.Interface("_id").(primitive.ObjectID).Hex(),
				"name": itemQM.String("name"),
			}
			startTime := itemQM.Int64("start_time")
			endTime := itemQM.Int64("end_time")
			fullTime := endTime - startTime
			costTime := time.Now().Unix() - startTime
			if fullTime > 0 && costTime > 0 {
				if costTime < fullTime {
					project["process"] = costTime * 100 / fullTime
				} else {
					project["process"] = 100
				}
			} else {
				project["process"] = 0
			}
			projectIds = append(projectIds, project)
		}
	}
	return projectIds
}
