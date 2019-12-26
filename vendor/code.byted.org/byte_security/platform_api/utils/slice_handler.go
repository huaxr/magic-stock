package utils

func StringsDistinct(src []string) (dst []string) {
	dst = make([]string, 0)
	for i := 0; i < len(src); i++ {
		repeat := false
		for j := i + 1; j < len(src); j++ {
			if src[i] == src[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			dst = append(dst, src[i])
		}
	}
	return
}
