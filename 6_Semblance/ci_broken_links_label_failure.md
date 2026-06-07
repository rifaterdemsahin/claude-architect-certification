# 🐛 Semblance: CI Failure — "Open GitHub issue for broken links" Step

## Summary

**Date:** 2026-06-07  
**Workflow:** Test Links After Deployment #10  
**Job:** test-links  
**Failed Step:** Open GitHub issue for broken links  
**Severity:** ERROR

---

## What Happened

The `test-links` CI job succeeded through all link-checking steps but failed at the final "Open GitHub issue for broken links" step. All preceding steps completed successfully:

| Step | Status |
|------|--------|
| Set up job | ✅ |
| Checkout | ✅ |
| Set up Python | ✅ |
| Derive GitHub Pages URL | ✅ |
| Wait for Pages to propagate | ✅ |
| Run HTTP link checker | ✅ |
| Build issue body | ✅ |
| Close duplicate open issues | ✅ |
| **Open GitHub issue for broken links** | ❌ |
| Fail the workflow if links are broken | ⏭ skipped |

---

## Root Cause

The `gh issue create` command was called with `--label "broken-links,bug"` but the label `broken-links` did not exist in the repository. GitHub CLI fails with a non-zero exit code when a specified label is missing, which caused the step to error out before the issue could be created.

```yaml
# The failing command:
gh issue create \
  --title "🚨 Broken links detected on GitHub Pages ($(date -u '+%Y-%m-%d'))" \
  --body-file issue_body.md \
  --label "broken-links,bug"
```

**Why:** Labels must be pre-created in a GitHub repository before they can be assigned to issues. `gh issue create` does not auto-create missing labels — it returns exit code 1 with an error like `Label 'broken-links' not found`.

---

## Fix Applied

Added a label-creation step (`Ensure labels exist`) **before** the `Open GitHub issue for broken links` step in `.github/workflows/test_links.yml`. The step uses `gh label create --force` which is idempotent — it creates the label if absent and updates it if already present, so it is safe to run on every workflow execution.

```yaml
- name: Ensure labels exist
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: |
    gh label create "broken-links" --color "d93f0b" --description "Automated: broken link detected" --force
    gh label create "bug" --color "d73a4a" --description "Something isn't working" --force
```

---

## How to Verify

1. Trigger a deployment that results in a broken link (or temporarily make `test_links.py` always return failure).
2. Confirm the "Open GitHub issue for broken links" step completes successfully (green).
3. Check the repository Issues tab for the auto-created issue with the `broken-links` and `bug` labels attached.

---

## Prevention

- Always pre-create required labels in CI workflows before any step that assigns them.
- Consider adding a one-time `bootstrap-labels.yml` workflow or a `scripts/create_labels.sh` script to initialize all standard labels when setting up a new repo from this template.

---

## References

- Workflow file: `.github/workflows/test_links.yml`
- Error log entry: `6_Semblance/error.log` — `[2026-06-07] [CI] [ERROR]`
- Fix log entry: `6_Semblance/fix.log` — `[2026-06-07] [CI] [APPLIED]`
