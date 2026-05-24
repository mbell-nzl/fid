# Reusable Audit Prompts

Use these with a fresh agent session, preferably with a different model family
from the implementation loop for high-stakes audits. Save outputs under
`project-docs/docs/reviews/` and triage findings through normal phase tasks.

## Architecture Audit

```text
You are a senior staff engineer conducting an architectural audit of a codebase
you did not write. Your goal is an honest assessment a human architect can act
on, not reassurance.

Context: <one-paragraph project summary, architecture summary, and relevant
files or repo access>.

Evaluate:
1. Overall structure: is there coherent layering, or is it tangled?
2. Abstraction quality: which abstractions earn their keep, and which are
   premature or leaky?
3. Dependency direction: do dependencies point toward core/domain rules, or are
   there violations?
4. State management: where does state live, and are there hidden globals or
   implicit ordering assumptions?
5. Change amplification: for a representative future feature, how many files
   would change and why?
6. Dead code and duplication: identify top candidates.

Output:
- Executive summary, most important first.
- Findings grouped by Critical / High / Medium / Low.
- For each finding: location, evidence, recommended fix, effort estimate.
- Explicit uncertainty list.
- Do not suggest rewrites unless genuinely warranted; prefer minimal-diff fixes.
```

## Security Audit

```text
You are a security engineer reviewing this codebase for vulnerabilities and
risky patterns.

Project summary: <summary>.
Threat model assumptions: <users, deployment, data sensitivity, likely
attackers>.

Evaluate:
1. Input validation and injection surfaces.
2. Authentication and authorization, if applicable.
3. Secret management and accidental credential exposure.
4. Sensitive data at rest and file/storage permissions.
5. Data in transit and transport security, if applicable.
6. Dependency risk and abandoned packages.
7. Insecure defaults and fail-open behavior.
8. Output encoding and injection into generated artifacts.
9. Crypto usage, if any.
10. Missing configuration that should exist for the stack.

Output:
- Findings by Critical / High / Medium / Low / Info.
- For each: CWE when useful, file/line, exploit sketch, remediation.
- Manual verification list for checks requiring runtime testing.
```

## Maintainability Audit

```text
You are a senior engineer inheriting this codebase. Assess how easy it will be
for a competent engineer unfamiliar with the project to make non-trivial changes
safely.

Evaluate:
1. Readability: naming, function/module length, comment quality.
2. Test quality: behavior coverage vs implementation trivia.
3. Documentation: can someone run and understand the project quickly?
4. Onboarding friction: what takes longest to understand?
5. Risky or hard-to-change areas.
6. Configuration complexity and defaults.
7. Error handling consistency.
8. Logging/observability consistency.

Output:
- Maintainability grade with justification.
- Top 10 improvements ranked by impact times ease.
- One-paragraph reading guide for the next engineer.
```

## Performance Audit

```text
You are a performance engineer reviewing this codebase for inefficiency and
scalability risk.

Expected scale: <scale>.
Budgets: <latency, throughput, memory, startup, bundle/binary size, etc.>.

Evaluate:
1. Algorithmic complexity in hot paths.
2. Storage/database access patterns.
3. Memory usage and large object lifetimes.
4. Network/API chattiness, timeouts, retries, and batching.
5. Concurrency and lock/contention risks.
6. Startup time and lazy-loading opportunities.
7. Client/UI rendering or bundle concerns, if applicable.

Output:
- For each finding: location, current cost estimate, suggested fix, expected
  improvement, effort.
- Suggested benchmarks or load scenarios.
- Flag anything that may look slow but is probably acceptable at expected scale.
```

## Code Complexity Audit

```text
Review this codebase for complexity hotspots that will become maintenance pain.
Complexity means cognitive load, not just cyclomatic count.

Evaluate:
1. Files/modules/functions/components that exceed the project's size budgets.
2. Functions with many parameters or deep nesting.
3. Modules with more than one reason to change.
4. Boolean flag parameters that substantially change behavior.
5. Clever or unusual code where straightforward code would do.
6. Hidden coupling via globals, shared state, or implicit ordering.
7. Dead abstractions: one-implementation interfaces, pass-through wrappers,
   factories for one thing.

Output:
- Top 15 hotspots ranked by future-change cost and likelihood.
- For each: location, what makes it complex, refactor suggestion, effort.
- Short list of simplifications the codebase would benefit from overall.
```

## Spec Drift Audit

```text
You have access to the product/technical/design/quality specs and the actual
codebase. Find divergence between specification and implementation.

For each category, list discrepancies:
1. Features in the spec that are not implemented.
2. Features in implementation that are not in the spec.
3. Behavior specified one way but implemented another.
4. Edge cases implementation handles but specs do not mention.
5. Edge cases specs mention but implementation does not handle.
6. Terminology drift between code, docs, UI, and tests.

Output:
- Table with discrepancy, location, evidence, recommended action (update spec,
  update code, discuss), and severity.
- Do not make changes; report only.
```

## Dependency And Supply-Chain Audit

```text
Review all third-party dependencies in this project, including development and
test dependencies.

For each dependency, evaluate:
1. Maintenance health: recent releases, activity, issue response.
2. Security track record: known vulnerabilities and response behavior.
3. Criticality: how hard it would be to remove or replace.
4. Version pinning and freshness.
5. License compatibility with the project policy.
6. Supply-chain risk signals such as unusual ownership changes, typosquatting,
   or packages doing too much.
7. Transitive dependency depth and risk.

Output:
- Dependencies review sorted by risk.
- Recommendations: remove, replace, pin, update, or accept.
- Required spec/decision updates for accepted risk.
```
