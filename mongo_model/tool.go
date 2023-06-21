package mongo_model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"
	"strconv"
	"strings"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
)

type ToolData struct {
	Id              bson.ObjectId `bson:"_id"`
	ToolNumber      int           `bson:"tool_number"`
	ToolName        string        `bson:"tool_name"`
	Tool            string        `bson:"tool"`
	TestPic         []postJson    `bson:"test_pic"`
	CategoryName    string        `bson:"category_name"`
	CategoryID      int           `bson:"category_id"`
	UseDetail       string        `bson:"use_detail"`
	ToolDetail      string        `bson:"tool_detail"`
	toolUrl         string        `bson:"tool_url"`
	toolLogo        []postJson    `bson:"tool_logo"`
	UseManual       []postJson    `bson:"use_manual"`
	UseManualLink   []string      `bson:"use_manual_link"`
	LinkPic         []postJson    `bson:"link_pic"`
	Script          []postJson    `bson:"script"`
	Brand           string        `bson:"brand"`
	ParamsJson      string        `bson:"params_json"`
	ParamsRemarks   string        `bson:"params_remarks"`
	SoftwareVersion string        `bson:"software_version"`
	HardwareVersion string        `bson:"hardware_version"`
	SpVersion       string        `bson:"sp_version"`
	SystemVersion   string        `bson:"system_version"`
	CreateUserName  string        `bson:"create_user_name"`
	Search          []string      `bson:"search"`
	CreateUserID    int           `bson:"create_user_id"`
	CreateTime      int           `bson:"create_time"`
	UpdateTime      int           `bson:"update_time"`
	Status          int           `bson:"status"`
	HistoryLog      []HistoryLog  `bson:"history_log"`
	Tag             string        `bson:"tag"`
	HistoryVersion  string        `bson:"history_version"` //变更历史迭代版本号
	IsCreate        int           `bson:"is_create"`       // 创建测试任务 1 勾选 2 不勾选
	IsJump          int           `bson:"is_jump"`         // 快速跳转 1 勾选 2 不勾选
}

type ToolDataResponse struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name"`        // 名称
	ToolNumber int           `bson:"tool_number"` // 名称
}

type postJson struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HistoryLog struct {
	Version       string `json:"version"`
	Date          int64  `json:"date"`
	UpdateUser    int    `json:"update_user"`
	UpdateContent string `json:"update_content"`
}

func (this *ToolData) GetToolCate() (*qmap.QM, error) {

	groupOperations := []bson.M{
		{"$match": bson.M{"status": 1}},
		{"$project": bson.M{"category_name": 1, "category_id": 1, "_id": 1, "tool_name": 1}},
		// {"$group": bson.M{"_id": bson.M{"cate_id":"$category_id","tool_name":"$tool_name"}, "count": bson.M{"$sum": 1}}},
	}
	list, err := mongo.NewMgoSession(common.MC_TOOL).QueryGet(groupOperations)
	result := []qmap.QM{}
	typeQM := qmap.QM{}
	if err == nil {
		for _, item := range *list {
			if typeQM[item["category_name"].(string)] == nil {
				typeQM[item["category_name"].(string)] = []qmap.QM{}
			}
			assetStruct := qmap.QM{
				"id":        item["_id"],
				"tool_name": item["tool_name"].(string),
			}
			typeQM[item["category_name"].(string)] = append(typeQM[item["category_name"].(string)].([]qmap.QM), assetStruct)
		}
		for typeName, item := range typeQM {
			itemSlice := qmap.QM{
				"category": typeName,
				"list":     item,
			}
			result = append(result, itemSlice)
		}

	}
	return &qmap.QM{"data": result}, err
}

