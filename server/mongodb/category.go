package mongodb

import (
	"context"
	"log"
	"server/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
)

type Category struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `json:"name"`
	Tags []Tag              `bson:"tags,omitempty" json:"tags,omitempty"`
}

var CategoryDb *mongo.Collection

func init() {
	registerDBData(Category{})
}

func (Category) initTable() {
	CategoryDb = db.Collection("Category")
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	if count, _ := CategoryDb.CountDocuments(ctx, bson.M{}); count != 4 {
		cate1, cate2, cate3, cate4 :=
			Category{Name: "Category1"},
			Category{Name: "Category2"},
			Category{Name: "Category3"},
			Category{Name: "Category4"}
		CategoryDb.InsertMany(ctx, []interface{}{cate1, cate2, cate3, cate4})
	}
}

func AddCategory(data Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	_, err := CategoryDb.InsertOne(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return util.Errorf("add Category %s failed to exec.", data.Name).WithCause(err).WithCode(codes.AlreadyExists)
		}
		return util.Errorf("add Category %s failed to exec.", data.Name).WithCause(err)
	}

	return nil
}

func GetAllCategories() ([]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	// Create a pipeline for aggregation
	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "Tag"},
			{"localField", "_id"},
			{"foreignField", "category"},
			{"as", "tags"},
		}}},
	}

	// Perform the aggregation
	cursor, err := CategoryDb.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
		return nil, util.Errorf("find all components failed").WithCause(err)
	}
	defer cursor.Close(context.Background())

	var datas []Category
	if err := cursor.All(context.Background(), &datas); err != nil {
		return nil, util.Errorf("get all categories failed").WithCause(err)
	}
	return datas, nil
}

func UpdateCategory(id string, data Category) error {
	update := bson.D{{"$set", data}}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return util.Errorf("update %s Tag failed", id).WithCause(err)
	}
	_, err = CategoryDb.UpdateByID(ctx, _id, update)
	if err != nil {
		return util.Errorf("update %s Tag failed", id).WithCause(err)
	}
	return nil
}

func DeleteCategory(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	_, err := CategoryDb.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return util.Errorf("delete %s Tag failed", id).WithCause(err).Log()
	}
	return nil
}
