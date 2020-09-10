package swarmctl

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ry4nz/sok/cmd"
	"github.com/ry4nz/sok/types"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <namespace>",
		Short: "Create Swarm namespaces",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(args[0])
		},
	}
	return cmd
}

func runInit(namespace string) error {
	namespaces, err := cmd.GetNamespaces()
	if os.IsNotExist(err) {
		return write(namespaces, namespace, cmd.GetConfigFile())
	}
	if err != nil {
		return err
	}

	return write(namespaces, namespace, cmd.GetConfigFile())
}
func write(namespaces types.Namespaces, namespace string, configFile string) error {

	namespaces.Namespaces = append(namespaces.Namespaces, types.Namespace{Name: namespace})

	file, _ := json.MarshalIndent(namespaces, "", " ")

	_ = ioutil.WriteFile(configFile, file, 0644)
	return nil
}
