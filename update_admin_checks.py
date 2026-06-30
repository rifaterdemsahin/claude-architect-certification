import os
import re

files_to_check = [
    '5_Symbols/production/postprod/edit_list.html',
    '5_Symbols/production/postprod/lower_thirds.html',
    '5_Symbols/production/preprod/bulk_image_generator.html',
    '5_Symbols/production/preprod/ivq.html',
    '5_Symbols/production/preprod/producer_checklist.html',
    '5_Symbols/production/preprod/scripts/index.html',
    '5_Symbols/supabase/ui/admin.html'
]

def update_file(path):
    if not os.path.exists(path):
        return
    with open(path, 'r') as f:
        content = f.read()
    
    # 1. Add btn-del to btn-danger (except close preview buttons)
    # 2. Add data-require-admin="true" to btn-danger
    def replace_btn(match):
        full_btn = match.group(0)
        if 'closePreview' in full_btn:
            return full_btn
        if 'data-require-admin' not in full_btn:
            full_btn = full_btn.replace('<button', '<button data-require-admin="true"')
        if 'btn-del' not in full_btn:
            full_btn = full_btn.replace('btn-danger', 'btn-danger btn-del')
        return full_btn
        
    content = re.sub(r'<button[^>]*class="[^"]*btn-danger[^"]*"[^>]*>', replace_btn, content)
    
    # Also find JS functions starting with delete, remove, unlink, clear
    funcs_to_protect = [
        'deleteVideo', 'deleteScene', 'removeSentence', 'deleteIvq', 
        'removeTask', 'deleteCodeRefUI', 'deleteSentLink', 'unlinkResearchUI',
        'clearAll'
    ]
    
    for fn in funcs_to_protect:
        # e.g. function deleteVideo(id) {
        pattern = r'(function\s+' + fn + r'\s*\([^)]*\)\s*\{)(?!\s*if\s*\(\s*window\.requireAdmin)'
        # Add the check
        repl = r"\1\n  if (window.requireAdmin && !window.requireAdmin('" + fn.replace('UI','') + r"')) return;"
        content = re.sub(pattern, repl, content, count=1) # only function definition
        
        # also async function
        pattern_async = r'(async\s+function\s+' + fn + r'\s*\([^)]*\)\s*\{)(?!\s*if\s*\(\s*window\.requireAdmin)'
        content = re.sub(pattern_async, repl, content, count=1)
        
    with open(path, 'w') as f:
        f.write(content)

for f in files_to_check:
    update_file(f)

print("Done")
