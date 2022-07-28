package gjson

import (
	"fmt"
	"strings"
)

type mdTitle struct {
	level string
	name  string
}

type stackState struct {
	// "#" -> 1, "##" -> 2, "###" -> 3
	level int
	// 同一个 level 的序列号，1.1 1.2 1.3 的序列号分别为 1、2、3
	serial int
}

func markdownParser(input []string) []mdTitle {
	// 用来记录解析 md 的上文，以 stack 的数据结构，包括了当前节点的所有的 parent 节点
	var stack = []stackState{}
	var result = []mdTitle{}

	for _, v := range input {
		fmt.Printf("ele: %s, stack: %+v\n", v, stack)

		level, name := parseTitle(v)
		// case 0. 初始状态
		if len(stack) == 0 {
			stack = append(stack, stackState{level: 1, serial: 1})
			result = append(result, mdTitle{level: "1", name: name})
			continue
		}
		n := len(stack)
		curFrame := stack[n-1]
		// case 1. 新的 markdown 与 stack 里面最新的 level 一致
		if level == curFrame.level {
			// 先删除当前的 stackFrame
			stack = stack[0 : n-1]
			stack = append(stack, stackState{level: level, serial: curFrame.serial + 1})
			result = append(result, mdTitle{level: stack2Level(stack), name: name})
		// case 2. 新的 markdown level 大于 stack 里面最新的 level，是更低的节点
		} else if level > curFrame.level {
			stack = append(stack, stackState{level: level, serial: 1})
			result = append(result, mdTitle{level: stack2Level(stack), name: name})
		// case 3. 新的 markdown level 大于 stack 里面最新的 level，是更高的节点
		} else {
			for {
				stack = stack[0 : n-1]
				n = len(stack)
				curFrame = stack[n-1]
				if curFrame.level == level {
					break
				}
				continue
			}
			// 删除这个同级的 stackFrame, stack 里面只保留 parent 节点
			stack = stack[0 : n-1]
			stack = append(stack, stackState{level: level, serial: curFrame.serial + 1})
			result = append(result, mdTitle{level: stack2Level(stack), name: name})
		}
	}
	return result
}

func parseTitle(in string) (int, string) {
	// "### a" -> 3, "a"
	strArray := strings.Split(in, " ")
	switch strArray[0] {
	case "#":
		return 1, strArray[1]
	case "##":
		return 2, strArray[1]
	case "###":
		return 3, strArray[1]
	default:
		panic("invalid title")
	}
}

func stack2Level(stack []stackState) string {
	// eg. 三个 stack 的 serial 分别为 1 3 2，结果就是 1.3.2
	serial := []int{}
	for _, v := range stack {
		serial = append(serial, v.serial)
	}
	switch len(serial) {
	default:
		panic("invalid stack")
	case 1:
		return fmt.Sprintf("%d", serial[0])
	case 2:
		return fmt.Sprintf("%d.%d", serial[0], serial[1])
	case 3:
		return fmt.Sprintf("%d.%d.%d", serial[0], serial[1], serial[2])
	}
}
