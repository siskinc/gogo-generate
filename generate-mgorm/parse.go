package main

var parseBasic = `
// Code generated by generate-mgorm DO NOT EDIT.
package {{.package}}

import (
	"time"
)

var _ = time.Now()

func (model *{{.struct}}) Save() (err error) {
	return {{.client}}.Save(model)
}

func (model *{{.struct}}) DeleteByID() (err error) {
	return {{.client}}.DeleteByID(model.{{.id}})
}

func (model *{{.struct}}) Delete(query interface{}) (err error) {
	return {{.client}}.DeleteAll(query)
}

func (model *{{.struct}}) FindByID() (err error) {
	query := bson.M{"_id": {{.struct}}.{{.id}}}
	return model.Find(query)
}

func (model *{{.struct}}) Find(query interface{}) (err error) {
	return {{.client}}.FindOne(model, query)
}

func (model *{{.struct}}) FindAll(query interface{}) (modelList []{{.struct}}, err error) {
	return {{.client}}.FindAll(modelList, query)
}

func (model *{{.struct}}) Find4Iter(query interface{}) (iter *mgo.Iter, err error) {
	return {{.client}}.FindAll4Iter(query)
}

func (model *{{.struct}}) FindPage(query interface{}, iPageSize, iPageIndex int, SortedStrs ...string) (err error) {
	return {{.client}}.FindPage(query)
}

func (model *{{.struct}}) UpdateByID() (err error) {
	objectID := {{.struct}}.{{.id}}
	query := bson.M{"_id": objectID}
	return {{.client}}.UpdateOne(query, model)
}

func (model *{{.struct}}) Update(query interface{}) (err error) {
	return {{.client}}.UpdateAll(query, model)
}

`

var parseSoftDelete = `
func (model *{{.struct}}) SoftDeleteByID() (err error) {
	query := bson.M{
		"_id": {{.struct}}.{{.id}},
		"{{.soft_delete_bson_name}}": false,
	}

	err = model.Count(query)
	if nil != err {
		return
	}

	err = model.Find(query)
	if nil != err {
		return
	}

	model.{{.soft_delete}} = true

	{{ if .soft_delete_at}}
	model.{{.soft_delete_at}} = time.Now().UTC()
	{{ end }}
	err = model.UpdateByID()
	return
}

`
