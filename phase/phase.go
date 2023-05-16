package phase

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type step struct {
	temp   float64
	volume float64
}

// 输入参数
const end int64 = 43784000 - 35000000
const timeStep int64 = 100
const target int = 500 // x步合并为一步

func Solve() {
	reader := bufio.NewReader(os.Stdin)
	ans := make([]step, end/timeStep)
	var x float64
	for i := int64(0); i < end/timeStep; i++ {
		fmt.Fscan(reader, &x)
		fmt.Fscan(reader, &ans[i].temp)
		fmt.Fscan(reader, &ans[i].volume)
		fmt.Fscan(reader, &x)
	}
	sort.Slice(ans, func(i, j int) bool {
		return ans[i].temp < ans[j].temp
	})
	final := make([]step, 0)
	curr, cnt := 0, 0
	sumT, sumV := 0.0, 0.0
	for i := int64(0); i < end/timeStep; i++ {
		sumT += ans[i].temp
		sumV += ans[i].volume
		curr += 1
		if curr == target {
			final = append(final, step{
				temp:   sumT / float64(target),
				volume: sumV / float64(target),
			})
			curr = 0
			cnt += 1
			sumT, sumV = 0.0, 0.0
		}
	}
	// output
	if f, err := os.OpenFile("./phase/out.txt", os.O_RDWR, 0666); err != nil {
		panic(err)
	} else {
		_, _ = f.Seek(0, 0)
		for i := 0; i < cnt; i++ {
			fmt.Fprintf(f, "%.3f %.1f\n", final[i].temp, final[i].volume)
		}
	}
}
