package main

import (
	"sort"
	"strings"
	"sync"
)

type ArgumentStat struct {
	Argument string
	Count    uint64
	Size     float64
}

type CommandStat struct {
	Command   string
	Count     uint64
	Size      float64
	Arguments []ArgumentStat
}

type RespResult struct {
	mu        sync.Mutex
	arguments map[string][]string
	counts    map[string]uint64
	sizes     map[string]float64
}

func NewRespResult() *RespResult {
	return &RespResult{
		arguments: make(map[string][]string),
		counts:    make(map[string]uint64),
		sizes:     make(map[string]float64),
	}
}

//Add save commans stats
func (r *RespResult) Add(command string, args string, size float64) {
	argslist, ok := r.arguments[command]
	if !ok {
		r.arguments[command] = []string{}
	}
	r.arguments[command] = append(argslist, args)
	r.counts[command] = r.counts[command] + 1
	r.sizes[command] = r.sizes[command] + size
}

//Commands export all commands as slice
func (r *RespResult) Commands() (list []string) {

	list = make([]string, 0, len(r.counts))
	for command := range r.counts {
		list = append(list, command)
	}

	return list
}

//Clear clear arguments of the command
func (r *RespResult) Clear(command string) {
	r.arguments[command] = []string{}
	r.counts[command] = 0
	r.sizes[command] = 0
}

//longestCommonPrefix I have get this code from leetcode
//ref: https://leetcode.com/problems/longest-common-prefix/discuss/272477/Golang-solution-(0ms-2.3mb)
func (r *RespResult) longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	ans := len(strs[0])
	for {
		if ans == 0 {
			return ""
		}
		prefix := strs[0][:ans]
		flag := true
		for i := 1; i < len(strs); i++ {
			if len(strs[i]) < ans || strs[i][:ans] != prefix {
				flag = false
				break
			}
		}
		if flag {
			return prefix
		}
		ans--
	}
}

func (r *RespResult) findSliceSegments(command string) (segments []int) {
	argslist, ok := r.arguments[command]
	if !ok {
		return
	}

	sort.Strings(argslist)

	matchedprefix, prev := "", ""
	matched := true
	segments = []int{0}

	for i := 0; i < len(argslist); i++ {
		currentItem := argslist[i]
		if i-1 < 0 {
			continue
		}

		if matchedprefix != "" {
			matched, matchedprefix = r.hasprefix(currentItem, matchedprefix, true)
			if !matched {
				segments = append(segments, i)
			}

			if i+1 >= len(argslist) {
				segments = append(segments, i+1)
			}
			continue
		}

		prev = argslist[i-1]
		matched, matchedprefix = r.hasprefix(currentItem, prev, false)
		if !matched {
			segments = append(segments, i)
		}
	}

	return segments
}

func (r *RespResult) hasprefix(current, next string, quick bool) (bool, string) {
	if !strings.HasPrefix(current, next[0:1]) {
		return false, ""
	}

	if quick {
		ok := strings.HasPrefix(current, next)
		if ok {
			return true, next
		}
		return false, ""
	}

	for i := len(next); i >= 0; i-- {
		if strings.HasPrefix(current, next[0:i]) {
			return true, next[0:i]
		}
	}

	return false, ""
}

func (r *RespResult) calculateCommandStats(command string) (stat CommandStat) {
	argslist, ok := r.arguments[command]
	if !ok {
		return
	}

	stat = CommandStat{Command: command}
	stat.Count = uint64(len(argslist))
	stat.Arguments = []ArgumentStat{}

	segments := r.findSliceSegments(command)
	for g := 0; g < len(segments); g++ {
		if g+1 >= len(segments) {
			continue
		}
		start := segments[g]
		stop := segments[g+1]
		segment := argslist[start:stop]
		argstat := ArgumentStat{}
		argstat.Count = uint64(len(segment))
		argstat.Argument = r.longestCommonPrefix(segment)

		stat.Arguments = append(stat.Arguments, argstat)
	}

	return stat
}
