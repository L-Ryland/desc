package mongodb

import (
	"context"
	"flag"
	"server/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
)

var tagIgnoreRef = flag.Bool("mongodb.tag.ignoreRef", true, "if ref abort delete tag")

var Tagdb *mongo.Collection

type Tag struct {
	Name     string
	Ref      int
	Order    int
	Category primitive.ObjectID `bson:"category,omitempty" json:"Category,omitempty"`
}

func init() {
	registerDBData(Tag{})
}

func (Tag) initTable() {
	Tagdb = db.Collection("Tag")

	// 创建索引选项
	indexOptions := options.Index().SetUnique(true)

	// 创建索引模型
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: indexOptions,
	}

	// 创建唯一性索引
	_, err := Tagdb.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		logrus.Errorf("create index for tag name err: %v", err)
	}
}

// 插入 Tag 表数据
func AddTag(data Tag) error {
	tags, _ := GetAllTags(bson.M{})
	data.Ref = 0
	data.Order = len(tags)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	_, err := Tagdb.InsertOne(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return util.Errorf("add Tag %s failed to exec.", data.Name).WithCause(err).WithCode(codes.AlreadyExists)
		}
		return util.Errorf("add Tag %s failed to exec.", data.Name).WithCause(err)
	}

	return nil
}

// 删除 Tag 表数据
func DeleteTag(name string) error {
	filter := bson.M{"name": name, "ref": 0}
	if *tagIgnoreRef {
		filter = bson.M{"name": name}
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()

	result := Tag{}
	err := Tagdb.FindOneAndDelete(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return util.Errorf("Tag with name %s not found or Ref not 0", name)
		}
		// 其他错误
		return util.Errorf("delete Tag with name %s failed", name).WithCause(err)
	}
	return nil
}

func GetTagByName(name string) (Tag, error) {
	filter := bson.M{"name": name}
	result := Tag{}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	err := Tagdb.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, util.Errorf("get %s Tag failed", name).WithCause(err).WithCode(codes.NotFound)
		}
		return result, util.Errorf("get %s Tag failed", name).WithCause(err)
	}
	return result, nil
}

func GetAllTags(filter bson.M) ([]Tag, error) {
	sortOptions := options.Find().SetSort(bson.D{{"order", 1}})
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	cursor, err := Tagdb.Find(ctx, filter, sortOptions)
	if err != nil {
		return nil, util.Errorf("find all components failed").WithCause(err)
	}
	defer cursor.Close(context.Background())

	var datas []Tag
	if err := cursor.All(context.Background(), &datas); err != nil {
		return nil, util.Errorf("get all tags failed").WithCause(err)
	}
	return datas, nil
}

// IncTagRef need session context
func IncTagRef(ctx context.Context, name string, num int) error {
	filter := bson.M{"name": name}
	update := bson.D{{"$inc", bson.M{"ref": num}}}
	var beforeData Tag
	err := Tagdb.FindOneAndUpdate(ctx, filter, update).Decode(&beforeData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			util.Errorf("inc %s Tag failed", name).WithCause(err).WithCode(codes.NotFound)
		}
		return util.Errorf("inc %s Tag failed", name).WithCause(err)
	}
	return nil
}

func UpdateTag(name string, tagData Tag) error {
	filter := bson.M{"name": name}
	update := bson.D{{"$set", tagData}}
	if tagData.Category.IsZero() {
		update = append(update, bson.E{"$unset", bson.M{"category": nil}})
		// update = bson.D{{"$set", tagData}, {"$unset", bson.M{"category": nil}}
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	err := Tagdb.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			util.Errorf("update %s Tag failed", name).WithCause(err).WithCode(codes.NotFound)
		}
		return util.Errorf("update %s Tag failed", name).WithCause(err)
	}
	return nil
}
