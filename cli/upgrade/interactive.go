// Copyright (c) 2022 Wireleap

package upgrade

import (
	"fmt"
	"os"
	"strings"
)

func Confirm(prompt string) bool {
	var (
		i      = 3
		answer = "n"
	)

	for j := 0; j < i; j++ {
		fmt.Printf("%s (y/N) ", prompt)
		fmt.Scanln(&answer)

		switch strings.TrimSpace(strings.ToLower(answer)) {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		}

		fmt.Printf("Response '%s' not understood. ", answer)
	}

	fmt.Printf("Failed to get an answer %d times. Aborting.\n", i)
	os.Exit(1)

	// never reached
	return false
}
