/**
 * Redirect from GitHub Pages static site to Fly.io live deployment
 *
 * This script redirects users accessing pages on the GitHub Pages static mirror
 * to the same page on the Fly.io live deployment. This ensures traffic flows to
 * the production version running the Go backend, while preserving paths.
 */

(function() {
  const staticHost = 'rifaterdemsahin.github.io';
  const liveHost = 'claude-architect-certification.fly.dev';
  const repoPath = '/claude-architect-certification';

  // Only redirect if accessed from GitHub Pages static mirror
  if (location.hostname === staticHost) {
    // Construct the live site URL preserving path and search
    // Note: GitHub Pages includes the repo name in the path, Fly.io does not.
    let newPath = location.pathname;
    if (newPath.startsWith(repoPath)) {
      newPath = newPath.substring(repoPath.length);
    }
    const liveSiteUrl = `https://${liveHost}${newPath}${location.search}`;
    location.replace(liveSiteUrl);
  }
})();

