package main

import (
	"github.com/openimsdk/wiseengage/v1/internal/rpc/customerservice"
	"github.com/openimsdk/wiseengage/v1/pkg/common/cmd"
)

func main() {
	cmd.Run(customerservice.Start)
}
