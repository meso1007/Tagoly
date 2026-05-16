#!/bin/bash

# Simple Go syntax validation without compilation
echo "Validating Go syntax..."

for file in cmd/tagoly/*.go internal/search/*.go internal/parser/*.go internal/git/*.go; do
    if [ -f "$file" ]; then
        # Check for basic syntax issues
        if grep -q "^func " "$file"; then
            echo "✓ $file: Contains function definitions"
        fi
        
        # Check for import statements
        if grep -q "^import\|^import (" "$file"; then
            echo "✓ $file: Contains imports"
        fi
        
        # Check balanced braces in implementation (not perfect but good enough)
        open_braces=$(grep -o "{" "$file" | wc -l)
        close_braces=$(grep -o "}" "$file" | wc -l)
        
        if [ "$open_braces" -eq "$close_braces" ]; then
            echo "✓ $file: Braces balanced ($open_braces)"
        else
            echo "✗ $file: Unbalanced braces (open: $open_braces, close: $close_braces)"
        fi
    fi
done

echo ""
echo "Validation complete!"
