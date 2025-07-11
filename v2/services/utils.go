package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	crud "github.com/dronm/crudifier"
	crudMd "github.com/dronm/crudifier/metadata"
	crudPg "github.com/dronm/crudifier/pg"
	crudTypes "github.com/dronm/crudifier/types"
	"github.com/dronm/ds/pgds"
	"github.com/jackc/pgx/v5"

	"github.com/dronm/gobizap/v2/database"
	"github.com/dronm/gobizap/v2/logger"
)

const defCollectionLimit = 5000

type DebugQueriesConfiger interface {
	GetDebugQueries() bool;
} 

var configer DebugQueriesConfiger

func InsertModel(ctx context.Context, db *pgds.PgProvider, model crudTypes.DbModel) (map[string]any, error) {
	poolConn, connID, err := db.GetPrimary()
	if err != nil {
		return nil, fmt.Errorf("GetPrimary() failed: %v", err)
	}
	defer db.Release(poolConn, connID)
	conn := poolConn.Conn()

	return InsertModelWithConn(ctx, conn, model)
}

// InsertModel insert model data to database ans returns server init field values and
// primary keys.
func InsertModelWithConn(ctx context.Context, conn *pgx.Conn, model crudTypes.DbModel) (map[string]any, error) {
	dbInsert := crudPg.NewPgInsert(model)
	if err := crud.PrepareInsertModel(dbInsert); err != nil {
		return nil, err
	}

	queryParams := make([]any, 0)
	queryText := dbInsert.SQL(&queryParams)

	if configer != nil && configer.GetDebugQueries() {
		logger.Logger.Debugf("InsertModel queryText: %s, params: %v", queryText, queryParams)
	}

	if err := conn.QueryRow(ctx, queryText, queryParams...).Scan(dbInsert.RetFieldValues()...); err != nil && !errors.Is(pgx.ErrNoRows, err) {
		return nil, err
	}

	return dbInsert.RetFields(), nil
}

func DeleteModel(ctx context.Context, db *pgds.PgProvider, keyModels []crudTypes.DbModel) (int64, error) {
	poolConn, connID, err := db.GetPrimary()
	if err != nil {
		return 0, fmt.Errorf("GetPrimary() failed: %v", err)
	}
	defer db.Release(poolConn, connID)
	conn := poolConn.Conn()

	return DeleteModelWithConn(ctx, conn, keyModels)
}

func DeleteModelWithConn(ctx context.Context, conn *pgx.Conn, keyModels []crudTypes.DbModel) (int64, error) {
	if len(keyModels) == 0 {
		return 0, fmt.Errorf("delete array is empty")
	}

	var filters crudPg.PgFilters
	for _, keyModel := range keyModels {
		if err := crud.ModelToDbFilters(keyModel, &filters, crudTypes.SQL_FILTER_OPERATOR_E, crudTypes.SQL_FILTER_JOIN_OR); err != nil {
			return 0, err
		}
	}

	dbDelete := crudPg.NewPgDelete(keyModels[0], filters)

	queryParams := make([]any, 0)
	queryText := dbDelete.SQL(&queryParams)

	if configer != nil && configer.GetDebugQueries() {
		logger.Logger.Debugf("DeleteModel queryText: %s, params: %v", queryText, queryParams)
	}

	cmd, err := conn.Exec(ctx, queryText, queryParams...)

	return cmd.RowsAffected(), err
}

func FetchModel(ctx context.Context, db *pgds.PgProvider,
	keyModel any, model crudTypes.DbModel,
) error {
	poolConn, connID, err := db.GetSecondary("")
	if err != nil {
		return fmt.Errorf("GetSecondary() failed: %v", err)
	}
	defer db.Release(poolConn, connID)
	conn := poolConn.Conn()

	return FetchModelWithConn(ctx, conn, keyModel, model)
}

