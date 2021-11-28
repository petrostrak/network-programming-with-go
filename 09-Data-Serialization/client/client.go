package client

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/petrostrak/network-programming-with-go/09-Data-Serialization/housework/v1"
)

var (
	addr, caCertFn string
)

func init() {
	flag.StringVar(&addr, "address", "localhost:34443", "server address")
	flag.StringVar(&caCertFn, "ca-cert", "cert.pem", "CA certificate")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			`Usage: %s [flags] [add chore, ...|complete #]
	add			add coma-separated chores
	complete	complete designated chore
	
Flags:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func list(ctx context.Context, client housework.RobotMaidClient) error {
	chores, err := client.List(ctx, new(housework.Empty))
	if err != nil {
		return err
	}

	if len(chores.Chores) == 0 {
		fmt.Println("You have nothing to do!")
		return err
	}

	fmt.Println("#\t[X]\tDescription")
	for i, chore := range chores.Chores {
		c := " "
		if chore.Complete {
			c = "X"
		}
		fmt.Printf("%d\t[%s]\t%s\n", i+1, c, chore.Description)
	}

	return nil
}

func add(ctx context.Context, client housework.RobotMaidClient, s string) error {
	chores := new(housework.Chores)

	for _, chore := range strings.Split(s, ",") {
		if desc := strings.TrimSpace(chore); desc != "" {
			chores.Chores = append(chores.Chores, &housework.Chore{
				Description: desc,
			})
		}
	}

	var err error
	if len(chores.Chores) > 0 {
		_, err = client.Add(ctx, chores)
	}

	return err
}

func complete(ctx context.Context, client housework.RobotMaidClient, s string) error {
	i, err := strconv.Atoi(s)
	if err == nil {
		_, err = client.Complete(ctx, &housework.CompleteRequest{ChoreNumber: int32(i)})
	}

	return err
}
