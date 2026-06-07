import os
import re
import urllib.parse
import sys

root_dir = "/Users/rifaterdemsahin/Projects/claude-architect-certification"
production_dir = os.path.join(root_dir, "5_Symbols/production")

print(f"🔍 Starting link checker in production folder: {production_dir}")

broken_links = []
scanned_files_count = 0
total_links_checked = 0

# Link regex patterns
href_pattern = re.compile(r'href=["\']([^"\']+)["\']', re.IGNORECASE)
src_pattern = re.compile(r'src=["\']([^"\']+)["\']', re.IGNORECASE)

def normalize_absolute_path(url_path):
    """
    Normalizes local file:/// URLs (both Mac and Windows formats) to local absolute paths.
    """
    if url_path.startswith("file:///"):
        path = url_path[8:]  # strip file:///
        
        # Handle Windows drive letters, e.g., C:/projects -> map to local equivalent if testing
        if len(path) > 2 and (path[1:3] == ":/" or path[1:3] == ":\\"):
            # If it's a C:/projects link, map it to our mac root dir for validation!
            win_prefix = "C:/projects/claude-architect-certification"
            win_prefix_alt = "C:/projects/delivery-pilot-template"
            normalized_path = path.replace("\\", "/")
            if normalized_path.lower().startswith(win_prefix.lower()):
                suffix = normalized_path[len(win_prefix):]
                return os.path.abspath(root_dir + suffix)
            elif normalized_path.lower().startswith(win_prefix_alt.lower()):
                suffix = normalized_path[len(win_prefix_alt):]
                return os.path.abspath(root_dir + suffix)
            return None  # Can't easily verify other arbitrary Windows drives
        
        # It's a standard unix path
        # Unquote URL-encoded characters (like %20)
        return urllib.parse.unquote(path)
    return None

for root, dirs, files in os.walk(production_dir):
    for file in files:
        if file.endswith(".html"):
            filepath = os.path.join(root, file)
            scanned_files_count += 1
            
            with open(filepath, "r", encoding="utf-8", errors="ignore") as f:
                content = f.read()
            
            # Find all href and src targets
            links = href_pattern.findall(content) + src_pattern.findall(content)
            
            for link in links:
                total_links_checked += 1
                clean_link = link.split('#')[0]

                if clean_link.startswith("file:///"):
                    normalized_abs = normalize_absolute_path(clean_link)
                    if normalized_abs:
                        if not os.path.exists(normalized_abs):
                            broken_links.append({
                                "file": filepath,
                                "link": link,
                                "resolved_path": normalized_abs,
                                "type": "Absolute Link"
                            })
                    else:
                        # Cannot resolve this absolute link (e.g. arbitrary Windows drive)
                        # We won't mark it broken unless it's supposed to be in our repo
                        pass
                else:
                    # It's a relative link
                    # Compute relative target relative to the directory containing the current HTML file
                    target_path = os.path.normpath(os.path.join(root, clean_link))
                    if not os.path.exists(target_path):
                        broken_links.append({
                            "file": filepath,
                            "link": link,
                            "resolved_path": target_path,
                            "type": "Relative Link"
                        })

print(f"\n📊 Scan Results:")
print(f"  Files Scanned: {scanned_files_count}")
print(f"  Links Checked: {total_links_checked}")
print(f"  Broken Links: {len(broken_links)}")

if broken_links:
    print("\n🚨 Broken Links Found:")
    for item in broken_links:
        print(f"  - File: {item['file']}")
        print(f"    Link: {item['link']}")
        print(f"    Resolved Path: {item['resolved_path']}")
        print(f"    Type: {item['type']}")
    sys.exit(1)
else:
    print("\n🎉 All links are fully intact!")
    sys.exit(0)
