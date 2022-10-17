package idgen

import "github.com/bwmarrin/snowflake"

type IdGen struct {
	node *snowflake.Node
}

func New(nodeId int64) (*IdGen, error) {
	n, err := snowflake.NewNode(nodeId)
	if err != nil {
		return nil, err
	}
	gen := &IdGen{
		node: n,
	}
	return gen, nil
}

func (g *IdGen) Generate() string {
	return g.node.Generate().String()
}
