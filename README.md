# Freelanging Rate Calculator

A modern, extensible web application for calculating freelancer rates, tracking expenses, and managing billing and schedules. Built for transparency, flexibility, and ease of use.

## Features
- **Expense Modeling:** Define and manage detailed expense templates and user-specific data in JSON, with schema validation.
- **Rate Calculation:** Calculate recommended rates based on your expenses, goals, and on-site vs. remote work.
- **Admin Interface:** Separate admin routes for managing data and settings.
- **Modern UI:** Responsive layout using semantic HTML, CSS, and a clean navigation structure.
- **Static Asset Embedding:** All static assets (CSS, JS, favicon) are embedded in the Go binary for easy deployment.

## Technology Stack
- **Go (Golang):** Main backend, using [chi](https://github.com/go-chi/chi) for routing and middleware.
- **Upstash/Redis:** (Planned/Optional) For session management, caching, or real-time features.
- **htmx:** For dynamic, modern HTML interactions without heavy JavaScript frameworks.
- **Alpine.js:** For lightweight, reactive UI enhancements where needed.
- **HTML Templates:** Go's `html/template` for server-side rendering.
- **JSON Schema:** For validating and documenting expense data structures.

## Project Structure
- `web/` — Go HTTP handlers and router setup
- `templates/` — HTML templates (layout, home, show, etc.)
- `static/` — Static assets (CSS, JS, favicon), embedded in the Go binary
- `data/` — Expense templates, user samples, and schemas
- `assets/` — (Legacy/Dev) Static assets before embedding
- `docs/` — Documentation and additional resources

## Getting Started
1. **Install Go 1.20+**
2. **Clone the repository:**
   ```sh
   git clone <your-repo-url>
   cd rate-calculator
   ```
3. **Run the app:**
   ```sh
   air
   ```
4. **Visit:** [http://localhost:3001](http://localhost:3001)

## Customization
- Edit `data/expenses/combined-expenses-template.json` to adjust the default expense model.
- Add or modify HTML templates in `templates/` for UI changes.
- Place your CSS/JS in `static/` for embedding.

## License
MIT

---

*Built with Go, htmx and a passion for freelancer empowerment.*
