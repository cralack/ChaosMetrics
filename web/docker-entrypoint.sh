#!/bin/sh
set -e

# Run setupDataDragon.sh script
/usr/local/bin/setupDataDragon.sh

# Execute Nginx entrypoint command
exec "$@"
