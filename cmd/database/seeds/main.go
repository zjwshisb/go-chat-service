package main

import (
	"ws/internal/databases"
)

func init()  {
	databases.Setup()
}
func main()  {
	setting()
	rules()
}