func (this *ToolData) Create(rawInfo *qmap.QM, UserID int, UserName string) (*ToolDataResponse, error) {

	toolName := rawInfo.MustString("name") // 工具名称
	// 添加工具之前首先对工具名称唯一性做判断
	checkSetParams := qmap.QM{
		"e_tool_name": toolName,
	}

	mongoClient := mongo.NewMgoSessionWithCond(common.MC_TOOL, checkSetParams)
	if _, err := mongoClient.GetOne(); err == nil {
		return nil, errors.New("该工具名称已经存在，请重新填写工具名称")
	} else {
		testPic := rawInfo.String("test_pic")
		testPicArr := []postJson{}
		json.Unmarshal([]byte(testPic), &testPicArr)

		// testPicArr := strings.Split(testPic, ",")
		categoryName := rawInfo.MustString("category_name") // 分类名称
		categoryID := rawInfo.MustInt("category_id")        // 分类ID
		useDetail := rawInfo.String("use_detail")           // 工具使用方法详情
		toolDetail := rawInfo.MustString("tool_detail")     // 工具介绍
		tag := rawInfo.String("tag")                        // 工具介绍
		useManualLink := rawInfo.String("use_manual_link")  // 使用手册Link
		useManualLinkArr := strings.Split(useManualLink, ",")
		useManual := rawInfo.String("use_manual") // 使用手册
		useManualArr := []postJson{}
		json.Unmarshal([]byte(useManual), &useManualArr)
		// useManualArr := strings.Split(useManual, ",")

		linkPic := rawInfo.String("link_pic") // 工具连接示意图
		linkPicArr := []postJson{}
		json.Unmarshal([]byte(linkPic), &linkPicArr)

		script := rawInfo.String("script") // 脚本
		scriptArr := []postJson{}
		json.Unmarshal([]byte(script), &scriptArr)

		isCreate := rawInfo.Int("is_create")
		isJump := rawInfo.Int("is_jump")

		toolUrl := rawInfo.String("tool_url") // 工具url

		toolLogo := rawInfo.String("logo") // 工具logo
		toolLogoArr := []postJson{}
		json.Unmarshal([]byte(toolLogo), &toolLogoArr)

		paramsJson := rawInfo.MustString("params_json") // 工具参数配置
		paramsJsonBase64 := base64.StdEncoding.EncodeToString([]byte(paramsJson))
		brand := rawInfo.MustString("brand")           // 工具品牌
		paramsRemarks := rawInfo.MustString("remarks") // 工具参数配置备注
		softV := rawInfo.String("software_version")    // 软件版本
		hardV := rawInfo.String("hardware_version")    // 硬件版本
		spV := rawInfo.MustString("sp_version")        // 工具总版本
		sysV := rawInfo.String("system_version")       // 工具总版本
		nowTime := custom_util.GetCurrentMilliSecond() / 1000
		toolNumber, _ := strconv.Atoi(custom_util.TimestampToStringNoSpace(nowTime))

		this.Id = bson.NewObjectId()
		this.ToolNumber = toolNumber
		this.ToolName = toolName
		this.Tool = rawInfo.String("tool")
		this.TestPic = testPicArr
		this.Brand = brand
		this.CategoryName = categoryName
		this.CategoryID = categoryID
		this.UseDetail = useDetail
		this.ToolDetail = toolDetail
		this.UseManual = useManualArr
		this.UseManualLink = useManualLinkArr
		this.LinkPic = linkPicArr
		this.Script = scriptArr
		this.ParamsJson = paramsJsonBase64
		this.ParamsRemarks = paramsRemarks
		this.SoftwareVersion = softV
		this.HardwareVersion = hardV
		this.SpVersion = spV
		this.SystemVersion = sysV
		this.CreateUserName = UserName
		this.CreateUserID = UserID
		this.Search = []string{toolName, categoryName, strconv.Itoa(categoryID), strconv.Itoa(toolNumber), toolDetail}
		this.CreateTime = int(nowTime)
		this.UpdateTime = int(nowTime)
		this.Status = 1
		this.Tag = tag
		this.toolUrl = toolUrl
		this.toolLogo = toolLogoArr
		this.IsCreate = isCreate
		this.IsJump = isJump

		//创建工具的时候变更历史记录默认版本为v1.0，之后每次变更+0.1
		historyLog := []HistoryLog{}
		nowHisLog := HistoryLog{
			Date:          nowTime,
			Version:       "v1.0",
			UpdateUser:    UserID,
			UpdateContent: "创建工具",
		}
		historyLog = append(historyLog, nowHisLog)
		this.HistoryLog = historyLog

		if err := mongo.NewMgoSession(common.MC_TOOL).Insert(this); err == nil {
			Rts := ToolDataResponse{
				Id:         this.Id,
				Name:       this.ToolName,
				ToolNumber: this.ToolNumber,
			}
			return &Rts, nil
		} else {
			return nil, err
		}
	}
}

