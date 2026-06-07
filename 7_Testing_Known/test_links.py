#!/usr/bin/env python3
import os
import re
import urllib.parse
import sys
import argparse
import json
import time
from urllib.request import urlopen, Request
from urllib.error import URLError, HTTPError
from html.parser import HTMLParser
from collections import deque

# ── CLI args ─────────────────────────────────────────────────────────────────
parser = argparse.ArgumentParser(description="Link checker — local or HTTP mode")
parser.add_argument("--mode", choices=["local", "http"], default="local",
                    help="local: scan HTML files on disk; http: crawl a live site")
parser.add_argument("--base-url", default="",
                    help="Base URL for HTTP mode (e.g. https://user.github.io/repo/)")
parser.add_argument("--output-json", default="",
                    help="Write broken-link results to this JSON file")
args = parser.parse_args()

broken_links = []

# ── LOCAL MODE ────────────────────────────────────────────────────────────────
if args.mode == "local":
    root_dir = os.environ.get("REPO_ROOT",
               os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
    production_dir = os.path.join(root_dir, "5_Symbols/production")

    print(f"🔍 Starting link checker in production folder: {production_dir}")

    scanned_files_count = 0
    total_links_checked = 0

    href_pattern = re.compile(r'href=["\']([^"\']+)["\']', re.IGNORECASE)
    src_pattern = re.compile(r'src=["\']([^"\']+)["\']', re.IGNORECASE)

    def normalize_absolute_path(url_path):
        if url_path.startswith("file:///"):
            path = url_path[8:]
            if len(path) > 2 and (path[1:3] == ":/" or path[1:3] == ":\\"):
                win_prefix = "C:/projects/claude-architect-certification"
                win_prefix_alt = "C:/projects/delivery-pilot-template"
                normalized_path = path.replace("\\", "/")
                if normalized_path.lower().startswith(win_prefix.lower()):
                    suffix = normalized_path[len(win_prefix):]
                    return os.path.abspath(root_dir + suffix)
                elif normalized_path.lower().startswith(win_prefix_alt.lower()):
                    suffix = normalized_path[len(win_prefix_alt):]
                    return os.path.abspath(root_dir + suffix)
                return None
            return urllib.parse.unquote(path)
        return None

    for root, dirs, files in os.walk(production_dir):
        for file in files:
            if file.endswith(".html"):
                filepath = os.path.join(root, file)
                scanned_files_count += 1

                with open(filepath, "r", encoding="utf-8", errors="ignore") as f:
                    content = f.read()

                links = href_pattern.findall(content) + src_pattern.findall(content)

                for link in links:
                    clean_link = link.split("?")[0].split("#")[0].strip()
                    if not clean_link:
                        continue
                    if any(clean_link.lower().startswith(p) for p in
                           ["http://", "https://", "mailto:", "tel:", "javascript:"]) \
                            or "${" in clean_link:
                        continue

                    total_links_checked += 1

                    if clean_link.startswith("file:///"):
                        normalized_abs = normalize_absolute_path(clean_link)
                        if normalized_abs and not os.path.exists(normalized_abs):
                            broken_links.append({
                                "file": filepath,
                                "link": link,
                                "resolved_path": normalized_abs,
                                "type": "Absolute Link",
                            })
                    else:
                        target_path = os.path.normpath(os.path.join(root, clean_link))
                        if not os.path.exists(target_path):
                            broken_links.append({
                                "file": filepath,
                                "link": link,
                                "resolved_path": target_path,
                                "type": "Relative Link",
                            })

    print(f"\n📊 Scan Results:")
    print(f"  Files Scanned:  {scanned_files_count}")
    print(f"  Links Checked:  {total_links_checked}")
    print(f"  Broken Links:   {len(broken_links)}")

# ── HTTP MODE ─────────────────────────────────────────────────────────────────
elif args.mode == "http":
    base_url = args.base_url.rstrip("/") + "/"
    if not base_url.startswith("http"):
        print("❌ --base-url must start with http:// or https://")
        sys.exit(2)

    parsed_base = urllib.parse.urlparse(base_url)
    origin = f"{parsed_base.scheme}://{parsed_base.netloc}"

    print(f"🌐 Starting HTTP link checker — base URL: {base_url}")

    class LinkExtractor(HTMLParser):
        def __init__(self):
            super().__init__()
            self.links = []

        def handle_starttag(self, tag, attrs):
            attrs_dict = dict(attrs)
            if tag == "a":
                href = attrs_dict.get("href", "")
                if href:
                    self.links.append(href)
            elif tag in ("img", "script"):
                src = attrs_dict.get("src", "")
                if src:
                    self.links.append(src)
            elif tag == "link":
                href = attrs_dict.get("href", "")
                if href:
                    self.links.append(href)

    def fetch_url(url, timeout=15):
        req = Request(url, headers={"User-Agent": "LinkChecker/1.0"})
        try:
            with urlopen(req, timeout=timeout) as resp:
                content_type = resp.headers.get("Content-Type", "")
                body = resp.read().decode("utf-8", errors="ignore") \
                    if "text/html" in content_type else ""
                return resp.status, body
        except HTTPError as e:
            return e.code, ""
        except URLError as e:
            return 0, ""

    def resolve(href, page_url):
        href = href.strip()
        if not href or href.startswith(("mailto:", "tel:", "javascript:", "#")):
            return None
        resolved = urllib.parse.urljoin(page_url, href.split("#")[0])
        parsed = urllib.parse.urlparse(resolved)
        # Only check links within the same origin
        if parsed.netloc and parsed.netloc != parsed_base.netloc:
            return None
        return resolved

    visited = set()
    queue = deque([base_url])
    pages_crawled = 0
    total_links_checked = 0

    while queue:
        url = queue.popleft()
        if url in visited:
            continue
        visited.add(url)

        status, body = fetch_url(url)
        pages_crawled += 1
        total_links_checked += 1
        time.sleep(0.1)  # be polite to the server

        if status not in range(200, 300):
            broken_links.append({
                "url": url,
                "status": status,
                "type": "HTTP Error",
                "found_on": url,
            })
            print(f"  ❌ {status}  {url}")
            continue

        # Parse links from HTML pages and queue internal ones
        if body:
            extractor = LinkExtractor()
            extractor.feed(body)
            for href in extractor.links:
                resolved = resolve(href, url)
                if resolved and resolved not in visited:
                    # Track where we found external resource links to report breakage
                    if resolved not in queue:
                        queue.append(resolved)

    print(f"\n📊 HTTP Scan Results:")
    print(f"  Pages Crawled:  {pages_crawled}")
    print(f"  URLs Checked:   {total_links_checked}")
    print(f"  Broken Links:   {len(broken_links)}")

# ── Results ───────────────────────────────────────────────────────────────────
if broken_links:
    print("\n🚨 Broken Links Found:")
    for item in broken_links:
        if args.mode == "local":
            print(f"  - File: {item['file']}")
            print(f"    Link: {item['link']}")
            print(f"    Resolved Path: {item['resolved_path']}")
            print(f"    Type: {item['type']}")
        else:
            print(f"  - URL:    {item['url']}")
            print(f"    Status: {item['status']}")

    if args.output_json:
        with open(args.output_json, "w") as f:
            json.dump(broken_links, f, indent=2)
        print(f"\n📄 Results written to {args.output_json}")

    sys.exit(1)
else:
    print("\n✅ All links are fully intact!")
    if args.output_json:
        with open(args.output_json, "w") as f:
            json.dump([], f)
    sys.exit(0)
