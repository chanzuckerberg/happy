package util

var localstackMode = false
var localstackEndpoint = "http://localhost:4566"

func IsLocalstackMode() bool {
	return localstackMode
}

func SetLocalstackMode(localstack bool) {
	localstackMode = localstack
}

func GetLocalstackEndpoint() string {
	return localstackEndpoint
}

func SetLocalstackEndpoint(endpoint string) {
	localstackEndpoint = endpoint
}
