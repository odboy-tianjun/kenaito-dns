package common

/*
 * @Description  自定义Set结构
 * @Author  https://www.odboy.cn
 * @Date  20241112
 */
import (
	"fmt"
	"strings"
)

type Set struct {
	m map[string]struct{}
}

func NewSet() *Set {
	return &Set{m: make(map[string]struct{})}
}

func (s *Set) Add(val string) {
	s.m[val] = struct{}{}
}

func (s *Set) Remove(val string) {
	delete(s.m, val)
}

func (s *Set) Contains(val string) bool {
	_, exists := s.m[val]
	return exists
}

func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) GetSimilarityValue(text string) string {
	for key := range s.m {
		if strings.HasPrefix(text, key+".") {
			return key
		}
	}
	return ""
}

func main() {
	set := NewSet()
	set.Add("bad")
	set.Add("man")
	fmt.Println(set.Contains("man")) // 输出: true
	fmt.Println(set.Size())          // 输出: 2
}