func FetchModelWithConn(ctx context.Context, conn *pgx.Conn,
	keyModel any, model crudTypes.DbModel,
) error {
	filters := crudPg.PgFilters{}

	dbSelect := crudPg.NewPgDetailSelect(model, &filters)

	if err := crud.PrepareFetchModel(keyModel, dbSelect); err != nil {
		return err
	}
	queryParams := make([]any, 0)
	queryText := dbSelect.SQL(&queryParams)

	if configer != nil && configer.GetDebugQueries() {
		logger.Logger.Debugf("FetchModel queryText: %s, params: %v", queryText, queryParams)
	}

	return conn.QueryRow(ctx, queryText, queryParams...).Scan(dbSelect.FieldValues()...)
}

// FetchCollectionModel returns data model and aggregation model.
func FetchCollectionModel[T crudTypes.DbAggModel, U any](ctx context.Context, db *pgds.PgProvider,
	model T, totModel *U, params crud.CollectionParams,
) ([]T, *U, error) {
	poolConn, connID, err := db.GetSecondary("")
	if err != nil {
		return nil, nil, fmt.Errorf("GetSecondary() failed: %v", err)
	}
	defer db.Release(poolConn, connID)
	conn := poolConn.Conn()

	return FetchCollectionModelWithConn(ctx, conn, model, totModel, params)
}

