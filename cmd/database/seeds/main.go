package main

import (
	"ws/app/databases"
)

func init()  {
	databases.Setup()
}
func main()  {
	rules()
}

