package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"github.com/cheggaaa/pb/v3"
	"time"
)

var (
	log          = logrus.New()
	dict, boards string
)

type Coordinates struct {
	X []string
	Y []string
}

//Stat represents statistics of each word
type WordStat struct {
	ID    string
	Text  string
	Clue  string
	Len   int
	Board string
	Coordinates
}

func (w WordStat) String() string {
	return fmt.Sprintf("ID: %s | Text: %s | Length: %d | Clue: %s | X: %s | y: %s", w.ID, w.Text, w.Len, w.Clue, w.X, w.Y)
}

func main() {
	var wg sync.WaitGroup
	start := time.Now()

	initFlags()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	/*files, err := ioutil.ReadDir(boards)
	if err != nil {
		log.Fatalf("Can't read the given directory: %s", err)
	}*/

	files, err := listOfBoards(boards)
	if err != nil {
		log.Fatal(err)
	}
	//log.Info(files1)

	log.Infof("%d boards found in the given directory", len(files))

	//New progress bar for counting processed files
	bar := pb.StartNew(len(files))

	for _, f := range files {
		//We process only xml files
		if filepath.Ext(f.Name()) != ".xml" {
			bar.Increment()
			log.Infof("Skipping incorrect file %s", f.Name())
			continue
		}
		wg.Add(1)
		fname := boards + string(os.PathSeparator) + f.Name()

		go worker(fname, &wg, bar)
	}

	wg.Wait()
	bar.Finish()

	err = report(db)
	if err != nil {
		log.Fatal(err)
	}

	//Close and remove database file
	err = cleanupDB(db)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()
	total := end.Sub(start)

	log.Infof("Processing took %f seconds", total.Seconds())
	log.Infof("All files are processed. Please, check %s file for report", repFile)
}

func worker(fname string, wg *sync.WaitGroup, bar *pb.ProgressBar) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Error(err)
	}
	if debug {
		log.Debugf("Opened file %s", fname)
	}

	b := &CrosswordCompiler{}
	err = xml.Unmarshal(data, b)
	if err != nil {
		log.Fatal(err)
	}
	if debug {
		log.Debugf("Unmarshalled XML file file %s", fname)
	}

	clues := getClues(b.RectangularPuzzle.Crossword.Grid.Cells)

	for _, w := range b.RectangularPuzzle.Crossword.Words {
		word := WordStat{}

		word.ID = w.ID
		word.Board = fname

		if clue, ok := clues[word.ID]; ok {
			word.Clue = clue
		}

		length, err := calcLength(w.X, w.Y)
		if err != nil {
			log.Fatal("Can't calculate length of the word")
		}
		if debug {
			log.Debugf("Length of the word %s is calculated: %d", w.Text, length)
		}

		word.Len = length

		coord, err := getCoordinates(w.X, w.Y)
		if err != nil {
			log.Fatal(err)
		}
		if debug {
			log.Debugf("Coordinated of the word %s calculated: %s", w.Text, coord)
		}
		word.Coordinates = coord

		text, err := findWords(b.RectangularPuzzle.Crossword.Grid.Cells, coord)
		if err != nil {
			log.Fatal(err)
		}

		word.Text = text

		if debug {
			log.Info(word)
		}

		//Save word to the database
		err = storeWord(db, word)
		if err != nil {
			log.Fatal(err)
		}
	}
	bar.Increment()
	wg.Done()
}

