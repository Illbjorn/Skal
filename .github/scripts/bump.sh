#!/usr/bin/env bash
##
# Performs an increment of +1 to one of the semantic version segments: major,
# minor or patch.
#
# Args:
#   $1 - The path to the VERSION file.
#   $2 - The version segment to increment.
##

# Collect Args.
s="$1"
v="$2"

# Read the current semantic version.
ver=$(cat "${v}")

# Split the read version string to segments.
major=$(echo "$ver" | grep -Po '^ *\K\d+')
minor=$(echo "$ver" | grep -Po '^ *\d+\.\K\d+')
patch=$(echo "$ver" | grep -Po '^ *\d+\.\d+\.\K\d+')

# Increment the requested version.
case "$s" in
  "major")
    ((major+=1))
    minor=0
    patch=0
    ;;

  "minor")
    ((minor+=1))
    patch=0
    ;;

  "patch")
    ((patch+=1))
    ;;

  *)
  echo "ERROR: Received invalid version string segment: '${s}'."
  exit 1
esac

# Assemble the new version.
new_ver="$major.$minor.$patch"

# Indicate the version bump.
echo "Bumping version $ver -> $new_ver."

# Increment the version.
echo "$major.$minor.$patch" > "${v}"
