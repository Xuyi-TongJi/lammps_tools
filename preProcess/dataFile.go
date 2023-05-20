package preProcess

import (
	. "atom/util"
	"errors"
	"fmt"
	"os"
)

const (
	Npg string = "npg"
	Tme string = "tme"
)

var (
	functions = make(map[string]func(*os.File) error)
	tme       = molecule{c: 5, h: 12, o: 3, count: 20}
	npg       = molecule{c: 5, h: 12, o: 2, count: 19}
)

func init() {
	functions[Npg] = NPG
	functions[Tme] = TME
}

func AddMol(input PreInput) {
	in, err := os.OpenFile(input.Path, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	fun, ext := functions[input.InputType]
	if !ext {
		panic("Invalid input type...")
	}
	if err := fun(in); err != nil {
		panic(err)
	}
}

type atom struct {
	id       int
	mol      int // 分子id
	atomType int
	q        float64
	x        float64
	y        float64
	z        float64
	lx       int
	ly       int
	lz       int
	name     string
}

type molecule struct {
	c     int
	o     int
	h     int
	count int // 原子总数
}

func NPG(in *os.File) error {
	molCnt := 19
	molId := 1
	num := 0
	var (
		output *os.File
		err    error
	)
	output, err = os.OpenFile("./npgAtom.txt", os.O_RDWR, 0666)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		output, _ = os.Create("./npgAtom.txt")
	} else if err != nil {
		return err
	}
	output.Truncate(0)
	output.Seek(0, 0)
	for true {
		curr := atom{}
		n, err := fmt.Fscanf(in, "%d %d %d %f %f %f %f %d %d %d # %s\n",
			&curr.id, &curr.mol, &curr.atomType, &curr.q, &curr.x, &curr.y, &curr.z, &curr.lx, &curr.ly, &curr.lz, &curr.name)
		if n == 0 || err != nil {
			if num != 0 {
				panic("Invalid input data file\n")
			}
			break
		}
		curr.mol = molId
		num += 1
		if num == molCnt {
			num = 0
			molId += 1
		}
		fmt.Fprintf(output, "   %d      %d   %d  %.6f    %.9f     %.9f    %.9f   %d   %d   %d # %s\n",
			curr.id, curr.mol, curr.atomType, curr.q, curr.x, curr.y, curr.z, curr.lx, curr.ly, curr.lz, curr.name)
	}
	return nil
}

// TME
// warning: 原始data文件中，TME的最后一个原子必须为H或C，不能为O
func TME(in *os.File) error {
	var (
		output *os.File
		err    error
	)
	output, err = os.OpenFile("./tmeAtom.txt", os.O_RDWR, 0666)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		output, _ = os.Create("./tmeAtom.txt")
	} else if err != nil {
		return err
	}
	output.Truncate(0)
	output.Seek(0, 0)
	molId := 1
	mole := molecule{0, 0, 0, 0}
	for true {
		curr := atom{}
		n, err := fmt.Fscanf(in, "%d %d %d %f %f %f %f %d %d %d # %s\n",
			&curr.id, &curr.mol, &curr.atomType, &curr.q, &curr.x, &curr.y, &curr.z, &curr.lx, &curr.ly, &curr.lz, &curr.name)
		if n == 0 || err != nil {
			if mole.count != 0 {
				panic("Invalid input data file\n")
			}
			break
		}
		switch curr.name[0] {
		case 'h':
			mole.h += 1
		case 'o':
			mole.o += 1
		case 'c':
			mole.c += 1
		default:
			panic("Invalid input data file\n")
		}
		curr.mol = molId
		mole.count += 1
		// npg
		if mole.count == npg.count {
			if mole.c == npg.c && mole.h == npg.h && mole.o == npg.o {
				mole = molecule{0, 0, 0, 0}
				molId += 1
			}
			// maybe tme
		} else if mole.count == tme.count {
			// tme
			if mole.c == tme.c && mole.h == tme.h && mole.o == tme.o {
				mole = molecule{0, 0, 0, 0}
				molId += 1
			} else {
				panic("Invalid input data file\n")
			}
		}
		// output
		fmt.Fprintf(output, "   %d      %d   %d  %.6f    %.9f     %.9f    %.9f   %d   %d   %d # %s\n",
			curr.id, curr.mol, curr.atomType, curr.q, curr.x, curr.y, curr.z, curr.lx, curr.ly, curr.lz, curr.name)
	}
	return nil
}
