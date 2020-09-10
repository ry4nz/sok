package swarmctl

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"

	"github.com/ry4nz/sok/cmd"
	"github.com/ry4nz/sok/constants"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

func NewJoinCommand(clientset kubernetes.Clientset) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join <node_name>",
		Short: "Join node to the active Swarm namespace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJoin(clientset, args[0])
		},
	}
	return cmd
}

func runJoin(clientset kubernetes.Clientset, nodename string) error {
	namespaces, err := cmd.GetNamespaces()
	if err != nil {
		return err
	}
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	node, err := cmd.GetNodeByUIDorName(nodes.Items, nodename)
	if err != nil {
		return err
	}
	for k, v := range node.Labels {
		if k == constants.SwarmNamespace {
			if namespaces.Active == v {
				return fmt.Errorf("node %s already joined swarm namespace %s, no-op", nodename, v)
			} else {
				return fmt.Errorf("node %s already joined swarm namespace %s, leave first before join another", nodename, v)
			}
		}
	}

	oldData, err := json.Marshal(node)
	if err != nil {
		return nil
	}
	node.Labels[constants.SwarmNamespace] = namespaces.Active
	newData, err := json.Marshal(node)
	if err != nil {
		return nil
	}
	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1.Node{})
	if err != nil {
		return nil
	}
	_, err = clientset.CoreV1().Nodes().Patch(context.TODO(), node.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("node %s joined swarm namespace %s\n", nodename, namespaces.Active)
	return nil
}
