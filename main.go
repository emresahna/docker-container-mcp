package main

import (
	"log"

	"github.com/emresahna/docker-container-mcp/resource"
	"github.com/emresahna/docker-container-mcp/tool"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Docker Container MCP",
		"0.0.0",
		server.WithLogging(),
	)

	containerCreateTool := mcp.NewTool("create_container",
		mcp.WithDescription("Create docker container"),
		mcp.WithString("image",
			mcp.Required(),
			mcp.Description("Image to create container."),
		),
	)

	containerStatusResourceTemplate := mcp.NewResourceTemplate(
		"status://{containerId}",
		"Docker Container Status",
		mcp.WithTemplateDescription("Gives information about container status"),
		mcp.WithTemplateMIMEType("text/plain"),
	)

	s.AddResourceTemplate(containerStatusResourceTemplate, resource.ContainerStatus)

	s.AddTool(containerCreateTool, tool.CreateContainerHandler)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
