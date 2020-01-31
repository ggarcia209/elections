/* dynamotypes.go defines the Table and Query objects, and functions for creating them.
In addition, it defines functions for creating DynamoDB AttributeValue objects and database keys in map format */

// Package dynamo contains controls and objects for DynamoDB CRUD operations
package dynamo

import (
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Table represents a table and holds basic information about it
type Table struct {
	TableName      string
	PrimaryKeyName string
	PrimaryKeyType string
	SortKeyName    string
	SortKeyType    string
}

// DbInfo holds different variables to be passed to db operation functions
type DbInfo struct {
	Svc        *dynamodb.DynamoDB
	Tables     map[string]*Table
	FailConfig *FailConfig
}

// SetSvc sets the Svc field of the DbInfo obj
func (d *DbInfo) SetSvc(svc *dynamodb.DynamoDB) {
	d.Svc = svc
}

// SetFailConfig sets the FailConfig field of the DbInfo obj
func (d *DbInfo) SetFailConfig(fc *FailConfig) {
	d.FailConfig = fc
}

// AddTable adds a new Table obj to the Tables field of the DbInfo obj
func (d *DbInfo) AddTable(t *Table) {
	d.Tables[t.TableName] = t
}

// InitDbInfo constructs a DbInfo object with default values
func InitDbInfo() *DbInfo {
	return &DbInfo{Svc: nil, Tables: make(map[string]*Table), FailConfig: nil}
}

// Query holds the search values for both the Partition and Sort Keys
// Query also holds data for updating a specific item in the UpdateFieldName column
type Query struct {
	PrimaryValue    interface{}
	SortValue       interface{}
	UpdateValue     interface{}
	UpdateFieldName string
}

// New creates a new query by setting the Partition Key and Sort Key values
func (q *Query) New(pv, sv interface{}) { q.PrimaryValue, q.SortValue = pv, sv }

// UpdateCurrent sets the update fields for the current item
func (q *Query) UpdateCurrent(uv interface{}, name string) {
	q.UpdateValue, q.UpdateFieldName = uv, name
}

// UpdateNew selects a new item for an update
func (q *Query) UpdateNew(pv, sv, uv interface{}, name string) {
	q.PrimaryValue, q.SortValue, q.UpdateValue, q.UpdateFieldName = pv, sv, uv, name
}

// Reset clears all fields
func (q *Query) Reset() {
	q.PrimaryValue, q.SortValue, q.UpdateValue, q.UpdateFieldName = nil, nil, nil, ""
}

// CreateNewTableObj creates a new Table struct
// The Table's key's Go types must be declared as strings
// ex: t := CreateNewTableObj("my_table", "Year", "int", "MovieName", "string")
func CreateNewTableObj(tableName, pKeyName, pType, sKeyName, sType string) *Table {
	typeMap := map[string]string{
		"[]byte":   "B",
		"[][]byte": "BS",
		"bool":     "BOOL",
		"list":     "L",
		"map":      "M",
		"int":      "N",
		"[]int":    "NS",
		"null":     "NULL",
		"string":   "S",
		"[]string": "SS",
	}

	pt := typeMap[pType]
	st := typeMap[sType]

	return &Table{tableName, pKeyName, pt, sKeyName, st}
}

// CreateNewQueryObj creates a new Query struct
// pval, sval == Primary/Partition key, Sort Key
func CreateNewQueryObj(pval, sval interface{}) *Query {
	return &Query{PrimaryValue: pval, SortValue: sval}
}

func createAV(val interface{}) *dynamodb.AttributeValue {
	if val == nil { // setNull
		av := &dynamodb.AttributeValue{}
		av.SetNULL(true)
		return av
	}
	if _, ok := val.([]byte); ok {
		av := &dynamodb.AttributeValue{}
		av.SetB(val.([]byte))
		return av
	}
	if _, ok := val.(bool); ok {
		av := &dynamodb.AttributeValue{}
		av.SetBOOL(val.(bool))
		return av
	}
	if _, ok := val.([][]byte); ok {
		av := &dynamodb.AttributeValue{}
		av.SetBS(val.([][]byte))
		return av
	}
	if _, ok := val.([]*dynamodb.AttributeValue); ok {
		av := &dynamodb.AttributeValue{}
		av.SetL(val.([]*dynamodb.AttributeValue))
		return av
	}
	if _, ok := val.(map[string]*dynamodb.AttributeValue); ok {
		av := &dynamodb.AttributeValue{}
		av.SetM(val.(map[string]*dynamodb.AttributeValue))
		return av
	}

	if _, ok := val.(int); ok {
		av := &dynamodb.AttributeValue{}
		av.SetN(strconv.Itoa(val.(int)))
		return av
	}
	if _, ok := val.([]int); ok {
		av := &dynamodb.AttributeValue{}

		ns := func(is []int) []*string {
			list := []*string{}
			for _, n := range is {
				str := strconv.Itoa(n)
				list = append(list, &str)
			}
			return list
		}(val.([]int))

		av.SetNS(ns)
		return av
	}
	if _, ok := val.(string); ok {
		av := &dynamodb.AttributeValue{}
		av.SetS(val.(string))
		return av
	}
	if _, ok := val.(string); ok {
		av := &dynamodb.AttributeValue{}
		av.SetS(val.(string))
		return av
	}
	return nil
}

// setType sets the the type of an AttributeValue object
/* func CreateAV(key interface{}, keyType string) *dynamodb.AttributeValue {
	av := *dynamodb.AttributeValue

	typeMap := map[string]func(key interface{}){
		"B":    func(key interface{}) { av.SetB(key.([]byte)) },
		"BOOL": func(key interface{}) { av.SetBOOL(key.(bool)) },
		"BS":   func(key interface{}) { av.SetBS(key.([][]byte)) },
		"L":    func(key interface{}) { av.SetL(key.([]*dynamodb.AttributeValue)) },
		"M":    func(key interface{}) { av.SetM(key.(map[string]*dynamodb.AttributeValue)) },
		"N":    func(key interface{}) { av.SetN(key.(string)) },
		"NS":   func(key interface{}) { av.SetNS(key.([]*string)) },
		"NULL": func(key interface{}) { av.SetNULL(key.(bool)) },
		"S":    func(key interface{}) { av.SetS(key.(string)) },
		"SS":   func(key interface{}) { av.SetSS(key.([]*string)) },
	}

	typeMap[keyType](key)

	return av
} */

// KeyMaker creates a map of Partition and Sort Keys
func keyMaker(q *Query, t *Table) map[string]*dynamodb.AttributeValue {
	keys := make(map[string]*dynamodb.AttributeValue)
	keys[t.PrimaryKeyName] = createAV(q.PrimaryValue)
	keys[t.SortKeyName] = createAV(q.SortValue)
	return keys
}
