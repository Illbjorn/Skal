package parse

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/illbjorn/fstr"
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func NewNode(tc *token.Collection, tk token.Token) *Node {
	n := new(Node)
	// Store token.
	if tk == nil {
		return n
	}
	n.Token = tk

	// Store source reference.
	if tc == nil {
		return n
	}

	return n
}

type Node struct {
	Token    token.Token `json:"-"`
	Parent   *Node       `json:"-"`
	Value    string      `json:"node_value,omitempty"`
	Children []*Node     `json:"children,omitempty"`
	Type     token.Type  `json:"node_type"`
}

/*------------------------------------------------------------------------------
 * Stringer Support
 *----------------------------------------------------------------------------*/

var tmplNode = `Node Value    : {value}
Node Type     : {type}
Node Children : {len}`

func (node Node) String() string {
	return fstr.Pairs(
		tmplNode,
		"value", node.Value,
		"type", node.Type.String(),
		"len", strconv.Itoa(len(node.Children)),
	)
}

/*------------------------------------------------------------------------------
 * Children Management
 *----------------------------------------------------------------------------*/

func (node *Node) AddChild(child *Node) *Node {
	if child == nil {
		return node
	}

	node.Children = append(node.Children, child)

	return node
}

func (node *Node) AddChildren(children []*Node) *Node {
	if len(children) == 0 {
		return node
	}

	node.Children = append(node.Children, children...)

	return node
}

/*------------------------------------------------------------------------------
 * Getters and Setters
 *----------------------------------------------------------------------------*/

func (node *Node) SetType(t token.Type) *Node {
	node.Type = t

	return node
}

func (node *Node) SetToken(tk token.Token) *Node {
	if tk == nil {
		return node
	}

	if tk.Type() > 0 && node.Type == 0 {
		node.Type = tk.Type()
	}

	node.Token = tk
	node.Value = tk.Value()

	return node
}

// Assigns the Token to node `a`.
//
// Sets the Token Type only if:
// a) It's not empty, and;
// b) One hasn't already been set.
func (node *Node) SetTokenOnly(tk token.Token) *Node {
	if tk.Type() > 0 && node.Type == 0 {
		node.Type = tk.Type()
	}

	node.Token = tk
	return node
}

/*------------------------------------------------------------------------------
 * Serialization Support
 *----------------------------------------------------------------------------*/

// Serializes Token `a` to JSON and writes the result to provided `path`.
func (node *Node) Serialize(path string) {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		sklog.CFatalF(
			"Failed to open output AST JSON file with error: {err}.",
			"err", err.Error(),
		)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			sklog.CFatalF(
				"Failed to close AST JSON file: {err}.",
				"err", err.Error(),
			)
		}
	}()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err = enc.Encode(node); err != nil {
		sklog.CFatalF(
			"Failed to encode the AST with error: {err}.",
			"err", err.Error(),
		)
	}
}
