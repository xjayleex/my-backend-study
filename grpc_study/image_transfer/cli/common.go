package cli

import (
	"fmt"
	"os"
)

func Trap(err error) {
	if err == nil {
		return
	}
	fmt.Println("Error: %+v\n", err)
	os.Exit(1)
}
