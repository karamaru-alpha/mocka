package testdata

import "fmt"

type Human interface {
	Say(name string) string
}

func Hello(myName fmt.Stringer, me Human) string {
	return me.Say(myName.String())
}
