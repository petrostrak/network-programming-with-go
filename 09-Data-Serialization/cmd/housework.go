package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/petrostrak/network-programming-with-go/09-Data-Serialization/housework"
	storage "github.com/petrostrak/network-programming-with-go/09-Data-Serialization/json"
)

var (
	datafile string
)

func init() {
	flag.StringVar(&datafile, " file", "housework.db", "data file")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			`Usage: %s [flags] [add chore, ...|complete #]
	add      add comma-separater chores
	complete complete designated chore
	
Flags:
	`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func load() ([]*housework.Chore, error) {
	if _, err := os.Stat(datafile); os.IsNotExist(err) {
		return make([]*housework.Chore, 0), nil
	}

	df, err := os.Open(datafile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := df.Close(); err != nil {
			fmt.Printf("closing data file: %v", err)
		}
	}()

	return storage.Load(df)
}

func flush(chores []*housework.Chore) error {
	df, err := os.Create(datafile)
	if err != nil {
		return err
	}
	defer func() {
		if err := df.Close(); err != nil {
			fmt.Printf("closing data file: %v", err)
		}
	}()

	return storage.Flush(df, chores)
}

func list() error {
	chores, err := load()
	if err != nil {
		return err
	}

	if len(chores) == 0 {
		fmt.Println("You're all caught up!")
		return nil
	}

	fmt.Println("#\t[X]\tDescription")
	for i, chore := range chores {
		c := " "
		if chore.Complete {
			c = "X"
		}
		fmt.Printf("%d\t[%s]\t%s\n", i+1, c, chore.Description)
	}

	return nil
}

func add(s string) error {
	chores, err := load()
	if err != nil {
		return err
	}

	for _, chore := range strings.Split(s, ",") {
		if desc := strings.TrimSpace(chore); desc != "" {
			chores = append(chores, &housework.Chore{
				Description: desc,
			})
		}
	}

	return flush(chores)
}

func main() {

}
