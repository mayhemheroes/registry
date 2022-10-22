package fuzz

import "strconv"
import "context"
import "github.com/apigee/registry/cmd/registry/controller"
import "github.com/apigee/registry/rpc"
import "github.com/apigee/registry/cmd/registry/core"
import "github.com/apigee/registry/cmd/registry/scoring"

func mayhemit(bytes []byte) int {

    var num int
    if len(bytes) < 1 {
        num = 0
    } else {
        num, _ = strconv.Atoi(string(bytes[0]))
    }

    switch num {
    case 0:
        core.NewLintFromZippedProtos("mayhem", bytes)
        return 0

    case 1:
        content := string(bytes)
        var test rpc.Manifest
        controller.ValidateManifest(content, &test)
        return 0

    case 2:
        content := string(bytes)
        ctx := context.Background()
        var test rpc.Artifact
        core.ExportVersionHistoryToSheet(ctx, content, &test)
        return 0

    case 3:
        content := string(bytes)
        var test rpc.ScoreDefinition
        scoring.ValidateScoreDefinition(content, &test)
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