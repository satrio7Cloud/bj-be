package qris

import (
	"fmt"
	"strconv"
)

type TLVNode struct {
	Tag      string
	Length   int
	Value    string
	Children []TLVNode
}

func ParseTLV(payload string, parseNested bool) ([]TLVNode, error) {
	nodes := make([]TLVNode, 0)
	cursor := 0

	for cursor < len(payload) {
		if cursor+4 > len(payload) {
			return nil, fmt.Errorf("%w: incomplete tlv field at position %d", ErrQRISInvalid, cursor)
		}

		tag := payload[cursor : cursor+2]
		lengthText := payload[cursor+2 : cursor+4]
		length, err := strconv.Atoi(lengthText)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid length for tag %s", ErrQRISInvalid, tag)
		}

		valueStart := cursor + 4
		valueEnd := valueStart + length
		if valueEnd > len(payload) {
			return nil, fmt.Errorf("%w: tag %s length exceeds payload", ErrQRISInvalid, tag)
		}

		node := TLVNode{Tag: tag, Length: length, Value: payload[valueStart:valueEnd]}
		if parseNested && isMerchantAccountInfo(tag) {
			if children, err := ParseTLV(node.Value, false); err == nil && SerializeTLV(children) == node.Value {
				node.Children = children
			}
		}

		nodes = append(nodes, node)
		cursor = valueEnd
	}

	return nodes, nil
}

func SerializeTLV(nodes []TLVNode) string {
	output := ""
	for _, node := range nodes {
		output += EncodeTLV(node.Tag, node.Value)
	}
	return output
}

func EncodeTLV(tag, value string) string {
	return fmt.Sprintf("%s%02d%s", tag, len(value), value)
}

func FindTag(nodes []TLVNode, tag string) *TLVNode {
	for index := range nodes {
		if nodes[index].Tag == tag {
			return &nodes[index]
		}
	}
	return nil
}

func RemoveTag(nodes []TLVNode, tag string) []TLVNode {
	filtered := make([]TLVNode, 0, len(nodes))
	for _, node := range nodes {
		if node.Tag != tag {
			filtered = append(filtered, node)
		}
	}
	return filtered
}

func UpsertTag(nodes []TLVNode, tag, value string, beforeTags map[string]bool) []TLVNode {
	updated := make([]TLVNode, 0, len(nodes)+1)
	for _, node := range nodes {
		if node.Tag == tag {
			node.Value = value
			node.Length = len(value)
		}
		updated = append(updated, node)
	}

	for _, node := range updated {
		if node.Tag == tag {
			return updated
		}
	}

	newNode := TLVNode{Tag: tag, Length: len(value), Value: value}
	for index, node := range updated {
		if beforeTags[node.Tag] {
			next := append([]TLVNode{}, updated[:index]...)
			next = append(next, newNode)
			next = append(next, updated[index:]...)
			return next
		}
	}

	return append(updated, newNode)
}

func isMerchantAccountInfo(tag string) bool {
	value, err := strconv.Atoi(tag)
	return err == nil && value >= 26 && value <= 51
}
