# HTML File Specifications

This directory contains auto-generated specifications for all HTML files in the project. You can edit these specs to document how the files should be built, their rationales, and acceptance criteria.

## 🔖 Versioning Convention

Every spec carries a `**Version**` field near the top. Bump it **manually** when you edit a spec:

| Bump | When | Example |
|------|------|---------|
| `0.1` | Initial auto-generated spec | new file |
| `0.11`, `0.12` … | Small update — typo, single field, minor clarification | tweak one heading |
| `0.2`, `0.3` … | Bigger update — new functionality, new tables, structural rewrite | add data-layer section |

> ♻️ **Re-generation is safe.** `4_Formula/generate_specs.py` **preserves** the existing version when it re-runs (it only stamps `0.1` on brand-new specs), so manual bumps and hand edits to the version line survive.

## 🧩 What every spec must contain (to recreate the page)

The generator emits these sections for each HTML file so the page can be rebuilt from the spec alone:

1. **🔖 Version** + versioning rule
2. **🎯 Purpose & Rationale**
3. **🧩 Functionality — Recreate the Page** — the interactive functions / behaviours
4. **🗄️ Data Layer — Tables & APIs Used** — Supabase/PostgREST tables + backend/external endpoints
5. **🖥️ UI — Core Layout & Containers** — container IDs (hover reveals the name for prompting) + key headings
6. **🏗️ Implementation Details** — stylesheets, shared scripts, inline constants
7. **✅ Acceptance Criteria**

Hand-authored specs (e.g. `5_Symbols_production_postprod_lower_thirds_spec.md` at `0.2`) may add table schemas, changelogs, and prose beyond what the generator extracts.

## All Specifications

