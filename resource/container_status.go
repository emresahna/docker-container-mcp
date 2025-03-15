package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func ContainerStatus(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	cotaninerId := extractContainerId(request.Params.URI)

	if cotaninerId == "" {
		return []mcp.ResourceContents{}, fmt.Errorf("Invalid container id.")
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return []mcp.ResourceContents{}, fmt.Errorf("Error occured while connecting docker client.")
	}

	c, err := cli.ContainerInspect(ctx, cotaninerId)
	if err != nil {
		return []mcp.ResourceContents{}, fmt.Errorf("Error while fetching container stats.")
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			Text:     fmt.Sprintf(`[%s] %s`, strings.ToUpper(c.State.Status), c.Name),
			MIMEType: "text/plain",
			URI:      request.Params.URI,
		},
	}, nil
}

func extractContainerId(uri string) string {
	s, ok := strings.CutPrefix(uri, "docker://")
	if !ok {
		return ""
	}
	return s
}
