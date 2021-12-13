package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/amery/dcfs/pkg/dcfs"
)

var listCmd = &cobra.Command{
	Use: "list [flags] <datadir>",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Not enough arguments")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		// instantiate
		fsys, err := dcfs.New(ctx, args[0])
		if err != nil {
			return err
		}

		// watch signals
		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

			for signum := range sig {
				switch signum {
				case syscall.SIGINT, syscall.SIGTERM:
					log.Println("Terminating...")
					cancel()
					return
				}
			}
		}()

		fsys.ForEachRecord(func(ctx context.Context, node *dcfs.NodeRecord) error {
			select {
			case <-ctx.Done():
				return errors.New("cancelled")
			default:
				log.Println(node)
				return nil
			}
		})
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