- [4_Formula_cost_calculator_gemini_generation_spec.md](4_Formula_cost_calculator_gemini_generation_spec.md)
- [5_Symbols_course_src_problem-server_templates_problem_spec.md](5_Symbols_course_src_problem-server_templates_problem_spec.md)
- [5_Symbols_course_src_templates_markdown_renderer_spec.md](5_Symbols_course_src_templates_markdown_renderer_spec.md)
- [5_Symbols_course_src_utils_markdown_viewer_spec.md](5_Symbols_course_src_utils_markdown_viewer_spec.md)
- [5_Symbols_menu_map_spec.md](5_Symbols_menu_map_spec.md)
- [5_Symbols_pipeline_spec.md](5_Symbols_pipeline_spec.md)
- [5_Symbols_production_postprod_ai_avatar_spec.md](5_Symbols_production_postprod_ai_avatar_spec.md)
- [5_Symbols_production_postprod_ai_broll_spec.md](5_Symbols_production_postprod_ai_broll_spec.md)
- [5_Symbols_production_postprod_ai_localization_spec.md](5_Symbols_production_postprod_ai_localization_spec.md)
- [5_Symbols_production_postprod_ai_script_gen_spec.md](5_Symbols_production_postprod_ai_script_gen_spec.md)
- [5_Symbols_production_postprod_ai_voiceover_spec.md](5_Symbols_production_postprod_ai_voiceover_spec.md)
- [5_Symbols_production_postprod_asset_checklist_spec.md](5_Symbols_production_postprod_asset_checklist_spec.md)
- [5_Symbols_production_postprod_audio_scoring_spec.md](5_Symbols_production_postprod_audio_scoring_spec.md)
- [5_Symbols_production_postprod_aws_genai_cert_spec.md](5_Symbols_production_postprod_aws_genai_cert_spec.md)
- [5_Symbols_production_postprod_cost_calculator_spec.md](5_Symbols_production_postprod_cost_calculator_spec.md)
- [5_Symbols_production_postprod_customer_discovery_spec.md](5_Symbols_production_postprod_customer_discovery_spec.md)
- [5_Symbols_production_postprod_edit_list_spec.md](5_Symbols_production_postprod_edit_list_spec.md)
- [5_Symbols_production_postprod_flywheel_spec.md](5_Symbols_production_postprod_flywheel_spec.md)
- [5_Symbols_production_postprod_gdrive_sync_spec.md](5_Symbols_production_postprod_gdrive_sync_spec.md)
- [5_Symbols_production_postprod_github_repos_spec.md](5_Symbols_production_postprod_github_repos_spec.md)
- [5_Symbols_production_postprod_graphics_generator_spec.md](5_Symbols_production_postprod_graphics_generator_spec.md)
- [5_Symbols_production_postprod_greenscreen_backgrounds_spec.md](5_Symbols_production_postprod_greenscreen_backgrounds_spec.md)
- [5_Symbols_production_postprod_image_generator_spec.md](5_Symbols_production_postprod_image_generator_spec.md)
- [5_Symbols_production_postprod_index_spec.md](5_Symbols_production_postprod_index_spec.md)
- [5_Symbols_production_postprod_linkedin_controversial_spec.md](5_Symbols_production_postprod_linkedin_controversial_spec.md)
- [5_Symbols_production_postprod_linkedin_messaging_spec.md](5_Symbols_production_postprod_linkedin_messaging_spec.md)
- [5_Symbols_production_postprod_lower_thirds_spec.md](5_Symbols_production_postprod_lower_thirds_spec.md)
- [5_Symbols_production_postprod_memory_palace_spec.md](5_Symbols_production_postprod_memory_palace_spec.md)
- [5_Symbols_production_postprod_music_sfx_score_spec.md](5_Symbols_production_postprod_music_sfx_score_spec.md)
- [5_Symbols_production_postprod_paid_vs_unpaid_certificates_spec.md](5_Symbols_production_postprod_paid_vs_unpaid_certificates_spec.md)
- [5_Symbols_production_postprod_post_production_checklist_spec.md](5_Symbols_production_postprod_post_production_checklist_spec.md)
- [5_Symbols_production_postprod_post_production_master_spec.md](5_Symbols_production_postprod_post_production_master_spec.md)
- [5_Symbols_production_postprod_production_shotlist_spec.md](5_Symbols_production_postprod_production_shotlist_spec.md)
- [5_Symbols_production_postprod_sales_and_marketing_plan_spec.md](5_Symbols_production_postprod_sales_and_marketing_plan_spec.md)
- [5_Symbols_production_postprod_visual_gallery_spec.md](5_Symbols_production_postprod_visual_gallery_spec.md)
- [5_Symbols_production_preprod_10x_certification_spec.md](5_Symbols_production_preprod_10x_certification_spec.md)
- [5_Symbols_production_preprod_bulk_image_generator_spec.md](5_Symbols_production_preprod_bulk_image_generator_spec.md)
- [5_Symbols_production_preprod_course_outline_spec.md](5_Symbols_production_preprod_course_outline_spec.md)
- [5_Symbols_production_preprod_customer_development_spec.md](5_Symbols_production_preprod_customer_development_spec.md)
- [5_Symbols_production_preprod_edit_scripts_spec.md](5_Symbols_production_preprod_edit_scripts_spec.md)
- [5_Symbols_production_preprod_explanations_spec.md](5_Symbols_production_preprod_explanations_spec.md)
- [5_Symbols_production_preprod_image_link_report_spec.md](5_Symbols_production_preprod_image_link_report_spec.md)
- [5_Symbols_production_preprod_index_spec.md](5_Symbols_production_preprod_index_spec.md)
- [5_Symbols_production_preprod_ivq_spec.md](5_Symbols_production_preprod_ivq_spec.md)
- [5_Symbols_production_preprod_log_viewer_spec.md](5_Symbols_production_preprod_log_viewer_spec.md)
- [5_Symbols_production_preprod_multimedia_learning_spec.md](5_Symbols_production_preprod_multimedia_learning_spec.md)
- [5_Symbols_production_preprod_planning_spec.md](5_Symbols_production_preprod_planning_spec.md)
- [5_Symbols_production_preprod_problem_spec.md](5_Symbols_production_preprod_problem_spec.md)
- [5_Symbols_production_preprod_producer_checklist_spec.md](5_Symbols_production_preprod_producer_checklist_spec.md)
- [5_Symbols_production_preprod_production_doctrine_spec.md](5_Symbols_production_preprod_production_doctrine_spec.md)
- [5_Symbols_production_preprod_research_audio_spec.md](5_Symbols_production_preprod_research_audio_spec.md)
- [5_Symbols_production_preprod_research_domain_specific_language_spec.md](5_Symbols_production_preprod_research_domain_specific_language_spec.md)
- [5_Symbols_production_preprod_research_images_spec.md](5_Symbols_production_preprod_research_images_spec.md)
- [5_Symbols_production_preprod_research_index_spec.md](5_Symbols_production_preprod_research_index_spec.md)
- [5_Symbols_production_preprod_research_infographic_generator_spec.md](5_Symbols_production_preprod_research_infographic_generator_spec.md)
- [5_Symbols_production_preprod_research_market_analysis_spec.md](5_Symbols_production_preprod_research_market_analysis_spec.md)
- [5_Symbols_production_preprod_research_notes_spec.md](5_Symbols_production_preprod_research_notes_spec.md)
- [5_Symbols_production_preprod_research_videos_spec.md](5_Symbols_production_preprod_research_videos_spec.md)
- [5_Symbols_production_preprod_risks_spec.md](5_Symbols_production_preprod_risks_spec.md)
- [5_Symbols_production_preprod_sanity_checklist_spec.md](5_Symbols_production_preprod_sanity_checklist_spec.md)
- [5_Symbols_production_preprod_scripts_generator_spec.md](5_Symbols_production_preprod_scripts_generator_spec.md)
- [5_Symbols_production_preprod_scripts_index_spec.md](5_Symbols_production_preprod_scripts_index_spec.md)
- [5_Symbols_production_preprod_self_learning_spec.md](5_Symbols_production_preprod_self_learning_spec.md)
- [5_Symbols_production_preprod_stats_spec.md](5_Symbols_production_preprod_stats_spec.md)
- [5_Symbols_production_preprod_strategy_spec.md](5_Symbols_production_preprod_strategy_spec.md)
- [5_Symbols_production_preprod_tell_show_do_apply_spec.md](5_Symbols_production_preprod_tell_show_do_apply_spec.md)
- [5_Symbols_production_preprod_theory_of_constraints_spec.md](5_Symbols_production_preprod_theory_of_constraints_spec.md)
- [5_Symbols_production_preprod_tools_admin_spec.md](5_Symbols_production_preprod_tools_admin_spec.md)
- [5_Symbols_production_preprod_tools_branding_spec.md](5_Symbols_production_preprod_tools_branding_spec.md)
- [5_Symbols_production_preprod_tools_database_analysis_spec.md](5_Symbols_production_preprod_tools_database_analysis_spec.md)
- [5_Symbols_production_preprod_tools_database_erd_spec.md](5_Symbols_production_preprod_tools_database_erd_spec.md)
- [5_Symbols_production_preprod_ways_of_working_spec.md](5_Symbols_production_preprod_ways_of_working_spec.md)
- [5_Symbols_production_prod_checklist_spec.md](5_Symbols_production_prod_checklist_spec.md)
- [5_Symbols_production_prod_footage_mapping_spec.md](5_Symbols_production_prod_footage_mapping_spec.md)
- [5_Symbols_production_prod_google_drive_folder_creator_spec.md](5_Symbols_production_prod_google_drive_folder_creator_spec.md)
- [5_Symbols_production_prod_google_drive_links_spec.md](5_Symbols_production_prod_google_drive_links_spec.md)
- [5_Symbols_production_prod_index_spec.md](5_Symbols_production_prod_index_spec.md)
- [5_Symbols_production_prod_screenshare_spec.md](5_Symbols_production_prod_screenshare_spec.md)
- [5_Symbols_production_prod_talking-heads_spec.md](5_Symbols_production_prod_talking-heads_spec.md)
- [5_Symbols_production_production_hub_spec.md](5_Symbols_production_production_hub_spec.md)
- [5_Symbols_production_publish_membership_spec.md](5_Symbols_production_publish_membership_spec.md)
- [5_Symbols_production_settings_spec.md](5_Symbols_production_settings_spec.md)
- [5_Symbols_supabase_ui_admin_spec.md](5_Symbols_supabase_ui_admin_spec.md)
- [5_Symbols_templates_axiom_errors_spec.md](5_Symbols_templates_axiom_errors_spec.md)
- [5_Symbols_templates_index_spec.md](5_Symbols_templates_index_spec.md)
- [5_Symbols_timeline_spec.md](5_Symbols_timeline_spec.md)
- [5_Symbols_tools_sitemap_spec.md](5_Symbols_tools_sitemap_spec.md)
- [home_spec.md](home_spec.md)
- [index_spec.md](index_spec.md)
- [markdown_renderer_spec.md](markdown_renderer_spec.md)
