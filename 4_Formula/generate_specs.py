import os
import glob
from bs4 import BeautifulSoup
import re

# --------------------------------------------------------------------------- #
# Versioning convention (see 4_Formula/specs/README.md):
#   0.1   = initial spec
#   0.11  = small update (typo, single field, minor clarification)
#   0.2   = bigger update (new functionality, new tables, structural change)
# The generator NEVER auto-bumps. It PRESERVES whatever version already lives in
# the spec file (so hand edits + manual bumps survive a re-generation) and only
# stamps 0.1 on brand-new specs.
# --------------------------------------------------------------------------- #
DEFAULT_VERSION = "0.1"

VERSION_RE = re.compile(r"\*\*Version\*\*:\s*`?([0-9]+\.[0-9]+)`?")


def read_existing_version(spec_path):
    """Preserve a manually-set version across re-generation."""
    if not os.path.exists(spec_path):
        return DEFAULT_VERSION
    try:
        with open(spec_path, "r", encoding="utf-8") as f:
            head = f.read(2000)
        m = VERSION_RE.search(head)
        if m:
            return m.group(1)
    except Exception:
        pass
    return DEFAULT_VERSION


def extract_tables_and_apis(inline_blob):
    """Pull the Supabase/PostgREST tables and REST endpoints out of the JS so the
    spec can be used to recreate the page's data layer."""
    # PostgREST-style helpers: api('scenes?...'), api(`videos?...`), .from('x')
    table_hits = re.findall(r"""api\(\s*[`'"]([a-z_][a-z0-9_]*)""", inline_blob)
    table_hits += re.findall(r"""\.from\(\s*['"]([a-z_][a-z0-9_]*)""", inline_blob)
    table_hits += re.findall(r"""rest/v1/([a-z_][a-z0-9_]*)""", inline_blob)
    tables = sorted(set(table_hits))

    # fetch('/api/...') style backend endpoints (exclude the PostgREST table proxy noise)
    ep_hits = re.findall(r"""fetch\(\s*[`'"]([^`'")\s]+)""", inline_blob)
    endpoints = sorted(set(e.split("?")[0] for e in ep_hits if e.startswith("/") or e.startswith("http")))
    return tables, endpoints


def create_spec(html_path, existing_version):
    with open(html_path, 'r', encoding='utf-8') as f:
        content = f.read()

    soup = BeautifulSoup(content, 'html.parser')

    title_tag = soup.find('title')
    title = title_tag.text.strip() if title_tag else os.path.basename(html_path)

    meta_desc = soup.find('meta', attrs={'name': 'description'})
    description = meta_desc['content'].strip() if meta_desc and meta_desc.has_attr('content') else 'No description provided.'

    # Extract scripts
    scripts = soup.find_all('script')
    script_srcs = [s['src'] for s in scripts if s.has_attr('src')]
    inline_scripts = [s.string for s in scripts if not s.has_attr('src') and s.string]
    inline_blob = " ".join(inline_scripts)

    # Extract stylesheets
    links = soup.find_all('link', rel='stylesheet')
    stylesheets = [l['href'] for l in links if l.has_attr('href')]

    # Structure (Headings)
    headings = []
    for h in soup.find_all(['h1', 'h2', 'h3']):
        headings.append(f"{h.name.upper()}: {h.text.strip()}")

    # Main content IDs/Classes to describe layout
    main_containers = []
    for tag in soup.find_all(['div', 'main', 'section', 'nav', 'header', 'footer']):
        tag_id = tag.get('id')
        tag_class = tag.get('class')
        if tag_id:
            main_containers.append(f"<{tag.name} id='{tag_id}'>")
        elif tag_class and len(tag_class) > 0:
            if tag.name in ['main', 'header', 'footer', 'nav'] or 'container' in tag_class or 'wrapper' in tag_class:
                main_containers.append(f"<{tag.name} class='{' '.join(tag_class)}'>")

    # Dedup
    main_containers = list(dict.fromkeys(main_containers))

    # Inline logic
    funcs = sorted(set(re.findall(r'function\s+(\w+)', inline_blob)))
    consts = sorted(set(re.findall(r'const\s+(\w+)', inline_blob)))

    # Data layer
    tables, endpoints = extract_tables_and_apis(inline_blob)

    # ----------------------------------------------------------------- build md
    spec_content = f"""# Spec: {title}

> 🔖 **Version**: `{existing_version}`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`{html_path}`

## 🎯 Purpose & Rationale
**Description**: {description}

*Rationale*: This file exists to serve as the `{title}` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `{os.path.dirname(html_path)}` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
"""
    if funcs:
        for fn in funcs:
            spec_content += f"- `{fn}()`\n"
    else:
        spec_content += "- Static page — no interactive JavaScript functions detected.\n"

    spec_content += "\n## 🗄️ Data Layer — Tables & APIs Used\n"
    spec_content += "**Database tables (Supabase / PostgREST):**\n"
    if tables:
        for t in tables:
            spec_content += f"- `{t}`\n"
    else:
        spec_content += "- None — this page does not read or write database tables.\n"
    spec_content += "\n**Backend / external endpoints:**\n"
    if endpoints:
        for e in endpoints:
            spec_content += f"- `{e}`\n"
    else:
        spec_content += "- None.\n"

    spec_content += """
## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
"""
    for css in stylesheets:
        spec_content += f"- `{css}`\n"
    if not stylesheets:
        spec_content += "- None external.\n"

    spec_content += "\n### 3. 🖥️ UI — Core Layout & Containers\nThe page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:\n"
    for c in main_containers[:15]:  # Limit to avoid huge lists
        spec_content += f"- `{c}`\n"
    if not main_containers:
        spec_content += "- No specific structural IDs found.\n"

    spec_content += "\n### 4. Key Headings\n"
    for h in headings[:10]:
        spec_content += f"- {h}\n"

    spec_content += "\n### 5. Scripts Required\nThe following JavaScript files must be loaded which are reusable shared scripts.\n"
    for src in script_srcs:
        spec_content += f"- `{src}`\n"
    if not script_srcs:
        spec_content += "- No external scripts.\n"

    if inline_scripts:
        spec_content += "\n**Inline Script Logic Includes:**\n"
        if funcs:
            spec_content += f"- Functions: {', '.join(funcs[:12])}\n"
        if consts:
            spec_content += f"- Constants/Variables: {', '.join(consts[:10])}\n"
        if not funcs and not consts:
            spec_content += "- Miscellaneous inline logic.\n"

    spec_content += """
## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
"""

    return spec_content


def main():
    specs_dir = "4_Formula/specs"
    os.makedirs(specs_dir, exist_ok=True)

    html_files = []
    for root, dirs, files in os.walk("."):
        if "node_modules" in root or ".git" in root or "_obsolete" in root:
            continue
        for file in files:
            if file.endswith(".html"):
                html_files.append(os.path.join(root, file))

    print(f"Found {len(html_files)} HTML files. Generating specs...")

    for html_file in html_files:
        try:
            # Normalize filename
            rel_path = os.path.relpath(html_file, ".")
            safe_name = rel_path.replace("/", "_").replace("\\", "_").replace(".html", "_spec.md")
            spec_path = os.path.join(specs_dir, safe_name)

            existing_version = read_existing_version(spec_path)
            spec_content = create_spec(html_file, existing_version)

            with open(spec_path, 'w', encoding='utf-8') as f:
                f.write(spec_content)
        except Exception as e:
            print(f"Failed to process {html_file}: {e}")

    print("Done!")


if __name__ == "__main__":
    main()
