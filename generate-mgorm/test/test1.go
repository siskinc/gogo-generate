//go:generate ./../generate-mgorm --type MogormTest --client Client
package test

import (
	"github.com/globalsign/mgo/bson"
	"github.com/siskinc/mgorm"
)

var Client *mgorm.MongoDBClient = &mgorm.MongoDBClient{}

//@def MogormTest struct comment1
//@def MogormTest struct comment2
//@def MogormTest struct comment3
//@def soft_delete SoftDeleted
//@def soft_delete_at SoftDeletedAt
type MogormTest struct {
	//das
	AA            bson.ObjectId `json:"aa" bson:"_id"` //@def aa
	SoftDeleted   bool          `bson:"soft_deleted"`
	SoftDeletedAt int64
}
