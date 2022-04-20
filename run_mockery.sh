#!/bin/bash
function mock {
  mockery --dir "$INPUT_DIR" --name "$NAME" --filename mock.go --output "$OUTPUT_DIR" --structname "$STRUCT_NAME" --outpkg "$OUT_PKG"
  rm -rf mocks
}

export NAME=Repository # Name of the interface which is going to be mocked
export STRUCT_NAME=RepositoryMock # The output struct name


export INPUT_DIR=domain/advert
export OUTPUT_DIR=internal/advert
export OUT_PKG=advert
mock

export INPUT_DIR=domain/user
export OUTPUT_DIR=internal/user
export OUT_PKG=user
mock
