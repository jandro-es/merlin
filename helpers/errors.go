package helpers

import "log"

func ExitOnFail(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
