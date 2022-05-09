package utils

func GenerateUid() string {
	return RandStringBytesMaskImprSrc(5)
}

func GeneratePassword() string {
	return RandStringBytesMaskImprSrc(16)
}
