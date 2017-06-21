package utils

import (
	"fmt"
)

// sends terminal codes to update title; generally to be used when running purrbot from tmux/screen
func UpdateXtermTitle(s string){
	fmt.Printf("\033]0;%s\007",s)
}
