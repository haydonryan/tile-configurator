#!/bin/bash -e

# Setup paths to use for this build.
export GOPATH=$PWD/go
export PATH=$GOPATH/bin:$PATH

# These need to match the parameters in your pipeline.yml
OUTPUT_DIR=$PWD/compiled-output
SOURCE_DIR=$PWD/source-code    


# Because this image doesn't have glide in it...
go get github.com/Masterminds/glide

# We start in the container default directory
# but there is no sourcecode 
cd ${SOURCE_DIR} 

# if Git cli is installed then get the version number from tags.
if [ -d ".git" ]; then
  DRAFT_VERSION=`versioning bump_patch`-`git rev-parse HEAD`
else
  DRAFT_VERSION="v0.0.0-local"
fi
echo "next version should be: ${DRAFT_VERSION}"


# Glide install to get our dependencies
glide install


# Copy the source code to the gopath since go will expect it there.
WORKING_DIR=$GOPATH/src/$GITHUB_URL
mkdir -p ${WORKING_DIR}
cp -R ${SOURCE_DIR}/* ${WORKING_DIR}/.
cd ${WORKING_DIR}

#
GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/${APP_NAME}-osx
md5sum ${OUTPUT_DIR}/${APP_NAME}-osx | cut -f1 -d' ' > ${OUTPUT_DIR}/${APP_NAME}-osx-shasum

GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/${APP_NAME}-linux
md5sum ${OUTPUT_DIR}/${APP_NAME}-linux | cut -f1 -d' ' > ${OUTPUT_DIR}/${APP_NAME}-linux-shasum

# Put name and tag into files in the output directory
echo ${DRAFT_VERSION} > ${OUTPUT_DIR}/name
echo ${DRAFT_VERSION} > ${OUTPUT_DIR}/tag