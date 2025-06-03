# Freelance Rate Calculator

A smart calculator for freelancers to determine optimal rates for on-site versus remote work scenarios.

## Overview

This application helps freelancers create accurate budget sheets and calculate realistic rates by:

1. Create a personal budget but wizard based by taking in the local where
   they come from, using possible templates they can choose from and then
   navigating through the wizard to input their base rates.
2. Calculate rates based on the selected location and expense model.
3. Calculate an on-site rate based on the client's location. Proposals will
   be done using online services to suggest the costs and how they will impact
   your rates.

A free version of the app will be available under [rate-calculator.frevara.app](https://rate-calculator.frevara.app).

## Technology Stack

- **Backend**: Go (Golang)
- **Frontend**: HTML + HTMX for interactive interfaces
- **Data Storage**: database, postgresql

## Getting Started

### Prerequisites

- Go 1.24+
- Postgresql

### Installation

```
git clone https://github.com/yourusername/rate-calculator.git
cd rate-calculator
go mod download
go run main.go
```

### Usage

1. Navigate to `http://localhost:8080` in your browser
2. Follow the wizard prompts to input your base rates
3. Adjust the suggested values from online services as needed
4. Toggle between on-site and remote scenarios to see different calculations
5. Generate and export your final budget sheet

## Development

### Project Structure

```
rate-calculator/
├── cmd/                # Application entrypoints
├── internal/           # Private application code
│   ├── calculator/     # Rate calculation logic
│   ├── services/       # External service integrations
│   └── web/            # Web handlers and HTMX templates
├── pkg/                # Public libraries
├── web/                # Frontend assets
│   ├── templates/      # HTML templates
│   ├── static/         # CSS, images, etc.
│   └── htmx/           # HTMX-specific components
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
- Various rate suggestion services for providing market data
