package gui

import (
	"fmt"
	"reflect"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Tree struct {
	*tview.TreeView
}

func NewTree() *Tree {
	t := &Tree{
		TreeView: tview.NewTreeView(),
	}

	t.SetBorder(true).SetTitle("json tree").SetTitleAlign(tview.AlignLeft)
	return t
}

func (t *Tree) UpdateView(g *Gui, i interface{}) {
	g.App.QueueUpdateDraw(func() {
		root := tview.NewTreeNode(".").SetChildren(t.AddNode(i))
		t.SetRoot(root).SetCurrentNode(root)
	})
}

func (t *Tree) AddNode(node interface{}) []*tview.TreeNode {
	var nodes []*tview.TreeNode

	switch node := node.(type) {
	case map[string]interface{}:
		for k, v := range node {
			newNode := t.NewNodeWithLiteral(k).
				SetColor(tcell.ColorMediumSlateBlue).SetReference(k)

			list, isList := v.([]interface{})
			if isList && len(list) > 0 {
				newNode.SetSelectable(true)
			}
			newNode.SetChildren(t.AddNode(v))
			nodes = append(nodes, newNode)
		}
	case []interface{}:
		for i, v := range node {
			switch n := v.(type) {
			case map[string]interface{}, []interface{}:
				if reflect.ValueOf(n).Len() > 0 {
					numberNode := tview.NewTreeNode(fmt.Sprintf("[%d]", i+1))
					numberNode.SetChildren(t.AddNode(v))
					nodes = append(nodes, numberNode)
				}
			default:
				nodes = append(nodes, t.AddNode(v)...)
			}
		}
	default:
		nodes = append(nodes, t.NewNodeWithLiteral(node))
	}
	return nodes
}

func (t *Tree) NewNodeWithLiteral(i interface{}) *tview.TreeNode {
	var text string
	node := tview.NewTreeNode("")
	switch v := i.(type) {
	case int32:
		text = fmt.Sprintf("%d", v)
	case int64:
		text = fmt.Sprintf("%d", v)
	case float32:
		text = fmt.Sprintf("%f", v)
	case float64:
		text = fmt.Sprintf("%f", v)
	case bool:
		text = fmt.Sprintf("%t", v)
	case nil:
		text = "null"
	case string:
		text = v
	}

	return node.SetText(text)
}

func (t *Tree) SetKeybindings() {
	t.SetSelectedFunc(func(node *tview.TreeNode) {
		if len(node.GetChildren()) > 0 {
			node.SetExpanded(!node.IsExpanded())
		}
	})
}
