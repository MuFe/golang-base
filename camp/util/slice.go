package utils

// InArray 数组是否包含某元素
func InArray(strif interface{}, arrayif interface{}) bool {
	switch strif.(type) {
	case string:
		str := strif.(string)
		array := arrayif.([]string)
		for _, j := range array {
			if j == str {
				return true
			}
		}
	case int:
		str := strif.(int)
		array := arrayif.([]int)
		for _, j := range array {
			if j == str {
				return true
			}
		}
	case int32:
		str := strif.(int32)
		array := arrayif.([]int32)
		for _, j := range array {
			if j == str {
				return true
			}
		}
	case int64:
		str := strif.(int64)
		array := arrayif.([]int64)
		for _, j := range array {
			if j == str {
				return true
			}
		}
	case uint:
		str := strif.(uint)
		array := arrayif.([]uint)
		for _, j := range array {
			if j == str {
				return true
			}
		}
	default:
		panic("wrong type")
	}
	return false
}

// RemoveDupsString 字符串数组去重
func RemoveDupsString(slice []string) []string {
	keys := make(map[string]bool)
	list := make([]string, 0)
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// RemoveDupsInt 整形数组去重
func RemoveDupsInt(intSlice []int) []int {
	keys := make(map[int]bool)
	list := make([]int, 0)
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// RemoveDupsInt32 整形32数组去重
func RemoveDupsInt32(intSlice []int32) []int32 {
	keys := make(map[int32]bool)
	list := make([]int32, 0)
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// RemoveDupsInt64 整形64数组去重
func RemoveDupsInt64(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := make([]int64, 0)
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// RemoveDupsInt64 整形64数组去重除0
func RemoveDupsInt64WithZero(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := make([]int64, 0)
	for _, entry := range intSlice {
		if entry == 0 {
			continue
		}
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// RemoveDupsUint 无符号整形数组去重
func RemoveDupsUint(intSlice []uint) []uint {
	keys := make(map[uint]bool)
	list := make([]uint, 0)
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// IntersectInt32 32位整形交集
func IntersectInt32(i, j []int32) (result []int32) {
	if len(i) == 0 || len(j) == 0 {
		return result
	}
	var m = map[int32]struct{}{}
	for _, v := range i {
		m[v] = struct{}{}
	}
	for _, v := range j {
		if _, ok := m[v]; ok {
			result = append(result, v)
		}
	}
	return result
}

// IntersectInt64 64位整形交集
func IntersectInt64(i, j []int64) (result []int64) {
	if len(i) == 0 || len(j) == 0 {
		return result
	}
	var m = map[int64]struct{}{}
	for _, v := range i {
		m[v] = struct{}{}
	}
	for _, v := range j {
		if _, ok := m[v]; ok {
			result = append(result, v)
		}
	}
	return result
}
