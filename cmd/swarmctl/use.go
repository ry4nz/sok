package swarmctl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ry4nz/sok/cmd"
	"github.com/ry4nz/sok/types"
	"github.com/spf13/cobra"
)

func NewUseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use <namespace>",
		Short: "Set Active Swarm namespaces",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUse(args[0])
		},
	}
	return cmd
}

func runUse(namespace string) error {
	namespaces, err := cmd.GetNamespaces()
	if err != nil {
		return err
	}

	for _, n := range namespaces.Namespaces {
		if n.Name == namespace {
			return activate(namespaces, namespace, cmd.GetConfigFile())
		}
	}
	return fmt.Errorf("namespace %s not found", namespace)
}
func activate(namespaces types.Namespaces, namespace string, configFile string) error {

	namespaces.Active = namespace

	file, _ := json.MarshalIndent(namespaces, "", " ")

	_ = ioutil.WriteFile(configFile, file, 0644)
	fmt.Printf("Use Swarm namespace: %s\n", namespace)
	return nil
}
