package commands

import (
	"fmt"
	"runtime/debug"
)

func HandlerVersion(s *State, cmd Command) error {

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("unknown")
	}

	fmt.Println(info.Main.Version)
	return nil
}
