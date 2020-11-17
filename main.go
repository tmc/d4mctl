// Command d4mctl offers command line manipulation of a docker-for-mac installation.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tmc/d4mctl/d4m"
)

func loadConf() *d4m.Settings {
	s, err := d4m.Load()
	if err != nil {
		fmt.Println("issue loading configuration:", err)
		os.Exit(1)
	}
	return s
}

var cmdk8s = &cobra.Command{
	Use:   "k8s [status|enable|disable]",
	Short: "get or set state of k8s-enabled in the local docker-for-mac installation",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := loadConf()
		switch args[0] {
		case "status":
			fmt.Println(s.KubernetesEnabled)
		case "enable":
			s.KubernetesEnabled = true
			if err := s.Write(); err != nil {
				fmt.Println("issue writing:", err)
				os.Exit(1)
			}
		case "disable":
			s.KubernetesEnabled = false
			if err := s.Write(); err != nil {
				fmt.Println("issue writing:", err)
				os.Exit(1)
			}
		}
	},
}

var cmdCpus = &cobra.Command{
	Use:   "cpus [n]",
	Short: "get or set number of cpus enabled in the local docker-for-mac installation",
	Run: func(cmd *cobra.Command, args []string) {
		s := loadConf()
		if len(args) == 0 {
			fmt.Println(s.Cpus)
		} else {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("issue interpreting n:", err)
				os.Exit(1)
			}
			s.Cpus = n
			if err := s.Write(); err != nil {
				fmt.Println("issue writing:", err)
				os.Exit(1)
			}
		}
	},
}

var cmdMem = &cobra.Command{
	Use:   "mem [n]",
	Short: "get or set amount of memory alllocated for the local docker-for-mac installation (in megabytes)",
	Run: func(cmd *cobra.Command, args []string) {
		s := loadConf()
		if len(args) == 0 {
			fmt.Println(s.MemoryMiB)
		} else {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("issue interpreting n:", err)
				os.Exit(1)
			}
			s.MemoryMiB = n
			if err := s.Write(); err != nil {
				fmt.Println("issue writing:", err)
				os.Exit(1)
			}
		}
	},
}

var cmdDump = &cobra.Command{
	Use:   "dump",
	Short: "dump current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		json.NewEncoder(os.Stdout).Encode(loadConf())
	},
}

var restartWait bool
var cmdRestart = &cobra.Command{
	Use:   "restart",
	Short: "restart docker-for-mac",
	Run: func(cmd *cobra.Command, args []string) {
		err := d4m.Restart(restartWait)
		// exit early if we complete without issue
		if err == nil {
			return
		}
		fmt.Println("issue restarting:", err)
		fmt.Println("retrying")
		if err = d4m.Restart(restartWait); err != nil {
			fmt.Println("issue restarting:", err)
			os.Exit(1)
		}
	},
}

func main() {
	rootCmd := &cobra.Command{Use: "d4mctl"}
	rootCmd.AddCommand(cmdDump)
	rootCmd.AddCommand(cmdk8s)
	rootCmd.AddCommand(cmdCpus)
	rootCmd.AddCommand(cmdMem)
	rootCmd.AddCommand(cmdRestart)
	cmdRestart.Flags().BoolVarP(&restartWait, "wait", "w", false, "If provided, wait until the docker engine is back up.")
	rootCmd.Execute()
}
