package fuzz

import "strconv"
import "github.com/apigee/registry/cmd/registry/controller"
import "github.com/apigee/registry/rpc"
import "github.com/apigee/registry/cmd/registry/core"

func mayhemit(bytes []byte) int {

    var num int
    if len(bytes) < 1 {
        num = 0
    } else {
        num, _ = strconv.Atoi(string(bytes[0]))
    }

    switch num {

    case 1:
        content := string(bytes)
        var test rpc.Manifest
        controller.ValidateManifest(content, &test)
        return 0

    default:
        core.UnzipArchiveToMap(bytes)
        return 0
    }
}

func Fuzz(data []byte) int {
    _ = mayhemit(data)
    return 0
}