package stringx

func Remove(strings []string, strs ...string) (res []string) {
	res = append([]string(nil), strings...)
	for _, target := range strs {
		for i, v := range res {
			if target == v {
				res = append(res[0:i], res[i+1:len(res)]...)
			}
		}
	}
	return res
}
