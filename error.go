package main

import ()

type Error interface {
	Error() string
}

func check(e error) {
	if e != nil {
		TraceActivity.Printf("Error : --%s--\n", e.Error())
		panic(e)
	}
}
