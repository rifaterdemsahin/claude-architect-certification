# 🐛 Error: VS Code Go Extension Cannot Find `go` Binary

## 📋 Summary

VS Code (via the `golang.go` extension) cannot locate the `go` binary even though
Go was installed with Homebrew and works correctly in the terminal.

---

## ❌ Error Message

```
Failed to find the "go" binary in either GOROOT() or PATH(
  /Users/rifaterdemsahin/.local/bin:
  /Users/rifaterdemsahin/.grok/bin:
  ...
  /opt/homebrew/bin:   ← Go IS here
  ...
).
Check PATH, or Install Go and reload the window.
If PATH isn't what you expected, see https://github.com/golang/vscode-go/issues/971
```

---

## 🔍 Root Cause

Homebrew on Apple Silicon installs binaries to `/opt/homebrew/bin`.
VS Code launches its extension host as a **GUI app**, not a login shell, so it reads
`launchctl` environment — not `~/.zshrc` or `~/.zprofile`.
The `$PATH` that VS Code sees is set by `launchctl` at login time and does **not**
include `/opt/homebrew/bin` unless it was explicitly added to the macOS GUI PATH.

The terminal works because `~/.zshrc` prepends `/opt/homebrew/bin` at every interactive
shell start — but VS Code never sources that file.

---

## ✅ Fix — Three Options (in order of preference)

### 🩹 Option 1 — Tell the extension where `go` lives (fastest, no restart)

Open VS Code **Settings** → search `Go: Go Root` / `go.goroot`, or set `go.alternateTools`:

```jsonc
// .vscode/settings.json  (project-level, committed)
{
  "go.goroot": "/opt/homebrew/opt/go/libexec",
  "go.alternateTools": {
    "go": "/opt/homebrew/bin/go"
  }
}
```

Save the file, then run **Go: Restart Language Server** from the Command Palette.

### 🩹 Option 2 — Add `/opt/homebrew/bin` to the GUI PATH (permanent, affects all apps)

```bash
sudo launchctl config user path \
  "/opt/homebrew/bin:/opt/homebrew/sbin:$(launchctl getenv PATH)"
```

Log out and back in. VS Code will now see the Homebrew prefix in its PATH.

### 🩹 Option 3 — Create a symlink in a path VS Code already sees

```bash
sudo ln -sf /opt/homebrew/bin/go /usr/local/bin/go
sudo ln -sf /opt/homebrew/bin/gofmt /usr/local/bin/gofmt
```

`/usr/local/bin` is in the default macOS GUI PATH so VS Code finds it immediately.
Reload VS Code window (`Cmd+Shift+P` → **Reload Window**).

---

## 🔬 Verification

After applying any fix:

1. Open VS Code Command Palette → **Go: Locate Configured Go Tools**
2. All tools should resolve without error.
3. Open `cmd/server/main.go` — hover over a type; IntelliSense should work.
4. Run **Terminal** in VS Code:

```bash
go version   # should print go1.26.x darwin/arm64
```

---

## 🗓 Timeline

| Date | Event |
|------|-------|
| 2026-06-12 | Homebrew installed; `go` works in terminal but VS Code extension fails |
| 2026-06-12 | Root cause identified: GUI PATH vs shell PATH divergence on macOS |
| 2026-06-12 | Fix applied: see above option used |

---

## 📚 References

- [golang/vscode-go#971](https://github.com/golang/vscode-go/issues/971) — canonical issue tracking this macOS PATH problem
- [Homebrew on Apple Silicon PATH docs](https://docs.brew.sh/Installation#installation)
