#!/bin/sh
alias redoc='sudo docker run --rm -v $(pwd):/ecum --workdir /ecum ghcr.io/redocly/redoc/cli'

# set working dir
cd ./internal/openapi

# merge openapi files
sudo docker run --rm -v "$PWD":/ecum -w /ecum node:20.0.0 npx openapi-merge-cli

# build openapi docs
redoc build openapi.yaml
if [[ $? != 0 ]]; then
    echo "Failed to generate HTML API documentation."
    exit 1
fi
mv redoc-static.html index.html
cd ../..
mkdir -p ./internal/docs
cp ./internal/openapi/index.html ./internal/docs
# Clean up artifacts
rm ./internal/openapi/index.html
rm ./internal/openapi/openapi.yaml
mkdir -p ./internal/generated

generate_docs()
{
    SERVICE_NAME=$1
    echo "started docs generation (service=$SERVICE_NAME)"
    OPENAPI_YAML_FILE="$SERVICE_NAME.yaml"
    if [ -f "./internal/openapi/$OPENAPI_YAML_FILE" ]; then
        echo "$OPENAPI_YAML_FILE exists. Processing next steps..."
    else 
        echo "$OPENAPI_YAML_FILE does not exist. Can not process. \
        You need to have openapi yaml file for \
        the service you want to have generated docs & code for router&types (pwd={$PWD})"
        exit 1
    fi
    OPENAPI_MERGE_FILE="$SERVICE_NAME-openapi-merge.json"
    if [ -f "./internal/openapi/$OPENAPI_MERGE_FILE" ]; then
        echo "$OPENAPI_MERGE_FILE exists. Processing next steps..."
    else 
        echo "$OPENAPI_MERGE_FILE does not exist. Can not process. \
        You need to have specific file for merging $SERVICE_NAME openapi yaml file with shared code. \
        Pattern for naming it is <service_name>-openapi-merge.json"
        exit 1
    fi
    # set working dir
    cd ./internal/openapi

    # merge openapi files
    sudo docker run --rm -v "$PWD":/ecum -w /ecum node:20.0.0 npx openapi-merge-cli --config ./$OPENAPI_MERGE_FILE
    echo "generated merged openapi yaml file (service=$SERVICE_NAME)"

    # build openapi docs
    OPENAPI_MERGED_YAML_FILE="$SERVICE_NAME-merged.yaml"
    redoc build $OPENAPI_MERGED_YAML_FILE
    if [[ $? != 0 ]]; then
        echo "Failed to generate HTML API documentation."
        exit 1
    fi
    DOCS_HTML_FILE="$SERVICE_NAME.html"
    mv redoc-static.html $DOCS_HTML_FILE
    cd ../..
    mkdir -p ./internal/docs
    mv ./internal/openapi/$DOCS_HTML_FILE ./internal/docs
    echo "generated HTML doc file and moved to ./internal/docs/ dir (service=$SERVICE_NAME)"

    GEN_ROUTER_DIR="./internal/generated/$SERVICE_NAME"
    mkdir -p $GEN_ROUTER_DIR
    oapi-codegen -o $GEN_ROUTER_DIR/router.go \
        -generate gorilla -package $SERVICE_NAME \
        --import-mapping ./shared-internal.yaml:github.com/ecumenos/ecumenos/internal/generated internal/openapi/$OPENAPI_MERGED_YAML_FILE
    oapi-codegen -o $GEN_ROUTER_DIR/types.go \
        -generate types,skip-prune -package $SERVICE_NAME \
        --import-mapping ./shared-internal.yaml:github.com/ecumenos/ecumenos/internal/generated internal/openapi/$OPENAPI_MERGED_YAML_FILE
    echo "generated router & types by openapi file (service=$SERVICE_NAME)"

    # Clean up artifacts
    echo "removed temporary files (service=$SERVICE_NAME)"
}

oapi-codegen -o internal/generated/types.go -generate types,skip-prune -package generated internal/openapi/shared-internal.yaml
generate_docs "zookeeperadmin"
generate_docs "zookeeper"
generate_docs "orbissocius"
generate_docs "orbissociusadmin"
generate_docs "pds"
generate_docs "pdsadmin"