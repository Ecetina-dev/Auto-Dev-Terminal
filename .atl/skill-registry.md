# Skill Registry - Auto-Dev-Terminal

Project: Auto-Dev-Terminal | Last Updated: 2026-03-27

## Available Skills

| Skill | Trigger Keywords | Location | Priority |
|-------|-----------------|----------|----------|
| sdd-init | sdd init, iniciar sdd, openspec init | file:///C:/Users/minec/.config/opencode/skills/sdd-init/SKILL.md | HIGH |
| sdd-explore | /sdd-explore | file:///C:/Users/minec/.config/opencode/skill/sdd-explore/SKILL.md | HIGH |
| sdd-propose | /sdd-propose | file:///C:/Users/minec/.config/opencode/skills/sdd-propose/SKILL.md | HIGH |
| sdd-spec | /sdd-spec | file:///C:/Users/minec/.config/opencode/skills/sdd-spec/SKILL.md | HIGH |
| sdd-design | /sdd-design | file:///C:/Users/minec/.config/opencode/skills/sdd-design/SKILL.md | HIGH |
| sdd-tasks | /sdd-tasks | file:///C:/Users/minec/.config/opencode/skills/sdd-tasks/SKILL.md | HIGH |
| sdd-apply | /sdd-apply | file:///C:/Users/minec/.config/opencode/skills/sdd-apply/SKILL.md | HIGH |
| sdd-verify | /sdd-verify | file:///C:/Users/minec/.config/opencode/skills/sdd-verify/SKILL.md | HIGH |
| sdd-archive | /sdd-archive | file:///C:/Users/minec/.config/opencode/skills/sdd-archive/SKILL.md | HIGH |
| sdd-state-tracking | sdd state, sdd tracking | file:///C:/Users/minec/.config/opencode/skills/sdd-state-tracking/SKILL.md | MEDIUM |
| go-testing | Go tests, Go test, testing Go | file:///C:/Users/minec/.config/opencode/skills/go-testing/SKILL.md | HIGH |
| go-testing-enhanced | Go tests enhanced | file:///C:/Users/minec/.config/opencode/skills/go-testing-enhanced/SKILL.md | HIGH |
| data-architect | database, DB, SQL, schema, model | file:///C:/Users/minec/.config/opencode/skills/data-architect/SKILL.md | MEDIUM |
| data-architect-enhanced | database enhanced | file:///C:/Users/minec/.config/opencode/skills/data-architect-enhanced/SKILL.md | MEDIUM |
| skill-creator | create skill, new skill, skill creation | file:///C:/Users/minec/.config/opencode/skills/skill-creator/SKILL.md | MEDIUM |
| skill-creator-enhanced | skill creation enhanced | file:///C:/Users/minec/.config/opencode/skill/skill-creator-enhanced/SKILL.md | MEDIUM |
| skill-registry | update skills, skill registry | file:///C:/Users/minec/.config/opencode/skills/skill-registry/SKILL.md | HIGH |
| skill-learning-core | memory, engram, persist learning | file:///C:/Users/minec/.config/opencode/skills/skill-learning-core/SKILL.md | MEDIUM |
| web-requirements-analyst | requirements, user stories, acceptance criteria | file:///C:/Users/minec/.config/opencode/skills/web-requirements-analyst/SKILL.md | MEDIUM |
| prd-creator | PRD, product requirements | file:///C:/Users/minec/.config/opencode/skills/prd-creator/SKILL.md | MEDIUM |
| typescript | TypeScript, types, interfaces | file:///C:/Users/minec/.config/opencode/skills/typescript/SKILL.md | HIGH |
| angular-architecture | Angular architecture | file:///C:/Users/minec/.config/opencode/skill/angular/architecture/SKILL.md | MEDIUM |
| angular-core | Angular core, Angular signals | file:///C:/Users/minec/.config/opencode/skills/angular/core/SKILL.md | MEDIUM |
| angular-forms | Angular forms, Reactive Forms | file:///C:/Users/minec/.config/opencode/skills/angular/forms/SKILL.md | MEDIUM |
| angular-performance | Angular performance, optimization | file:///C:/Users/minec/.config/opencode/skills/angular/performance/SKILL.md | MEDIUM |
| react-19 | React 19, React Compiler | file:///C:/Users/minec/.config/opencode/skills/react-19/SKILL.md | MEDIUM |
| react-native | React Native, Expo, mobile | file:///C:/Users/minec/.config/opencode/skills/react-native/SKILL.md | MEDIUM |
| nextjs-15 | Next.js 15, App Router | file:///C:/Users/minec/.config/opencode/skills/nextjs-15/SKILL.md | MEDIUM |
| zustand-5 | Zustand, state management | file:///C:/Users/minec/.config/opencode/skills/zustand-5/SKILL.md | MEDIUM |
| tailwind-4 | Tailwind CSS 4 | file:///C:/Users/minec/.config/opencode/skill/tailwind-4/SKILL.md | MEDIUM |
| playwright | Playwright, E2E testing | file:///C:/Users/minec/.config/opencode/skills/playwright/SKILL.md | MEDIUM |
| pytest | pytest, Python testing | file:///C:/Users/minec/.config/opencode/skills/pytest/SKILL.md | MEDIUM |
| django-drf | Django, DRF, REST API | file:///C:/Users/minec/.config/opencode/skills/django-drf/SKILL.md | MEDIUM |
| spring-boot-3 | Spring Boot 3 | file:///C:/Users/minec/.config/opencode/skills/spring-boot-3/SKILL.md | MEDIUM |
| java-21 | Java 21, records, sealed types | file:///C:/Users/minec/.config/opencode/skills/java-21/SKILL.md | MEDIUM |
| hexagonal-architecture-layers-java | Hexagonal architecture Java | file:///C:/Users/minec/.config/opencode/skills/hexagonal-architecture-layers-java/SKILL.md | MEDIUM |
| electron | Electron, desktop app | file:///C:/Users/minec/.config/opencode/skills/electron/SKILL.md | MEDIUM |
| elixir-antipatterns | Elixir, Phoenix antipatterns | file:///C:/Users/minec/.config/opencode/skills/elixir-antipatterns/SKILL.md | LOW |
| ai-sdk-5 | Vercel AI SDK, AI chat | file:///C:/Users/minec/.config/opencode/skills/ai-sdk-5/SKILL.md | MEDIUM |
| zod-4 | Zod validation | file:///C:/Users/minec/.config/opencode/skills/zod-4/SKILL.md | MEDIUM |
| github-pr | GitHub PR, pull request | file:///C:/Users/minec/.config/opencode/skill/github-pr/SKILL.md | MEDIUM |
| jira-epic | Jira epic, large feature | file:///C:/Users/minec/.config/opencode/skills/jira-epic/SKILL.md | LOW |
| jira-task | Jira task, ticket | file:///C:/Users/minec/.config/opencode/skills/jira-task/SKILL.md | LOW |

## Auto-Load Rules

| Context | Skill to Load |
|---------|--------------|
| Go tests, Bubbletea TUI | go-testing |
| Creating new AI skills | skill-creator |
| Database design, models | data-architect |
| Web requirements analysis | web-requirements-analyst |

## Project Notes

- This is a new CLI tool project for automating terminal setup
- Tech stack (TypeScript/Go) to be determined during exploration phase
- Use `/sdd-explore` to investigate and determine optimal stack