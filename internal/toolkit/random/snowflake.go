package random

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/snowflake"
)

type Node struct {
	node *snowflake.Node
}

func NewSnowflakeNode(node int64) (*Node, error) {
	n, err := snowflake.NewNode(node)
	if err != nil {
		return nil, fmt.Errorf("node creation err = %w", err)
	}

	return &Node{node: n}, nil
}

func (n *Node) GenerateInt64() int64 {
	return n.node.Generate().Int64()
}

type GetEntityByIDCallback[T interface{}] func(ctx context.Context, id int64) (*T, error)

func GetSnowflakeID[T interface{}](ctx context.Context, n int64, cb GetEntityByIDCallback[T]) (int64, error) {
	node, err := NewSnowflakeNode(n)
	if err != nil {
		return 0, fmt.Errorf("node creation err = %w", err)
	}

	for i := 0; i < 10; i++ {
		id := node.GenerateInt64()
		a, err := cb(ctx, id)
		if err != nil || a != nil {
			continue
		}
		return id, nil
	}

	return 0, errors.New("can not generate unique id for 10 times of try")
}
