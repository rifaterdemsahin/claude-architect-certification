path = '5_Symbols/production/preprod/research/domain_specific_language.html'
with open(path, 'r') as f:
    content = f.read()

content = content.replace(
    '<button class="btn btn-primary" onclick="saveEntry()">', 
    '<button data-require-admin="true" class="btn btn-primary btn-del" onclick="saveEntry()">',
    1
)
with open(path, 'w') as f:
    f.write(content)
