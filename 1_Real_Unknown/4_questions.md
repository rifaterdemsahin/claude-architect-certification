# ❓ Open Questions

> **Stage 1: Real Unknown** — Document the specific questions and unknowns that this project must answer.

---

## 📋 Active Unknowns

| Question | Owner / Agent | Target Stage for Resolution | Resolution Notes / Link |
| :--- | :--- | :--- | :--- |
| **Q1:** What is the preferred hosting platform? | Human / Dev | `2_Environment` | GitHub Pages via GitHub Actions |
| **Q2:** How will user authentication be handled? | Claude / Gemini | `4_Formula` | YouTube Join button handles membership + payment; no custom auth needed |
| **Q3:** What are the performance metrics to meet? | Human | `1_Real_Unknown` | YouTube Analytics: video retention rate is the primary KPI. Target: sustained retention indicating course value |
| **Q4:** What is the break-even point for Erdem? | Human | `1_Real_Unknown` | **£5,000 net/month** to replace contracting income. At ~50% tax, gross must be **£10,000/month**. With life buffer, realistic target is **£15,000 gross/month** |
| **Q5:** What payment/monetisation system will be used? | Human | `2_Environment` | YouTube Join button + YouTube membership tiers. No custom payment gateway needed |
| **Q6:** How many product sections are there? | Human | `1_Real_Unknown` | **Two sections:** (1) **Content** — Claude AI Architect Certification Course; (2) **Delivery Pilot** — the debug/framework system |
| **Q7:** How is the system tested? | Claude / CI | `7_Testing_Known` | GitHub Actions on every push; post-deployment link checks; active GitHub Issues used as failure reports |

---

## 💰 Break-Even Model

| Level | Amount (GBP/month) |
| :--- | :--- |
| Net income target (replace contracting) | £5,000 |
| Gross required (50% tax rate) | £10,000 |
| Realistic target with life buffer | **£15,000** |

YouTube membership pricing and subscriber count needed to hit £15k gross will be tracked in `4_Formula/revenue_model.md`.

---

## 🎯 Two Product Sections

### 1. Content — Claude AI Architect Certification Course
- Delivered via YouTube (videos, retention analytics)
- Membership gated via YouTube Join button
- Performance measured by: view count, watch time, retention %, member count

### 2. Delivery Pilot — Debug Framework
- The 7-stage self-learning system used to build the course
- Showcased as a reusable template for other projects
- Tested via GitHub Actions + post-deployment checks + GitHub Issues

---

## 🧪 Testing Strategy

| Layer | Tool | What it checks |
| :--- | :--- | :--- |
| CI on push | GitHub Actions | Build succeeds, no broken HTML |
| Post-deployment | Link checker workflow | All public URLs return 200 |
| Failures reported | GitHub Issues (auto-created) | Visible to team without manual monitoring |
| YouTube | YouTube Analytics | Retention, watch time, member growth |

---

## 📌 Instructions
1. Document questions **before** writing code.
2. Update the "Resolution Notes" column as soon as a decision is made or implemented.
3. Move fully resolved questions to [1_Real_Unknown/_obsolete/questions.md](_obsolete/questions.md) if the log gets too cluttered.
