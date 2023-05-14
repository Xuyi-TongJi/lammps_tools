package heatFlux

import (
	"bufio"
	"fmt"
	"os"
)

// input variable
const chunk int = 20      // 块的个数
const timeStep int = 1000 // 输出时间间隔
const tot int = 2780000   // 总采样步数

func Solve() {
	reader := bufio.NewReader(os.Stdin)
	data := make([][]float64, tot/timeStep)
	for i := 0; i < tot/timeStep; i++ {
		data[i] = make([]float64, chunk)
	}
	var (
		x, y  int
		z, zz float64
	) // 占位符
	for i := 0; i < tot/timeStep; i++ {
		fmt.Fscanf(reader, "%d %d %f\n", &x, &y, &z)
		for j := 0; j < chunk; j++ {
			fmt.Fscanf(reader, "%d %f %f %f\n", &x, &z, &zz, &data[i][j])
		}
	}
	f, _ := os.OpenFile("./heatFlux/out.txt", os.O_RDWR, 0666)
	for i := 0; i < tot/timeStep; i++ {
		for j := 0; j < chunk; j++ {
			fmt.Fprintf(f, "%.03f", data[i][j])
			if j != chunk-1 {
				fmt.Fprintf(f, " ")
			}
		}
		fmt.Fprintf(f, "\n")
	}
	// 计算平均
	af, _ := os.OpenFile("./heatFlux/avg.txt", os.O_RDWR, 0666)
	for i := 0; i < chunk; i++ {
		sum := 0.0
		for j := 0; j < tot/timeStep; j++ {
			sum += data[j][i]
		}
		avg := sum / float64(tot/timeStep)
		fmt.Fprintf(af, "%.3f\n", avg)
	}
}
