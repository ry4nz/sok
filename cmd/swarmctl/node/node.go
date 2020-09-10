package node

import (
	"github.com/ry4nz/sok/cmd"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

// NewNodeCommand returns a cobra command for `node` subcommands
func NewNodeCommand(clientset kubernetes.Clientset) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage Swarm nodes",
		RunE:  cmd.ShowHelp(),
	}
	cmd.AddCommand(
		newListCommand(clientset),
	)
	return cmd
}
