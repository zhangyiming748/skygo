package pg

func FindById(id int, modelPtr interface{}) (bool, error) {
	has, err := GetSession().ID(id).Get(modelPtr)
	return has, err
}

func DeleteById(id int, modelPtr interface{}) (bool, error) {
	count, err := GetSession().ID(id).Delete(modelPtr)
	if count == 0 {
		return false, err
	} else {
		return true, err
	}
}