func FetchCollectionModelWithConn[T crudTypes.DbAggModel, U any](ctx context.Context, conn *pgx.Conn,
	model T, totModel *U, params crud.CollectionParams,
) ([]T, *U, error) {
	// max limit if not defined in gui
	if params.Count == 0 {
		params.Count = defCollectionLimit
	}

	dbSelect := crudPg.NewPgSelect(model, &crudPg.PgFilters{}, &crudPg.PgSorters{}, &crudPg.PgLimit{})
	if err := crud.PrepareFetchModelCollection(dbSelect, params); err != nil {
		return nil, nil, fmt.Errorf("crud.PrepareFetchModelCollection(): %v", err)
	}

	queryParams := make([]any, 0)
	queryText, totQueryText := dbSelect.CollectionSQL(&queryParams)

	if configer != nil && configer.GetDebugQueries() {
		logger.Logger.Debugf("FetchCollectionModel queryText: %s, params: %v", queryText, queryParams)
		if totQueryText != "" {
			logger.Logger.Debugf("FetchCollectionModel totQueryText: %s, params: %v", totQueryText, queryParams)
		}
	}

	rows, err := conn.Query(ctx, queryText, queryParams...)
	if err != nil {
		return nil, nil, fmt.Errorf("collection conn.Query() failed: %v", err)
	}
	defer rows.Close()

	modelType := reflect.TypeOf(model).Elem()
	dataModelResult := make([]T, 0)

	for rows.Next() {
		row := reflect.New(modelType).Interface()

		// TODO: take only columns with json tags!!!
		rowVal := reflect.ValueOf(row).Elem()
		// rowFields := make([]any, rowVal.NumField())
		var rowFields []any
		for i := 0; i < rowVal.NumField(); i++ {
			fieldTag := modelType.Field(i).Tag.Get(crudMd.FieldAnnotationName)
			if fieldTag == "" || fieldTag == "-" {
				continue
			}
			field := rowVal.Field(i)
			// json tag must be present
			if field.CanSet() {
				// rowFields[i] = field.Addr().Interface()
				rowFields = append(rowFields, field.Addr().Interface())
			}
		}

		// Scan the row values into the struct fields
		if err := rows.Scan(rowFields...); err != nil {
			return nil, nil, fmt.Errorf("collection rows.Scan() failed: %v", err)
		}
		dataModelResult = append(dataModelResult, row.(T))
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	// total query
	if totQueryText == "" || totModel == nil {
		return dataModelResult, nil, nil
	}

	totRows, err := conn.Query(ctx, totQueryText, queryParams...)
	if err != nil {
		return nil, nil, fmt.Errorf("total conn.Query() failed: %v", err)
	}
	defer totRows.Close()

	if totRows.Next() {
		totResult := reflect.New(reflect.TypeOf(totModel).Elem()).Interface().(*U)

		// take only fields with json tags
		rowVal := reflect.ValueOf(totResult).Elem()
		rowFields := make([]any, rowVal.NumField())
		for i := 0; i < rowVal.NumField(); i++ {
			field := rowVal.Field(i)
			if field.CanSet() {
				rowFields[i] = field.Addr().Interface()
			}
		}

		// Scan the row values into the struct fields
		if err := totRows.Scan(rowFields...); err != nil {
			return nil, nil, fmt.Errorf("total rows.Scan() failed: %v", err)
		}
		// totModel = totResult
		reflect.ValueOf(totModel).Elem().Set(reflect.ValueOf(totResult).Elem())
	}

	return dataModelResult, totModel, nil
}

func UpdateModel(ctx context.Context, db *pgds.PgProvider,
	keyModel any, model crudTypes.DbModel,
) (int64, error) {
	poolConn, connID, err := db.GetPrimary()
	if err != nil {
		return 0, fmt.Errorf("GetPrimary() failed: %v", err)
	}
	defer db.Release(poolConn, connID)
	conn := poolConn.Conn()

	return UpdateModelWithConn(ctx, conn, keyModel, model)
}

// UpdateModel update date in model table.
// It returns number of actually affected rows.
func UpdateModelWithConn(ctx context.Context, conn *pgx.Conn, keyModel any, model crudTypes.DbModel) (int64, error) {
	dbUpdate := crudPg.NewPgUpdate(model)
	if err := crud.PrepareUpdateModel(keyModel, dbUpdate); err != nil {
		return 0, err
	}

	queryParams := make([]any, 0)
	queryText := dbUpdate.SQL(&queryParams)

	if configer != nil && configer.GetDebugQueries() {
		logger.Logger.Debugf("UpdateModel queryText: %s, params: %v", queryText, queryParams)
		// logger.Logger.Debugf("queryParams[0]:",string(queryParams[0].([]byte)))
	}

	cmd, err := conn.Exec(ctx, queryText, queryParams...)
	if err != nil {
		return 0, err
	}

	return cmd.RowsAffected(), nil
}

// tags: sql:"false" f:"fieldName" json:"fieldName"
// If "sql" tag is set to false, then field is ignored.
// If "f" tag present then it is treated as a field name.
// If "json" tag present and it is not "-" then it is treated as a firld name.
func AddStructFieldsToList(v reflect.Value, fields *[]any, fieldIDs *strings.Builder, fieldPrefix string) error {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			break
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Anonymous {
			if err := AddStructFieldsToList(v.Field(i), fields, fieldIDs, fieldPrefix); err != nil {
				return err
			}
		} else if sql, ok := t.Field(i).Tag.Lookup("sql"); !ok || sql != "false" {
			var fieldId string
			if fieldId, ok = t.Field(i).Tag.Lookup("f"); !ok {
				// no f tag
				if fieldId, ok = t.Field(i).Tag.Lookup("json"); !ok || fieldId == "-" {
					// no json
					continue
				}
			} else if fieldId == "-" {
				continue
			}

			value_field := v.Field(i)
			*fields = append(*fields, value_field.Addr().Interface())

			if fieldIDs.Len() > 0 {
				fieldIDs.WriteString(",")
			}
			fieldIDs.WriteString(fieldPrefix + fieldId)
		}
	}
	return nil
}

// Returns:
//
//	struct fields,
//	list of field IDs for select query
//	error if any
func MakeStructRowFields(resultStruct any, fieldPrefix string) ([]any, string, error) {
	fields := make([]any, 0)
	var field_ids strings.Builder
	AddStructFieldsToList(reflect.ValueOf(resultStruct), &fields, &field_ids, fieldPrefix)
	return fields, field_ids.String(), nil
}

func GetMd5(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil || !os.IsNotExist(err) {
		return true
	}
	return false
}

func PublishEvent(ctx context.Context, sessionID string, eventID string, params string) error {
	poolConn, connID, err := database.DB.GetPrimary()
	if err != nil {
		return fmt.Errorf("GetPrimary() failed: %v", err)
	}
	defer database.DB.Release(poolConn, connID)
	conn := poolConn.Conn()

	return PublishEventWithConn(ctx, conn, sessionID, eventID, params)
}

func PublishEventWithConn(ctx context.Context, conn *pgx.Conn, sessionID string, eventID string, params string) error {
	_, err := conn.Exec(ctx, `SELECT pg_notify($1, $2)`, eventID, params)
	return err
}

