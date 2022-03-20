package stop

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"strconv"
	"ws/app/sys"
)

func NewStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop the server",
		Run: func(cmd *cobra.Command, args []string) {
			if !sys.IsRunning() {
				log.Fatalln("service is not running")
			}
			pid := sys.GetPid()
			if pid == 0 {
				log.Fatalln("service is not running")
			}
			e := exec.Command("kill", "15", strconv.Itoa(pid))
			_, err := e.Output()
			if err != nil {
				log.Fatalln(err)
			} else {
				fmt.Println("service had been stopped")
			}
		},
	}
	return cmd
}
