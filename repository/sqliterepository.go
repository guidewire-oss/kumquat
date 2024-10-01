package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type SQLiteRepository struct {
	Db                 *sql.DB
	StoredKinds        map[string]bool
	UseQueryObjectPool bool
}

var tableErrorRegexp = regexp.MustCompile(`^no such table: (.*)$`)

func (r *SQLiteRepository) Query(query string) (ResultSet, error) {
	slog.Debug("Running query", "query", query)

	type objectKey struct {
		kind      string
		namespace string
		name      string
	}

	var rows *sql.Rows
	var err error

	for rows, err = r.Db.Query(query); err != nil; rows, err = r.Db.Query(query) {
		m := tableErrorRegexp.FindStringSubmatch(err.Error())

		if m != nil {
			fmt.Println("Table not found, creating table", "table", m[1])
			missingTable := m[1]
			err := r.createTable(missingTable)

			if err != nil {
				return ResultSet{}, fmt.Errorf("unable to create empty table to run query: %w", err)
			}
		} else {
			return ResultSet{}, fmt.Errorf("error running query: %w", err)
		}
	}

	defer rows.Close() //nolint:errcheck

	/*
	 * Process Columns
	 */
	columnNames, err := rows.Columns()
	if err != nil {
		return ResultSet{}, fmt.Errorf("error getting query result column information: %w", err)
	}

	/*
	 * Allocate temporary storage to Scan() into
	 */
	columnValues := make([]any, len(columnNames))
	columnValuePtrs := make([]any, len(columnNames))

	for i := 0; i < len(columnNames); i++ {
		columnValuePtrs[i] = &columnValues[i]
	}

	/*
	 * Process Rows
	 */
	var results = make([]map[string]*Resource, 0)
	var pool = make(map[objectKey]*Resource)

	for rows.Next() {
		err = rows.Scan(columnValuePtrs...)
		if err != nil {
			return ResultSet{}, fmt.Errorf("error while scanning query result: %w", err)
		}

		var result = make(map[string]*Resource)

		for i, columnName := range columnNames {
			var parsed any

			err = json.Unmarshal([]byte(columnValues[i].(string)), &parsed)
			if err != nil {
				return ResultSet{}, fmt.Errorf("error while unmarshaling column '%s' to JSON: %w",
					strings.Trim(columnName, "'"), err)
			}

			switch v := parsed.(type) {
			case map[string]any:
				res, err := MakeResource(v)
				if err != nil {
					return ResultSet{}, fmt.Errorf("error retrieving resource from column '%s': %w",
						strings.Trim(columnName, "'"), err)
				}

				resPtr := &res
				if r.UseQueryObjectPool {
					key := objectKey{
						kind:      res.Kind(),
						namespace: res.Namespace(),
						name:      res.Name(),
					}

					if found, ok := pool[key]; ok {
						resPtr = found
					} else {
						pool[key] = resPtr
					}
				}

				result[columnName] = resPtr
			default:
				return ResultSet{}, fmt.Errorf("expected JSON object in column '%s' but got %T",
					strings.Trim(columnName, "'"), v)
			}
		}

		results = append(results, result)
	}

	err = rows.Err()
	if err != nil {
		return ResultSet{}, fmt.Errorf("error while iterating query result: %w", err)
	}

	resultset := ResultSet{
		Names:   columnNames,
		Results: results,
	}

	return resultset, nil
}

func (r *SQLiteRepository) Close() error {
	return r.Db.Close()
}

func (r *SQLiteRepository) createTable(table string) error {
	if r.StoredKinds[table] {
		return fmt.Errorf("table already exists: %s", table)
	}

	_, err := r.Db.Exec( /* sql */ `CREATE TABLE "` + table +
		`" (namespace TEXT NOT NULL, name TEXT NOT NULL, data TEXT NOT NULL, PRIMARY KEY (namespace, name)) STRICT`)

	if err != nil {
		return fmt.Errorf("unable to create table: %w", err)
	}

	r.StoredKinds[table] = true

	return nil
}

func (r *SQLiteRepository) Upsert(resource Resource) error {
	byteJSON, err := json.Marshal(resource.Content())

	if err != nil {
		return fmt.Errorf("unable to encode resource as JSON: %w", err)
	}

	table := resource.Kind() + "." + resource.Group()
	contentJSON := string(byteJSON)

	if !r.StoredKinds[table] {
		err := r.createTable(table)

		if err != nil {
			return err
		}
	}

	_, err = r.Db.Exec( /* sql */ `INSERT INTO "`+table+`" (namespace, name, data) VALUES (?,?,?)
		ON CONFLICT(namespace, name) DO UPDATE SET data=excluded.data`,
		resource.Namespace(), resource.Name(), contentJSON)

	if err != nil {
		return fmt.Errorf("unable to upsert resource: %w", err)
	}

	return nil
}

// a function that builds data columns from the resource and check if a row with the same value for data column exists
func (r *SQLiteRepository) CheckIfResourceExists(resource Resource) (bool, error) {
	byteJSON, err := json.Marshal(resource.Content())

	if err != nil {
		return false, fmt.Errorf("unable to encode resource as JSON: %w", err)
	}

	table := resource.Kind() + "." + resource.Group()
	fmt.Println("Checking if resource with same data already exists",
		"table", table, "namespace", resource.Namespace(), "name", resource.Name())
	contentJSON := string(byteJSON)

	if !r.StoredKinds[table] {
		err := r.createTable(table)

		if err != nil {
			return false, err
		}
	}

	var count int
	err = r.Db.QueryRow( /* sql */ `SELECT COUNT(*) FROM "`+table+`" WHERE data = ?`, contentJSON).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("unable to check if resource exists: %w", err)
	}

	return count > 0, nil
}

func (r *SQLiteRepository) ExtractTableNamesFromQuery(query string) []string {
	// Extract table names from query
	tableNames := make([]string, 0)
	tableNameSet := make(map[string]struct{})

	// Find all table names in the query, including subqueries and quoted table names
	tableNameRegexp := regexp.MustCompile(`(?i)(?:FROM|JOIN|CROSS JOIN)\s+["]?(\w+(\.\w+)*|[\w.]+)["]?`)
	tableNameMatches := tableNameRegexp.FindAllStringSubmatch(query, -1)

	// Extract table names from the matches and add to set to avoid duplicates
	for _, match := range tableNameMatches {
		tableName := match[1]
		if _, exists := tableNameSet[tableName]; !exists {
			tableNames = append(tableNames, tableName)
			tableNameSet[tableName] = struct{}{}
		}
	}

	return tableNames
}
func (r *SQLiteRepository) DropTable(table string) error {
	if !r.StoredKinds[table] {
		return fmt.Errorf("table does not exist: %s", table)
	}
	_, err := r.Db.Exec( /* sql */ `DROP TABLE "` + table + `"`)

	if err != nil {
		return fmt.Errorf("unable to drop table: %w", err)
	}

	delete(r.StoredKinds, table)
	log.Log.Info("Table dropped", "tableName", table)

	return nil
}

func NewSQLiteRepository() (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	repo := &SQLiteRepository{
		Db:          db,
		StoredKinds: make(map[string]bool),
	}

	return repo, nil
}
