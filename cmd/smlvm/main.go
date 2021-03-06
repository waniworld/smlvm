package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"

	"shanhu.io/smlvm/arch"
	"shanhu.io/smlvm/dasm"
	"shanhu.io/smlvm/debug"
	"shanhu.io/smlvm/image"
)

func run(
	img []byte, conf *arch.Config, ncycle int, printStatus bool,
) (int, error) {

	// create a single core machine
	m := arch.NewMachine(conf)
	secs, err := image.Read(bytes.NewReader(img))
	if err != nil {
		return 0, err
	}

	if err := m.LoadSections(secs); err != nil {
		return 0, err
	}

	ret, exp := m.Run(ncycle)
	if printStatus {
		m.PrintCoreStatus()
	}

	if exp != nil && !arch.IsHalt(exp) {
		fmt.Println(exp)
		err := arch.FprintStack(os.Stdout, m, exp)
		if err != nil {
			log.Fatal(err)
		}
	}

	if exp == nil {
		return ret, nil
	}
	return ret, exp
}

func main() {
	doDasm := flag.Bool("d", false, "do dump")
	printDebug := flag.Bool("debug", false, "print debug symbols")
	ncycle := flag.Int("n", 100000, "max cycles to execute")
	memSize := flag.Int("m", 0, "memory size; 0 for full 4GB")
	printStatus := flag.Bool("s", false, "print status after execution")
	bootArg := flag.Uint("arg", 0, "boot argument, a uint32 number")
	romRoot := flag.String("rom", "", "rom root path")
	randSeed := flag.Int64("seed", 0, "random seed, 0 for using the time")
	initSP := flag.Int64("initsp", 0, "init stack pointer")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("need exactly one input file\n")
	}

	fname := args[0]

	if *doDasm {
		f, err := os.Open(fname)
		defer f.Close()

		err = dasm.DumpImage(f, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	} else if *printDebug {
		f, err := os.Open(fname)
		defer f.Close()

		secs, err := image.Read(f)
		if err != nil {
			log.Fatal(err)
		}

		for _, sec := range secs {
			if sec.Type != image.Debug {
				continue
			}

			tab, err := debug.UnmarshalTable(sec.Bytes)
			if err != nil {
				log.Fatal(err)
			}

			tab.PrintTo(os.Stdout)
		}
	} else {
		bs, err := ioutil.ReadFile(fname)
		if err != nil {
			log.Fatal(err)
		}

		if *bootArg > math.MaxUint32 {
			log.Fatalf("boot arg(%d) is too large", *bootArg)
		}
		conf := &arch.Config{
			MemSize:  uint32(*memSize),
			ROM:      *romRoot,
			RandSeed: *randSeed,
			BootArg:  uint32(*bootArg),
			InitSP:   uint32(*initSP),
		}

		n, e := run(bs, conf, *ncycle, *printStatus)
		fmt.Printf("(%d cycles)\n", n)
		if e != nil {
			if !arch.IsHalt(e) {
				fmt.Println(e)
			}
		} else {
			fmt.Println("(end of time)")
		}
	}
}
