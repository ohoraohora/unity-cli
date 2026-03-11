# unity-cli

Command-line tool to control Unity Editor from AI coding assistants (Claude Code, Cursor, etc.)

**No MCP protocol. No Python relay. Just a CLI.**

## Install

```bash
go install github.com/fedtop/unity-cli@latest
```

Or download a binary from [Releases](https://github.com/fedtop/unity-cli/releases).

## Unity Setup

Add the Connector package to your Unity project:

```
// Packages/manifest.json
"com.fedtop.unity-cli-connector": "https://github.com/fedtop/unity-cli.git?path=unity-connector"
```

The Connector auto-starts an HTTP server when Unity opens. The CLI discovers it automatically.

## Usage

```bash
# Editor control
unity-cli editor play --wait
unity-cli editor stop
unity-cli editor refresh --compile

# ECS queries
unity-cli query entities --world server --component Health
unity-cli query inspect --world server --index 42
unity-cli query singleton --world server --component GamePhaseState

# Game flow
unity-cli game connect
unity-cli game load --index 0
unity-cli game phase Playing
unity-cli game spawn-bots 10

# Execute C# in Unity
unity-cli exec "Debug.Log(Time.time)"

# Console logs
unity-cli console --lines 50

# Custom tools
unity-cli tool list
unity-cli tool call my_custom_tool --params '{"key":"value"}'
```

## How It Works

```
unity-cli ──HTTP POST──→ Unity Editor (localhost:8090)
                          ├── CommandRouter
                          ├── ToolDiscovery
                          └── [UnityCliTool] handlers
```

1. Unity Connector opens HTTP server on `localhost:8090`
2. Registers itself in `~/.unity-cli/instances.json`
3. CLI reads instance file, sends `POST /command` with JSON
4. Connector dispatches to the matching tool handler
5. CLI prints result to stdout

## Writing Custom Tools

```csharp
using UnityCliConnector;
using Newtonsoft.Json.Linq;

[UnityCliTool(Description = "Returns the current scene name")]
public static class GetCurrentScene
{
    // Command name: auto-derived as "get_current_scene"

    public class Parameters
    {
        [ToolParameter("Include build index", Required = false)]
        public bool IncludeBuildIndex { get; set; }
    }

    public static object HandleCommand(JObject parameters)
    {
        var scene = UnityEngine.SceneManagement.SceneManager.GetActiveScene();
        return new SuccessResponse("Current scene", new
        {
            name = scene.name,
            path = scene.path,
        });
    }
}
```

Call it:
```bash
unity-cli tool call get_current_scene
```

## Options

| Flag | Description |
|------|-------------|
| `--port <N>` | Override Unity instance port |
| `--project <path>` | Select instance by project path |
| `--json` | Raw JSON output |
| `--timeout <ms>` | Request timeout (default: 120000) |
| `--help` | Show help |
| `--version` | Show version |

## Compared to MCP

| | MCP | unity-cli |
|---|-----|-----------|
| **Install** | Python + uv + FastMCP + config JSON | `go install` or download binary |
| **Dependencies** | Python runtime, WebSocket | None (single binary) |
| **Protocol** | JSON-RPC 2.0 over stdio + WebSocket | HTTP POST |
| **Setup** | Generate MCP config, restart AI tool | Just works |
| **Compatibility** | MCP-compatible clients only | Anything with a shell |

## License

MIT
