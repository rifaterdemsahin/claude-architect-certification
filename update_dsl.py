import re
import os

path = '5_Symbols/production/preprod/research/domain_specific_language.html'
with open(path, 'r') as f:
    content = f.read()

# Protect deleteEntry button
content = content.replace(
    '<button class="btn-sm btn-delete"', 
    '<button data-require-admin="true" class="btn-sm btn-delete btn-del"'
)

# Protect JS functions
funcs = ['saveEntry', 'deleteEntry']
for fn in funcs:
    pattern_async = r'(async\s+function\s+' + fn + r'\s*\([^)]*\)\s*\{)(?!\s*if\s*\(\s*window\.requireAdmin)'
    repl = r"\1\n  if (window.requireAdmin && !window.requireAdmin('" + fn + r"')) return;"
    content = re.sub(pattern_async, repl, content, count=1)
    
with open(path, 'w') as f:
    f.write(content)
with open(path, 'r') as f:
    content = f.read()

content = content.replace(
    '<button class="btn-primary" onclick="saveEntry()">', 
    '<button data-require-admin="true" class="btn-primary" onclick="saveEntry()">'
)
with open(path, 'w') as f:
    f.write(content)
