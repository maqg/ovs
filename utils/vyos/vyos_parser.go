package vyos

import (
	"bufio"
	"fmt"
	"octlink/ovs/utils"
	"strings"
)

// Parser for Vyos
type Parser struct {
	parsed bool
	Tree   *ConfigTree
}

type role int

const (
	// Root role
	Root role = iota

	// RootAttribute for role
	RootAttribute

	// KeyValue defs
	KeyValue

	// Close defs
	Close

	// Ignore defs
	Ignore

	// Value defs
	Value
)

var (
	// UnitTest if true
	UnitTest = false
)

func matchToken(words []string) (int, role, []string, string) {
	ws := make([]string, 0)
	next := 0

	// find until \n
	for next = 0; next < len(words); next++ {
		w := words[next]
		if w == "\n" {
			break
		}

		ws = append(ws, w)
	}

	length := len(ws)
	if length == 2 && ws[length-1] == "{" {
		return next, Root, []string{ws[0]}, ""
	} else if length > 2 && ws[length-1] == "{" {
		return next, RootAttribute, ws[:length-1], ""
	} else if length >= 2 && ws[length-1] != "{" && ws[length-1] != "}" {
		return next, KeyValue, []string{ws[0]}, strings.Join(ws[1:], " ")
	} else if length == 1 && ws[0] == "}" {
		return next, Close, nil, ""
	} else if length == 1 && ws[0] != "{" && ws[0] != "}" {
		return next, Value, nil, ws[0]
	} else if length == 0 {
		return next + 1, Ignore, nil, ""
	} else {
		panic(fmt.Errorf("unable to parser the words: %s", strings.Join(words, " ")))
	}
}

// GetValue from parser
func (parser *Parser) GetValue(key string) (string, bool) {
	c := parser.Tree.Get(key)
	if c == nil {
		return "", false
	}
	return c.Value(), true
}

// Parse for parser
func (parser *Parser) Parse(text string) *ConfigTree {
	parser.parsed = true

	words := make([]string, 0)
	for _, s := range strings.Split(text, "\n") {
		scanner := bufio.NewScanner(strings.NewReader(s))
		scanner.Split(bufio.ScanWords)
		ws := make([]string, 0)
		for scanner.Scan() {
			ws = append(ws, scanner.Text())
		}
		ws = append(ws, "\n")
		words = append(words, ws...)
	}

	offset := 0
	tree := &ConfigTree{Root: &ConfigNode{}}
	tree.Root.tree = tree
	tstack := &utils.Stack{}

	currentNode := tree.Root
	for i := 0; i < len(words); i += offset {
		o, role, keys, value := matchToken(words[i:])
		offset = o
		if role == Root {
			tstack.Push(currentNode)
			currentNode = currentNode.addNode(keys[0])
		} else if role == KeyValue {
			currentNode.addNode(keys[0]).addNode(value)
		} else if role == Value {
			currentNode.addNode(value)
		} else if role == RootAttribute {
			tstack.Push(currentNode)

			for _, key := range keys {
				if n := currentNode.getNode(key); n == nil {
					currentNode = currentNode.addNode(key)
				} else {
					currentNode = n
				}
			}
		} else if role == Close {
			currentNode = tstack.Pop().(*ConfigNode)
		}
	}

	//txt, _ := json.Marshal(parser.data)
	//fmt.Println(string(txt))

	//fmt.Println(tree.String())
	parser.Tree = tree
	return tree
}

// ConfigurationSourceFunc for configure source
var ConfigurationSourceFunc = func() string {
	bash := utils.Bash{
		Command: "/bin/cli-shell-api showCfg",
		NoLog:   true,
	}

	_, o, _, _ := bash.RunWithReturn()
	bash.PanicIfError()
	return o
}

// ShowConfiguration for configure display
func ShowConfiguration() string {
	return ConfigurationSourceFunc()
}

// NewParserFromShowConfiguration for new parser
func NewParserFromShowConfiguration() *Parser {
	p := &Parser{}
	p.Parse(ConfigurationSourceFunc())
	return p
}

// NewParserFromConfiguration for new parser
func NewParserFromConfiguration(text string) *Parser {
	p := &Parser{}
	p.Parse(text)
	return p
}

// ConfigNode for config node
type ConfigNode struct {
	name          string
	children      []*ConfigNode
	childrenIndex map[string]*ConfigNode
	parent        *ConfigNode
	tree          *ConfigTree
}

// Children for config node
func (n *ConfigNode) Children() []*ConfigNode {
	return n.children
}

