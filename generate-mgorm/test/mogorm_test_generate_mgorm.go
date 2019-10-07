// Code generated by generate-mgorm DO NOT EDIT.
package test

import (
	"time"
)

var _ = time.Now()

func (model *MogormTest) Save() (err error) {

	model.CreateAt = time.Now().UTC()

	return Client.Save(model)
}

func (model *MogormTest) DeleteByID() (err error) {
	return Client.DeleteByID(model.AA)
}

func (model *MogormTest) Delete(query interface{}) (err error) {
	return Client.DeleteAll(query)
}

func (model *MogormTest) FindByID() (err error) {
	query := bson.M{
		"_id": MogormTest.AA,

		"soft_deleted": false,
	}
	return model.Find(query)
}

func (model *MogormTest) Find(query interface{}) (err error) {

	mQuery := query.(map[string]interface{})
	mQuery["soft_deleted"] = false

	return Client.FindOne(model, query)
}

func (model *MogormTest) FindAll(query interface{}) (modelList []MogormTest, err error) {

	mQuery := query.(map[string]interface{})
	mQuery["soft_deleted"] = false

	return Client.FindAll(modelList, query)
}

func (model *MogormTest) Find4Iter(query interface{}) (iter *mgo.Iter, err error) {

	mQuery := query.(map[string]interface{})
	mQuery["soft_deleted"] = false

	return Client.FindAll4Iter(query)
}

func (model *MogormTest) FindPage(query interface{}, iPageSize, iPageIndex int, SortedStrs ...string) (err error) {

	mQuery := query.(map[string]interface{})
	mQuery["soft_deleted"] = false

	return Client.FindPage(query)
}

func (model *MogormTest) UpdateByID() (err error) {
	objectID := MogormTest.AA
	query := bson.M{
		"_id": MogormTest.AA,

		"soft_deleted": false,
	}

	model.UpdateAt = time.Now().UTC()

	return Client.UpdateOne(query, model)
}

func (model *MogormTest) Update(query interface{}) (err error) {

	mQuery := query.(map[string]interface{})
	mQuery["soft_deleted"] = false

	model.UpdateAt = time.Now().UTC()

	return Client.UpdateAll(query, model)
}

func (model *MogormTest) SoftDeleteByID() (err error) {
	query := bson.M{
		"_id":          MogormTest.AA,
		"soft_deleted": false,
	}

	err = model.Count(query)
	if nil != err {
		return
	}

	err = model.Find(query)
	if nil != err {
		return
	}

	model.SoftDeleted = true

	model.SoftDeletedAt = time.Now().UTC()

	err = model.UpdateByID()
	return
}
