package mongo_model

type TestSubClass struct {
	Name string `bson:"sub_class_name" json:"sub_class_name"` //测试子类名称
}

type TestClass struct {
	Name     string         `bson:"class_name" json:"class_name"` //测试分类名称
	SubClass []TestSubClass `bson:"sub_class" json:"sub_class"`   //包含的测试子类
}

type EvaluateTestCaseTemplate struct {
	Id       string      `bson:"_id,omitempty"  json:"_id"`  // 模板ID
	Name     string      `bson:"name" json:"name"`           // 模板名字
	Describe string      `bson:"objective" json:"objective"` // 模板描述
	Class    []TestClass `bson:"class" json:"class"`         // 测试分类
}
