package adventofcode2022

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type CmdName int

const (
	NotSupported CmdName = iota
	LS
	CD
)

func GetCmdName(cmd string) CmdName {
	switch cmd {
	case "ls":
		return LS
	case "cd":
		return CD
	default:
		return NotSupported
	}
}

type Command struct {
	CMD    CmdName
	Args   []string
	Output []string
}

type CommandQueue []*Command

func ToCmdQueue(ir InputReader) (CommandQueue, error) {
	content, err := ir.GetInput()
	if err != nil {
		return nil, err
	}

	var cmdQueue CommandQueue = []*Command{}
	for i := 0; i < len(content); i++ {
		line := content[i]
		command, err := parseCmd(line)
		if err != nil {
			return nil, err
		}
		switch command.CMD {
		case LS:
			parsedOut, parsedLinesCnt := parseLSOutput(i+1, content)
			i += parsedLinesCnt
			command.Output = parsedOut
			cmdQueue = append(cmdQueue, command)
		case CD:
			cmdQueue = append(cmdQueue, command)
		case NotSupported:
			return nil, fmt.Errorf("found unsuppoerted command: %v", command.CMD)
		default:
			return nil, fmt.Errorf("unexpected behaviour")
		}
	}
	return cmdQueue, nil
}

func parseLSOutput(lineN int, lines []string) ([]string, int) {
	res := []string{}
	parsedCnt := 0
	for i := lineN; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "$") {
			break
		}
		parsedCnt++
		res = append(res, line)
	}
	return res, parsedCnt
}

func parseCmd(line string) (*Command, error) {
	cmdRe := regexp.MustCompile(`(\$) (\w+)\s?(.+)?`)
	parsed := cmdRe.FindStringSubmatch(line)
	if parsed[1] != "$" {
		return nil, fmt.Errorf("%v not a command", line)
	}
	cmd := GetCmdName(parsed[2])
	args := []string{}
	if len(parsed) == 4 {
		parsedArgs := strings.Split(parsed[3], "\\s")
		args = append(args, parsedArgs...)
	}
	return &Command{CMD: cmd, Args: args}, nil
}

type NodeType int

const (
	File NodeType = iota
	Dir
)

type Tree struct {
	Type     NodeType
	Size     int
	Name     string
	Parent   *Tree
	Children map[string]*Tree
}

func (t *Tree) populateDirSizes() int {
	if t.Type == File {
		return t.Size
	}
	for _, ch := range t.Children {
		size := ch.populateDirSizes()
		t.Size += size
	}
	return t.Size
}

const LS_OUTPUT_DIR_PREFIX = "dir"

func parseLsOutput(output []string, node *Tree) error {
	for _, item := range output {
		if err := parseLsOutputLine(item, node); err != nil {
			return err
		}
	}
	return nil
}

// Expected ls output either
// dir <dir name>
// <file size> <file name>
func parseLsOutputLine(item string, node *Tree) error {
	splitted := strings.Split(item, " ")
	if len(splitted) != 2 {
		return fmt.Errorf("got: %v, wanted either: [dir <dirname>] or [<size> <filename>]", item)
	}
	if splitted[0] == LS_OUTPUT_DIR_PREFIX {
		return nil
	}
	size, err := strconv.Atoi(splitted[0])
	if err != nil {
		return err
	}
	file := splitted[1]
	node.Children[file] = &Tree{
		Type:     File,
		Size:     size,
		Name:     file,
		Parent:   node,
		Children: map[string]*Tree{},
	}
	return nil
}

const (
	ROOT_DIR = "/"
	GO_UP    = ".."
)

func buildDirTreeByCmdOuque(cq CommandQueue) (*Tree, error) {
	var root *Tree = &Tree{
		Type:     Dir,
		Name:     ROOT_DIR,
		Children: map[string]*Tree{},
		Parent:   nil,
	}
	var currentNode *Tree = nil
	for _, cmd := range cq {
		switch cmd.CMD {
		case LS:
			if err := parseLsOutput(cmd.Output, currentNode); err != nil {
				return nil, err
			}
		case CD:
			switch cmd.Args[0] {
			case ROOT_DIR:
				currentNode = root
			case GO_UP:
				// skip if we already in root dir
				if currentNode.Parent == nil {
					continue
				}
				currentNode = currentNode.Parent
			default:
				ch, ok := currentNode.Children[cmd.Args[0]]
				if !ok {
					ch = &Tree{
						Name:     cmd.Args[0],
						Type:     Dir,
						Parent:   currentNode,
						Children: map[string]*Tree{},
					}
					currentNode.Children[cmd.Args[0]] = ch
				}
				currentNode = ch
			}
		default:
			return nil, fmt.Errorf("unexpected cmd: %v", cmd)
		}
	}
	root.populateDirSizes()
	return root, nil
}

func (t *Tree) sumSizesByCondition(cond func(*Tree) bool) int {
	if t.Type == File {
		return 0
	}
	sum := 0
	for _, ch := range t.Children {
		if cond(ch) {
			sum += ch.Size
		}
		sum += ch.sumSizesByCondition(cond)
	}
	return sum
}

func (t *Tree) findAllByCondition(cond func(*Tree) bool) []*Tree {
	if t.Type == File {
		return []*Tree{}
	}
	found := []*Tree{}
	for _, ch := range t.Children {
		if cond(ch) {
			found = append(found, ch)
		}
		found = append(found, ch.findAllByCondition(cond)...)
	}
	return found
}

func Task7_1(ir InputReader, cnvrInpt func(InputReader) (CommandQueue, error)) (string, error) {
	cmdQueue, err := cnvrInpt(ir)
	if err != nil {
		return "", err
	}
	root, err := buildDirTreeByCmdOuque(cmdQueue)
	if err != nil {
		return "", err
	}
	sum := root.sumSizesByCondition(func(t *Tree) bool { return t.Type == Dir && t.Size <= 100000 })
	return fmt.Sprintf("Result: %v", sum), nil
}

const (
	MAX_SIZE                = 70000000
	MIN_EXPECTED_FREE_SPACE = 30000000
)

func Task7_2(ir InputReader, cnvrInpt func(InputReader) (CommandQueue, error)) (string, error) {
	cmdQueue, err := cnvrInpt(ir)
	if err != nil {
		return "", err
	}
	root, err := buildDirTreeByCmdOuque(cmdQueue)
	if err != nil {
		return "", err
	}
	needToCleanUp := MIN_EXPECTED_FREE_SPACE - (MAX_SIZE - root.Size)
	found := root.findAllByCondition(func(t *Tree) bool { return t.Type == Dir && t.Size >= needToCleanUp })
	sizes := []int{}
	for _, f := range found {
		sizes = append(sizes, f.Size)
	}
	sort.Ints(sizes)
	return fmt.Sprintf("Result: %v", sizes[0]), nil
}
