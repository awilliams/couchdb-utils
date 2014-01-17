package main

import (
	"fmt"
	"github.com/awilliams/cobra"
	"github.com/awilliams/couchdb-utils/api"
	"github.com/awilliams/couchdb-utils/util"
	"os"
)

func checkError(err error) {
	if err != nil {
		util.PrintError(err)
		os.Exit(1)
	}
}

func handleResult(result *api.Result) {
	if GlobalConfig.Debug {
		util.PrettyPrint(*result)
	}
}

var couchdb *api.Couchdb

func Couchdb() *api.Couchdb {
	c := couchdb
	if c == nil {
		c, err := api.New(GlobalConfig.Host)
		checkError(err)
		c.ResultHandler = handleResult
		couchdb = c
	}
	return couchdb
}

func parseDatabases(args []string) api.Databases {
	if len(args) < 1 {
		databases, err := Couchdb().GetDatabases()
		checkError(err)
		return databases
	}
	dbs := make(api.Databases, len(args))
	for i, s := range args {
		dbs[i] = api.Database{&s}
	}
	return dbs
}

var GlobalConfig = struct {
	Host    string
	Verbose bool
	Debug   bool
}{
	"",
	false,
	false,
}

var cli = &cobra.Command{
	Use:   CouchdbUtils.Name,
	Short: CouchdbUtils.Name,
	Long:  CouchdbUtils.Description + " See help below for more information.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number of " + CouchdbUtils.Name,
	Run: func(cmd *cobra.Command, args []string) {
		util.PrettyPrint(CouchdbUtils)
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Print basic server info",
	Long:  "Print basic server info. See help for more options.\nhttp://docs.couchdb.org/en/latest/api/misc.html#get",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := Couchdb().GetServer()
		checkError(err)
		util.PrettyPrint(server)
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats [(<part1> <part2>)]",
	Short: "Print server stats (optionally only a certain section eg: couchdb request_time).",
	Long:  "Print server stats (optionally only a certain section eg: couchdb request_time). See help for more options.\nhttp://docs.couchdb.org/en/latest/api/server/common.html#stats",
	Run: func(cmd *cobra.Command, args []string) {
		var sectionA string
		var sectionB string
		if len(args) != 0 && len(args) != 2 {
			checkError(fmt.Errorf("Must provide 2 parts of stat\nSee http://docs.couchdb.org/en/latest/api/server/common.html?highlight=stats#stats"))
		}
		if len(args) == 2 {
			sectionA = args[0]
			sectionB = args[1]
		}
		stats, err := Couchdb().GetStats(sectionA, sectionB)
		checkError(err)
		util.PrettyPrint(stats)
	},
}

var activeTasksCmd = &cobra.Command{
	Use:   "activetasks [<type>]",
	Short: "Print active tasks (optionally filtering by type)",
	Long:  "Print active tasks (optionally filtering by type). Types include:\n 'indexer'\n 'replication'\nSee help for more options.\nhttp://docs.couchdb.org/en/latest/api/server/common.html#active-tasks",
	Run: func(cmd *cobra.Command, args []string) {
		var taskFilter string
		if len(args) > 0 {
			taskFilter = args[0]
		}

		activeTasks, err := Couchdb().GetActiveTasks()
		checkError(err)
		if taskFilter != "" {
			util.PrettyPrint(activeTasks.ByType(taskFilter))
		} else {
			util.PrettyPrint(activeTasks)
		}
	},
}

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Print information about authenticated user",
	Long:  "Print information about authenticated user. See help for more options.\nhttp://docs.couchdb.org/en/latest/api/server/authn.html?highlight=session#get--_session",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := Couchdb().GetSession()
		checkError(err)
		util.PrettyPrint(session)
	},
}

var databaseListCmd = &cobra.Command{
	Use:   "databases",
	Short: "Print all databases",
	Long:  "Print all databases. See help for more options.\nhttp://docs.couchdb.org/en/latest/api/server/common.html#all-dbs",
	Run: func(cmd *cobra.Command, args []string) {
		databases, err := Couchdb().GetDatabases()
		checkError(err)
		util.PrettyPrint(databases)
	},
}

var databaseListViewsCmd = &cobra.Command{
	Use:   "views [<db>...]",
	Short: "Print all views (optionally filtering by database(s))",
	Long:  "Print all views (optionally filtering by database(s)). See help for more options.",
	Run: func(cmd *cobra.Command, args []string) {
		dbs := parseDatabases(args)
		for _, db := range dbs {
			views, err := Couchdb().GetViews(db)
			checkError(err)
			util.PrettyPrint(views)
		}
	},
}

