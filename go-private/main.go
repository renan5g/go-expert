package main

import (
	"fmt"

	"github.com/renan5g/go-expert-utils/pkg/events"
)

func main() {
	ed := events.NewEventDispatcher()
	fmt.Println(ed)
}
