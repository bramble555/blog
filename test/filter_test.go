package test

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/bramble555/blog/pkg/filter"
)

func loadWordsFromFile(t testing.TB, path string) []string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open sensitive words file error: %v", err)
	}
	defer f.Close()

	var words []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, line)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan sensitive words file error: %v", err)
	}
	return words
}

func maskWord(word string) string {
	return strings.Repeat("*", utf8.RuneCountInString(word))
}

func buildBenchmarkData(t testing.TB, wordCount int) ([]string, string) {
	t.Helper()
	baseWords := loadWordsFromFile(t, "../filter.txt")
	words := make([]string, 0, wordCount)
	words = append(words, baseWords...)

	r := rand.New(rand.NewSource(1))
	for len(words) < wordCount {
		n := r.Intn(999999)
		words = append(words, fmt.Sprintf("badword%06d", n))
	}

	var sb strings.Builder
	sb.Grow(64 * 1024)
	for i := 0; i < 200; i++ {
		sb.WriteString("这是正常内容。This is normal content. ")
		sb.WriteString("PORN and Violence and Spam. ")
		sb.WriteString("这里有色情和暴力以及广告。 ")
		if wordCount > len(baseWords) {
			sb.WriteString(words[len(baseWords)+i%minInt(wordCount-len(baseWords), 100)])
			sb.WriteString(" ")
			sb.WriteString(strings.ToUpper(words[len(baseWords)+i%minInt(wordCount-len(baseWords), 100)]))
			sb.WriteString(" ")
		}
	}
	return words[:wordCount], sb.String()
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestChineseFiltering 测试中文敏感词过滤
func TestChineseFiltering(t *testing.T) {
	sf := filter.NewSensitiveFilter()

	// 添加中文敏感词
	words := []string{"色情", "暴力", "敏感词"}
	sf.Build(words)

	tests := []struct {
		name     string
		input    string
		expected string
		hasSens  bool
	}{
		{
			name:     "单个中文敏感词",
			input:    "这是一条色情内容",
			expected: "这是一条**内容",
			hasSens:  true,
		},
		{
			name:     "多个中文敏感词",
			input:    "色情和暴力都是敏感词",
			expected: "**和**都是***",
			hasSens:  true,
		},
		{
			name:     "无敏感词",
			input:    "这是正常内容",
			expected: "这是正常内容",
			hasSens:  false,
		},
		{
			name:     "敏感词在开头",
			input:    "敏感词在这里",
			expected: "***在这里",
			hasSens:  true,
		},
		{
			name:     "敏感词在结尾",
			input:    "这里有敏感词",
			expected: "这里有***",
			hasSens:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasSens, result := sf.Filter(tt.input)
			if hasSens != tt.hasSens {
				t.Errorf("Filter() hasSens = %v, want %v", hasSens, tt.hasSens)
			}
			if result != tt.expected {
				t.Errorf("Filter() result = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestEnglishFiltering 测试英文敏感词过滤
func TestEnglishFiltering(t *testing.T) {
	sf := filter.NewSensitiveFilter()

	// 添加英文敏感词
	words := []string{"porn", "violence", "spam"}
	sf.Build(words)

	tests := []struct {
		name     string
		input    string
		expected string
		hasSens  bool
	}{
		{
			name:     "单个英文敏感词",
			input:    "This is porn content",
			expected: "This is **** content",
			hasSens:  true,
		},
		{
			name:     "大小写混合",
			input:    "This is PORN and Violence",
			expected: "This is **** and ********",
			hasSens:  true,
		},
		{
			name:     "多个英文敏感词",
			input:    "spam and violence are bad",
			expected: "**** and ******** are bad",
			hasSens:  true,
		},
		{
			name:     "无敏感词",
			input:    "This is normal content",
			expected: "This is normal content",
			hasSens:  false,
		},
		{
			name:     "部分匹配不应触发",
			input:    "sporting event",
			expected: "sporting event",
			hasSens:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasSens, result := sf.Filter(tt.input)
			if hasSens != tt.hasSens {
				t.Errorf("Filter() hasSens = %v, want %v", hasSens, tt.hasSens)
			}
			if result != tt.expected {
				t.Errorf("Filter() result = %v, want %v", result, tt.expected)
			}
		})
	}
}

func BenchmarkDFAFilterVsNaiveReplace(b *testing.B) {
	sizes := []int{20, 200, 2000, 10000}
	for _, size := range sizes {
		size := size
		b.Run(fmt.Sprintf("words=%d", size), func(b *testing.B) {
			words, text := buildBenchmarkData(b, size)

			b.Run("DFA", func(b *testing.B) {
				sf := filter.NewSensitiveFilter()
				sf.Build(words)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, _ = sf.Filter(text)
				}
			})

			b.Run("NaiveReplaceCaseFold", func(b *testing.B) {
				lowerText := strings.ToLower(text)
				lowerWords := make([]string, len(words))
				masks := make([]string, len(words))
				for i, w := range words {
					lowerWords[i] = strings.ToLower(w)
					masks[i] = maskWord(w)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					s := lowerText
					for j := 0; j < len(lowerWords); j++ {
						s = strings.Replace(s, lowerWords[j], masks[j], -1)
					}
					_ = s
				}
			})
		})
	}
}

// TestConcurrentSafety 测试并发安全性
func TestConcurrentSafety(t *testing.T) {
	sf := filter.NewSensitiveFilter()

	// 初始化一些敏感词
	initialWords := []string{"色情", "暴力", "porn", "violence"}
	sf.Build(initialWords)

	// 并发测试参数
	numReaders := 50  // 读协程数量
	numWriters := 10  // 写协程数量
	iterations := 100 // 每个协程的迭代次数

	// 用于同步的 WaitGroup
	var wg sync.WaitGroup

	// 启动多个读协程（模拟用户发评论时的过滤操作）
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			texts := []string{
				"这是色情内容",
				"This is porn",
				"正常评论",
				"Normal comment",
				"暴力和violence混合",
			}
			for j := 0; j < iterations; j++ {
				text := texts[j%len(texts)]
				// 调用 Filter 方法（读操作）
				_, _ = sf.Filter(text)
				// 调用 Contains 方法（读操作）
				_ = sf.Contains(text)
			}
		}(i)
	}

	// 启动多个写协程（模拟动态更新敏感词库）
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			newWords := []string{
				"新敏感词",
				"newbadword",
				"广告",
				"spam",
			}
			for j := 0; j < iterations; j++ {
				word := newWords[j%len(newWords)]
				// 调用 AddWord 方法（写操作）
				sf.AddWord(word)
			}
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()

	// 验证过滤器仍然正常工作
	hasSens, _ := sf.Filter("这是色情内容")
	if !hasSens {
		t.Error("Filter should detect sensitive word after concurrent operations")
	}
}
