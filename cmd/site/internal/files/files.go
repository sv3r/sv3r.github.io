package files

import (
	"log"
	"os"
)

func Copy(src, dst string) {
	input, err := os.ReadFile(src)
	must(err)

	must(os.WriteFile(dst, input, 0o644))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
