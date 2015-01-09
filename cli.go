package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/awilliams/couchdb-utils/couch"
	"github.com/spf13/cobra"
)

var (
	Context struct {
		server  *string
		verbose *bool
		uri     *bool
	}

	delimiter = "\t"
)

const (
	URI  = "uri"
	View = ""
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "couchdb-utils",
		Short: "couchdb-utils fetches data from a CouchDB server",
		Long:  "couchdb-utils fetches data from a CouchDB server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	Context.server = rootCmd.PersistentFlags().StringP("server", "s", "http://localhost:5984", "CouchDB URL")
	Context.uri = rootCmd.PersistentFlags().BoolP("uri", "", false, "output format, separted by commas")
	Context.verbose = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	databasesCmd := &cobra.Command{
		Use:   "databases",
		Short: "list databases",
		Long:  "List databases",
		Run: func(cmd *cobra.Command, args []string) {
			var output bytes.Buffer
			if err := databaseList(&output, *Context.uri); err != nil {
				log.Fatal(err)
			}
			os.Stdout.Write(output.Bytes())
		},
	}
	replicatorsCmd := &cobra.Command{
		Use:   "replicators",
		Short: "replicators databases",
		Long:  "Replicators databases",
		Run: func(cmd *cobra.Command, args []string) {
			var output bytes.Buffer
			if err := replicatorList(&output, *Context.uri); err != nil {
				log.Fatal(err)
			}
			os.Stdout.Write(output.Bytes())
		},
	}
	replicateCmd := &cobra.Command{
		Use:   "replicate",
		Short: "replicate",
		Long:  "Replicate",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := inputFromArgsOrStdin(args)
			if err != nil {
				log.Fatal(err)
			}
			var output bytes.Buffer
			if err := replicate(input, &output); err != nil {
				log.Fatal(err)
			}
			os.Stdout.Write(output.Bytes())
		},
	}
	/*
		databaseListCmd := &cobra.Command{
			Use: "list",
			Run: func(_ *cobra.Command, _ []string) {
				var output bytes.Buffer
				if err := databaseList(&output, *Context.databaseListURI); err != nil {
					log.Fatal(err)
				}
				os.Stdout.Write(output.Bytes())
			},
		}
		Context.databaseListURI = databaseListCmd.Flags().BoolP("uri", "", false, "print full database URI")

		databaseReplicateCmd := &cobra.Command{
			Use: "replicate",
			Run: func(cmd *cobra.Command, args []string) {
			},
		}

		viewCmd := &cobra.Command{
			Use:   "view [command]",
			Short: "view commands",
			Long:  "View commands",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		viewListCmd := &cobra.Command{
			Use: "list",
			Run: func(cmd *cobra.Command, args []string) {
				var output bytes.Buffer
				if err := viewListAll(&output, *Context.viewListURI); err != nil {
					log.Fatal(err)
				}
				os.Stdout.Write(output.Bytes())
			},
		}
		Context.viewListURI = viewListCmd.Flags().BoolP("uri", "", false, "print full view URI")

		viewFilterCmd := &cobra.Command{
			Use: "filter",
			Run: func(cmd *cobra.Command, args []string) {
				if !stdinHasInput() {
					cmd.Help()
					return
				}
				var output bytes.Buffer
				if err := viewFilter(&output); err != nil {
					log.Fatal(err)
				}
				os.Stdout.Write(output.Bytes())
			},
		}

		designCmd := &cobra.Command{
			Use:   "design [command]",
			Short: "design commands",
			Long:  "Design commands",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		designListCmd := &cobra.Command{
			Use: "list",
			Run: func(_ *cobra.Command, _ []string) {
				var output bytes.Buffer
				if err := designListAll(&output, *Context.designListURI); err != nil {
					log.Fatal(err)
				}
				os.Stdout.Write(output.Bytes())
			},
		}
		Context.designListURI = designListCmd.Flags().BoolP("uri", "", false, "print full design URI")

		databaseCmd.AddCommand(databaseListCmd, databaseReplicateCmd)
		viewCmd.AddCommand(viewListCmd, viewFilterCmd)
		designCmd.AddCommand(designListCmd)
	*/
	rootCmd.AddCommand(databasesCmd, replicatorsCmd, replicateCmd)
	rootCmd.Execute()
}

