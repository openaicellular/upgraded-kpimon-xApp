package main

import (
	//"gerrit.o-ran-sc.org/r/scp/ric-app/kpimon/control"
	"example.com/kpimon/control"
)

func main() {
	c := control.NewControl()
	c.Run()
}

