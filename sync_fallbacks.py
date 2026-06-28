import json

with open('navigation_config.json', 'r') as f:
    nav_str = f.read()

def replace_fallback(filename):
    with open(filename, 'r') as f:
        content = f.read()
    
    start_str = "let navigationData = {"
    start_idx = content.find(start_str)
    if start_idx == -1:
        print(f"Could not find start in {filename}")
        return
        
    brace_count = 0
    in_str = False
    escape = False
    
    for i in range(start_idx + len("let navigationData = "), len(content)):
        char = content[i]
        
        if in_str:
            if escape:
                escape = False
            elif char == '\\':
                escape = True
            elif char == '"':
                in_str = False
        else:
            if char == '"':
                in_str = True
            elif char == '{':
                brace_count += 1
            elif char == '}':
                brace_count -= 1
                if brace_count == 0:
                    end_idx = i
                    break
    
    end_idx = content.find(';', end_idx)
    
    new_content = content[:start_idx] + "let navigationData = " + nav_str.strip() + ";" + content[end_idx+1:]
    
    with open(filename, 'w') as f:
        f.write(new_content)
    print(f"Updated {filename}")

replace_fallback('index.html')
replace_fallback('markdown_renderer.html')
