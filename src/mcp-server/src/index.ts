import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { SSEServerTransport } from "@modelcontextprotocol/sdk/server/sse.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";
import sqlite3 from "sqlite3";
import { promisify } from "util";
import dotenv from "dotenv";
import express from "express";

dotenv.config();

// Initialize an isolated, in-memory enterprise database for demonstration/testing
const db = new sqlite3.Database(":memory:");
const dbAll = (sql: string, params: any[] = []): Promise<any[]> => {
  return new Promise((resolve, reject) => {
    db.all(sql, params, (err, rows) => {
      if (err) reject(err);
      else resolve(rows);
    });
  });
};

const dbRun = (sql: string, params: any[] = []): Promise<void> => {
  return new Promise((resolve, reject) => {
    db.run(sql, params, (err) => {
      if (err) reject(err);
      else resolve();
    });
  });
};

// Seed the mockup enterprise schema (Simulating an active CRM/ERP infrastructure)
async function seedDatabase() {
  await dbRun(`
    CREATE TABLE IF NOT EXISTS enterprise_inventory (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      sku TEXT UNIQUE,
      region TEXT,
      stock_level INTEGER,
      classification TEXT
    )
  `);
  
  await dbRun(`
    INSERT OR IGNORE INTO enterprise_inventory (sku, region, stock_level, classification)
    VALUES 
    ('CLD-35-SONNET', 'eu-west-1', 450, 'Tier-1'),
    ('CLD-35-HAIKU', 'us-east-1', 1200, 'Tier-2'),
    ('MCP-GW-01', 'eu-west-2', 85, 'Critical')
  `);
}

// Instantiate the Model Context Protocol Server
const server = new Server(
  {
    name: "enterprise-data-bridge",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {}, // Expose tool capabilities to Claude
    },
  }
);

// Register Tool Discovery Lifecycle Hook
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "query_inventory",
        description: "Securely queries regional enterprise stock levels and classifications. Read-only access enforced.",
        inputSchema: {
          type: "object",
          properties: {
            region: {
              type: "string",
              description: "The targeted cloud infrastructure region (e.g., eu-west-1, us-east-1)",
            },
          },
          required: ["region"],
        },
      },
    ],
  };
});

// Register Tool Execution Handler (With strict boundary constraints)
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  if (request.params.name !== "query_inventory") {
    throw new Error(`Unsupported tool requested: ${request.params.name}`);
  }

  const region = request.params.arguments?.region as string;
  
  try {
    // Architectural Guardrail: Enforcing parameterized inputs to mitigate SQL injection at the agent boundary
    const rows = await dbAll(
      "SELECT sku, region, stock_level, classification FROM enterprise_inventory WHERE region = ?",
      [region]
    );

    return {
      content: [
        {
          type: "text",
          text: JSON.stringify({ status: "SUCCESS", data: rows }, null, 2),
        },
      ],
    };
  } catch (error: any) {
    return {
      isError: true,
      content: [
        {
          type: "text",
          text: `Architectural Fault Triggered: ${error.message}`,
        },
      ],
    };
  }
});

// Launch server using standard input/output transport layers or SSE depending on configuration
async function run() {
  await seedDatabase();

  const useSse = process.env.TRANSPORT === "sse" || !!process.env.PORT;

  if (useSse) {
    const app = express();
    const port = process.env.PORT || 8080;
    
    // Store active transport mapping for routing messages
    let activeTransport: SSEServerTransport | null = null;

    app.get("/sse", async (req, res) => {
      console.error("🔒 New SSE connection received on enterprise gateway");
      activeTransport = new SSEServerTransport("/messages", res);
      await server.connect(activeTransport);
    });

    app.post("/messages", express.json(), async (req, res) => {
      if (activeTransport) {
        await activeTransport.handlePostMessage(req, res);
      } else {
        res.status(400).send("No active SSE transport session. Connect to /sse first.");
      }
    });

    app.get("/health", (req, res) => {
      res.json({ status: "healthy", server: "enterprise-data-bridge" });
    });

    app.listen(port, () => {
      console.error(`🔒 Enterprise MCP SSE Server listening on port ${port}`);
    });
  } else {
    // Default Stdio connection for local run configurations
    const transport = new StdioServerTransport();
    await server.connect(transport);
    console.error("🔒 Enterprise MCP Data Bridge active on STDIO transport layer.");
  }
}

run().catch((error) => {
  console.error("Fatal startup crash in MCP Server:", error);
  process.exit(1);
});
