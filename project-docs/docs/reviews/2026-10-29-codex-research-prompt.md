# Research brief: the state of agentic software development (late 2026)

You are conducting a thorough, source-backed research survey of the leading
workflows, tools and guidance for agentic software development: agents write
all or most of the code and review most of it, while engineering quality stays
high. Be exhaustive. Use subagents to research areas in parallel, then
synthesise their findings yourself. Prefer primary sources (official docs,
repositories, authors' own posts, talks, papers) over commentary. Record the
date you observed each fact (for example star counts and release versions).

## Why this research exists

The reader runs two personal projects and will compare your findings against a
revised plan for them:

- **fid**: a local-first, agent-first CLI issue tracker (Go, SQLite with an
  append-only journal). The revised plan makes it a work ledger with
  verifiable receipts. It tracks epics, tasks and subtasks, and records proof
  that checks ran on a specific commit, that an independent reviewer
  reviewed the work, and that a human approved protected decisions. It
  refuses "done" without the evidence its risk level requires. It has a
  polished human CLI for drop-in control and `--json` for agents.
- **Agentic Engineering Kit**: a personal engineering-method package. It is
  moving from role-based subagent prompts (Tech Lead, Developer, QA Reviewer,
  QA Lead, an LLM "Phase Runner" orchestrator) towards a small skills plugin.
  The plugin would add one fresh-context reviewer subagent, risk lanes
  (Routine, Standard, Assured) and a deterministic shell loop instead of an
  LLM orchestrator.

The author's previous workflow became "agent waterfall": heavy phase
ceremonies, validation contracts, many markdown status files, more prose than
code. Your report should help them see what the field has learned, where the
consensus is, and where it is still contested.

## Sources to cover (minimum)

Verify each identification. If a name below is ambiguous or wrong, say so and
cover the most likely intended source.

1. **Anthropic**: its guidance on an agentic software development lifecycle
   (search for "Agentic SDLC" and related Anthropic publications); Claude
   Code best practices; "Building effective agents"; posts on context
   engineering and on harnesses for long-running agents; Agent Skills (the
   format and the open standard); Claude Code plugins and marketplaces,
   hooks, subagents, headless mode and the Agent SDK; the official
   `anthropics/skills` repository; and Anthropic's own reported internal
   practices for engineering with Claude.
2. **Matt Pocock's** skills repository and his writing on agentic workflows.
3. **Addy Osmani's** writing and any skills or workflow repositories he
   publishes.
4. **pstack**: the author means pstack specifically, not gstack. Find the
   pstack skills or workflow project for agentic development and cover it in
   depth. If more than one project uses the name, cover the one relevant to
   agentic coding and note the others.
5. **The recent Anthropic hackathon winner** whose Claude Code configuration
   repo gained a large number of GitHub stars (possibly
   `everything-claude-code`). Confirm the identity and cover it.
6. **Other highly starred skills and workflow repos**: for example
   `obra/superpowers`, and anything comparable that is widely adopted now.
   Rank them by adoption, with star counts and dates.
7. **Spec-driven and planning-heavy methods**: GitHub Spec Kit, Kiro,
   BMAD Method, Tessl and similar, plus critiques of them (for example
   Thoughtworks/Martin Fowler site articles on spec-driven development).
8. **Agent task trackers and state tools**: Beads (Steve Yegge), Backlog.md,
   Task Master and any others. Include Yegge's orchestration work, such as
   Gas Town, if relevant.
9. **Loop and orchestration patterns**: the "Ralph" loop (Geoffrey Huntley),
   parallel agents in git worktrees, background and cloud agents, and
   multi-agent orchestration frameworks. What works for unattended runs?
10. **OpenAI Codex** guidance, `AGENTS.md`, Codex skills, and how teams run
    one method across Codex and Claude Code.
11. **Practitioners**: Simon Willison, Kent Beck ("augmented coding"),
    Thoughtworks Technology Radar entries, and other respected engineers
    writing from real use.
12. **Evidence on outcomes**: the DORA AI reports, METR studies, vendor and
    independent productivity and quality studies, and published failure
    analyses. Separate measured evidence from anecdote.
13. **AI code review**: Claude Code review, Codex review, CodeRabbit,
    Greptile and similar, including independent or cross-model review
    patterns.
14. **Verification and provenance**: tools or practices that prove an agent
    actually ran tests or gates. Include adjacent supply-chain attestation
    ideas (SLSA, in-toto, Sigstore) where they inform "receipts" for agent
    work.

Add anything else that is leading, reputable and relevant. Leave out
marketing-only content unless adoption makes it significant.

## Research questions to answer explicitly

For each question, give the consensus (if any), the main alternatives, the
evidence quality, and notable dissent.

1. **Skills vs subagents vs role personas**: how is workflow knowledge best
   encoded now? When do subagents earn their keep (context isolation,
   independent review), and when are role personas an anti-pattern?
2. **Orchestration**: LLM orchestrators vs deterministic scripts, loops and
   hooks. What keeps long or unattended runs on track?
3. **Durable state across sessions**: markdown files vs issue trackers vs
   purpose-built tools. What do agents actually use well?
4. **Planning artifacts**: how much spec or plan prose is right? Where should
   it live? How do teams avoid documentation bloat ("agent waterfall")?
5. **Validation and acceptance**: acceptance criteria, validation contracts,
   eval-driven development, test-first with agents. How is "done" defined
   and enforced?
6. **Verification of agent claims**: how teams stop agents from claiming that
   tests passed or work is finished when it isn't. Is anyone doing
   receipts, attestations or hook-enforced gates?
7. **Review**: independent review agents, cross-model review, how findings are
   tracked and resolved, and how much human review remains.
8. **Human in the loop**: when to interrupt a human, human drop-in or
   take-over UX, approval of risky changes, risk tiers or lanes.
9. **Quality outcomes**: what measurably improves quality (TDD, strict gates,
   small diffs, review) and what doesn't?
10. **Measurement**: how practitioners measure whether a workflow is good
    (human minutes, tokens, defect escape, rework).
11. **Portability**: running the same method across Claude Code, Codex and
    others, and the state of standards (Agent Skills, AGENTS.md, MCP).
12. **Trends and direction**: what changed in 2026, what is fading, and what
    is likely next.

## Output

Write one self-contained HTML file to:

`project-docs/docs/reviews/2026-10-29-agentic-workflow-landscape.html`

Requirements:

- Inline CSS only: no external scripts, no build step, readable offline in a
  desktop browser. Clean typography, a sticky table of contents, and anchor
  links on every section. Use inline SVG or HTML/CSS diagrams where they help
  (for example a landscape map, or a spectrum from heavy to light process).
- **Section 1, "Executive summary for the reader"**: one to two screens of
  plain language. The ten most important takeaways, what has become
  standard, what has fallen out of favour, and the three to five ideas most
  worth adopting.
- **Section 2, landscape tables**: tools, repos and methods with what each
  is, adoption (stars and date), its stance on the research questions, and
  links.
- **Section 3, one subsection per research question**, with findings,
  evidence quality and citations.
- **Section 4, deep dives**: one per major source in the list above. Cover
  what it recommends, distinctive ideas, strengths, weaknesses, and how it
  is used in practice.
- **Section 5, trends and open debates.**
- **Section 6, implications**: a neutral list of patterns, practices and
  tools most relevant to the fid and kit plans described above (state
  tracking, receipts or evidence, review findings, validation, risk lanes,
  human drop-in UX, skills packaging, deterministic loops). Flag where the
  field agrees with, contradicts or goes beyond those plans. Don't redesign
  them; another reviewer will do the comparison.
- **Section 7, sources**: every source with its URL, author, publication date
  and access date. Mark each one as primary or secondary.
- Make the report as long as it needs to be. Thoroughness matters more than
  brevity, but keep the executive summary short.
- Mark confidence levels where the evidence is thin, and never invent a
  source, quote, star count or date. If you can't verify something, say so.