func report(db *bolt.DB) (err error) {
	rep, err := os.OpenFile(repFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer rep.Close()
	wr := csv.NewWriter(rep)

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("Words"))

		b.ForEach(func(k, v []byte) error {
			var w []WordStat
			err := json.Unmarshal(v, &w)
			if err != nil {
				return err
			}
			if len(w) > 0 {
				var raw []string
				var files []string
				for _, f := range w {
					files = append(files, f.Board)
				}
				if debug {
					log.Infof("Word: %s | Count: %d | Boards: %s", string(k), len(files), files)
				}
				raw = append(raw, string(k))
				raw = append(raw, strconv.Itoa(len(files)))
				raw = append(raw, strings.Join(files, " "))

				err = wr.Write(raw)
				if err != nil {
					log.Fatal(err)
				}
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return err
	}
	//Flush all changes to the io.Writer
	wr.Flush()

	return nil
}

func calcLength(x, y string) (val int, err error) {
	var split = func(s string) (int, error) {
		coord := strings.Split(x, "-")
		if len(coord) == 1 {
			return 1, nil
		}

		a, err := strconv.Atoi(coord[0])
		if err != nil {
			return 0, err
		}
		b, err := strconv.Atoi(coord[1])
		if err != nil {
			return 0, err
		}

		l := (b - a) + 1
		return l, nil
	}

	if len(x) == 1 && len(y) == 1 {
		return 1, nil
	}

	if len(x) > 1 {
		l, err := split(x)
		if err != nil {
			return 0, err
		}
		return l, nil
	}

	if len(y) > 1 {
		l, err := split(x)
		if err != nil {
			return 0, err
		}
		return l, nil
	}

	return 0, nil
}

//getCoordinates gets the coordinates of the word on the board
func getCoordinates(x, y string) (coord Coordinates, err error) {
	var xList, yList, xN, yN []string

	xList = strings.Split(x, "-")
	if strings.Contains(x, "-") && len(xList) == 1 {
		return coord, errors.New("can't extract X coordinates")
	}

	if len(xList) > 1 {
		xN, err = makeRange(xList[0], xList[1])
		if err != nil {
			return coord, errors.New("can't extract X coordinates")
		}
	} else {
		xN = xList
	}

	yList = strings.Split(y, "-")
	if strings.Contains(y, "-") && len(yList) == 1 {
		return coord, errors.New("can't extract Y coordinates")
	}

	if len(yList) > 1 {
		yN, err = makeRange(yList[0], yList[1])
		if err != nil {
			return coord, errors.New("can't extract Y coordinates")
		}
	} else {
		yN = yList
	}

	coord.X = xN
	coord.Y = yN

	return coord, nil
}

func findWords(cells []Cell, coord Coordinates) (text string, err error) {
	var word []string

	var letters [][]string

	//log.Info(coord)
	//log.Info(letters)

	if len(coord.X) == 1 && len(coord.Y) == 1 {
		letters = append(letters, []string{coord.X[0], coord.Y[0]})
	}

	if len(coord.X) > 1 {
		//log.Info("Horizontal")
		for _, v := range coord.X {
			letters = append(letters, []string{v, coord.Y[0]})
		}
	}

	if len(coord.Y) > 1 {
		//log.Info("Vertical")
		for _, v := range coord.Y {
			letters = append(letters, []string{coord.X[0], v})
		}
	}

	for _, c := range cells {
		if c.Type == "clue" {
			continue
		}

		for _, l := range letters {
			if c.X == l[0] && c.Y == l[1] {
				word = append(word, c.Solution)
			}
		}
	}

	text = strings.Join(word, "")

	return text, nil
}

func listOfBoards(path string) (boards []os.FileInfo, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			boards = append(boards, info)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}

	return boards, nil
}

func getClues(cells []Cell) (clues map[string]string) {
	clues = make(map[string]string)
	for _, c := range cells {
		if c.Type == "clue" {
			for _, cl := range c.Clue {
				clues[cl.Word] = cl.Text
			}
		}
	}

	if debug {
		log.Debugf("Clues are processed")
	}

	return clues
}

func makeRange(min, max string) (ran []string, err error) {
	var mi, ma int
	mi, err = strconv.Atoi(min)
	if err != nil {
		return nil, err
	}

	ma, err = strconv.Atoi(max)
	if err != nil {
		return nil, err
	}

	a := make([]int, ma-mi+1)
	for i := range a {
		a[i] = mi + i
	}

	for v := range a {
		ran = append(ran, strconv.Itoa(a[v]))
	}

	return ran, nil
}
