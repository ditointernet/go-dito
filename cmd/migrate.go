package cmd

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ditointernet/go-dito/lib/env"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/spf13/cobra"
)

// Migration is a structure that encapsulates migration's dependencies
type Migration struct {
	migrate *migrate.Migrate
}

var (
	databaseURI     string
	prodEnvironment bool
	dbDriver        string
	db              *sql.DB
	migration       Migration
)

func newMigration(db *sql.DB, source string) (Migration, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return Migration{}, err
	}

	m, err := migrate.NewWithDatabaseInstance(source, "postgres", driver)
	if err != nil {
		return Migration{}, err
	}

	return Migration{migrate: m}, nil
}

func init() {
	databaseURI = env.GetString("DATABASE_URI")
	prodEnvironment = env.GetBool("PROD_ENVIRONMENT", false)

	if prodEnvironment {
		dbDriver = "cloudsqlpostgres"
	} else {
		dbDriver = "postgres"
	}

	db = newDBConnection(dbDriver, databaseURI)

	var err error
	migration, err = newMigration(db, "file://./migrations")
	if err != nil {
		fmt.Printf("error when creating migration: %s\n", err.Error())
		os.Exit(1)
	}
}

var (
	migrateCmd = &cobra.Command{
		Use:     "migrate",
		Version: "0.0.1",
	}

	upCmd = &cobra.Command{
		Use:  "up",
		Long: "looks at the currently active migration version and will migrate up",
		Run: func(cmd *cobra.Command, args []string) {
			steps, err := cmd.Flags().GetUint("steps")
			if err != nil {
				fmt.Printf("error when fetching flag steps: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("migrating up steps=%v\n", steps)

			err = migration.runMigrationUp(steps)
			if err != nil {
				fmt.Printf("error when migrating database: %v\n", err)
				os.Exit(1)
			}
		},
	}

	downCmd = &cobra.Command{
		Use:  "down",
		Long: "",
		Run: func(cmd *cobra.Command, args []string) {
			steps, err := cmd.Flags().GetUint("steps")
			if err != nil {
				fmt.Printf("error when fetching flag steps: %v\n", err)
				os.Exit(1)
			}

			reader := bufio.NewReader(os.Stdin)

			for {
				fmt.Println("The 'down' operation may results data loss. Do you confirm [y/n]? ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				switch input {
				case "n":
					os.Exit(0)
				case "y":
					fmt.Printf("migrating down steps=%v\n", steps)
					err = migration.runMigrationDown(steps)
					if err != nil {
						fmt.Printf("error when migrating database: %v\n", err)
						os.Exit(1)
					}
				default:
					fmt.Println("Invalid option.")
					os.Exit(0)
				}
			}
		},
	}

	seedCmd = &cobra.Command{
		Use:  "seed",
		Long: "a flag that states if it will run the local database seed",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("seeding the local database...")

			var seedFiles []string
			err := filepath.Walk("./seed", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				if matched, err := filepath.Match("*.sql", filepath.Base(path)); err != nil {
					return err
				} else if matched {
					seedFiles = append(seedFiles, path)
				}
				return nil
			})
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			for _, seedFile := range seedFiles {
				fmt.Printf("running %s...\n", seedFile)
				sql, err := ioutil.ReadFile(seedFile)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				_, err = db.Exec(string(sql))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				fmt.Printf("ran %s\n", seedFile)
			}
		},
	}
)

func newDBConnection(dbDriver string, databaseURI string) *sql.DB {
	db, err := sql.Open(dbDriver, databaseURI)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	db.SetMaxOpenConns(10)

	fmt.Println("got connected to database")

	return db
}

func (c Migration) runMigrationUp(steps uint) error {
	if steps > 0 {
		return c.runMigrationSteps(int(steps))
	}

	version, dirty, _ := c.migrate.Version()
	if dirty {
		return errors.New("database is in dirty state, solve version in 'schema_migrations' table manually")
	}

	if err := c.migrate.Up(); err != nil {
		return c.handleMigrationError(err, version)
	}

	return nil
}

func (c Migration) runMigrationDown(steps uint) error {
	if steps > 0 {
		return c.runMigrationSteps(int(steps) * -1)
	}

	version, dirty, _ := c.migrate.Version()
	if dirty {
		return errors.New("database is in dirty state, solve version in 'schema_migrations' table manually")
	}

	if err := c.migrate.Down(); err != nil {
		return c.handleMigrationError(err, version)
	}

	return nil
}

func (c Migration) runMigrationSteps(steps int) error {
	err := c.migrate.Steps(steps)
	if err != nil {
		return err
	}

	return nil
}

func (c Migration) handleMigrationError(err error, previousVersion uint) error {
	if err.Error() == "no change" {
		return nil
	}
	if previousVersion > 0 {
		e := c.migrate.Force(int(previousVersion))
		if e != nil {
			return fmt.Errorf("%s: %s", err.Error(), e.Error())
		}
		return fmt.Errorf("error on migration, database version was reverted to version %d: %s", previousVersion, err.Error())
	}

	return err
}
