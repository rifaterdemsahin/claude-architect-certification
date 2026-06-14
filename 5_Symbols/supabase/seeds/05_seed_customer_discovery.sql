-- ============================================================
-- 🔍 Customer Discovery — Checklist Items
-- ============================================================

INSERT INTO checklist_items (phase, item_name, item_desc, sort_order) VALUES
-- Part 1: The Foundational Mindset (The Pyramid)
('Customer Discovery', '1. 🕵️ Customer Discovery', 'Going out to talk to real people to validate assumptions.', 10),
('Customer Discovery', '2. 🩹 Solving a Problem or Pain Point', 'Using those insights to address a genuine friction point.', 20),
('Customer Discovery', '3. 💎 Providing Value to Users', 'Turning that solution into something so useful that customers care about it.', 30),
('Customer Discovery', '4. 🏆 Successful Business', 'The ultimate peak achieved only after the bottom three layers are solid.', 40),

-- Part 2: Phase 1: Formulate Hypotheses
('Customer Discovery', 'Phase 1: 1. Product Hypothesis', 'What are you building, and what features do you think it needs?', 110),
('Customer Discovery', 'Phase 1: 2. Customer & Problem Hypothesis', 'Who do you think has this problem, and how severe is it?', 120),
('Customer Discovery', 'Phase 1: 3. Distribution & Pricing Hypothesis', 'How will you reach these customers, and how much will they pay?', 130),
('Customer Discovery', 'Phase 1: 4. Demand Creation Hypothesis', 'How will you acquire and grow your user base?', 140),
('Customer Discovery', 'Phase 1: 5. Market Type Hypothesis', 'Are you entering an existing market, cloning one, or creating a brand new one?', 150),
('Customer Discovery', 'Phase 1: 6. Competitive Hypothesis', 'Who are the current competitors, and how will you beat them?', 160),

-- Part 2: Phase 2: Test the Problem Hypothesis
('Customer Discovery', 'Phase 2: 7. First Contacts', 'Reach out to potential customers who fit your target persona.', 210),
('Customer Discovery', 'Phase 2: 8. Problem Presentation', 'Share the problem (not your solution yet!) to see if it resonates with them.', 220),
('Customer Discovery', 'Phase 2: 9. Customer Understanding', 'Listen deeply to how they currently cope with this issue.', 230),
('Customer Discovery', 'Phase 2: 10. Market Knowledge', 'Use these conversations to build a broader understanding of the market landscape.', 240),

-- Part 2: Phase 3: Test the Product Hypothesis
('Customer Discovery', 'Phase 3: 11. First Reality Check', 'Take your initial feedback and align it with your product concept.', 310),
('Customer Discovery', 'Phase 3: 12. Product Presentation', 'Show a mockup, prototype, or MVP (Minimum Viable Product) to the customers.', 320),
('Customer Discovery', 'Phase 3: 13. More Customer Visits', 'Gather deep feedback from a wider pool of users on the prototype.', 330),
('Customer Discovery', 'Phase 3: 14. Second Reality Check', 'Analyze whether the product actually solves the problem well enough for them to buy it.', 340),

-- Part 2: Phase 4: Verify and Decide
('Customer Discovery', 'Phase 4: 15. Verify the Product', 'Does the product actually work for the user?', 410),
('Customer Discovery', 'Phase 4: 16. Verify the Problem', 'Is this a must-solve problem that customers are willing to pay for?', 420),
('Customer Discovery', 'Phase 4: 17. Verify the Business Model', 'Can you make money, scale, and acquire customers sustainably?', 430),
('Customer Discovery', 'Phase 4: 18. Iterate or Exit', 'Based on the data, either pivot, proceed, or stop.', 440);