// ChildNodeKeys for child node keys
func (n *ConfigNode) ChildNodeKeys() []string {
	keys := make([]string, 0)
	for k := range n.childrenIndex {
		keys = append(keys, k)
	}
	return keys
}

// String for config node
func (n *ConfigNode) String() string {
	stack := &utils.Stack{}
	p := n
	for {
		if p == nil {
			return func() string {
				sl := stack.Slice()
				ss := make([]string, len(sl))
				for i, s := range sl {
					ss[i] = s.(string)
				}
				return strings.TrimSpace(strings.Join(ss, " "))
			}()
		}

		stack.Push(p.name)
		p = p.parent
	}
}

func (n *ConfigNode) isValueNode() bool {
	return n.childrenIndex == nil && n.children == nil
}

func (n *ConfigNode) isKeyNode() bool {
	if len(n.children) != 1 {
		return false
	}

	c := n.children[0]
	return c.isValueNode()
}

// Values for config node
func (n *ConfigNode) Values() []string {
	values := make([]string, 0)
	for _, c := range n.children {
		if c.isValueNode() {
			values = append(values, c.name)
		}
	}
	return values
}

// ValueSize for config node
func (n *ConfigNode) ValueSize() int {
	return len(n.Values())
}

// Value for config node
func (n *ConfigNode) Value() string {
	values := n.Values()
	utils.Assert(len(values) != 0, fmt.Sprintf("the node[%s] doesn't have any value", n.String()))
	utils.Assert(len(values) == 1, fmt.Sprintf("the node[%s] has more than one value%s", n.String(), values))
	return values[0]
}

// Size for config node
func (n *ConfigNode) Size() int {
	return len(n.children)
}

// Delete for config node
func (n *ConfigNode) Delete() {
	n.tree.Delete(n.String())
}

func (n *ConfigNode) deleteSelf() *ConfigNode {
	return n.parent.deleteNode(n.name)
}

func (n *ConfigNode) deleteNode(name string) *ConfigNode {
	delete(n.childrenIndex, name)
	nsl := make([]*ConfigNode, 0)
	for _, c := range n.children {
		if c.name != name {
			nsl = append(nsl, c)
		}
	}
	n.children = nsl
	return n
}

// Getf for config node
func (n *ConfigNode) Getf(f string, args ...interface{}) *ConfigNode {
	if args != nil {
		return n.Get(fmt.Sprintf(f, args...))
	}
	return n.Get(f)
}

// Get for config node
func (n *ConfigNode) Get(config string) *ConfigNode {
	cs := strings.Split(config, " ")
	current := n

	for _, c := range cs {
		current = current.getNode(c)
		if current == nil {
			return nil
		}
	}

	return current
}

func (n *ConfigNode) getNode(name string) *ConfigNode {
	return n.childrenIndex[name]
}

func (n *ConfigNode) addNode(name string) *ConfigNode {
	if c, ok := n.childrenIndex[name]; ok {
		return c
	}

	utils.Assertf(n.tree != nil, "node[%s] has tree == nil", n.String())
	newNode := &ConfigNode{
		name: name,
		tree: n.tree,
	}

	if n.children == nil {
		n.children = make([]*ConfigNode, 0)
	}
	n.children = append(n.children, newNode)

	if n.childrenIndex == nil {
		n.childrenIndex = make(map[string]*ConfigNode)
	}
	n.childrenIndex[name] = newNode
	newNode.parent = n
	return newNode
}

// ConfigTree for vyos
type ConfigTree struct {
	Root           *ConfigNode
	changeCommands []string
}

// HasChanges judge changes of command tree
func (t *ConfigTree) HasChanges() bool {
	return len(t.changeCommands) != 0
}

// Apply changes for command tree
func (t *ConfigTree) Apply(asVyosUser bool) {
	if UnitTest {
		fmt.Println(strings.Join(t.changeCommands, "\n"))
		return
	}

	if len(t.changeCommands) == 0 {
		logger.Debugf("[Vyos Configuration] no changes to apply")
		return
	}

	if asVyosUser {
		RunVyosScriptAsUserVyos(strings.Join(t.changeCommands, "\n"))
	} else {
		RunVyosScript(strings.Join(t.changeCommands, "\n"), nil)
	}
}

func (t *ConfigTree) init() {
	if t.changeCommands == nil {
		t.changeCommands = make([]string, 0)
	}
	if t.Root == nil {
		t.Root = &ConfigNode{
			children:      make([]*ConfigNode, 0),
			childrenIndex: make(map[string]*ConfigNode),
			tree:          t,
		}
	}
}

