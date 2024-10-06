package repository_test

import (
	"fmt"
	"kumquat/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyQuery(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	rs, err := r.Query(`VACUUM`) // Rebuild indexes, doesn't return anything
	require.NoError(t, err)
	assert.Empty(t, rs.Names)
	assert.Empty(t, rs.Results)
}

func TestJSONLiteral(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	_, err = r.Query( /* sql */ `SELECT '{}' AS test`)
	assert.ErrorContains(t, err, "missing apiVersion")

	rs, err := r.Query(`SELECT ` +
		`'{"apiVersion":"guidewire.com/v1beta1","kind":"Test","metadata":{"namespace":"testns","name":"test"}}'` +
		` AS test`,
	)
	require.NoError(t, err)
	assert.ElementsMatch(t, rs.Names, []string{"test"})
	require.Len(t, rs.Results, 1)
	require.Len(t, rs.Results[0], 1)
	assert.Contains(t, rs.Results[0], "test")
	assert.Equal(t, "Test", rs.Results[0]["test"].Content()["kind"])
	assert.Equal(t, "testns", rs.Results[0]["test"].Content()["metadata"].(map[string]any)["namespace"])
}

func TestNonJSONLiteral(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	_, err = r.Query( /* sql */ `SELECT "1" AS test`)
	assert.ErrorContains(t, err, "expected JSON object in column 'test'")

	_, err = r.Query( /* sql */ `SELECT '["hello"]' AS test`)
	assert.ErrorContains(t, err, "expected JSON object in column 'test'")

	_, err = r.Query( /* sql */ `SELECT '[1]' AS test`)
	assert.ErrorContains(t, err, "expected JSON object in column 'test'")
}

func TestInsert(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	res, err := repository.MakeResource(map[string]any{
		"apiVersion": "guidewire.com/v1beta1",
		"kind":       "Example",
		"metadata": map[string]any{
			"name":      "alpha",
			"namespace": "examples",
		},
	})
	require.NoError(t, err)
	err = r.Upsert(res)
	require.NoError(t, err)

	rs, err := r.Query(
		/* sql */ `SELECT example.data AS e FROM "Example.guidewire.com" AS example WHERE example.name = 'alpha'`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, rs.Names, []string{"e"})
	require.Len(t, rs.Results, 1)
	assert.Len(t, rs.Results[0], 1)

	res2 := rs.Results[0]["e"]
	assert.NotNil(t, res2)
	assert.NotZero(t, res2)
	assert.Equal(t, "examples", res2.Namespace())
	assert.Equal(t, "examples", res2.Content()["metadata"].(map[string]any)["namespace"])
}

func TestInsertResourceHavingUnmarshallableData(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	res, err := repository.MakeResource(map[string]any{
		"apiVersion": "guidewire.com/v1beta1",
		"kind":       "Example",
		"metadata": map[string]any{
			"name":      "alpha",
			"namespace": "examples",
		},
		// can't marshal a channel
		"data": make(chan int),
	})
	require.NoError(t, err)

	err = r.Upsert(res)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "unable to encode resource as JSON")
}

func TestUpdate(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	res, err := repository.MakeResource(map[string]any{
		"apiVersion": "guidewire.com/v1beta1",
		"kind":       "Example",
		"metadata": map[string]any{
			"name":      "alpha",
			"namespace": "examples",
		},
		"spec": map[string]any{
			"released": false,
		},
	})
	require.NoError(t, err)
	err = r.Upsert(res)
	require.NoError(t, err)

	res.Content()["spec"].(map[string]any)["released"] = true
	err = r.Upsert(res)
	require.NoError(t, err)

	rs, err := r.Query(
		/* sql */ `SELECT example.data AS e FROM "Example.guidewire.com" AS example WHERE example.name = 'alpha'`)
	assert.NoError(t, err)
	assert.ElementsMatch(t, rs.Names, []string{"e"})
	require.Len(t, rs.Results, 1)
	assert.Len(t, rs.Results[0], 1)

	res2 := rs.Results[0]["e"]
	assert.NotNil(t, res2)
	assert.NotZero(t, res2)
	assert.True(t, res2.Content()["spec"].(map[string]any)["released"].(bool))
}

func TestQueryOfMissingTable(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	rs, err := r.Query(
		/* sql */ `SELECT example.data AS e FROM "Example.guidewire.com" AS example WHERE example.name = 'test'`)
	assert.NoError(t, err)
	assert.Empty(t, rs.Results)
	assert.ElementsMatch(t, rs.Names, []string{"e"})
}

func TestExtractingTableNames(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	tableNames := r.ExtractTableNamesFromQuery( /* sql */ `
	SELECT *
	FROM employees
	WHERE employee_id < 50
	ORDER BY last_name ASC;
	`)
	assert.ElementsMatch(t, tableNames, []string{"employees"})
}

func TestExtractingTableNamesWithDoubleQuotes(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck
	tableNames := r.ExtractTableNamesFromQuery( /* sql */ `
	SELECT persistentvolumeclaim.data AS pvc, persistentvolume.data AS pv, pod.data AS p
	FROM "PersistentVolumeClaim.core" AS persistentvolumeclaim
	JOIN "PersistentVolume.core" AS persistentvolume
	ON persistentvolumeclaim.data ->> '$.spect.volumeName' = persistentvolume.name
	CROSS JOIN "Pod.core" AS pod
	WHERE persistentvolumeclaim.namespace = 'default' AND pod.namespace = 'default'
	`)
	assert.ElementsMatch(t, tableNames, []string{"PersistentVolumeClaim.core", "PersistentVolume.core", "Pod.core"})
}
func TestExtractingTableNamesFromQueriesWIthSubQuery(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	tableNames := r.ExtractTableNamesFromQuery( /* sql */ `
	SELECT *
	FROM (
		SELECT *
		FROM employees
		WHERE employee_id < 50
		ORDER BY last_name ASC
	);
	`)
	assert.ElementsMatch(t, tableNames, []string{"employees"})
}

