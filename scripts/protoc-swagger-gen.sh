#!/usr/bin/env bash

set -eo pipefail

mkdir -p ./tmp-swagger-gen

# move the vendor folder to a temp dir so that go list works properly
temp_dir="f29ea6aa861dc4b083e8e48f67cce"
if [ -d vendor ]; then
  mv ./vendor ./$temp_dir
fi

# Get the path of the cosmos-sdk repo from go/pkg/mod
gogo_proto_dir=$(go list -f '{{ .Dir }}' -m github.com/gogo/protobuf)
google_api_dir=$(go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway)
cosmos_proto_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-proto)
cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
alliance_dir=$(go list -f '{{ .Dir }}' -m github.com/terra-money/alliance)
wasm_dir=$(go list -f '{{ .Dir }}' -m github.com/CosmWasm/wasmd)
ibc_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/ibc-go/v6)

# move the vendor folder back to ./vendor
if [ -d $temp_dir ]; then
  mv ./$temp_dir ./vendor
fi

proto_dirs=$(find "$cosmos_sdk_dir"/proto "$alliance_dir"/proto "$wasm_dir"/proto "$ibc_dir"/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    protoc  \
    -I "$cosmos_sdk_dir/proto" \
    -I "$alliance_dir/proto" \
    -I "$ibc_dir/proto" \
    -I "$wasm_dir"/proto \
    -I "$gogo_proto_dir" \
    -I "$google_api_dir"/third_party/googleapis \
    -I "$cosmos_proto_dir/proto" \
    -I "third_party/proto" \
      "$query_file" \
    --swagger_out ./tmp-swagger-gen \
    --swagger_opt logtostderr=true \
    --swagger_opt fqn_for_swagger_name=true \
    --swagger_opt simple_operation_ids=true
  fi
done

npm install -g swagger-combine
swagger-combine ./client/docs/config.json -o ./client/docs/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
rm -rf ./tmp-swagger-gen