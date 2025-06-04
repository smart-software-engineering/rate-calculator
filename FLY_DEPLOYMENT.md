# Deploying Rate Calculator to Fly.io

This guide shows how to deploy the Rate Calculator to Fly.io.

## Prerequisites

1. Install [Fly CLI](https://fly.io/docs/flyctl/install/)
2. Sign up for a Fly.io account: `fly auth signup`
3. Login: `fly auth login`

## Deployment Steps

### 1. Initialize Fly App

```bash
# Create a new Fly app (or use existing fly.toml)
fly launch --no-deploy

# Or if using the included fly.toml:
fly apps create rate-calculator
```

### 2. Set Required Secrets

The app requires a session key for cookie encryption:

```bash
# Generate and set a secure session key
fly secrets set COOKIE_SESSION_KEY=$(openssl rand -base64 32)
```

### 3. Deploy

```bash
# Deploy the application
fly deploy
```

### 4. Access Your App

```bash
# Open the deployed app in your browser
fly open
```

## Configuration Details

### Environment Variables

- **`PORT`**: Automatically set by Fly.io (the app reads this and binds to the correct port)
- **`COOKIE_SESSION_KEY`**: Required secret for session management (set via `fly secrets`)

### Resource Allocation

The app is configured for minimal resource usage:
- **Memory**: 256MB 
- **CPU**: 1 shared CPU
- **Auto-scaling**: Suspends when idle, starts automatically on requests

### Custom Domain (Optional)

To use a custom domain:

```bash
# Add your domain
fly certs add yourdomain.com

# Update DNS to point to Fly.io
# Create CNAME record: yourdomain.com -> your-app.fly.dev
```

## Monitoring

```bash
# View logs
fly logs

# Check app status
fly status

# Scale if needed
fly scale count 2  # Run 2 instances
fly scale memory 512  # Use 512MB RAM
```

## Troubleshooting

### App Won't Start
- Check logs: `fly logs`
- Ensure `COOKIE_SESSION_KEY` is set: `fly secrets list`

### Port Issues
The app automatically detects the PORT environment variable provided by Fly.io. No manual configuration needed.

### SSL/TLS Issues
Fly.io automatically provides HTTPS. The app is configured to work correctly behind Fly.io's proxy.

## Cost Optimization

The default configuration is optimized for minimal costs:
- Suspends when idle (no charges for idle time)
- Minimal resource allocation
- Single instance (can scale up if needed)

For high-traffic scenarios, consider:
- Increasing memory allocation
- Running multiple instances across regions
- Using dedicated CPU instances 