func (t *ConfigTree) has(config ...string) bool {
	if t.Root == nil || t.Root.children == nil {
		return false
	}

	current := t.Root
	for _, c := range config {
		current = current.childrenIndex[c]
		if current == nil {
			return false
		}
	}

	return true
}

// Has judgement for config tree
func (t *ConfigTree) Has(config string) bool {
	return t.has(strings.Split(config, " ")...)
}

// AttachFirewallToInterface to add firewall config for interface
func (t *ConfigTree) AttachFirewallToInterface(ethname, direction string) {
	t.Setf("interfaces ethernet %v firewall %s name %v.%v", ethname, direction, ethname, direction)
}

// FindFirewallRuleByDescription to find firewall config by description
func (t *ConfigTree) FindFirewallRuleByDescription(ethname, direction, des string) *ConfigNode {
	rs := t.Getf("firewall name %v.%v rule", ethname, direction)

	if rs == nil {
		return nil
	}

	for _, r := range rs.children {
		if d := r.Get("description"); d != nil && d.Value() == des {
			return r
		}
	}

	return nil
}

// SetFirewallDefaultAction to set firewall's default action
func (t *ConfigTree) SetFirewallDefaultAction(ethname, direction, action string) {
	utils.Assertf(action == "drop" || action == "reject" || action == "accept", "action must be drop or reject or accept, but %s got", action)
	t.Setf("firewall name %s.%s default-action %v", ethname, direction, action)
}

// SetFirewallOnInterface to set firewall on interface
func (t *ConfigTree) SetFirewallOnInterface(ethname, direction string, rules ...string) int {
	if direction != "in" && direction != "out" && direction != "local" {
		panic(fmt.Sprintf("the direction can only be [in, out, local], but %s get", direction))
	}

	currentRuleNum := -1
	for i := 1; i <= 9999; i++ {
		if c := t.Getf("firewall name %s.%s rule %v", ethname, direction, i); c == nil {
			currentRuleNum = i
			break
		}
	}

	if currentRuleNum == -1 {
		panic(fmt.Sprintf("No firewall rule number found for the interface %s.%s. You have set more than 9999 rules???", ethname, direction))
	}

	for _, rule := range rules {
		t.Setf("firewall name %v.%v rule %v %s", ethname, direction, currentRuleNum, rule)
	}

	return currentRuleNum
}

// SetFirewallWithRuleNumber to set firewall with rule number
func (t *ConfigTree) SetFirewallWithRuleNumber(ethname, direction string, number int, rules ...string) {
	if direction != "in" && direction != "out" && direction != "local" {
		panic(fmt.Sprintf("the direction can only be [in, out, local], but %s get", direction))
	}

	for _, rule := range rules {
		t.Setf("firewall name %v.%v rule %v %s", ethname, direction, number, rule)
	}
}

// SetDnat to set dnat by config tree
func (t *ConfigTree) SetDnat(rules ...string) int {
	currentRuleNum := -1

	for i := 1; i <= 9999; i++ {
		if c := t.Getf("nat destination rule %v", i); c == nil {
			currentRuleNum = i
			break
		}
	}

	if currentRuleNum == -1 {
		panic("No rule number avaible for dnat. You have set more than 9999 rules???")
	}

	for _, rule := range rules {
		t.Setf("nat destination rule %v %s", currentRuleNum, rule)
	}

	return currentRuleNum
}

// FindDnatRuleDescription to find dnat rule's description
func (t *ConfigTree) FindDnatRuleDescription(des string) *ConfigNode {
	rs := t.Get("nat destination rule")
	if rs == nil {
		return nil
	}

	for _, r := range rs.children {
		if d := r.Get("description"); d != nil && d.Value() == des {
			return r
		}
	}

	return nil
}

// FindSnatRuleDescription to find snat's rule decription
func (t *ConfigTree) FindSnatRuleDescription(des string) *ConfigNode {
	rs := t.Get("nat source rule")

	if rs == nil {
		return nil
	}

	for _, r := range rs.children {
		if d := r.Get("description"); d != nil && d.Value() == des {
			return r
		}
	}

	return nil
}

// SetSnatWithRuleNumber to find snat's rule Number
func (t *ConfigTree) SetSnatWithRuleNumber(ruleNum int, rules ...string) {
	for _, rule := range rules {
		t.Setf("nat source rule %v %s", ruleNum, rule)
	}
}