func TestExtractingTableNamesFromNestedSubqueries(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	tableNames := r.ExtractTableNamesFromQuery( /* sql */ `
    SELECT a.*
    FROM (
        SELECT b.*
        FROM (
            SELECT c.*
            FROM employees AS c
            WHERE c.employee_id < 50
        ) AS b
    ) AS a
    `)
	assert.ElementsMatch(t, tableNames, []string{"employees"})
}

func TestExtractingTableNamesFromMultipleJoins(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	tableNames := r.ExtractTableNamesFromQuery( /* sql */ `
    SELECT a.*, b.*, c.*
    FROM "TableA" AS a
    JOIN "TableB" AS b ON a.id = b.a_id
    LEFT JOIN "TableC" AS c ON b.id = c.b_id
    `)
	assert.ElementsMatch(t, tableNames, []string{"TableA", "TableB", "TableC"})
}

func TestExtractingTableNamesFromComplexQuery(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	tableNames := r.ExtractTableNamesFromQuery( /* sql */ `
    SELECT e.*, d.*
    FROM employees AS e
    JOIN (
        SELECT department_id, department_name
        FROM departments
        WHERE location_id = 1700
    ) AS d ON e.department_id = d.department_id
    WHERE e.salary > (
        SELECT AVG(salary)
        FROM employees
        WHERE department_id = e.department_id
    )
    `)
	assert.ElementsMatch(t, tableNames, []string{"employees", "departments"})
}

func TestExtractingTableNamesFromQueryWithMultipleSubqueries(t *testing.T) {
	r, err := repository.NewSQLiteRepository()
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	tableNames := r.ExtractTableNamesFromQuery( /* sql */ `
    SELECT e.*, d.*, l.*
    FROM employees AS e
    JOIN departments AS d ON e.department_id = d.department_id
    JOIN locations AS l ON d.location_id = l.location_id
    WHERE e.salary > (
        SELECT AVG(salary)
        FROM employees
        WHERE department_id = e.department_id
    )
    AND e.employee_id IN (
        SELECT employee_id
        FROM job_history
        WHERE job_id = 'IT_PROG'
    )
    `)
	assert.ElementsMatch(t, tableNames, []string{"employees", "departments", "locations", "job_history"})
}

// func TestExtractingTableNamesFromQueryWithCTE(t *testing.T) {
// 	r, err := repository.NewSQLiteRepository()
// 	require.NoError(t, err)
// 	defer r.Close() //nolint:errcheck

// 	tableNames, err := r.ExtractTableNamesFromQuery( /* sql */ `
//     WITH EmployeeCTE AS (
//         SELECT employee_id, first_name, last_name, department_id
//         FROM employees
//         WHERE salary > 5000
//     ),
//     DepartmentCTE AS (
//         SELECT department_id, department_name
//         FROM departments
//     )
//     SELECT e.*, d.*
//     FROM EmployeeCTE AS e
//     JOIN DepartmentCTE AS d ON e.department_id = d.department_id
//     `)
// 	require.NoError(t, err)
// 	assert.ElementsMatch(t, tableNames, []string{"employees", "departments"})
// }

func BenchmarkQueryPerformance(b *testing.B) {
	const DB_ENTRIES = 500

	r, err := repository.NewSQLiteRepository()
	require.NoError(b, err)
	defer r.Close() //nolint:errcheck

	// Populate repository with some resources
	for i := 0; i < DB_ENTRIES; i++ {
		res, err := repository.MakeResource(map[string]any{
			"apiVersion": "guidewire.com/v1beta1",
			"kind":       "Example",
			"metadata": map[string]any{
				"name":      fmt.Sprintf("%04d", i),
				"namespace": "examples",
			},
		})
		require.NoError(b, err)
		err = r.Upsert(res)
		require.NoError(b, err)
	}

	benchmarkQueryExpecting := func(query string, expected int) func(*testing.B) {
		return func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rs, err := r.Query(query)
				if err != nil {
					b.Fatal(err)
				}

				if len(rs.Results) != expected {
					b.Fatalf("unexpected result: %v", rs)
				}
			}
		}
	}

	// Find the time to query the first resource
	q := fmt.Sprintf(
		`SELECT example.data AS e FROM "Example.guidewire.com" AS example WHERE example.name = '%04d'`,
		0,
	)
	b.Run("QueryFirst", benchmarkQueryExpecting(q, 1))

	// Find the time to query the last resource
	q = fmt.Sprintf(
		`SELECT example.data AS e FROM "Example.guidewire.com" AS example WHERE example.name = '%04d'`,
		DB_ENTRIES-1,
	)
	b.Run("QueryLast", benchmarkQueryExpecting(q, 1))

	// Find the time to query a non-existent resource
	q = `SELECT example.data AS e FROM "Example.guidewire.com" AS example WHERE example.name = 'missing'`
	b.Run("QueryMissing", benchmarkQueryExpecting(q, 0))

	// Get Cartesian product of table with itself
	q = `SELECT a.data, b.data FROM "Example.guidewire.com" AS a CROSS JOIN "Example.guidewire.com" AS b`
	b.Run("QueryCartesianProduct", benchmarkQueryExpecting(q, DB_ENTRIES*DB_ENTRIES))
}
