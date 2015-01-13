package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/awilliams/couchdb-utils/couch"
	"github.com/spf13/cobra"
)

// container for global flags
var Context struct {
	server  *string
	uri     *bool
	logHTTP *bool
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "couchdb-utils",
		Short: "couchdb-utils",
		Long:  "couchdb-utils - https://github.com/awilliams/couchdb-utils",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	// global flags
	Context.server = rootCmd.PersistentFlags().StringP("server", "s", "http://localhost:5984", "CouchDB URL")
	Context.uri = rootCmd.PersistentFlags().BoolP("uri", "", false, "output format, separted by commas")
	Context.logHTTP = rootCmd.PersistentFlags().BoolP("log", "l", false, "log all http requests and responses to stderr")

	listDatabasesCmd := &cobra.Command{
		Use:   "databases",
		Short: "List databases",
		Long:  "List databases",
		Run: func(_ *cobra.Command, _ []string) {
			run(func(client *couch.Client, output *Output) error {
				return listDatabases(client, output, *Context.uri)
			})
		},
	}

	listReplicatorsCmd := &cobra.Command{
		Use:   "replicators",
		Short: "List replicators",
		Long:  "List replicators",
		Run: func(_ *cobra.Command, _ []string) {
			run(func(client *couch.Client, output *Output) error {
				return listReplicators(client, output, *Context.uri)
			})
		},
	}

	listViewsCmd := &cobra.Command{
		Use:   "views",
		Short: "List views",
		Long:  "List views",
		Run: func(_ *cobra.Command, args []string) {
			input, err := inputFromArgsOrStdin(args)
			checkErr(err)
			run(func(client *couch.Client, output *Output) error {
				return listViews(client, output, input, *Context.uri)
			})
		},
	}

	delDocsCmd := &cobra.Command{
		Use:   "deldocs",
		Short: "Delete documents",
		Long:  "Delete documents. Reads from STDIN or args",
		Run: func(_ *cobra.Command, args []string) {
			input, err := inputFromArgsOrStdin(args)
			checkErr(err)
			run(func(client *couch.Client, output *Output) error {
				return deleteDocs(client, output, input)
			})
		},
	}

	var (
		pullContinuous bool
		createTarget   bool
	)
	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Start pull replication",
		Long:  "Start pull replication",
		Run: func(_ *cobra.Command, args []string) {
			input, err := inputFromArgsOrStdin(args)
			checkErr(err)
			run(func(client *couch.Client, output *Output) error {
				return pull(client, output, input, pullContinuous, createTarget)
			})
		},
	}
	pullCmd.Flags().BoolVarP(&pullContinuous, "continuous", "", true, "set continuous")
	pullCmd.Flags().BoolVarP(&createTarget, "create-target", "", true, "set createTarget")

	refreshViewsCmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh views",
		Long:  "Refresh views",
		Run: func(_ *cobra.Command, args []string) {
			input, err := inputFromArgsOrStdin(args)
			checkErr(err)
			run(func(client *couch.Client, output *Output) error {
				return refreshViews(client, output, input, *Context.uri)
			})
		},
	}

	rootCmd.AddCommand(
		listDatabasesCmd,
		listReplicatorsCmd,
		listViewsCmd,
		delDocsCmd,
		pullCmd,
		refreshViewsCmd,
	)
	rootCmd.Execute()
}

// run wraps a function by creating the couch client and output, calling function, then printing output
func run(f func(*couch.Client, *Output) error) {
	c, err := newClient()
	checkErr(err)
	output := &Output{buf: new(bytes.Buffer)}
	checkErr(f(c, output))
	output.Flush()
}

type Output struct {
	buf *bytes.Buffer
}

// delimiter for tabular output
const delimiter = "\t"

func (o *Output) Println(line ...string) {
	o.buf.WriteString(strings.Join(line, delimiter) + "\n")
}

func (o *Output) Flush() {
	if o.buf.Len() == 0 {
		return
	}
	_, err := os.Stdout.Write(o.buf.Bytes())
	checkErr(err)
}

func newClient() (*couch.Client, error) {
	host := strings.TrimSpace(*Context.server)
	if host == "" {
		return nil, fmt.Errorf("--server cannot be blank")
	}
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}
	client, err := couch.NewClient(host)
	if err != nil {
		return nil, err
	}
	client.LogHTTP = *Context.logHTTP
	return client, nil
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func inputFromArgsOrStdin(args []string) ([]string, error) {
	if len(args) == 1 && args[0] == "-" {
		return stdinLines()
	}
	return args, nil
}

func stdinLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	return lines, err
}
