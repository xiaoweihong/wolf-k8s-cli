package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"wolf-k8s-cli/configs"
)

var rootCmd = &cobra.Command{
	Use:   "wolf-k8s-cli",
	Short: "战狼k8s工具",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(monadoCmd)
	if home := homedir.HomeDir(); home != "" {
		configs.Kubeconfig = rootCmd.PersistentFlags().String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		configs.Kubeconfig = rootCmd.PersistentFlags().String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
}
