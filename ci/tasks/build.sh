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

# Glide install to get our dependencies
glide install

# Copy the source code to the gopath since go will expect it there.
WORKING_DIR=$GOPATH/src/$GITHUB_URL
mkdir -p ${WORKING_DIR}
cp -R ${SOURCE_DIR}/* ${WORKING_DIR}/.
cd ${WORKING_DIR}

#
GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/${APP_NAME}-osx
GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/${APP_NAME}-linux


# Put name and tag into files in the output directory
# echo ${DRAFT_VERSION} > ${OUTPUT_DIR}/name
# echo ${DRAFT_VERSION} > ${OUTPUT_DIR}/tag

echo "testversionname" > ${OUTPUT_DIR}/name
echo "testversiontag" > ${OUTPUT_DIR}/tag