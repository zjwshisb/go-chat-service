package stop

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"strconv"
	"strings"
	"ws/app"
)

func NewStopCommand() *cobra.Command  {
	cmd := &cobra.Command{
		Use:                        "stop",
		Short: "stop the server",
		Run: func(cmd *cobra.Command, args []string) {
			pid := app.GetPid()
			if pid == 0 {
				fmt.Println("serve not running")
			} else {
				cmd := 	exec.Command("ps")
				out, err := cmd.Output()
				if err != nil {
					fmt.Println(err)
				}
				if strings.Contains(string(out), strconv.Itoa(pid)) {
					cmd := exec.Command("kill", "15", strconv.Itoa(pid))
					_, err := cmd.Output()
					if err != nil{
						fmt.Println(err)
					}
					fmt.Println("serve had stop")

				} else {
					fmt.Println("server not running")
				}
			}
		},
	}
	return cmd
}	