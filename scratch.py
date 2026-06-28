import os
import json

with open('navigation_config.json', 'r') as f:
    nav_content = f.read()

html_files = []
for root, dirs, files in os.walk('.'):
    if '/_' in root or '/node_modules' in root or '/.git' in root:
        continue
    for file in files:
        if file.endswith('.html'):
            path = os.path.join(root, file).replace('./', '')
            html_files.append(path)

missing = []
for h in html_files:
    if h not in nav_content:
        missing.append(h)

for m in missing:
    print(m)
