package ecode

// rest code: [1000, 1999]
const (
	// EcodeRESTPanic error code of rest panic
	EcodeRESTPanic = 1000
	// EcodeRESTDBNotFound db not found
	EcodeRESTDBNotFound = 1010
	// EcodeRESTDBMigrateErr db migrate failed
	EcodeRESTDBMigrateErr = 1011

	// EcodeRESTAlertWechatErr failed in sending message to wechat
	EcodeRESTAlertWechatErr = 1100
)
