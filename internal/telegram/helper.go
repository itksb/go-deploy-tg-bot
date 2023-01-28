package telegram

func sliceContains(s []int64, id int64) bool {
	for _, v := range s {
		if v == id {
			return true
		}
	}

	return false
}
