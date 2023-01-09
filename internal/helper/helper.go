package helper

func CalculateMode(isCdn, isSni string) (int, int) {
	if isCdn != "" && isSni != "" {
		return 1, 0
	} else if isCdn != "" {
		return 1, 1
	} else if isSni != "" {
		return 0, 0
	} else {
		return 0, 1
	}
}
