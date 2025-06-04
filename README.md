# Freelance Rate Calculator

A smart calculator for freelancers to determine optimal rates for on-site versus remote work scenarios.

## Overview

This application helps freelancers create accurate budget sheets and calculate realistic rates by:

1. Create a personal budget through a wizard that uses customizable templates as a starting point
2. Categorize costs into "must-have", "wishful", and "realistic" budgets
3. Calculate rates based on working days, accounting for holidays, education days, and vacations
4. Generate on-site vs. remote work rate comparisons

A free version of the app will be available under [rate-calculator.frevara.app](https://rate-calculator.frevara.app).

## Technology Stack

- **Backend**: Go (Golang)
- **Frontend**: HTML + HTMX for interactive interfaces
- **Data Storage**: Firestore (with architecture supporting storage replacements)

## Features

### Cost Templates
- Pre-defined and user-created cost templates
- Templates can be private or public
- Categories of expenses with monthly or yearly inputs
- Costs categorized as "must-have", "wishful", or "realistic"

### Working Schedule
- Customizable working days selection
- Country-specific templates for holidays and non-working days
- Automatic calculation of available working hours

### Rate Calculation
- Remote work rate calculations
- On-site rate calculations with travel expense considerations
- Different rate views based on budget categories

## Getting Started

### Prerequisites

- Go 1.20+
- Google Cloud account for Firestore

### Installation

```
git clone https://github.com/yourusername/rate-calculator.git
cd rate-calculator
go mod download
go run cmd/app/main.go
```

### Development with Hot Reload

For development, you can use Air for hot reloading:

```
# Install Air (if not already installed)
go install github.com/cosmtrek/air@latest

# Run the application with Air
./dev.sh
```

Or manually:

```
air
```

This will automatically rebuild and restart the application when changes are detected in your source files.

### Configuration

Create a `config.yaml` file with your Firestore credentials:

```yaml
storage:
  firestore:
    projectID: "your-project-id"
    collection: "rate-calculator"
```

### Usage

1. Navigate to `http://localhost:8080` in your browser
2. Select or create a cost template
3. Fill in your personal expense information
4. Adjust working days and non-working periods
5. View calculated rates for remote and on-site work

## Development

### Project Structure

```
rate-calculator/
├── cmd/                # Application entrypoints
│   └── app/            # Main application
├── internal/           # Private application code
│   ├── calculator/     # Rate calculation logic
│   ├── server/         # HTTP server with embedded assets
│   │   ├── static/     # Embedded static files (CSS, JS)
│   │   └── templates/  # Embedded HTML templates
│   ├── storage/        # Storage interface and implementations
│   │   ├── firestore/  # Firestore implementation
│   │   └── memory/     # In-memory implementation for testing
│   ├── model/          # Data models
│   └── web/            # Web handlers and HTMX templates
├── web/                # Frontend assets
│   ├── template/       # HTML templates
│   └── app/            # Static assets
│       ├── css/        # CSS files
│       └── js/         # JavaScript files
├── .air.toml           # Air configuration for hot reload
├── dev.sh              # Development script for running with Air
└── README.md           # This file
```

### Running Tests

```
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0) - see the LICENSE file for details. This ensures that if someone modifies the code and offers it as a service over a network, they must make their modifications available to users of that service.

## Acknowledgments

- HTMX for enabling rich interactions with minimal JavaScript
- Firebase Firestore for flexible data storage