package main

import (
	"dev.ideatip.appendr/logv1"
	"dev.ideatip.appendr/logv2"
	"dev.ideatip.appendr/logv3"
)

func main() {
	logv1.ExampleUsage()
	println("----------------------------------------------------")
	logv2.ExampleUsage()
	println("----------------------------------------------------")
	logv3.ExampleUsage()
}
