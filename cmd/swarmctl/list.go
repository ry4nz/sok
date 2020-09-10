package swarmctl

import (
	"fmt"

	"github.com/ry4nz/sok/cmd"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Swarm namespaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList()
		},
	}
	return cmd
}

func runList() error {
	namespaces, err := cmd.GetNamespaces()
	if err != nil {
		return err
	}

	fmt.Printf("Swarm Namespaces: %d\n", len(namespaces.Namespaces))
	for _, n := range namespaces.Namespaces {
		if n.Name == namespaces.Active {
			fmt.Printf("%s *\n", n.Name)
		} else {
			fmt.Println(n.Name)
		}
	}
	return nil
}