// SetSnatWithStartRuleNumber to set snat with rule number start
func (t *ConfigTree) SetSnatWithStartRuleNumber(startNum int, rules ...string) int {
	currentRuleNum := -1

	for i := startNum; i <= 9999; i++ {
		if c := t.Getf("nat source rule %v", i); c == nil {
			currentRuleNum = i
			break
		}
	}

	if currentRuleNum == -1 {
		panic("No rule number avaible for source nat. You have set more than 9999 rules???")
	}

	for _, rule := range rules {
		t.Setf("nat source rule %v %s", currentRuleNum, rule)
	}

	return currentRuleNum
}

// SetSnat for config node
func (t *ConfigTree) SetSnat(rules ...string) int {
	return t.SetSnatWithStartRuleNumber(1, rules...)
}

// SetWithoutCheckExisting set the config without checking any existing config with the same path
// usually used for set multi-value keys
func (t *ConfigTree) SetWithoutCheckExisting(config string) {
	t.changeCommands = append(t.changeCommands, fmt.Sprintf("$SET %s", config))
}

// SetfWithoutCheckExisting set the config without checking any existing config with the same path
// usually used for set multi-value keys
func (t *ConfigTree) SetfWithoutCheckExisting(f string, args ...interface{}) {
	if args != nil {
		t.SetWithoutCheckExisting(fmt.Sprintf(f, args...))
	} else {
		t.SetWithoutCheckExisting(f)
	}
}

// Setf if existing value is different from the config
// delete the old one and set the new one
func (t *ConfigTree) Setf(f string, args ...interface{}) bool {
	if args != nil {
		return t.Set(fmt.Sprintf(f, args...))
	}
	return t.Set(f)
}

// Set if existing value is different from the config
// delete the old one and set the new one
func (t *ConfigTree) Set(config string) bool {
	t.init()
	cs := strings.Split(config, " ")
	key := strings.Join(cs[:len(cs)-1], " ")
	value := cs[len(cs)-1]
	keyNode := t.Get(key)
	if keyNode != nil && keyNode.ValueSize() > 0 {
		// the key found
		cvalue := keyNode.Value()
		if value != cvalue {
			keyNode.deleteNode(cvalue)
			keyNode.addNode(value)
			// the value is changed, delete the old one
			t.changeCommands = append(t.changeCommands, fmt.Sprintf("$DELETE %s", key))
			t.changeCommands = append(t.changeCommands, fmt.Sprintf("$SET %s", config))
			return true
		} // the value is unchanged
		return false
	}
	// the key not found
	current := t.Root
	for _, c := range cs {
		current = current.addNode(c)
	}
	t.changeCommands = append(t.changeCommands, fmt.Sprintf("$SET %s", config))
	return true
}

// Getf for config tree
func (t *ConfigTree) Getf(f string, args ...interface{}) *ConfigNode {
	if args != nil {
		return t.Get(fmt.Sprintf(f, args...))
	}
	return t.Get(f)
}

// Get config node from config tree
func (t *ConfigTree) Get(config string) *ConfigNode {
	t.init()
	return t.Root.Get(config)
}

// Deletef delete config node from config tree
func (t *ConfigTree) Deletef(f string, args ...interface{}) bool {
	if args != nil {
		return t.Delete(fmt.Sprintf(f, args...))
	}
	return t.Delete(f)
}

// Delete from config tree
func (t *ConfigTree) Delete(config string) bool {
	n := t.Get(config)
	if n == nil {
		return false
	}

	n.deleteSelf()
	t.changeCommands = append(t.changeCommands, fmt.Sprintf("$DELETE %s", config))
	return true
}

// CommandsAsString to do command as string
func (t *ConfigTree) CommandsAsString() string {
	return strings.Join(t.changeCommands, "\n")
}

// Commands to change
func (t *ConfigTree) Commands() []string {
	return t.changeCommands
}

// ConfigTree to string
func (t *ConfigTree) String() string {
	if t.Root == nil {
		return ""
	}

	strs := make([]string, 0)
	for _, n := range t.Root.children {
		path := utils.Stack{}

		var pathBuilder func(node *ConfigNode)
		pathBuilder = func(node *ConfigNode) {
			if node.children == nil {
				path.Push(node.name)
				strs = append(strs, func() string {
					sl := path.ReverseSlice()
					ss := make([]string, len(sl))
					for i, s := range sl {
						ss[i] = s.(string)
					}
					return strings.Join(ss, " ")
				}())
				path.Pop()

				return
			}

			path.Push(node.name)
			for _, cn := range node.children {
				pathBuilder(cn)
			}
			path.Pop()
		}

		pathBuilder(n)
	}

	return strings.Join(strs, "\n")
}
