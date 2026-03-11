using System;

namespace UnityCliConnector
{
    /// <summary>
    /// Alias for backward compatibility with existing [McpForUnityTool] tools.
    /// Projects can keep their existing attribute and it will be discovered.
    /// </summary>
    [AttributeUsage(AttributeTargets.Class)]
    public class McpForUnityToolAttribute : Attribute
    {
        public string Description { get; set; } = "";
        public string Name { get; set; }
        public string Group { get; set; } = "";
        public bool AutoRegister { get; set; } = true;
        public bool RequiresPolling { get; set; }
        public string PollAction { get; set; }
        public bool StructuredOutput { get; set; }
    }
}
