package sql

import (
	"database/sql"
)

// CreateLogin connects to the SQL Database to create a login with the provided
// credentials
func (c Connector) CreateLogin(username string, password string) error {
	cmd := `DECLARE @sql nvarchar(max)
					SET @sql = 'CREATE LOGIN ' + QuoteName(@username) + ' ' +
										 'WITH PASSWORD = ' + QuoteName(@password, '''')
					EXEC (@sql)`
	return c.Execute(cmd, sql.Named("username", username), sql.Named("password", password))
}

// DeleteLogin connects to the SQL Database and removes a login with the provided
// username, if it exists. If it does not exist, this is a noop.
func (c Connector) DeleteLogin(username string) error {
	cmd := `DECLARE @sql nvarchar(max);
					SET @sql = 'IF EXISTS (SELECT 1 FROM [master].[sys].[server_principals] WHERE [name] = ' + QuoteName(@username, '''') + ') ' +
										 'DROP LOGIN ' + QuoteName(@username);
					EXEC (@sql)`
	return c.Execute(cmd, sql.Named("username", username))
}

// Login represents an SQL Server Login
type Login struct {
	Username    string
	PrincipalID int64
}

// GetLogin reads a login from the SQL Database, if it exists. If it does not,
// no error is returned, but the returned Login is nil
func (c Connector) GetLogin(username string) (*Login, error) {
	var principalID int64 = -1

	err := c.Query(
		"SELECT principal_id FROM [master].[sys].[server_principals] WHERE [name] = @username",
		func(r *sql.Rows) error {
			for r.Next() {
				err := r.Scan(&principalID)
				if err != nil {
					return err
				}
			}
			return nil
		},
		sql.Named("username", username),
	)

	if err != nil {
		return nil, err
	}
	if principalID != -1 {
		return &Login{Username: username, PrincipalID: principalID}, nil
	}
	return nil, nil
}

// UpdateLogin updates the password of a login, if it exists.
func (c Connector) UpdateLogin(username string, password string) error {
	cmd := `DECLARE @sql nvarchar(max)
					SET @sql = 'IF EXISTS (SELECT 1 FROM [master].[sys].[server_principals] WHERE [name] = ' + QuoteName(@username, '''') + ') ' +
										 'ALTER LOGIN ' + QuoteName(@username) + ' ' +
										 'WITH PASSWORD = ' + QuoteName(@password, '''')
					EXEC (@sql)`

	return c.Execute(cmd, sql.Named("username", username), sql.Named("password", password))
}
