package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/youngwoocho02/unity-cli/internal/client"
)

type toolParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
}

type toolSchema struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Group       string      `json:"group"`
	Parameters  []toolParam `json:"parameters"`
}

var builtinTools = []toolSchema{
	{
		Name: "editor", Description: "Control editor state: play, stop, pause, refresh.", Group: "built-in",
		Parameters: []toolParam{
			{Name: "action", Description: "play, stop, pause, refresh", Type: "String", Required: true},
			{Name: "wait_for_completion", Description: "Wait for action to complete", Type: "Boolean"},
		},
	},
	{
		Name: "console", Description: "Read or clear Unity console logs.", Group: "built-in",
		Parameters: []toolParam{
			{Name: "filter", Description: "Comma-separated log types: error, warning, log (default: error,warning)", Type: "String"},
			{Name: "lines", Description: "Maximum number of log entries to return", Type: "Int32"},
			{Name: "filter_text", Description: "Filter log messages containing this text", Type: "String"},
			{Name: "stacktrace", Description: "Stack trace mode: none, short, full (default: none)", Type: "String"},
			{Name: "clear", Description: "Clear console", Type: "Boolean"},
		},
	},
	{
		Name: "exec", Description: "Execute arbitrary C# code at runtime.", Group: "built-in",
		Parameters: []toolParam{
			{Name: "code", Description: "C# code to execute", Type: "String", Required: true},
			{Name: "usings", Description: "Additional using directives (comma-separated)", Type: "String"},
		},
	},
	{
		Name: "menu", Description: "Execute a Unity menu item by path.", Group: "built-in",
		Parameters: []toolParam{
			{Name: "menu_path", Description: "Menu item path (e.g. File/Save Project)", Type: "String", Required: true},
		},
	},
	{
		Name: "screenshot", Description: "Capture scene/game view as PNG.", Group: "built-in",
		Parameters: []toolParam{
			{Name: "view", Description: "scene (default), game", Type: "String"},
			{Name: "width", Description: "Image width (default: 1920)", Type: "Int32"},
			{Name: "height", Description: "Image height (default: 1080)", Type: "Int32"},
			{Name: "output_path", Description: "Output path (default: Screenshots/screenshot.png)", Type: "String"},
		},
	},
	{
		Name: "reserialize", Description: "Force reserialize Unity assets.", Group: "built-in",
		Parameters: []toolParam{
			{Name: "path", Description: "Single asset path", Type: "String"},
			{Name: "paths", Description: "Multiple asset paths", Type: "String[]"},
		},
	},
	{
		Name: "profiler", Description: "Control Unity Profiler.", Group: "built-in",
		Parameters: []toolParam{
			{Name: "action", Description: "hierarchy, enable, disable, status, clear", Type: "String", Required: true},
			{Name: "depth", Description: "Recursive depth (default: 1)", Type: "Int32"},
			{Name: "root", Description: "Set root by name (substring match)", Type: "String"},
			{Name: "frames", Description: "Average over last N frames", Type: "Int32"},
			{Name: "from", Description: "Start frame index", Type: "Int32"},
			{Name: "to", Description: "End frame index", Type: "Int32"},
			{Name: "parent", Description: "Drill into item by ID", Type: "Int32"},
			{Name: "min", Description: "Minimum time (ms) filter", Type: "Single"},
			{Name: "sort", Description: "Sort by: total, self, calls", Type: "String"},
			{Name: "max", Description: "Max children per level (default: 30)", Type: "Int32"},
			{Name: "frame", Description: "Specific frame index", Type: "Int32"},
			{Name: "thread", Description: "Thread index (0=main)", Type: "Int32"},
		},
	},
	{
		Name: "test", Description: "Run EditMode/PlayMode tests.", Group: "built-in",
		Parameters: []toolParam{
			{Name: "mode", Description: "EditMode (default) or PlayMode", Type: "String", Required: true},
			{Name: "filter", Description: "Filter by test name", Type: "String"},
		},
	},
	{
		Name: "refresh_unity", Description: "Refresh assets and optionally compile.", Group: "built-in",
		Parameters: []toolParam{
			{Name: "mode", Description: "if_dirty (default) or force", Type: "String"},
			{Name: "scope", Description: "all (default) or specific path", Type: "String"},
			{Name: "compile", Description: "none (default) or request", Type: "String"},
		},
	},
}

func listCmd(project string, port int, timeout int) error {
	all := make([]toolSchema, len(builtinTools))
	copy(all, builtinTools)

	// Try to fetch custom tools from Unity
	inst, err := client.DiscoverInstance(project, port)
	if err == nil {
		if err := waitForAlive(inst.Port, 3000); err == nil {
			resp, err := client.Send(inst, "list_tools", map[string]interface{}{}, timeout)
			if err == nil && resp.Success {
				var remote []toolSchema
				if json.Unmarshal(resp.Data, &remote) == nil {
					// Add custom tools (skip built-in duplicates)
					names := map[string]bool{
						"manage_editor": true,
						"refresh_unity": true,
						"run_tests":     true,
					}
					for _, t := range builtinTools {
						names[t.Name] = true
					}
					for _, t := range remote {
						if !names[t.Name] {
							if t.Group == "" {
								t.Group = "custom"
							}
							all = append(all, t)
						}
					}
				}
			}
		}
	}

	b, _ := json.MarshalIndent(all, "", "  ")
	fmt.Println(string(b))
	return nil
}
