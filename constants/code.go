package constants

const CODE_OK = 0
const CODE_UNHANDLED_ERROR = -1

// 登录相关
const (
	CODE_NOT_LOGIN = 100_0000 + iota
	CODE_TOKEN_NOT_VALID
	CODE_USERNAME_IS_REGISTERED
	CODE_USERNAME_OR_PASSWORD_ERROR
	CODE_REGISTER_PARAM_NOT_VALID
	CODE_LOGIN_PARAM_NOT_VALID
)

// 目录
const (
	CODE_CREATE_DIR_PARAM_NOT_VALID = 200_0000 + iota
	CODE_QUERY_DIR_INFO_WITH_EMPTY_RES
)

// 文件
const (
	CODE_FILENAME_HAS_BEEN_USED = 300_0000 + iota
	CODE_FILE_NOT_EXIST
)

// 通用
const (
	CODE_PARAMS_NOT_VALID = 900_0000 + iota
)
