#!/bin/sh
echo "Deleting things..."
seq 1 5 | xargs -I{} sh -c "date && sleep 10"
echo "Done deleting things..."