#!/bin/bash

# Project renaming script for burn-after-read template
# Usage: ./scripts/rename-project.sh <old-name> <new-name>
# Example: ./scripts/rename-project.sh burn-after-read my-awesome-project

set -e

OLD_NAME=${1:-"burn-after-read"}
NEW_NAME=${2}

if [ -z "$NEW_NAME" ]; then
    echo "Usage: $0 <old-name> <new-name>"
    echo "Example: $0 burn-after-read my-awesome-project"
    exit 1
fi

echo "Renaming project from '$OLD_NAME' to '$NEW_NAME'..."

# Files to exclude from replacement (binary files, .git, etc.)
EXCLUDE_PATTERN="\.(git|jpg|jpeg|png|gif|ico|pdf|exe|bin|so|dylib)(/|$)"

# Find all text files and replace the project name
find . -type f \
    ! -path "./.git/*" \
    ! -path "./vendor/*" \
    ! -path "./node_modules/*" \
    ! -path "./bin/*" \
    ! -name "*.jpg" ! -name "*.jpeg" ! -name "*.png" ! -name "*.gif" \
    ! -name "*.ico" ! -name "*.pdf" ! -name "*.exe" ! -name "*.bin" \
    ! -name "*.so" ! -name "*.dylib" \
    -exec grep -l "$OLD_NAME" {} \; 2>/dev/null | \
while IFS= read -r file; do
    echo "Processing: $file"
    
    # Skip binary files by checking if file is text
    if file "$file" | grep -q "text\|empty"; then
        # Use different sed syntax for macOS and Linux
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            sed -i "" "s|$OLD_NAME|$NEW_NAME|g" "$file"
        else
            # Linux
            sed -i "s|$OLD_NAME|$NEW_NAME|g" "$file"
        fi
    else
        echo "Skipping binary file: $file"
    fi
done

echo "Project renamed successfully!"
echo "Don't forget to:"
echo "1. Update the git remote URL if needed"
echo "2. Run 'go mod tidy' to update dependencies"
echo "3. Update any documentation with the new project name"