package conns

import (
	"fmt"
	"github.com/spf13/cobra"
	 "ws/app/rpc/client"
)

type Reply struct {
	Name string
}

func NewConnsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connection",
		Short: "show the connections",
		Run: func(cmd *cobra.Command, args []string) {
			c := client.ClientTotal(1, "admin")
			fmt.Println(c)
		},
	}
	return cmd
}
