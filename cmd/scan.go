package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/amery/dcfs/pkg/scan"
)

var scanCmd = &cobra.Command{
	Use: "scan <dir>",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()

		scanner, err := scan.NewScanner(ctx)
		if err != nil {
			return err
		}
		defer scanner.Close()

		// watch signals
		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

			for signum := range sig {
				switch signum {
				case syscall.SIGINT, syscall.SIGTERM:
					log.Println("Terminating...")

					err := fmt.Errorf("Terminated by %s", signum)
					scanner.Cancel(err)
					return
				}
			}
		}()

		if len(args) == 0 {
			args = []string{"."}
		}

		for _, s := range args {

			fsys, dir, name, err := scanner.Split(s)
			if err != nil {
				scanner.Cancel(err)
				break
			}

			scanner.Spawn(func(ctx context.Context) {
				scanner.Scan(ctx, fsys, dir, name)
			})
		}

		return scanner.Wait()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