func (this *ToolData) Edit(rawInfo *qmap.QM, UserID int, UserName string) (*ToolDataResponse, error) {

	masterID := bson.ObjectIdHex(rawInfo.MustString("id"))
	params := qmap.QM{
		"e__id": masterID,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_TOOL, params)
	if rts, err := mongoClient.GetOne(); err == nil {
		toolName := rawInfo.MustString("name") // 工具名称
		// 编辑工具时 首先检测除了本身外，该工具是否存在
		checkSetParams := qmap.QM{
			"ne__id":      masterID,
			"e_tool_name": toolName,
		}

		if _, checkErr := mongo.NewMgoSessionWithCond(common.MC_TOOL, checkSetParams).GetOne(); checkErr == nil {
			return nil, errors.New("该工具名称已经存在，请重新填写工具名称")
		} else {
			testPic := rawInfo.String("test_pic") // 测试工具图片
			testPicArr := []postJson{}
			json.Unmarshal([]byte(testPic), &testPicArr)
			// testPicArr := strings.Split(testPic, ",")
			categoryName := rawInfo.MustString("category_name") // 分类名称
			categoryID := rawInfo.MustInt("category_id")        // 分类ID
			useDetail := rawInfo.String("use_detail")           // 工具使用方法详情
			toolDetail := rawInfo.MustString("tool_detail")     // 工具介绍

			useManualLink := rawInfo.String("use_manual_link") // 使用手册Link
			useManualLinkArr := strings.Split(useManualLink, ",")

			useManual := rawInfo.String("use_manual") // 使用手册
			useManualArr := []postJson{}
			json.Unmarshal([]byte(useManual), &useManualArr)
			// useManualArr := strings.Split(useManual, ",")

			linkPic := rawInfo.String("link_pic") // 工具连接示意图
			linkPicArr := []postJson{}
			json.Unmarshal([]byte(linkPic), &linkPicArr)

			script := rawInfo.String("script") // 脚本
			scriptArr := []postJson{}
			json.Unmarshal([]byte(script), &scriptArr)

			isCreate := rawInfo.Int("is_create")
			isJump := rawInfo.Int("is_jump")

			toolUrl := rawInfo.String("tool_url") // 工具url

			toolLogo := rawInfo.String("logo") // 工具logo
			toolLogoArr := []postJson{}
			json.Unmarshal([]byte(toolLogo), &toolLogoArr)

			paramsJson := rawInfo.MustString("params_json") // 工具参数配置
			paramsJsonBase64 := base64.StdEncoding.EncodeToString([]byte(paramsJson))
			brand := rawInfo.MustString("brand")                // 工具品牌
			paramsRemarks := rawInfo.MustString("remarks")      // 工具参数配置备注
			softV := rawInfo.String("software_version")         // 软件版本
			hardV := rawInfo.String("hardware_version")         // 硬件版本
			spV := rawInfo.MustString("sp_version")             // 工具总版本
			sysV := rawInfo.String("system_version")            // 工具总版本
			historyCon := rawInfo.MustString("history_content") // 变更历史记录的变更内容
			historyVersion := rts.String("history_version")     // 变更历史记录的版本号
			newHistoryVersion := addHistoryVersion(historyVersion)

			nowTime := custom_util.GetCurrentMilliSecond() / 1000

			historyLog := []HistoryLog{}

			newHistoryLogVersion := ""

			tag := rawInfo.String("tag")
			if dbHisLog, has := rts.TrySlice("history_log"); has {
				for _, val := range dbHisLog {
					oneLog := val.(map[string]interface{})
					tmpDate, _ := oneLog["date"].(int64)
					tmpCon, _ := oneLog["updatecontent"].(string)
					tmpU, _ := oneLog["updateuser"].(int)
					tmpV, _ := oneLog["version"].(string)
					updateHisLog := HistoryLog{
						Date:          tmpDate,
						Version:       tmpV,
						UpdateUser:    tmpU,
						UpdateContent: tmpCon,
					}
					//获取最后一次的历史变更记录版本号
					historyLog = append(historyLog, updateHisLog)
					newHistoryLogVersion = tmpV
				}
			}
			nowHisLog := HistoryLog{
				Date:          nowTime,
				Version:       addHistoryVersion(newHistoryLogVersion), // 变更内容的版本号 如v1.2
				UpdateUser:    UserID,
				UpdateContent: historyCon,
			}
			historyLog = append(historyLog, nowHisLog)
			update := bson.M{
				"$set": bson.M{
					"tool_name":        toolName,
					"test_pic":         testPicArr,
					"brand":            brand,
					"category_name":    categoryName,
					"category_id":      categoryID,
					"use_detail":       useDetail,
					"tool_detail":      toolDetail,
					"use_manual":       useManualArr,
					"use_manual_link":  useManualLinkArr,
					"link_pic":         linkPicArr,
					"script":           scriptArr,
					"params_json":      paramsJsonBase64,
					"params_remarks":   paramsRemarks,
					"software_version": softV,
					"hardware_version": hardV,
					"sp_version":       spV,
					"system_version":   sysV,
					"search":           []string{toolName, categoryName, strconv.Itoa(categoryID), strconv.Itoa(10001)},
					"update_time":      int(nowTime),
					"history_log":      historyLog,
					"tag":              tag,
					"history_version":  newHistoryVersion,
					"is_create":        isCreate,
					"is_jump":          isJump,
					"tool_url":         toolUrl,
					"tool_logo":        toolLogoArr,
				},
			}
			if err := mongoClient.Update(bson.M{"_id": masterID}, update); err != nil {
				return nil, err
			} else {
				Rts := ToolDataResponse{
					Id:         masterID,
					Name:       toolName,
					ToolNumber: rts.Int("tool_number"),
				}
				return &Rts, nil
			}
		}
	} else {
		return nil, errors.New("tool not found")
	}
}