var databaseRefreshViewsCmd = &cobra.Command{
	Use:   "refreshviews [<db>...] [--verbose]",
	Short: "Refresh views (optionally filtering by database(s))",
	Long:  "Refresh all views (optionally filtering by database(s)).\nThis is done by requesting a random view from each design doc with stale=update_after. If verbose, the command will print out the views which were requested.\nSee help for more options.",
	Run: func(cmd *cobra.Command, args []string) {
		dbs := parseDatabases(args)
		for _, db := range dbs {
			views, err := Couchdb().GetViews(db)
			checkError(err)
			refreshedViews, errors := Couchdb().RefreshViews(views)
			if len(errors) != 0 {
				for _, err := range errors {
					util.PrintError(err)
				}
				os.Exit(1)
			}
			if GlobalConfig.Verbose && len(refreshedViews) > 0 {
				util.PrettyPrint(db, refreshedViews)
			}
		}
	},
}

var replicatorBaseCmd = &cobra.Command{
	Use:   "rep <command>...",
	Short: "Replication subcommands",
	Long:  "Replication subcommands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Must provide a subcommand")
		cmd.Usage()
	},
}

var replicatorsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all replicators",
	Run: func(cmd *cobra.Command, args []string) {
		replicators, err := Couchdb().GetReplicators()
		checkError(err)
		util.PrettyPrint(*replicators)
	},
}

var deleteReplicatorConf struct {
	All bool
}
var deleteReplicatorCmd = &cobra.Command{
	Use:   "stop (<id>... | --all) [--verbose]",
	Short: "Stop replicating given id(s) or all",
	Run: func(cmd *cobra.Command, args []string) {
		if deleteReplicatorConf.All {
			replicators, err := Couchdb().DeleteAllReplicators()
			checkError(err)
			if GlobalConfig.Verbose {
				fmt.Printf("Stopped %d replicators\n", len(*replicators))
				util.PrettyPrint(replicators)
			}
			return
		}
		if len(args) == 0 {
			checkError(fmt.Errorf("Must provide at least 1 replicator id or use --all to delete all"))
		}
		for _, id := range args {
			err := Couchdb().DeleteReplicator(id)
			checkError(err)
		}
	},
}

var replicateConf api.ReplicationConfig
var replicateCmd = &cobra.Command{
	Use:   "start <source> <target> [--create --continuous]",
	Short: "Configure replication from source to target",
	Long:  "Configure replication from source to target. See help for more options. Verbose option will display response.\nhttp://docs.couchdb.org/en/latest/api/misc.html#post-replicate",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			checkError(fmt.Errorf("Must provide source and target."))
		} else {
			replicateConf.Source = args[0]
			replicateConf.Target = args[1]
		}
		session, err := Couchdb().GetSession()
		checkError(err)
		replicateConf.UserCtx = session.UserCtx
		err = Couchdb().Replicate(replicateConf)
		checkError(err)
	},
}

var replicateHostConf api.ReplicationConfig
var replicateHostCmd = &cobra.Command{
	Use:   "host <remote_host> [--create --continuous --verbose]",
	Short: "Replicates all databases in remote host",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			checkError(fmt.Errorf("Must provide remote database"))
		}
		remoteCouch, err := api.New(args[0])
		remoteCouch.ResultHandler = handleResult
		checkError(err)
		databases, err := Couchdb().ReplicateHost(remoteCouch, replicateHostConf)
		checkError(err)
		if GlobalConfig.Verbose {
			fmt.Printf("Replicated %d databases\n", len(*databases))
			util.PrettyPrint(databases)
		}
	},
}

func executeCli() {
	cli.PersistentFlags().StringVarP(&GlobalConfig.Host, "host", "h", "http://localhost:5984", "Couchdb server url (http://user:password@host:port)")
	cli.PersistentFlags().BoolVarP(&GlobalConfig.Verbose, "verbose", "v", false, "chatty output")
	cli.PersistentFlags().BoolVarP(&GlobalConfig.Debug, "debug", "d", false, "print http requests")

	replicateCmd.Flags().BoolVarP(&replicateConf.Cancel, "delete", "", false, "cancel replication")
	replicateCmd.Flags().BoolVarP(&replicateConf.CreateTarget, "create", "", true, "create target database if doesn't exist")
	replicateCmd.Flags().BoolVarP(&replicateConf.Continuous, "continuous", "", false, "make the replication continuous")
	replicateCmd.Flags().StringVarP(&replicateConf.ID, "id", "", "", "replicator id, required if persistent")

	replicateHostCmd.Flags().BoolVarP(&replicateHostConf.CreateTarget, "create", "", true, "create target database if doesn't exist")
	replicateHostCmd.Flags().BoolVarP(&replicateHostConf.Continuous, "continuous", "", true, "make the replication continuous")

	deleteReplicatorCmd.Flags().BoolVarP(&deleteReplicatorConf.All, "all", "", false, "delete all replicators")

	replicatorBaseCmd.AddCommand(replicatorsListCmd, replicateCmd, deleteReplicatorCmd, replicateHostCmd)
	cli.AddCommand(versionCmd, serverCmd, statsCmd, activeTasksCmd, sessionCmd, databaseListCmd, databaseListViewsCmd, databaseRefreshViewsCmd, replicatorBaseCmd)

	cli.Execute()
}
