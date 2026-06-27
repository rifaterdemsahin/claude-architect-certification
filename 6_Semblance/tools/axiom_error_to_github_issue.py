import os
import requests
import json
import datetime

# Environment Variables
AXIOM_TOKEN = os.getenv("AXIOM_TOKEN")
AXIOM_DATASET = os.getenv("AXIOM_DATASET", "videoproduction")
AXIOM_QUERY_URL = os.getenv("AXIOM_QUERY_URL", "https://api.axiom.co")
OPENROUTER_API_KEY = os.getenv("OPENROUTER_API_KEY")
GITHUB_TOKEN = os.getenv("GITHUB_TOKEN")
GITHUB_REPOSITORY = os.getenv("GITHUB_REPOSITORY")

def query_axiom_errors():
    if not AXIOM_TOKEN:
        print("Missing AXIOM_TOKEN")
        return []
    
    # We query errors from the last 24 hours
    headers = {
        "Authorization": f"Bearer {AXIOM_TOKEN}",
        "Content-Type": "application/json"
    }
    
    # APL Query for the dataset, looking for ERROR or FATAL logs
    query = f"['{AXIOM_DATASET}'] | where _time > now(-24h) | where severity == 'ERROR' or severity == 'FATAL' or level == 'error' | limit 10"
    
    url = f"{AXIOM_QUERY_URL}/v1/datasets/_apl?format=legacy"
    
    try:
        response = requests.post(url, headers=headers, json={"apl": query})
        response.raise_for_status()
        data = response.json()
        
        matches = data.get("matches", [])
        return matches
    except Exception as e:
        print(f"Error querying Axiom: {e}")
        if hasattr(e, 'response') and e.response is not None:
            print(f"Response: {e.response.text}")
        return []

def analyze_error_with_openrouter(error_match):
    if not OPENROUTER_API_KEY:
        print("Missing OPENROUTER_API_KEY")
        return "No analysis provided because OPENROUTER_API_KEY is missing."
        
    url = "https://openrouter.ai/api/v1/chat/completions"
    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json",
        "HTTP-Referer": "https://github.com/rifaterdemsahin/claude-architect-certification", 
        "X-Title": "Axiom Error Analyzer"
    }
    
    # Serialize the log securely
    error_json = json.dumps(error_match, indent=2)
    
    prompt = (
        "You are an expert Full-Stack Developer and DevOps Engineer.\n"
        "Analyze the following server-side error log retrieved from Axiom.\n"
        "Provide a root cause analysis and a concrete, actionable fix that a local agent or developer can implement.\n\n"
        "Error Log:\n"
        f"```json\n{error_json}\n```"
    )
    
    payload = {
        "model": "anthropic/claude-3-opus",
        "messages": [
            {"role": "system", "content": "You are a debugging assistant."},
            {"role": "user", "content": prompt}
        ]
    }
    
    try:
        response = requests.post(url, headers=headers, json=payload)
        response.raise_for_status()
        data = response.json()
        return data["choices"][0]["message"]["content"]
    except Exception as e:
        print(f"Error analyzing with OpenRouter: {e}")
        if hasattr(e, 'response') and e.response is not None:
            print(f"Response: {e.response.text}")
        return f"Failed to analyze error with OpenRouter. Error: {e}"

def create_github_issue(error_match, analysis_text):
    if not GITHUB_TOKEN or not GITHUB_REPOSITORY:
        print("Missing GITHUB_TOKEN or GITHUB_REPOSITORY")
        return
        
    url = f"https://api.github.com/repos/{GITHUB_REPOSITORY}/issues"
    headers = {
        "Authorization": f"token {GITHUB_TOKEN}",
        "Accept": "application/vnd.github.v3+json"
    }
    
    # Try to extract a meaningful title from the log
    description = error_match.get("description", "")
    message = error_match.get("message", "")
    error_message = error_match.get("error", "")
    
    error_summary = description if description else message
    if not error_summary:
        error_summary = error_message
    if not error_summary:
        error_summary = "Unknown Server Error"
        
    # Truncate summary for title
    title = f"🚨 Axiom Server Error: {error_summary[:60]}{'...' if len(error_summary) > 60 else ''}"
    
    body = (
        f"## 💥 Server Error Detected\n\n"
        f"An error was detected in Axiom logs in the past 24 hours.\n\n"
        f"### 🤖 AI Analysis (Claude 3 Opus via OpenRouter)\n"
        f"{analysis_text}\n\n"
        f"### 📜 Raw Log Metadata\n"
        f"```json\n{json.dumps(error_match, indent=2)}\n```\n\n"
        f"**Goal**: A local agent should pull this, check the analysis, implement the fix, and close this issue."
    )
    
    payload = {
        "title": title,
        "body": body,
        "labels": ["bug", "axiom-error", "ai-analyzed"]
    }
    
    try:
        response = requests.post(url, headers=headers, json=payload)
        response.raise_for_status()
        print(f"Successfully created GitHub Issue: {response.json().get('html_url')}")
    except Exception as e:
        print(f"Error creating GitHub issue: {e}")
        if hasattr(e, 'response') and e.response is not None:
            print(f"Response: {e.response.text}")

def main():
    print("Starting Axiom Error Analyzer...")
    errors = query_axiom_errors()
    
    if not errors:
        print("No errors found in the last 24 hours.")
        return
        
    print(f"Found {len(errors)} errors. Analyzing the most recent distinct error...")
    
    # Process the first one to avoid issue spam
    target_error = errors[0]
    
    analysis = analyze_error_with_openrouter(target_error)
    create_github_issue(target_error, analysis)
    
if __name__ == "__main__":
    main()
