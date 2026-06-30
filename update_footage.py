import re
path = '5_Symbols/production/prod/footage_mapping.html'
with open(path, 'r') as f:
    content = f.read()

content = content.replace(
    '<button class="btn-delete"', 
    '<button data-require-admin="true" class="btn-delete btn-del"'
)

fn = 'removeMapping'
pattern = r'(function\s+' + fn + r'\s*\([^)]*\)\s*\{)(?!\s*if\s*\(\s*window\.requireAdmin)'
repl = r"\1\n  if (window.requireAdmin && !window.requireAdmin('" + fn + r"')) return;"
content = re.sub(pattern, repl, content, count=1)
    
with open(path, 'w') as f:
    f.write(content)
