import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";
import * as dotenv from "dotenv";

// Load environment variables
dotenv.config();

const server = new Server(
  {
    name: "enterprise-secure-gateway",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// Mock database representing system infrastructure states
const cloudTopologies = [
  {
    id: "aws-zdr-vpc",
    name: "AWS PrivateLink Zero-Data Retention VPC",
    description: "Multi-AZ VPC routing Claude API requests securely via private endpoints with no public transit.",
    tags: ["aws", "security", "vpc", "zdr"]
  },
  {
    id: "multi-agent-orchestrator",
    name: "Deterministic Multi-Agent State Machine",
    description: "Orchestration layer containing hop limits and state checks to prevent looping.",
    tags: ["patterns", "multi-agent", "orchestration"]
  }
];

// Register list of available tools
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "retrieve_architecture_patterns",
        description: "Retrieve approved enterprise-grade cloud reference architecture patterns.",
        inputSchema: {
          type: "object",
          properties: {
            category: {
              type: "string",
              description: "Category of architecture (e.g., 'security', 'vpc', 'multi-agent')",
            },
          },
          required: ["category"],
        },
      },
      {
        name: "verify_zdr_compliance",
        description: "Verify if a specific cloud topology violates Zero-Data Retention rules.",
        inputSchema: {
          type: "object",
          properties: {
            topologyId: {
              type: "string",
              description: "The ID of the topology to audit (e.g., 'aws-zdr-vpc')",
            },
            enableLogging: {
              type: "boolean",
              description: "Flag indicating whether payload logging is turned on (MUST be false for ZDR)",
            },
          },
          required: ["topologyId", "enableLogging"],
        },
      },
    ],
  };
});

// Handle tool execution requests
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  try {
    if (name === "retrieve_architecture_patterns") {
      const category = (args?.category as string).toLowerCase();
      const filtered = cloudTopologies.filter(item => 
        item.tags.includes(category)
      );

      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              status: "success",
              results: filtered.length > 0 ? filtered : "No matching patterns found.",
            }, null, 2),
          },
        ],
      };
    }

    if (name === "verify_zdr_compliance") {
      const topologyId = args?.topologyId as string;
      const enableLogging = args?.enableLogging as boolean;

      const target = cloudTopologies.find(t => t.id === topologyId);
      if (!target) {
        return {
          content: [
            {
              type: "text",
              text: JSON.stringify({
                status: "error",
                message: `Topology with ID '${topologyId}' not found.`,
              }),
            },
          ],
          isError: true,
        };
      }

      // Security check: ZDR is breached if logging is enabled on payloads
      if (enableLogging) {
        return {
          content: [
            {
              type: "text",
              text: JSON.stringify({
                status: "failed",
                compliance: "NON-COMPLIANT",
                reason: "Payload logging is enabled. To satisfy Zero-Data Retention (ZDR), payload logging must be disabled.",
              }, null, 2),
            },
          ],
        };
      }

      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              status: "success",
              compliance: "COMPLIANT",
              details: "Topology conforms to standard security guardrails and data retention rules.",
            }, null, 2),
          },
        ],
      };
    }

    throw new Error(`Tool not found: ${name}`);
  } catch (error: any) {
    return {
      content: [
        {
          type: "text",
          text: JSON.stringify({
            status: "error",
            message: error.message || "An unknown error occurred during execution.",
          }),
        },
      ],
      isError: true,
    };
  }
});

// Initialize server transport
async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error("Enterprise Secure Gateway MCP Server running on stdio");
}

main().catch((error) => {
  console.error("Fatal error in MCP Server:", error);
  process.exit(1);
});
