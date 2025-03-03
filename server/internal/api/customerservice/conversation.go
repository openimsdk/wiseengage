package customerservice

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/wiseengage/v1/pkg/protocol/customerservice"
)

func (x *API) StartConsultation(c *gin.Context) {
	a2r.Call(c, customerservice.CustomerserviceClient.StartConsultation, x.client)
}

func (x *API) RegisterCustomer(c *gin.Context) {
	a2r.Call(c, customerservice.CustomerserviceClient.ChangeConversationRole, x.client)
}
