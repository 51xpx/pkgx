package rescue

import (
	"fmt"
)

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if err := recover(); err != nil {
		fmt.Println("[rescue]recover from panic", err)
	}
}
