Here is the complete generic Google Drive storage structure for production, formatted in Markdown for easy copying and implementation.

---

# Comprehensive Google Drive Storage Structure Guide for Video & Course Production

This standardized, scalable folder structure is designed for media and digital course production within Google Drive. Inspired by sequential project lifecycles, this system ensures that production teams, content creators, editors, and marketers can seamlessly collaborate without asset loss or version confusion.

## 1. Core Principles of Drive Organization

To maintain a clean and reliable production pipeline, always adhere to these three rules:

1. **Sequential Numbering:** Prefix root folders with double digits (e.g., `00`, `01`, `02`) to force chronological or logical sorting regardless of alphabetical order.
2. **Immutability of Raw Assets:** Never modify or overwrite original raw recordings. Always save modified versions as separate files or inside designated edit folders.
3. **Centralized Documentation:** Keep tracking sheets, scripts, and guidelines at the root level using a central file (like a `README`) for immediate visibility.

---

## 2. Root Directory Architecture

The table below illustrates the top-level directory structure designed to take a course or video series from initial concept to public launch.

| Folder / File Name | Purpose & Contents | Primary Owners |
| --- | --- | --- |
| **`00-pre-production`** | Scripts, outlines, research, storyboards, budgets, and scheduling timelines. | Producers, Writers |
| **`01-course-intro`** | Assets, timelines, and project files for the introductory module or welcome video. | Editors, Instructors |
| **`02-module-01-[Topic_Name]`** | Core content folder for the first instructional module (scalable for multiple modules). | Instructors, Editors |
| **`03-module-02-[Topic_Name]`** | Core content folder for the second instructional module. | Instructors, Editors |
| **`04-module-03-[Topic_Name]`** | Core content folder for the third instructional module. | Instructors, Editors |
| **`05-course-outro`** | Concluding remarks, call-to-actions, wrap-up materials, and next-step guides. | Instructors, Editors |
| **`06-promo-video`** | Trailers, marketing clips, social media cuts, and promotional assets. | Marketing Team |
| **`README.md`** | Central dashboard containing project links, status updates, and team contacts. | Project Manager |

---

## 3. Detailed Sub-Folder Breakdown

Each root folder requires a standardized internal schema so team members can navigate any section intuitively.

### 3.1 Inside `00-pre-production`

* **`00_scripts_and_outlines`**: Final approved scripts, teleprompter text, and conceptual outlines.
* **`01_research_and_briefs`**: Background documentation, source literature, competitor analyses, and technical references.
* **`02_planning_and_schedules`**: Production calendars, booking confirmations, call sheets, and milestones trackers.
* **`03_brand_and_style_guides`**: Color palettes, typography specifications, intro/outro animations, and lower-third graphic templates.

### 3.2 Inside Content Modules (e.g., `02-module-01-[Topic_Name]`)

Replicate this exact internal structure for every core module or chapter folder created:

* **`01_raw_recordings`**: Unedited video files directly from camera cards, screen captures, or studio feeds.
* `A-Roll/` (Main speaker footage)
* `B-Roll/` (Supplemental context footage)


* **`02_audio_stems`**: High-quality external audio recordings, voiceovers, sound effects (SFX), and ambient tracks.
* **`03_graphics_and_slides`**: Presentation decks (Google Slides/PowerPoint), diagrams, illustrations, and static overlay assets.
* **`04_project_files`**: Video editing project files (e.g., Adobe Premiere Pro, DaVinci Resolve, Final Cut Pro templates).
* **`05_review_exports`**: Draft versions exported for internal stakeholder feedback, quality assurance, and legal review.
* **`06_final_deliverables`**: High-bitrate master files ready for deployment to the learning management system (LMS) or video host.

### 3.3 Inside `06-promo-video`

* **`social_media_cuts`**: Vertical 9:16 and square 1:1 aspect ratio cutdowns for platforms like YouTube Shorts, Instagram Reels, and LinkedIn.
* **`thumbnails_and_banners`**: Designed image assets, title graphics, and variations optimized for different distribution networks.
* **`copywriting_and_meta`**: Video descriptions, optimized titles, tag lists, and email announcement copy templates.

---

## 4. File Naming Conventions & Version Control

To prevent critical tracking failures and accidentally deleting work, enforce the following structured naming system for all uploaded files:

```text
Format:  [ProjectCode]_[Module/Section]_[AssetDescription]_[YYYYMMDD]_v[X.X].[Extension]

Example 1: SEC101_Mod01_IntroLecture_20260607_v1.0.mp4
Example 2: SEC101_PreProd_CourseScript_20260515_v2.1.gdoc

```

* **Version Control Rules:** Major updates increment the integer digit (`v1.0` to `v2.0`), while minor edits or feedback iterations increment the decimal digit (`v1.0` to `v1.1`).

---

## 5. Access Control and Collaboration Settings

Maintain cloud infrastructure security by configuring appropriate Google Drive permissions at the root level:

1. **Full Editors:** Core production staff (Producers, Video Editors, Animators) should have edit access to the entire root directory.
2. **Restricted Editors (Folder-Level):** Instructors or external subject matter experts should only have edit access to their specific module folders and pre-production draft folders.
3. **Viewers:** Executive stakeholders, marketing copywriters, and legal compliance officers should generally be restricted to "Viewer" or "Commenter" access to prevent accidental deletion or relocation of assets.