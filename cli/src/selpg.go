package main

import(
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type selpgArgs struct{
	startPage int
	endPage int
	inFilename string
	pageLen int
	pageType int
	printDest string
}

var progname string 
var maxSize = 1 << 32 - 1

func main(){
	var sa selpgArgs
	ac := len(os.Args)
	av := os.Args
	progname = os.Args[0]
	
	sa.startPage = -1
	sa.endPage = -1
	sa.inFilename = ""
	sa.pageLen = 72
	sa.pageType = 'l'
	sa.printDest = ""

	processArgs(ac, av, &sa)
	processInput(sa)
}

func processArgs(ac int, av []string, psa *selpgArgs){
	var s1 string
	var s2 string 
	var argno int 
	var i int

	if ac < 3{
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
		usage()
		os.Exit(1)
	}
	
	s1 = av[1] 
	if s1[0] != '-' || s1[1] != 's'{
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstartPage\n", progname)
		usage()
		os.Exit(2)
	}
	
	i,_ = strconv.Atoi(s1[2:])
	if  i < 1 || i > maxSize{
		fmt.Fprintf(os.Stderr, "%s: invalid start page %s\n", progname, s1[2:])
		usage()
		os.Exit(3)
	}
	psa.startPage = i

	s1 = av[2] 
	if s1[0] != '-' || s1[1] != 'e'{
		fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -eend_page\n", progname)
		usage()
		os.Exit(4)
	}
	i,_ = strconv.Atoi(s1[2:])
	if i < 1 || i > maxSize || i < psa.startPage {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %s\n", progname, s1[2:])
		usage()
		os.Exit(5)
	}
	psa.endPage = i

	argno = 3
	for argno <= ac - 1 && av[argno][0] == '-'{
		s1 = av[argno] 
		switch s1[1]{
			case 'l':
				s2 = s1[2:] 
				i,_ = strconv.Atoi(s2)
				if  i < 1 || i > maxSize{
					fmt.Fprintf(os.Stderr, "%s: invalid page length %s\n", progname, s2)
					usage()
					os.Exit(6)
				}
				psa.pageLen = i
				argno++
			case 'f':
				
				if s1[0] != '-' || s1[1] != 'f'{
					fmt.Fprintf(os.Stderr, "%s: option should be \"-f\"\n", progname)
					usage()
					os.Exit(7)
				}
				psa.pageType = 'f'
				argno++
			case 'd':
				s2 = s1[2:]
				
				if len(s2) < 1{
					fmt.Fprintf(os.Stderr,
					"%s: -d option requires a printer destination\n", progname)
					usage()
					os.Exit(8)
				}
				psa.printDest = s2
				argno++
			default:
				fmt.Fprintf(os.Stderr, "%s: unknown option %s\n", progname, s1)
				usage()
				os.Exit(9)
		} 
	} 

	if argno <= ac - 1{
		psa.inFilename = av[argno]
		
		_, err := os.Stat(psa.inFilename)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: input file \"%s\" does not exist\n", progname, psa.inFilename)
			os.Exit(10)
		}
	}
}

func processInput(sa selpgArgs) {
	var fin *os.File
	var fout *os.File 
	var s1 string
	var lineCtr int
	var pageCtr int

	if sa.inFilename == "" {
		fin = os.Stdin;
	} else{
		fin,_ = os.Open(sa.inFilename)
		if fin == nil {
			fmt.Fprintf(os.Stderr, "%s: could not open input file \"%s\"\n",
			progname, sa.inFilename)
			os.Exit(11)
		}
	}
	
	if sa.printDest == ""{
		fout = os.Stdout
	} else{
		fin,_ = os.Open(s1)
		if fin == nil {
			fmt.Fprintf(os.Stderr, "%s: could not open pipe to \"%s\"\n",
			progname, s1)
			os.Exit(12)
		}
	}
    
    var result1 string
	
	if (sa.pageType == 'l'){
		lineCtr = 0
		pageCtr = 1
		reader := bufio.NewReader(fin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			lineCtr++
			if lineCtr > sa.pageLen {
				pageCtr++
				lineCtr = 1
			}
			if pageCtr >= sa.startPage && pageCtr <= sa.endPage {
				if sa.printDest == "" {
					fmt.Fprintf(fout, "%s", line)
				} else {
					result1 += line
				}
			}
		}
	} else{
		pageCtr = 1
		reader := bufio.NewReader(fin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil { 
				break
			}
			for _, v := range line{
				if v == '\f' {
					pageCtr ++
				}
				if pageCtr >= sa.startPage && pageCtr <= sa.endPage{
					if sa.printDest == "" {
						fmt.Fprintf(fout, "%s", line)
					} else {
						result1 += line
					}
				}
			}
		}
	}

	if pageCtr < sa.startPage {
		fmt.Fprintf(os.Stderr,
		"%s: startPage (%d) greater than total pages (%d), no output written\n", progname, sa.startPage, pageCtr)
	}else if pageCtr < sa.endPage{
		fmt.Fprintf(os.Stderr, "%s: end_page (%d) greater than total pages (%d), less output than expected\n", progname, sa.endPage, pageCtr)
	}
	fin.Close()
	fout.Close()
}

func usage(){
	fmt.Fprintf(os.Stderr,
	"\nUSAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ][ -ddest ] [ in_filename ]\n", progname)
}

