package golanglibs

import (
	"github.com/bwmarrin/snowflake"
)

func getSnowflakeID(nodeNumber ...int) int64 {
	var node *snowflake.Node
	var err error
	if len(nodeNumber) != 0 {
		node, err = snowflake.NewNode(Int64(nodeNumber[0]))
	} else {
		node, err = snowflake.NewNode(1)
	}
	panicerr(err)
	return node.Generate().Int64()
}
