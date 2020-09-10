package node

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/ry4nz/sok/constants"

	"github.com/docker/cli/cli/command/formatter"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	nodeIDHeader         = "ID"
	hostnameHeader       = "HOSTNAME"
	availabilityHeader   = "AVAILABILITY"
	managerStatusHeader  = "MANAGER STATUS"
	engineVersionHeader  = "ENGINE VERSION"
	swarmNamespaceHeader = "SWARM NAMESPACE"
)

type nodeContext struct {
	formatter.HeaderContext
	n corev1.Node
}

func (c *nodeContext) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(c)
}

func (c *nodeContext) ID() string {
	return string(c.n.UID)
}

func (c *nodeContext) Hostname() string {
	return c.n.Name
}

func (c *nodeContext) Status() string {
	for _, cond := range c.n.Status.Conditions {
		if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
			return "Ready"
		}
	}
	return "NotReady"
}

func (c *nodeContext) Availability() string {
	return "active"
}

func (c *nodeContext) ManagerStatus() string {
	var statuses []string
	for k, _ := range c.n.Labels {
		if strings.HasPrefix(k, constants.NodeRole) {
			statuses = append(statuses, strings.TrimPrefix(k, constants.NodeRole))
		}
	}
	return strings.Join(statuses, ",")
}

func (c *nodeContext) EngineVersion() string {
	return c.n.Status.NodeInfo.ContainerRuntimeVersion
}

func (c *nodeContext) SwarmNamespace() string {
	for k, v := range c.n.Labels {
		if k == constants.SwarmNamespace {
			return v
		}
	}
	return ""
}

func newListCommand(clientset kubernetes.Clientset) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Aliases: []string{"list"},
		Short:   "List nodes in the swarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(clientset)
		},
	}

	return cmd
}

func runList(clientset kubernetes.Clientset) error {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	nodesCtx := formatter.Context{
		Output: buf,
		Format: "table {{.ID}} \t{{.Hostname}}\t{{.Status}}\t{{.Availability}}\t{{.ManagerStatus}}\t{{.EngineVersion}}\t{{.SwarmNamespace}}",
	}
	err = FormatWrite(nodesCtx, nodes.Items)
	fmt.Println(buf.String())
	return err
}

// FormatWrite writes the context
func FormatWrite(ctx formatter.Context, nodes []corev1.Node) error {
	render := func(format func(subContext formatter.SubContext) error) error {
		for _, node := range nodes {
			nodeCtx := &nodeContext{n: node}
			if err := format(nodeCtx); err != nil {
				return err
			}
		}
		return nil
	}
	nodeCtx := nodeContext{}
	nodeCtx.Header = formatter.SubHeaderContext{
		"ID":             nodeIDHeader,
		"Hostname":       hostnameHeader,
		"Status":         formatter.StatusHeader,
		"Availability":   availabilityHeader,
		"ManagerStatus":  managerStatusHeader,
		"EngineVersion":  engineVersionHeader,
		"SwarmNamespace": swarmNamespaceHeader,
	}
	return ctx.Write(&nodeCtx, render)
}
