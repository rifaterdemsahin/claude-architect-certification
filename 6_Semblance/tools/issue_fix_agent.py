import os
import json
import subprocess
import requests

OPENROUTER_API_KEY = os.getenv("OPENROUTER_API_KEY")

def run_cmd(cmd):
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"Error running: {cmd}\nOutput: {result.stderr}")
    return result.stdout.strip(), result.returncode

def fetch_open_issues():
    stdout, rc = run_cmd("gh issue list --label 'axiom-error' --state open --json number,title,body")
    if rc != 0:
        return []
    return json.loads(stdout)

def ask_openrouter_for_fix(issue_body):
    if not OPENROUTER_API_KEY:
        print("Missing OPENROUTER_API_KEY")
        return None
        
    url = "https://openrouter.ai/api/v1/chat/completions"
    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json"
    }
    
    prompt = (
        "You are an autonomous AI coding agent. Your goal is to apply the fix described in the issue below to the codebase.\n"
        "Return ONLY a JSON array containing the files you want to create or modify, with their FULL updated file contents.\n"
        "Format your output strictly as a JSON array (no markdown blocks, no conversational text).\n"
        'Example: [{"file_path": "cmd/server/main.go", "content": "package main\\n..."}]\n\n'
        f"Issue Description:\n{issue_body}\n"
    )
    
    payload = {
        "model": "anthropic/claude-opus-4.6",
        "max_tokens": 4000,
        "messages": [
            {"role": "system", "content": "You output only valid JSON."},
            {"role": "user", "content": prompt}
        ]
    }
    
    try:
        response = requests.post(url, headers=headers, json=payload)
        response.raise_for_status()
        content = response.json()["choices"][0]["message"]["content"]
        
        # Clean up markdown if the LLM still wrapped it
        if content.startswith("```json"):
            content = content[7:]
        if content.startswith("```"):
            content = content[3:]
        if content.endswith("```"):
            content = content[:-3]
            
        return json.loads(content.strip())
    except Exception as e:
        print(f"Error asking OpenRouter: {e}")
        if hasattr(e, 'response') and e.response is not None:
            print(f"Response: {e.response.text}")
        return None

def apply_fixes(fixes):
    for fix in fixes:
        file_path = fix.get("file_path")
        content = fix.get("content")
        
        if not file_path or not content:
            continue
            
        # Ensure directory exists
        os.makedirs(os.path.dirname(file_path) or ".", exist_ok=True)
        
        with open(file_path, "w") as f:
            f.write(content)
        print(f"Patched {file_path}")

def process_issue(issue):
    number = issue["number"]
    title = issue["title"]
    print(f"Processing Issue #{number}: {title}")
    
    fixes = ask_openrouter_for_fix(issue["body"])
    
    if not fixes:
        print(f"Failed to generate fix for issue #{number}")
        return
        
    apply_fixes(fixes)
    
    # Check if there are git changes
    status_out, _ = run_cmd("git status --porcelain")
    if not status_out:
        print(f"No changes made for issue #{number}")
        return
        
    # Commit and close
    run_cmd("git add .")
    run_cmd(f'git commit -m "🤖 Auto-fix issue #{number}: {title}"')
    
    # Verify the build locally if applicable (optional safety step)
    # run_cmd("go build ./...") 
    
    run_cmd(f'gh issue comment {number} --body "✅ This issue has been automatically fixed by the Issue Fix Agent. Changes have been committed."')
    run_cmd(f'gh issue close {number} -r completed')
    
    print(f"Successfully closed issue #{number}")

def main():
    print("Starting Issue Fix Agent...")
    issues = fetch_open_issues()
    if not issues:
        print("No open axiom-error issues found.")
        return
        
    for issue in issues:
        process_issue(issue)
        
    # Push changes if any commits were made
    print("Pushing applied fixes to remote...")
    run_cmd("git push")

if __name__ == "__main__":
    main()
