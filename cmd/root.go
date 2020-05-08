package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ms2my",
	Short: "A tool for transfer sql server database to mysql execute file",
	Long: `A tool for transfer mssql to mysql,but you should custom pk,
and not support binrary`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if password == "" {
			fmt.Println("no such flag --password")
			return
		}
		if dbname == "" {
			fmt.Println("no such flag --dbname")
			return
		}
		if tablename == "" {
			fmt.Println("no such flag --tablename")
			return
		}
		var err error
		mssql, err = NewMssql(username, dbname, ip, password, port)
		if err != nil {
			fmt.Println(err)
			return
		}
		tables := strings.Split(tablename, ",")
		err = os.MkdirAll("mysql", 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, v := range tables {
			f, err := os.Create("mysql/" + v + ".sql")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			createsql, err := generateCreate(v)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = f.WriteString(createsql)
			if err != nil {
				fmt.Println(err)
				return
			}
			//f.WriteString("truncate table " + strings.ToLower(tablename) + ";\n")
			insertsql, err := generateInsert(v)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = f.WriteString(insertsql)
			if err != nil {
				fmt.Println(err)
				return
			}
			_ = f.Sync()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	cfgFile   string
	port      int
	tablename string
	username  string
	password  string
	ip        string
	dbname    string
)

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().IntVar(&port, "port", 1433, "mssql port")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "sa", "mssql username")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "mssql password")
	rootCmd.PersistentFlags().StringVarP(&tablename, "tablename", "t", "", "table names,you yan use , join mult")
	rootCmd.PersistentFlags().StringVarP(&ip, "ip", "i", "127.0.0.1", "ip")
	rootCmd.PersistentFlags().StringVarP(&dbname, "dbname", "d", "", "dbname")

	//_ = viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	//_ = viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	//_ = viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	//_ = viper.BindPFlag("tablename", rootCmd.PersistentFlags().Lookup("tablename"))
	//_ = viper.BindPFlag("ip", rootCmd.PersistentFlags().Lookup("ip"))
	//_ = viper.BindPFlag("dbname", rootCmd.PersistentFlags().Lookup("dbname"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		//home, err := homedir.Dir()
		//if err != nil {
		//  fmt.Println(err)
		//  os.Exit(1)
		//}

		// Search config in home directory with name ".ms2my" (without extension).
		//viper.AddConfigPath("")
		viper.AddConfigPath(".")
		viper.SetConfigName("config.yaml")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