func (this *ToolData) Del(rawInfo *qmap.QM) (*ToolDataResponse, error) {

	masterID := bson.ObjectIdHex(rawInfo.MustString("id"))
	params := qmap.QM{
		"e__id": masterID,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_TOOL, params)
	if err := mongoClient.One(&this); err == nil {
		nowTime := custom_util.GetCurrentMilliSecond() / 1000
		this.UpdateTime = int(nowTime)
		this.Status = 0
		mongoClient := mongo.NewMgoSession(common.MC_TOOL)
		if err := mongoClient.Update(bson.M{"_id": masterID}, this); err != nil {
			return nil, err
		} else {
			Rts := ToolDataResponse{
				Id:         masterID,
				Name:       this.ToolName,
				ToolNumber: this.ToolNumber,
			}
			return &Rts, nil
		}
	} else {
		return nil, errors.New("tool not found")
	}
}

func (this *ToolData) GetOne(rawInfo *qmap.QM) (*qmap.QM, error) {

	masterID := bson.ObjectIdHex(rawInfo.MustString("id"))
	params := qmap.QM{
		"e__id":    masterID,
		"e_status": 1,
	}
	if len(params) == 0 {
		return nil, errors.New("Provide at least one id and name!")
	}
	if rts, err := mongo.NewMgoSessionWithCond(common.MC_TOOL, params).GetOne(); err == nil {
		return rts, nil
	} else {
		return nil, err
	}
}

// 增加“变更历史”的版本号
func addHistoryVersion(HistoryVersion string) string {
	if HistoryVersion == "" {
		return "v1.0"
	}
	strs := strings.Split(HistoryVersion, "v")
	var v = strs[1]
	val, _ := strconv.ParseFloat(v, 64)
	val += 0.1
	num := fmt.Sprintf("%.1f", val)
	finall := strings.Join([]string{"v", num}, "")
	return finall

}
func (this *ToolData) UpdateTag(rawInfo *qmap.QM, UserID int, UserName string) (*ToolDataResponse, error) {

	masterID := bson.ObjectIdHex(rawInfo.MustString("id"))
	params := qmap.QM{
		"e__id": masterID,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_TOOL, params)
	if rts, err := mongoClient.GetOne(); err == nil {
		toolName := rawInfo.MustString("name") // 工具名称
		tag := rawInfo.String("tag")

		update := bson.M{
			"$set": bson.M{
				"tag":       tag,
				"tool_name": toolName,
			},
		}
		if err := mongoClient.Update(bson.M{"_id": masterID}, update); err != nil {
			return nil, err
		} else {
			Rts := ToolDataResponse{
				Id:         masterID,
				Name:       toolName,
				ToolNumber: rts.Int("tool_number"),
			}
			return &Rts, nil
		}
	} else {
		return nil, errors.New("tool not found")
	}
}
