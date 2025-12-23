# jcaldwell-labs Organization Guidelines

This document establishes standards for repository organization and documentation across jcaldwell-labs projects.

## Repository Structure

### Root Directory (Minimal)

Keep the root directory clean with only essential files:

```
fintrack/
├── README.md              # Project overview and quick start
├── LICENSE                # MIT license
├── .gitignore             # Git ignore patterns
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── Makefile               # Build automation
├── llms.txt               # AI discoverability file
├── CONTRIBUTING.md        # Contribution guidelines
├── TESTING.md             # Test strategy
├── SECURITY.md            # Security guidelines
├── CHANGELOG.md           # Version history
└── CLAUDE.md              # AI assistant context
```

### Documentation Structure

```
docs/
├── README.md              # Documentation index
├── guides/                # How-to guides
├── tutorials/             # Step-by-step learning
└── examples/              # Sample configs and scripts

.github/
├── planning/              # Internal planning (roadmap, backlogs)
├── workflows/             # CI/CD workflows
├── ISSUE_TEMPLATE/        # Issue templates
└── ORG-GUIDELINES.md      # This file
```

## README Standards

Every README must include:

1. **Badges** - License, language version, PR status
2. **Value Proposition** - "Why [project]?" section explaining the problem solved
3. **Quick Start** - Installation and first use in under 5 minutes
4. **Demo** - Visual demonstration or try-it-yourself commands
5. **Features** - Current capabilities and roadmap
6. **Comparison** - How it differs from alternatives
7. **Documentation Links** - Clear navigation to detailed docs
8. **Community** - Issues, PRs, discussions links

### Quality Bar

- Can a stranger understand what this does in 30 seconds?
- Can they try it in less than 5 minutes?
- Is the value proposition clear?

## AI Discoverability

### llms.txt File

Every project should have an `llms.txt` file in the root containing:

- Project name and tagline
- What the project does (2-3 sentences)
- Key capabilities (bullet list)
- Quick start commands
- Common usage patterns
- Architecture overview
- Repository URL
- License

This enables AI assistants to better understand and assist with the project.

## Documentation Organization

### User-Facing (docs/)

- **guides/** - Task-oriented how-to content
- **tutorials/** - Learning-oriented step-by-step content
- **examples/** - Reference code and configurations

### Internal (hidden)

- **.github/planning/** - Roadmaps, backlogs, sprint planning
- Keep separate from user documentation
- Not discoverable from main navigation

## GitHub Configuration

### Topics

Add 5-10 relevant topics to the repository:
- Language/framework (go, cli, terminal)
- Domain (finance, budgeting, personal-finance)
- Features (privacy, local-first, unix)

### Description

Concise summary (70-120 characters) with keywords:
> Terminal-based personal finance tracking - privacy-first, scriptable, Unix philosophy

### Templates

Use issue and PR templates for consistent contributions.

## Project Maturity Levels

| Level | Criteria |
|-------|----------|
| L1 - Experimental | README exists, builds, basic docs |
| L2 - Alpha | Tests, CI/CD, contributing guide |
| L3 - Beta | Full docs, llms.txt, organized structure |
| L4 - Production | Comparison table, tutorials, community |

## Checklist

Before considering a project "polished":

- [ ] README has badges and value proposition
- [ ] llms.txt created
- [ ] docs/ organized with README index
- [ ] .github/planning/ contains roadmap
- [ ] Repository has topics and description
- [ ] Issue templates configured
- [ ] Contributing guidelines documented

## Resources

- [shields.io](https://shields.io) - Badge generation
- [asciinema](https://asciinema.org) - Terminal recordings
- [my-grid reference](https://github.com/jcaldwell-labs/my-grid) - Example implementation
