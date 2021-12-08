package cmd

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"

	"github.com/amery/dcfs/pkg/fuse"
	"github.com/spf13/cobra"
)

var (
	traceFile string
)

var mountCmd = &cobra.Command{
	Use: "mount [flags] <datadir> <mountpoint>",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Not enough arguments")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// trace
		if traceFile != "" {
			if f, err := os.Create(traceFile); err != nil {
				log.Println("Failed to create trace output file", err)
			} else {
				log.Printf("Writing Tracing to %q", traceFile)

				defer f.Close()
				trace.Start(f)
				defer trace.Stop()
			}
		}

		// fuse daemon
		m, err := fuse.New(args[0], args[1])
		if err != nil {
			return err
		}

		// watch signals
		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

			for signum := range sig {
				switch signum {
				case syscall.SIGHUP:
					// reload
					if err := m.Reload(); err != nil {
						log.Println("Failed to reload:", err)
					}
				case syscall.SIGINT, syscall.SIGTERM:
					// terminate
					log.Println("Terminating...")
					if err := m.Abort(); err != nil {
						log.Println("Failed to terminate:", err)
						continue
					}

					return
				}
			}
		}()

		defer m.Close()
		return m.Serve()
	},
}

func init() {
	flags := mountCmd.Flags()
	flags.StringVarP(&traceFile, "trace-file", "T", "", "Trace output")

	rootCmd.AddCommand(mountCmd)
}
