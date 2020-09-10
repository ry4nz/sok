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

func NewLeaveCommand(clientset kubernetes.Clientset) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leave <node_name>",
		Short: "leave node from the active Swarm namespace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLeave(clientset, args[0])
		},
	}
	return cmd
}

func runLeave(clientset kubernetes.Clientset, nodename string) error {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	node, err := cmd.GetNodeByUIDorName(nodes.Items, nodename)
	if err != nil {
		return err
	}
	for k, _ := range node.Labels {
		if k == constants.SwarmNamespace {

			oldData, err := json.Marshal(node)
			if err != nil {
				return nil
			}
			delete(node.Labels, constants.SwarmNamespace)
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
			fmt.Printf("node %s left swarm namespace\n", nodename)
			return nil
		}
	}
	return fmt.Errorf("node %s doesn not belong to any swarm namespace", nodename)

}
