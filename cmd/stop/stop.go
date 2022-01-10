package stop

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"strconv"
	"ws/app"
)

func NewStopCommand() *cobra.Command  {
	cmd := &cobra.Command{
		Use:                        "stop",
		Short: "stop the server",
		Run: func(cmd *cobra.Command, args []string) {
			if !app.IsRunning() {
				log.Fatalln("serve is not running")
			}
 			pid := app.GetPid()
			e := exec.Command("kill", "15", strconv.Itoa(pid))
			_, err := e.Output()
			if err != nil{
				log.Fatalln(err)
			}  else {
				fmt.Println("serve is stop")
			}
		},
	}
	return cmd
}	