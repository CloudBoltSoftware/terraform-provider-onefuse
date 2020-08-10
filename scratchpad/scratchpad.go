package main

import (
	"fmt"
	"github.com/cloudboltsoftware/terraform-provider-onefuse/onefuse"
	"strings"
)

func main() {
	config := onefuse.NewConfig("https", "localhost", "443", "admin", "admin", false)
	client := config.NewOneFuseApiClient()

	_ := client.DoSomething()
}
