(function () {
  var BASE_URL = 'https://rifaterdemsahin.github.io/claude-architect-certification';
  var path = location.pathname.replace(/\/claude-architect-certification/, '').split('?')[0];
  if (path.endsWith('/')) path += 'index.html';

  var pageMeta = {
    '/index.html': {
      title: 'Claude AI Certification for Architects — Enterprise MCP & Multi-Agent Systems',
      description: 'Enterprise-grade Claude AI certification course covering MCP servers, prompt caching, multi-agent topologies, VPC PrivateLink, and system architecture for AI architects.',
      keywords: 'Claude AI, AI certification, MCP server, multi-agent systems, prompt caching, enterprise AI, system architecture, Anthropic, AI architect'
    },
    '/home.html': {
      title: 'Claude AI Certification for Architects — Home',
      description: 'Enterprise Claude AI Certification for Architects. Learn custom MCP, VPC PrivateLink, multi-agent topologies, and cost-reduction prompt caching.',
      keywords: 'Claude AI, AI certification, MCP server, multi-agent systems, prompt caching, enterprise AI'
    },
    '/markdown_renderer.html': {
      title: 'Document Renderer — Claude AI Certification',
      description: 'Render and view markdown documentation for the Claude AI Certification for Architects course. Browse architecture guides, topology docs, and reference materials.',
      keywords: 'markdown renderer, documentation viewer, Claude AI, course docs, architecture documentation'
    },
    '/5_Symbols/pipeline.html': {
      title: 'Production Pipeline — Claude AI Certification for Architects',
      description: 'End-to-end production pipeline for the Claude AI Certification course. Track progress across pre-production, production, post-production, and publishing stages.',
      keywords: 'production pipeline, course production, video pipeline, content workflow, Claude AI certification'
    },
    '/5_Symbols/timeline.html': {
      title: 'Project Timeline — Claude AI Certification for Architects',
      description: 'Course project timeline showing milestones, deadlines, and progress tracking for the Claude AI Certification for Architects masterclass.',
      keywords: 'project timeline, course milestones, certification timeline, AI course schedule'
    },
    '/5_Symbols/production_hub.html': {
      title: 'Production Hub — Claude AI Certification for Architects',
      description: 'Central production hub for the Claude AI Certification course. Manage pre-production, production, post-production workflows and track course creation progress.',
      keywords: 'production hub, course production, content management, AI certification workflow'
    },
    '/5_Symbols/course_outline.html': {
      title: '📚 Course Outline — Claude Architect Certification',
      description: 'Complete course outline for the Claude AI Certification for Architects. Browse modules, videos, and learning paths for enterprise AI system design.',
      keywords: 'course outline, Claude AI modules, certification curriculum, AI architect training'
    }
  };

  var pageTitles = {
    '/5_Symbols/production/preprod/index.html': 'Pre-Production Hub — Claude AI Certification',
    '/5_Symbols/production/preprod/planning.html': 'Planning Hub — Claude AI Certification',
    '/5_Symbols/production/preprod/sanity_checklist.html': 'Master Sanity Checklist — Claude AI Certification',
    '/5_Symbols/production/preprod/producer_checklist.html': 'Producer Checklist — Claude AI Certification',
    '/5_Symbols/production/preprod/edit_scripts.html': 'Script Editor — Claude AI Certification',
    '/5_Symbols/production/preprod/explanations.html': 'AI Architect Explanations — Claude AI Certification',
    '/5_Symbols/production/preprod/course_outline.html': 'Course Outline — Claude Architect Certification',
    '/5_Symbols/production/preprod/bulk_image_generator.html': 'Bulk Image Generator — Claude AI Certification',
    '/5_Symbols/production/preprod/log_viewer.html': 'Log Viewer — Claude AI Certification',
    '/5_Symbols/production/preprod/stats.html': 'Statistics — Claude AI Certification',
    '/5_Symbols/production/preprod/problem.html': 'Problem Statement — Claude Architect Certification',
    '/5_Symbols/production/preprod/ivq.html': 'Interactive Video Quiz — Claude AI Certification',
    '/5_Symbols/production/preprod/research/index.html': 'Research Hub — Claude AI Certification',
    '/5_Symbols/production/preprod/research/videos.html': 'Video Research — Claude AI Certification',
    '/5_Symbols/production/preprod/research/notes.html': 'Research Notes — Claude AI Certification',
    '/5_Symbols/production/preprod/research/images.html': 'Image Research — Claude AI Certification',
    '/5_Symbols/production/preprod/research/audio.html': 'Audio Research — Claude AI Certification',
    '/5_Symbols/production/preprod/research/domain_specific_language.html': 'Domain Specific Language — Claude AI Certification',
    '/5_Symbols/production/preprod/research/market_analysis.html': 'Market Analysis — Claude AI Certification',
    '/5_Symbols/production/preprod/research/infographic_generator.html': 'Infographic Generator — Claude AI Certification',
    '/5_Symbols/production/preprod/scripts/index.html': 'Scripts Hub — Claude AI Certification',
    '/5_Symbols/production/preprod/scripts/generator.html': 'Script Generator — Claude AI Certification',
    '/5_Symbols/production/preprod/tools/admin.html': 'Admin Tools — Claude AI Certification',
    '/5_Symbols/production/prod/index.html': 'Production — Claude AI Certification',
    '/5_Symbols/production/prod/checklist.html': 'Production Checklist — Claude AI Certification',
    '/5_Symbols/production/prod/footage_mapping.html': 'Footage Mapping — Claude AI Certification',
    '/5_Symbols/production/postprod/index.html': 'Post-Production — Claude AI Certification',
    '/5_Symbols/production/postprod/production_shotlist.html': 'Production Shot List — Claude AI Certification',
    '/5_Symbols/production/postprod/post_production_checklist.html': 'Post Production Checklist — Claude AI Certification',
    '/5_Symbols/production/postprod/edit_list.html': 'Edit List — Claude AI Certification',
    '/5_Symbols/production/postprod/audio_scoring.html': 'Audio Scoring — Claude AI Certification',
    '/5_Symbols/production/postprod/asset_checklist.html': 'Asset Checklist — Claude AI Certification',
    '/5_Symbols/production/postprod/visual_gallery.html': 'Visual Asset Gallery — Claude AI Certification',
    '/5_Symbols/production/postprod/lower_thirds.html': 'Lower Thirds Manager — Claude AI Certification',
    '/5_Symbols/production/postprod/linkedin_messaging.html': 'LinkedIn Messaging — Claude AI Certification',
    '/5_Symbols/production/postprod/linkedin_controversial.html': 'Controversial Post Playbook — Claude AI Certification',
    '/5_Symbols/production/postprod/image_generator.html': 'AI Image Generator — Claude AI Certification',
    '/5_Symbols/production/postprod/flywheel.html': 'The Flywheel — Claude AI Certification',
    '/5_Symbols/production/postprod/customer_discovery.html': 'Customer Discovery — Claude AI Certification',
    '/5_Symbols/production/publish/membership.html': 'Join Members — Claude AI Certification',
    '/5_Symbols/production/settings.html': 'Production Settings — Claude AI Certification',
    '/5_Symbols/supabase/ui/admin.html': 'Supabase Admin — Claude AI Certification',
    '/5_Symbols/templates/index.html': 'Claude AI Certification for Architects',
    '/5_Symbols/templates/axiom_errors.html': 'Axiom Errors Admin — Claude AI Certification',
    '/5_Symbols/course_src/shared-templates/markdown_renderer.html': 'Markdown Renderer — Claude AI Certification',
    '/5_Symbols/course_src/shared-utils/markdown_viewer.html': 'Markdown Document Viewer — Claude AI Certification',
    '/5_Symbols/tools/sitemap.html': 'Project Sitemap — Claude AI Certification'
  };

  var pageDescriptions = {
    '/5_Symbols/production/preprod/index.html': 'Pre-production hub for the Claude AI Certification course. Manage scripts, outlines, milestones, and all pre-production activities.',
    '/5_Symbols/production/preprod/planning.html': 'Planning hub for the Claude AI Certification. Track critical path, milestones, pipeline tasks, and production planning.',
    '/5_Symbols/production/preprod/sanity_checklist.html': 'Master sanity checklist for the Claude AI Certification course. Validate all production assets and pre-production tasks.',
    '/5_Symbols/production/preprod/producer_checklist.html': 'Producer checklist for the Claude AI Certification course. Track daily recording tasks, video production, and course creation progress.',
    '/5_Symbols/production/preprod/edit_scripts.html': 'Edit and manage course scripts for the Claude AI Certification for Architects masterclass.',
    '/5_Symbols/production/preprod/explanations.html': 'AI Architect explanations and technical deep-dives for the Claude AI Certification course concepts.',
    '/5_Symbols/production/preprod/course_outline.html': 'Browse the complete course outline for Claude AI Certification for Architects. View all modules, sections, and learning objectives.',
    '/5_Symbols/production/preprod/bulk_image_generator.html': 'Generate bulk images for course production using AI. Create visual assets for the Claude AI Certification course.',
    '/5_Symbols/production/preprod/log_viewer.html': 'View and analyze application logs for the Claude AI Certification project. Debug and monitor system activity.',
    '/5_Symbols/production/preprod/stats.html': 'View project statistics and analytics for the Claude AI Certification course production pipeline.',
    '/5_Symbols/production/preprod/problem.html': 'Define the core problem being solved by the Claude AI Certification for Architects course. Enterprise AI system architecture challenges.',
    '/5_Symbols/production/preprod/ivq.html': 'Interactive Video Quiz manager for the Claude AI Certification course. Create and manage video-based assessments.',
    '/5_Symbols/production/preprod/research/index.html': 'Central research hub for the Claude AI Certification course. Access images, audio, video, notes, and market analysis.',
    '/5_Symbols/production/preprod/research/videos.html': 'Video research assets and references for the Claude AI Certification for Architects course production.',
    '/5_Symbols/production/preprod/research/notes.html': 'Research notes and technical references for the Claude AI Certification course development.',
    '/5_Symbols/production/preprod/research/images.html': 'Image research assets, mockups, and visual references for the Claude AI Certification course.',
    '/5_Symbols/production/preprod/research/audio.html': 'Audio research assets, voiceovers, and sound effects for the Claude AI Certification course.',
    '/5_Symbols/production/preprod/research/domain_specific_language.html': 'Domain Specific Language (DSL) research and definitions for the Claude AI Certification course.',
    '/5_Symbols/production/preprod/research/market_analysis.html': 'Market analysis for AI certification and enterprise AI training. Research findings for the Claude AI Certification course.',
    '/5_Symbols/production/preprod/research/infographic_generator.html': 'Generate infographics for the Claude AI Certification course. Create visual data representations and diagrams.',
    '/5_Symbols/production/preprod/scripts/index.html': 'Scripts hub for the Claude AI Certification course. Create, edit, and manage video narration scripts.',
    '/5_Symbols/production/preprod/scripts/generator.html': 'AI-powered script generator for the Claude AI Certification course. Generate narration and video scripts.',
    '/5_Symbols/production/preprod/tools/admin.html': 'Admin tools and utilities for managing the Claude AI Certification course production pipeline.',
    '/5_Symbols/production/prod/index.html': 'Production phase hub for the Claude AI Certification course. Manage recordings, voiceovers, and footage.',
    '/5_Symbols/production/prod/checklist.html': 'Production checklist for roll and voiceover recordings in the Claude AI Certification course.',
    '/5_Symbols/production/prod/footage_mapping.html': 'Footage and research mapping for the Claude AI Certification course. Organize video assets and source material.',
    '/5_Symbols/production/postprod/index.html': 'Post-production hub for the Claude AI Certification course. Manage EDL, asset checklists, scene review, and final editing.',
    '/5_Symbols/production/postprod/production_shotlist.html': 'Production shot list and asset management for Module 1 of the Claude AI Certification course.',
    '/5_Symbols/production/postprod/post_production_checklist.html': 'Post-production checklist for the Claude AI Certification course. Track editing tasks and quality assurance.',
    '/5_Symbols/production/postprod/edit_list.html': 'Edit Decision List (EDL) manager for the Claude AI Certification course. Track scene edits and transitions.',
    '/5_Symbols/production/postprod/audio_scoring.html': 'Audio scoring, sound effects, and background music management for the Claude AI Certification course post-production.',
    '/5_Symbols/production/postprod/asset_checklist.html': 'Production asset checklist for Module 1, Section 1 of the Claude AI Certification course.',
    '/5_Symbols/production/postprod/visual_gallery.html': 'Visual asset gallery for post-production. Browse and manage images, diagrams, and visual elements.',
    '/5_Symbols/production/postprod/lower_thirds.html': 'Create, edit, and preview lower thirds graphics for the Claude AI Certification course videos.',
    '/5_Symbols/production/postprod/linkedin_messaging.html': 'LinkedIn messaging templates and outreach scripts for promoting the Claude AI Certification.',
    '/5_Symbols/production/postprod/linkedin_controversial.html': 'Controversial post playbook and LinkedIn engagement strategies for the Claude AI Certification.',
    '/5_Symbols/production/postprod/image_generator.html': 'AI image generator for creating visual assets and thumbnails for the Claude AI Certification course.',
    '/5_Symbols/production/postprod/flywheel.html': 'The Flywheel business model for the Claude AI Certification. Learn, build, certify, publish, and sell.',
    '/5_Symbols/production/postprod/customer_discovery.html': 'Customer discovery framework and workflow from pyramid to product-market fit for the certification.',
    '/5_Symbols/production/publish/membership.html': 'Join the Claude AI Certification membership. Unlock full access to the enterprise AI architect masterclass and community.',
    '/5_Symbols/production/settings.html': 'Production settings and configuration for the Claude AI Certification course management system.',
    '/5_Symbols/supabase/ui/admin.html': 'Supabase database admin interface for managing the Claude AI Certification course data and seed operations.',
    '/5_Symbols/templates/index.html': 'Claude AI Certification for Architects — Enterprise Systems & Integration Masterclass Companion Workspace.',
    '/5_Symbols/templates/axiom_errors.html': 'Axiom error monitoring and admin dashboard for the Claude AI Certification production system.',
    '/5_Symbols/course_src/shared-templates/markdown_renderer.html': 'Render and view markdown documentation within the course source. Browse technical references and guides.',
    '/5_Symbols/course_src/shared-utils/markdown_viewer.html': 'View course documentation, design topologies, and architecture diagrams for the Claude AI Certification.',
    '/5_Symbols/tools/sitemap.html': 'Complete sitemap of all project pages for the Claude AI Certification for Architects. Browse the full site structure.'
  };

  var pageKeywords = {
    'preprod': 'pre-production, course planning, script editing, content outline, AI certification production',
    'research': 'research assets, market analysis, competitive research, course development, AI training research',
    'scripts': 'video scripts, narration scripts, script editing, AI script generation, course content',
    'prod': 'video production, audio recording, footage management, course recording, content production',
    'postprod': 'post-production, video editing, EDL, audio scoring, asset management, course finishing',
    'publish': 'course publishing, membership, content delivery, certification launch, audience growth',
    'tools': 'production tools, admin utilities, sitemap, course management, AI certification tools'
  };

  function ensureTags() {
    if (document.querySelector('meta[name="seo-injected"]')) return;
    var existingTitle = document.title;

    var metaKey = path;
    var meta = pageMeta[metaKey] || {};

    var finalTitle = pageTitles[metaKey] || meta.title || existingTitle || 'Claude AI Certification for Architects';

    var finalDescription = meta.description || pageDescriptions[metaKey] || (
      'Claude AI Certification for Architects — ' + (finalTitle.replace(/—.*$/, '').trim() || 'Course Production')
    );

    var segment = path.includes('/preprod/') ? 'preprod'
      : path.includes('/research/') ? 'research'
      : path.includes('/scripts/') ? 'scripts'
      : path.includes('/prod/') ? 'prod'
      : path.includes('/postprod/') ? 'postprod'
      : path.includes('/publish/') ? 'publish'
      : path.includes('/tools/') ? 'tools'
      : '';
    var fallbackKeywords = 'Claude AI, AI certification, enterprise AI, system architecture, MCP, multi-agent';
    var finalKeywords = meta.keywords || (segment ? pageKeywords[segment] : '') || fallbackKeywords;

    var canonical = BASE_URL + path;

    var head = document.head || document.querySelector('head') || document.documentElement;
    function addMeta(name, value, property) {
      if (!value) return;
      var existing = property
        ? head.querySelector('meta[property="' + property + '"]')
        : head.querySelector('meta[name="' + name + '"]');
      if (existing) return;
      var el = document.createElement('meta');
      if (property) {
        el.setAttribute('property', property);
      } else {
        el.setAttribute('name', name);
      }
      el.setAttribute('content', value);
      head.appendChild(el);
    }

    function addLink(rel, href) {
      if (!href) return;
      var existing = head.querySelector('link[rel="' + rel + '"]');
      if (existing) return;
      var el = document.createElement('link');
      el.setAttribute('rel', rel);
      el.setAttribute('href', href);
      head.appendChild(el);
    }

    addMeta('description', finalDescription);
    addMeta('keywords', finalKeywords);
    addMeta('robots', 'index, follow');

    addMeta('og:title', finalTitle, 'og:title');
    addMeta('og:description', finalDescription, 'og:description');
    addMeta('og:type', 'website', 'og:type');
    addMeta('og:url', canonical, 'og:url');
    addMeta('og:image', BASE_URL + '/favicon.png', 'og:image');
    addMeta('og:site_name', 'Claude AI Certification for Architects', 'og:site_name');

    addMeta('twitter:card', 'summary_large_image', 'twitter:card');
    addMeta('twitter:title', finalTitle, 'twitter:title');
    addMeta('twitter:description', finalDescription, 'twitter:description');
    addMeta('twitter:image', BASE_URL + '/favicon.png', 'twitter:image');

    addLink('canonical', canonical);

    var existingLD = head.querySelector('script[type="application/ld+json"]');
    if (!existingLD) {
      var ld = {
        '@context': 'https://schema.org',
        '@type': 'Course',
        'name': finalTitle,
        'description': finalDescription,
        'url': canonical,
        'provider': {
          '@type': 'Person',
          'name': 'Rifat Erdem Sahin',
          'url': 'https://www.linkedin.com/in/rifaterdemsahin'
        },
        'educationalLevel': 'Advanced',
        'teaches': 'Claude AI system architecture, MCP server deployment, multi-agent orchestration, prompt caching optimization, enterprise AI security'
      };
      if (path === '/index.html') {
        ld['@type'] = 'Course';
        ld['name'] = 'Claude AI Certification for Architects';
      }
      var script = document.createElement('script');
      script.setAttribute('type', 'application/ld+json');
      script.textContent = JSON.stringify(ld);
      head.appendChild(script);
    }

    addMeta('seo-injected', 'true');
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', ensureTags);
  } else {
    ensureTags();
  }
})();