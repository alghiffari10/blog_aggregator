package main

import (
	"fmt"
	"runtime/debug"
)

func handlerVersion(s *state, cmd command) error {

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("unknown")
	}

	fmt.Println(info.Main.Version)
	return nil
}
