#!/usr/bin/env bash

# Function to check if a container has port 5432 mapped
container_uses_port_5432() {
    local container_name=$1
    if podman inspect --format '{{range $p, $conf := .NetworkSettings.Ports}}{{if eq $p "5432/tcp"}}true{{end}}{{end}}' "$container_name" 2>/dev/null | grep -q "true"; then
        return 0  # Container uses port 5432
    else
        return 1  # Container does not use port 5432
    fi
}

# Check if there are other containers using port 5432 (excluding rate-calculator-db)
for container in $(podman ps --format "{{.Names}}"); do  # Only check running containers
    if [ "$container" != "rate-calculator-db" ] && container_uses_port_5432 "$container"; then
        echo "Stopping container '$container' using port 5432..."
        podman stop "$container"
    fi
done

# Check if the container is already running
if podman ps --filter "name=rate-calculator-db" --filter "status=running" | grep -q rate-calculator-db; then
    echo "Container rate-calculator-db is already running."
    exit 0
fi

# Check if the container exists in any state other than running
if podman ps -a --filter "name=rate-calculator-db" | grep -q rate-calculator-db; then
    # Container exists but is not running, try to start it
    if podman start rate-calculator-db; then
        echo "Started existing rate-calculator-db container."
        exit 0
    else
        # If starting fails, remove it and create a new one
        echo "Failed to start existing container. Removing it..."
        podman rm -f rate-calculator-db
    fi
fi

# If we got here, we need to create a new container
podman pull postgres
podman run --name rate-calculator-db \
    -e POSTGRES_DB=rate_calculator_db \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -p 5432:5432 \
    -d postgres
echo "Created and started new rate-calculator-db container."