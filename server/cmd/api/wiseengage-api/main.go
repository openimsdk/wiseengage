package main

import (
	customerapi "github.com/openimsdk/wiseengage/v1/internal/api/customerservice"
	"github.com/openimsdk/wiseengage/v1/pkg/common/cmd"
)

func main() {
	cmd.Run(customerapi.Start)
}
