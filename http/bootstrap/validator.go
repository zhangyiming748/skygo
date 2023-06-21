package bootstrap

// func InitValidator() {
// 	v, ok := binding.Validator.Engine().(*validator.Validate)
//     if ok {
//         // 自定义验证方法
//         if err := v.RegisterValidation("checkMobile", checkMobile); err != nil {
//         	panic(err)
// 		}
//         if err := v.RegisterValidation("FieldTypeValid", FieldTypeValid); err != nil {
//         	panic(err)
// 		}
// 		if err := v.RegisterValidation("TaskDataConditionValid", form.TaskDataConditionValid); err != nil {
// 			panic(err)
// 		}
// 		// 规则表单的样例字段，必须是一个json数组对应的字符串
// 		if err := v.RegisterValidation("RuleSampleContentValid", form.RuleSampleContentValid); err != nil {
// 			panic(err)
// 		}
//     } else {
//     	panic("validator init failed")
// 	}
// }
//
// func checkMobile(fl validator.FieldLevel) bool {
// 	 mobile := strconv.Itoa(int(fl.Field().Uint()))
//     re := `^1[3456789]\d{9}$`
//     r := regexp.MustCompile(re)
//     return r.MatchString(mobile)
// }
//
// func FieldTypeValid(fl validator.FieldLevel) bool {
// 	if v, ok := fl.Field().Interface().(int); ok {
// 		_, ok := common.FieldTypeMap[v]
// 		return ok
// 	}
// 	return false
// }
