package log

import "go-pkg/config"

func A() {
	config.Echo()
	config.SetOutput("world")
}
