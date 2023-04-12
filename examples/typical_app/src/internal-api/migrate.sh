#!/bin/sh
echo "Migrating things..."
seq 1 5 | xargs -I{} sh -c "date && sleep 10"
echo "Done migrating things..."