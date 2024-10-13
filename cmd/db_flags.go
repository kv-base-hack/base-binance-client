package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	postgresHostFlag     = "postgres-host"
	postgresPortFlag     = "postgres-port"
	postgresUserFlag     = "postgres-user"
	postgresPasswordFlag = "postgres-password"
	postgresDatabaseFlag = "postgres-database"
)

// NewPostgreSQLFlags creates new cli flags for PostgreSQL client.
func NewPostgreSQLFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    postgresHostFlag,
			EnvVars: []string{"POSTGRES_HOST"},
		},
		&cli.StringFlag{
			Name:    postgresUserFlag,
			EnvVars: []string{"POSTGRES_USER"},
		},
		&cli.StringFlag{
			Name:    postgresPasswordFlag,
			EnvVars: []string{"POSTGRES_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    postgresDatabaseFlag,
			EnvVars: []string{"POSTGRES_DB"},
		},
		&cli.Int64Flag{
			Name:    postgresPortFlag,
			EnvVars: []string{"POSTGRES_PORT"},
		},
	}
}

// NewDBFromContext creates a DB instance from cli flags configuration.
func NewDBFromContext(c *cli.Context) (*gorm.DB, error) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		c.String(postgresHostFlag),
		c.String(postgresUserFlag),
		c.String(postgresPasswordFlag),
		c.String(postgresDatabaseFlag),
		c.Int(postgresPortFlag),
	)
	return gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

// DatabaseNameFromContext return database name
func DatabaseNameFromContext(c *cli.Context) string {
	return c.String(postgresDatabaseFlag)
}
