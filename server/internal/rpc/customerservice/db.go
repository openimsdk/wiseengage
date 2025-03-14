package customerservice

import "github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"

func setNotNil[V any](kv map[string]any, key string, value *V) bool {
	if value == nil {
		return false
	}
	kv[key] = *value
	return true
}

func UpdateAgent(req *customerservice.UpdateAgentReq) map[string]any {
	res := make(map[string]any)
	setNotNil(res, "nickname", req.Nickname)
	setNotNil(res, "face_url", req.FaceURL)
	setNotNil(res, "status", req.Status)
	setNotNil(res, "start_msg", req.StartMsg)
	setNotNil(res, "end_msg", req.EndMsg)
	setNotNil(res, "timeout_msg", req.TimeoutMsg)
	return res
}
