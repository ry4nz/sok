package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"

	"github.com/ry4nz/sok/types"
	"github.com/spf13/cobra"
)

// ShowHelp shows the command help.
func ShowHelp() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
func GetConfigFile() string {
	homedir, _ := os.UserHomeDir()
	return filepath.Join(homedir, ".swarmctl")
}

func GetNamespaces() (types.Namespaces, error) {
	var namespaces types.Namespaces

	configFile := GetConfigFile()
	_, err := os.Stat(configFile)
	if err != nil {
		return namespaces, err
	}

	jsonFile, err := os.Open(configFile)
	if err != nil {
		return namespaces, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return namespaces, err
	}
	jsonFile.Close()

	err = json.Unmarshal(byteValue, &namespaces)
	return namespaces, err
}

func GetNodeByUIDorName(nodes []corev1.Node, uidOrName string) (corev1.Node, error) {
	for _, n := range nodes {
		if string(n.UID) == uidOrName || n.Name == uidOrName {
			return n, nil
		}
	}
	return corev1.Node{}, errors.New("unable to find the node")
}
