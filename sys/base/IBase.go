package base

type ICdRequest struct {
	ctx  string
	m    string
	c    string
	a    string
	dat  string
	args string
}

type JWT struct {
	jwtToken   string
	checked    bool
	checkTime  int
	authorized bool
}

type ISessResp struct {
	cd_token string
	userId   int
	jwt      JWT
	ttl      int
	initUuid string
	initTime string
}

type IRespInfo struct {
	messages []string
	code     string
	app_msg  string
}

type IServerConfig struct {
	usePush       bool
	usePolling    bool
	useCacheStore bool
}

type IAppState struct {
	success bool
	info    IRespInfo
	sess    ISessResp
	cache   string
	sConfig IServerConfig
}

type ICdResponse struct {
	app_state IAppState
	data      string
}
