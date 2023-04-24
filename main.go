package main

import (
	"bufio"
	"fmt"
	"os"
)

type atom struct {
	id     int
	aType  int     // 原子类型
	mId    int     // 分子ID
	charge float64 // 电荷
	x      float64
	y      float64
	z      float64
	dx     int
	dy     int
	dz     int
	aStr   string
}

var charges map[int]float64

const (
	atomCnt = 19
	n       = 4000 // 原子总数
)

func init() {
	charges = make(map[int]float64, 0)
	charges[1] = -0.57  // OH
	charges[2] = 0.054  // C2 -> C4
	charges[3] = 0.0    // C -> C0
	charges[4] = -0.159 // C3 -> C3
	charges[5] = 0.41   // HO
	charges[6] = 0.53   // HC
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	atoms := make([]atom, n)
	for i := 0; i < n; i++ {
		fmt.Fscanf(reader, "%d", &atoms[i].id)
		fmt.Fscanf(reader, "%d", &atoms[i].mId) // 占位
		fmt.Fscanf(reader, "%d", &atoms[i].aType)
		fmt.Fscanf(reader, "%f", &atoms[i].charge) // 占位
		fmt.Fscanf(reader, "%f", &atoms[i].x)
		fmt.Fscanf(reader, "%f", &atoms[i].y)
		fmt.Fscanf(reader, "%f", &atoms[i].z)
		fmt.Fscanf(reader, "%d %d %d", &atoms[i].dx, &atoms[i].dy, &atoms[i].dz)
		fmt.Fscanf(reader, "%s", &atoms[i].aStr) // #
		fmt.Fscanf(reader, "%s\n", &atoms[i].aStr)
	}
	// write to file
	f, _ := os.OpenFile("atom.txt", os.O_RDWR, 0666)
	_, _ = f.Seek(0, 0)
	writer := bufio.NewWriter(f)
	mId, cnt := 1, 0
	for i := 0; i < n; i++ {
		// 原子ID
		fmt.Fprint(writer, "   ")
		fmt.Fprintf(writer, "%d", atoms[i].id)
		// 分子ID
		fmt.Fprint(writer, "   ")
		fmt.Fprintf(writer, "%d", mId)
		// 原子类型
		fmt.Fprint(writer, "   ")
		fmt.Fprintf(writer, "%d", atoms[i].aType)
		// 电荷
		fmt.Fprint(writer, "   ")
		fmt.Fprintf(writer, "%.6f", atoms[i].charge)
		// X
		fmt.Fprint(writer, "   ")
		fmt.Fprintf(writer, "%.9f", atoms[i].x)
		// Y
		fmt.Fprint(writer, "   ")
		fmt.Fprintf(writer, "%.9f", atoms[i].y)
		// Z
		fmt.Fprint(writer, "   ")
		fmt.Fprintf(writer, "%.9f", atoms[i].z)
		fmt.Fprintf(writer, " %d %d %d", atoms[i].dx, atoms[i].dy, atoms[i].dz)
		fmt.Fprintf(writer, " # %s\r\n", atoms[i].aStr)
		cnt += 1
		if cnt == atomCnt {
			cnt = 0
			mId++
		}
	}
}
