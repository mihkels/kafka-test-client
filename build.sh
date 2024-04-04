#!/usr/bin/env bash

BUILD_YEAR=$(date +"%Y")
BUILD_MONTH=$(date +"%m")
BUILD_DAY=$(date +"%d")

echo "All tags:"
git tag -l

TAGS=$(git tag -l | grep -E "^v${BUILD_YEAR}\.${BUILD_MONTH}\.${BUILD_DAY}\.[0-9]+$")
echo "Tags after grep:"
echo "$TAGS"

SORTED_TAGS=$(echo "$TAGS" | sort -V)
echo "Tags after sort:"
echo "$SORTED_TAGS"

LATEST_TAG=$(echo "$SORTED_TAGS" | tail -n 1)
echo "Latest tag:"
echo "$LATEST_TAG"

if [ -n "$LATEST_TAG" ]; then
    START_NUMBER=$(echo "$LATEST_TAG" | awk -F '.' '{print $NF}')
else
    START_NUMBER=0
fi

echo "Start number: $START_NUMBER"

#if [ -f build_number.txt ]; then
#    START_NUMBER=$(cat build_number.txt)
#else
#    START_NUMBER=0
#fi

# shellcheck disable=SC2004
MONTH_BUILD_NUMBER=$(($START_NUMBER + 1))

REPOSITORY='mihkels'
IMAGE='kafka-tester'

# List of directories
DIRS=("python" "rust" "java" "golang")

# Generate a hash of the .dockerignore file
DOCKERIGNORE_HASH=$(shasum -a 256 .dockerignore | awk '{ print $1 }')

for DIR in "${DIRS[@]}"; do
    if [ -f "${DIR}/.dockerignore" ]; then
        TARGET_DOCKERIGNORE_HASH=$(shasum -a 256 "${DIR}/.dockerignore" | awk '{ print $1 }')
        if [ "$DOCKERIGNORE_HASH" != "$TARGET_DOCKERIGNORE_HASH" ]; then
            cp .dockerignore "${DIR}/"
        fi
    else
        cp .dockerignore "${DIR}/"
    fi

    for DOCKERFILE in "$DIR"/Dockerfile.*; do
        echo "Building $DOCKERFILE"
        EXTENSION="${DOCKERFILE#"$DIR"/Dockerfile.}"
        IMAGE_NAME="${REPOSITORY}/${IMAGE}:${EXTENSION}-${BUILD_YEAR}.${BUILD_MONTH}.${BUILD_DAY}.${MONTH_BUILD_NUMBER}-${DIR}"
        docker buildx build --load --cache-from=type=gha --cache-to=type=gha,mode=max --build-arg BASE_DIR="${DIR}" --progress plain --platform=linux/amd64 -t "$IMAGE_NAME" -f "$DOCKERFILE" .
#        docker buildx build --build-arg BASE_DIR="${DIR}" --progress plain --platform=linux/amd64 -t "$IMAGE_NAME" -f "$DOCKERFILE" .
        docker push "$IMAGE_NAME"
        echo "$IMAGE_NAME"
        echo "Done building $DOCKERFILE"
    done
done

echo "BUILD_NUMBER=v${BUILD_YEAR}.${BUILD_MONTH}.${BUILD_DAY}.${MONTH_BUILD_NUMBER}" >> $GITHUB_ENV
# echo $MONTH_BUILD_NUMBER > build_number.txt