func databaseList(output *bytes.Buffer, uri bool) error {
	client, err := couch.NewClient(*Context.server)
	if err != nil {
		return err
	}
	dbs, err := couch.AllDBs(client)
	if err != nil {
		return err
	}
	var s string
	for _, db := range dbs {
		if uri {
			s = db.URI()
		} else {
			s = db.String()
		}
		output.WriteString(s + "\n")
	}
	return nil
}

func replicatorList(output *bytes.Buffer, uri bool) error {
	client := newClient()
	replicators, err := couch.ListReplicators(client)
	if err != nil {
		return err
	}
	var s string
	for _, replicator := range replicators {
		if uri {
			s = replicator.URI()
		} else {
			var cont string
			if replicator.Continuous {
				cont = "continuous"
			} else {
				cont = "noncontinuous"
			}
			s = join(replicator.Source, replicator.Target, replicator.ReplicationState, cont, replicator.String())
		}
		output.WriteString(s + "\n")
	}
	return nil
}

func replicate(sources []string, output *bytes.Buffer) error {
	sourceDBs, err := couch.DatabasesFromURLs(sources)
	if err != nil {
		return err
	}
	//_ := newClient()
	for _, db := range sourceDBs {
		output.WriteString(join(db.URI(), db.String()) + "\n")
	}
	return nil
}

func designListAll(output *bytes.Buffer, uri bool) error {
	client, err := couch.NewClient(*Context.server)
	if err != nil {
		return err
	}
	dbs, err := couch.AllDBs(client)
	if err != nil {
		return err
	}
	for _, db := range dbs {
		docs, err := couch.ListDesignDocs(db)
		if err != nil {
			return err
		}
		var s string
		for _, doc := range docs {
			if uri {
				s = doc.URI()
			} else {
				s = join(doc.String(), db.String())
			}
			output.WriteString(s + "\n")
		}
	}
	return nil
}

func viewListAll(output *bytes.Buffer, uri bool) error {
	client, err := couch.NewClient(*Context.server)
	if err != nil {
		return err
	}
	dbs, err := couch.AllDBs(client)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	views := make([]string, len(dbs), len(dbs))
	for i, db := range dbs {
		wg.Add(1)
		go func(db *couch.Database, i int) {
			out, err := viewList(db, uri)
			if err != nil {
				log.Fatal(err)
			}
			views[i] = out
			wg.Done()
		}(db, i)
	}
	wg.Wait()
	output.WriteString(strings.Join(views, ""))
	return nil
}

func viewList(db *couch.Database, uri bool) (string, error) {
	var output string
	docs, err := couch.ListDesignDocs(db)
	if err != nil {
		return output, err
	}
	viewCount := 0
	for _, doc := range docs {
		viewCount += len(doc.Views)
	}
	lines := make([]string, viewCount, viewCount)
	i := 0
	for _, doc := range docs {
		for _, view := range doc.Views {
			if uri {
				lines[i] = view.URI()
			} else {
				lines[i] = join(view.String(), doc.String(), db.String())
			}
			i++
		}
	}
	//sort.Strings(lines)
	if len(lines) > 0 {
		output = strings.Join(lines, "\n") + "\n"
	}
	return output, nil
}

func viewFilter(output *bytes.Buffer) error {
	_, err := couch.NewClient(*Context.server)
	if err != nil {
		return err
	}
	designPath := "_design"
	viewPath := "_view"
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var db, design, view string
		if strings.HasPrefix(scanner.Text(), "http") {
			u, err := url.Parse(scanner.Text())
			if err != nil {
				return err
			}
			designIndex := strings.Index(u.Path, designPath)
			viewIndex := strings.Index(u.Path, viewPath)
			if designIndex == -1 || viewIndex == -1 {
				return fmt.Errorf("unable to parse %s as view URI", scanner.Text())
			}
			db = strings.Trim(u.Path[0:designIndex], "/")
			design = strings.Trim(u.Path[designIndex+len(designPath):viewIndex], "/")
			view = strings.Trim(u.Path[viewIndex+len(viewPath):], "/")

		} else {
			components := strings.Split(scanner.Text(), "\t")
			if len(components) != 3 {
				return fmt.Errorf("unable to parse %s as view", scanner.Text())
			}
			db = components[0]
			design = components[1]
			view = components[2]
		}
		output.WriteString(db + delimiter + design + delimiter + view + "\n")
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
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

func stdinHasInput() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func join(args ...string) string {
	return strings.Join(args, delimiter)
}

func newClient() *couch.Client {
	host := strings.TrimSpace(*Context.server)
	if host == "" {
		log.Fatal("--server cannot be blank")
	}
	client, err := couch.NewClient(host)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
