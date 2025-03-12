package customerapi

import (
	"net"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestName(t *testing.T) {
	l, err := net.Listen("tcp", ":10006")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	t.Log("addr", l.Addr())
	var api API
	r := gin.New()
	// http://127.0.0.1:10006/callback/openim/callbackAfterSendGroupMsgCommand?key=eyJzZW5kSUQiOiI1MzE4NTQzODIyIiwicmVjdklEIjoiIiwiZ3JvdXBJRCI6IjUzODQwNzI4NSJ9
	r.Group("/callback").POST("/openim/*command", api.OpenIMCallback)
	t.Log(r.RunListener(l))
}
