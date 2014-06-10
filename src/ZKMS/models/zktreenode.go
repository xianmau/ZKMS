package models

type ZkTreeNode struct {
	Text     string `json:"text"`
	Children bool   `json:"children"`
}
