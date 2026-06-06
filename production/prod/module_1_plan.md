# 🎬 Module 1 Production: Screen Action & Speaking Scripts

## 🎥 Section 1: The Hook (0:00 - 2:30)
* **Visual Action:** Open your web browser to the live repository dashboard: `https://rifaterdemsahin.github.io/claude-architect-certification/` (or display the local `index.html` file in browser). Activate your webcam in picture-in-picture mode and hover your cursor over the "Enterprise Masterclass" badges.
* **Speaking Script:**
  > *"Welcome to Module 1. Moving LLMs from a playground prototype or an experimental chat window into a secure, scalable production system requires systems engineering—not just clever prompting. Today, we are breaking down the official Claude Architecture Blueprint. We are moving beyond basic API calls to understand model topologies, stateless middleware boundaries, and how the Model Context Protocol changes enterprise integration forever. All the code, blueprints, and architecture files we cover are live right here on our project hub."*

---

## 🎥 Section 2: The Mechanics (2:30 - 7:45)
* **Visual Action:** Present an architecture diagram (or draw on-screen) illustrating a stateless payload loop. Highlight how standard single-turn requests compare to context payloads that pass system instructions, history, and schemas.
* **Speaking Script:**
  > *"To build resilient systems, we must accept one core rule: the Claude API is entirely stateless. It retains zero memory of previous interactions. If you don't pass the context back, it doesn't exist. Look at this architecture mapping. In an enterprise configuration, every turn requires a stateless middleware layer to coordinate history, append corporate governance boundaries, and pass structural schemas back into the context window. We're going to optimize this using Anthropic's 200k token limits, showing how to structure system vs. user space safely."*

---

## 🎥 Section 3: The Tech Stack Setup (7:45 - 12:00)
* **Visual Action:** Switch your screen to VS Code displaying `/src/mcp-server`. Expand the `/src` folder, open `index.ts`, open the terminal panel, run `cd src/mcp-server && npm run build`, and then execute `npm start` to show the silent diagnostics message.
* **Speaking Script:**
  > *"Let's look at the actual engineering layout in the IDE. Inside our project directory, under `src/mcp-server`, we’ve initialized a clean, container-ready TypeScript and SQLite foundation. When I run our build and initialize the process, notice that it outputs an isolated diagnostic indicator directly to standard error. This is critical because Model Context Protocol servers communicate over standard input and output streams (`stdio`). If your initialization code outputs random text or logs to standard out, it corrupts the JSON-RPC communication layer that Claude uses to execute local tools. This clean, read-only baseline is exactly where we will implement our custom enterprise database bridge in the very next module."*

---

## 🎥 Section 4: Call to Action (12:00 - End)
* **Visual Action:** Open your root GitHub repository page: `github.com/rifaterdemsahin/claude-architect-certification` and highlight the `git clone` block in the README.
* **Speaking Script:**
  > *"To follow along with the hands-on code deployments, the security compliance matrices, and the multi-agent routing engines we are spinning up, pull down this repository right now. Run `git clone` on your local machine, check out the directory layouts, and get your environment ready. In Module 2, we are taking this local setup and deploying it straight to live edge microVM production servers on Fly.io. Hit subscribe, clone the architecture blueprint, and let’s build production-grade AI."*
