package conns

import (
	"context"
	"github.com/pterm/pterm"
	"github.com/smallnest/rpcx/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"ws/app/rpc/server/connection"
)

type Reply struct {
	Name string
}

func NewConnsCommand() *cobra.Command  {
	cmd := &cobra.Command{
		Use:                        "connection",
		Short: "show the connections",
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetBool("App.Cluster") {

			} else {
				d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:" + viper.GetString("Rpc.Port") , "")
				// #2
				xclient := client.NewXClient("Connection", client.Failtry, client.RandomSelect, d, client.DefaultOption)

				reply := &connection.Reply{}
				a := &connection.Args{}
				// #3
				call, err := xclient.Go(context.Background(), "Show", a, reply, nil)
				if err != nil {
					log.Fatalf("failed to call: %v", err)
				}

				replyCall := <-call.Done
				if replyCall.Error != nil {
					log.Fatalf("failed to call: %v", replyCall.Error)
				} else {
					table := make([][]string, 0)
					table = append(table, []string{
						"uid", "created_time",
					})
					for _, d := range reply.Data {
						table = append(table, []string{
							d["uid"], d["created_at"],
						})
					}
					pterm.DefaultTable.WithHasHeader().WithBoxed(true).WithData(table).Render()
				}
			}

		},
	}
	return cmd
}