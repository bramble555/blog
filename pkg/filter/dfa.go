package filter

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

// TrieNode 前缀树节点
type TrieNode struct {
	children map[rune]*TrieNode // 子节点映射
	isEnd    bool               // 是否为敏感词结尾
}

// NewTrieNode 创建新的前缀树节点
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEnd:    false,
	}
}

// SensitiveFilter 敏感词过滤器 (基于DFA算法，线程安全)
type SensitiveFilter struct {
	root *TrieNode    // 前缀树根节点
	mu   sync.RWMutex // 读写锁，保证并发安全
}

// NewSensitiveFilter 创建新的敏感词过滤器
func NewSensitiveFilter() *SensitiveFilter {
	return &SensitiveFilter{
		root: NewTrieNode(),
	}
}

// AddWord 添加单个敏感词到前缀树（线程安全）
func (sf *SensitiveFilter) AddWord(word string) {
	if word == "" {
		return
	}

	// 写锁保护
	sf.mu.Lock()
	defer sf.mu.Unlock()

	// 转换为小写以支持大小写不敏感匹配
	word = strings.ToLower(word)
	node := sf.root

	// 遍历每个字符（支持中英文）
	for _, char := range word {
		if _, exists := node.children[char]; !exists {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}

	// 标记为敏感词结尾
	node.isEnd = true
}

// LoadFromFile 从文件加载敏感词
func (sf *SensitiveFilter) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		sf.AddWord(line)
	}

	return scanner.Err()
}

// Filter 检测并替换文本中的敏感词（线程安全，内存优化）
// 返回值: (是否包含敏感词, 替换后的文本)
func (sf *SensitiveFilter) Filter(text string) (bool, string) {
	if text == "" {
		return false, text
	}

	// 读锁保护
	sf.mu.RLock()
	defer sf.mu.RUnlock()

	// 转换为小写进行匹配
	lowerText := strings.ToLower(text)
	runes := []rune(lowerText)
	length := len(runes)

	// 第一遍扫描：检测是否包含敏感词，记录所有匹配位置
	type matchPos struct {
		start int
		end   int
	}
	var matches []matchPos

	for i := 0; i < length; i++ {
		node, matchEnd := sf.root, -1
		// DFA 状态机匹配
		for j := i; j < length; j++ {
			if child, exists := node.children[runes[j]]; exists {
				node = child
				if node.isEnd {
					matchEnd = j + 1 // 记录匹配结束位置
				}
			} else {
				break
			}
		}

		// 如果找到敏感词，记录位置
		if matchEnd > i {
			matches = append(matches, matchPos{start: i, end: matchEnd})
			i = matchEnd - 1
		}
	}

	// 如果没有敏感词，直接返回原字符串（内存优化：避免不必要的分配）
	if len(matches) == 0 {
		return false, text
	}

	// 第二遍：只有在检测到敏感词时才进行替换
	originalRunes := []rune(text)
	for _, match := range matches {
		for k := match.start; k < match.end; k++ {
			// 所有字符统一替换为 *
			originalRunes[k] = '*'
		}
	}

	return true, string(originalRunes)
}

// Contains 检测文本是否包含敏感词（不替换，线程安全）
func (sf *SensitiveFilter) Contains(text string) bool {
	if text == "" {
		return false
	}

	// 读锁保护
	sf.mu.RLock()
	defer sf.mu.RUnlock()

	lowerText := strings.ToLower(text)
	runes := []rune(lowerText)
	length := len(runes)

	for i := 0; i < length; i++ {
		node := sf.root
		j := i

		for j < length {
			char := runes[j]
			if child, exists := node.children[char]; exists {
				node = child
				j++
				if node.isEnd {
					return true
				}
			} else {
				break
			}
		}
	}

	return false
}

// Build 为了测试,批量构建敏感词前缀树
func (sf *SensitiveFilter) Build(words []string) {
	for _, word := range words {
		sf.AddWord(word)
	}
}
