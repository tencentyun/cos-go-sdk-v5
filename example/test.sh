#!/usr/bin/env bash

run() {
    go run "$@"

    if [ $? -ne 0 ]
    then
        exit 3
    fi
}

echo '###### service ####'
run ./service/get.go


echo '##### bucket ####'

run ./bucket/delete.go
run ./bucket/put.go
run ./bucket/putACL.go
run ./bucket/putCORS.go
run ./bucket/putLifecycle.go
run ./bucket/putTagging.go
run ./bucket/get.go
run ./bucket/getACL.go
run ./bucket/getCORS.go
run ./bucket/getLifecycle.go
run ./bucket/getTagging.go
run ./bucket/getLocation.go
run ./bucket/head.go
run ./bucket/listMultipartUploads.go
run ./bucket/delete.go
run ./bucket/deleteCORS.go
run ./bucket/deleteLifecycle.go
run ./bucket/deleteTagging.go


echo '##### object ####'

run ./bucket/putCORS.go
run ./object/put.go
run ./object/uploadFile.go
run ./object/putACL.go
run ./object/append.go
run ./object/get.go
run ./object/head.go
run ./object/getAnonymous.go
run ./object/getACL.go
run ./object/listParts.go
run ./object/options.go
run ./object/initiateMultipartUpload.go
run ./object/uploadPart.go
run ./object/completeMultipartUpload.go
run ./object/abortMultipartUpload.go
run ./object/delete.go
run ./object/deleteMultiple.go
run ./object/copy.go
