package tool

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/mark3labs/mcp-go/mcp"
)

func CreateContainerHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	imageName := request.Params.Arguments["image"].(string)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error while initalizing docker: %v", err)), nil
	}

	var selectedImage image.InspectResponse
	selectedImage, err = cli.ImageInspect(ctx, imageName)
	if errdefs.IsNotFound(err) {
		pulledImage, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error while pulling image: %v", err)), nil
		}
		defer pulledImage.Close()

		_, err = io.Copy(io.Discard, pulledImage)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error while pulling image: %v", err)), nil
		}

		selectedImage, err = cli.ImageInspect(ctx, imageName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error while inspecting pulled image: %v", err)), nil
		}

	} else if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error while inspecting image: %v", err)), nil
	}

	containerCreated, err := cli.ContainerCreate(ctx,
		&container.Config{Image: selectedImage.ID}, nil, nil, nil, fmt.Sprintf("%s-%d", imageName, time.Now().Unix()))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error while creating container: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("Container created with %s image and %s containerID.", imageName, containerCreated.ID),
	), nil
}
