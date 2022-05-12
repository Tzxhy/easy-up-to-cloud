package utils

func GenerateUid() string {
	return RandStringBytesMaskImprSrc(5)
}
func GenerateDid() string {
	return RandStringBytesMaskImprSrc(8)
}
func GenerateFid() string {
	return RandStringBytesMaskImprSrc(10)
}

func GenerateGid() string {
	return RandStringBytesMaskImprSrc(5)
}

func GenerateRid() string {
	return RandStringBytesMaskImprSrc(8)
}

func GeneratePassword() string {
	return RandStringBytesMaskImprSrc(16)
}
