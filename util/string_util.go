package util

func RemoveStringFromList(list []string, token string) []string {
	var newList []string
	for _, v := range list {
		if v != token {
			newList = append(newList, v)
		}
	}
	return newList
}
