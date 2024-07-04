# SmoothStart

Free onboarding platform for your team!

Built with Go, Echo, HTMX, Templ, PostgreSQL, Redis and Docker.

# Functionalities

- Setup custom onbording plan for each role
- Manage and control the onboarding process of each new team member
- Setup custom notifications

# Usage

1. Clone smoothstart repo:
```bash
git clone git@github.com:Fishmansky/SmoothStart.git
```

2. Copy env_example to .env
```bash
cp env_example .env
```
and set variables inside it.

3. Build it
```bash
make build
```

4. Run it
```bash
make run
```

App will run at http://localhost:8080
