package swarmctl

import (
	"github.com/ry4nz/sok/cmd"
	"github.com/ry4nz/sok/cmd/swarmctl/node"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

func CmdSwarmctl(clientset kubernetes.Clientset) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "node",
		RunE: cmd.ShowHelp(),
	}
	cmd.AddCommand(
		node.NewNodeCommand(clientset),
	)
	return cmd
}
