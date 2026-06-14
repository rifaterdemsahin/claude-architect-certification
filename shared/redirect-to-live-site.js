/**
 * Redirect from Fly.io staging to GitHub Pages live site
 *
 * This script redirects users accessing pages on the Fly.io staging deployment
 * to the same page on the GitHub Pages live site. This ensures traffic flows to
 * the canonical, production version while preserving all path and query parameters.
 *
 * Usage: Add this line in the <head> of any HTML file:
 *   <script src="../../shared/redirect-to-live-site.js"></script>
 *
 * Or for pages in different directory depths, adjust the relative path accordingly.
 */

(function() {
  const stagingHost = 'claude-architect-certification.fly.dev';
  const liveHost = 'rifaterdemsahin.github.io';
  const repoPath = '/claude-architect-certification';

  // Only redirect if accessed from Fly.io staging
  if (location.hostname === stagingHost) {
    // Construct the live site URL preserving path and search
    const liveSiteUrl = `https://${liveHost}${repoPath}${location.pathname}${location.search}`;
    location.replace(liveSiteUrl);
  }
})();
