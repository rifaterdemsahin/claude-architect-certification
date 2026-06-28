import json
import copy

with open('navigation_config.json', 'r') as f:
    nav = json.load(f)

def find_menu(menu_array, label_substring):
    for item in menu_array:
        if label_substring in item.get('label', ''):
            return item
        if 'children' in item:
            res = find_menu(item['children'], label_substring)
            if res: return res
    return None

def add_to_menu(parent_label, new_item):
    parent = find_menu(nav['projectMenu'], parent_label)
    if parent and 'children' in parent:
        # Check if already exists
        exists = any(child.get('url') == new_item['url'] for child in parent['children'])
        if not exists:
            parent['children'].append(new_item)

# 1. Customer Discovery Interviews
add_to_menu("1. 🧭 Strategy & Research", {
    "label": "🎤 Customer Discovery Interviews",
    "url": "5_Symbols/production/preprod/customer_discovery_interviews.html",
    "description": "User interviews for customer discovery."
})

# 2. LinkedIn Controversial
add_to_menu("11. 🤝 LinkedIn Outreach", {
    "label": "Controversial Post Playbook",
    "url": "5_Symbols/production/postprod/linkedin_controversial.html"
})

# 3. Visual Gallery
add_to_menu("2. 🎨 Visuals & Graphics", {
    "label": "🖼️ Visual Asset Gallery",
    "url": "5_Symbols/production/postprod/visual_gallery.html"
})

# 4. Explanations
add_to_menu("2. 🎨 Method & Doctrine", {
    "label": "🧠 AI Architect Explanations",
    "url": "5_Symbols/production/preprod/explanations.html"
})

# 5. Environment
add_to_menu("🗄️ Data & Backend", {
    "label": "🌍 Loaded Environment",
    "url": "5_Symbols/production/preprod/environment.html"
})

# 6. Stats
add_to_menu("🗄️ Data & Backend", {
    "label": "📊 Database Stats",
    "url": "5_Symbols/production/preprod/stats.html"
})

# 7. Strategy
add_to_menu("1. 🧭 Strategy & Research", {
    "label": "♟️ Certification Strategy",
    "url": "5_Symbols/production/preprod/strategy.html"
})

# 8. Bulk Image Generator
add_to_menu("🖼 Images", {
    "label": "🖼️ Bulk Image Generator",
    "url": "5_Symbols/production/preprod/bulk_image_generator.html"
})

# 9. Product Market Fit
add_to_menu("1. 🧭 Strategy & Research", {
    "label": "🎯 Product-Market Fit",
    "url": "5_Symbols/production/preprod/product_market_fit.html"
})

# 10. Edit Scripts
add_to_menu("5. 🎬 Script", {
    "label": "✍️ Script Editor",
    "url": "5_Symbols/production/preprod/edit_scripts.html"
})

# 11. Asset Checklist
add_to_menu("3. 📋 Shots & Tracking", {
    "label": "✅ Asset Checklist",
    "url": "5_Symbols/production/postprod/asset_checklist.html"
})

# Hubs
add_to_menu("🎬 Preprod", {
    "label": "🗂️ Pre-Production Hub",
    "url": "5_Symbols/production/preprod/index.html"
})
add_to_menu("🎥 Production", {
    "label": "🗂️ Production Hub",
    "url": "5_Symbols/production/prod/index.html"
})
add_to_menu("📦 Post Prod", {
    "label": "🗂️ Post-Production Hub",
    "url": "5_Symbols/production/postprod/index.html"
})
add_to_menu("🛠️ Tools", {
    "label": "⚙️ Settings",
    "url": "5_Symbols/production/settings.html"
})

with open('navigation_config.json', 'w') as f:
    json.dump(nav, f, indent=2)
print("Updated navigation_config.json")
