package mongodb

import (
	"context"
	"server/util"
	"strings"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
)

var WebDatadb *mongo.Collection
var WebDataNum int

type WebData struct {
	ID          int `bson:"_id"`
	Name        string
	Url         string
	Tags        []string
	Description string
}

func init() {
	registerDBData(WebData{})
}

func (WebData) initTable() {
	WebDatadb = db.Collection("webData")

	// 创建索引选项
	indexOptions := options.Index().SetUnique(true)

	// 创建索引模型
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},
		Options: indexOptions,
	}

	// 创建唯一性索引
	_, err := WebDatadb.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		logrus.Errorf("create index for webData url name err: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	// 创建一个排序条件，按降序排列
	sort := bson.D{{"_id", -1}}
	// 设置查询选项，仅返回一条文档
	options := options.FindOne().SetSort(sort)
	// 执行查询
	var result bson.M
	if err := WebDatadb.FindOne(ctx, bson.D{}, options).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			WebDataNum = 0
		}
	} else {
		// 获取最大 ID 值
		WebDataNum = int(result["_id"].(int32))
	}
}

// 插入 WebData 表数据
func AddWebData(data WebData) (num int, err error) {
	// session, err := WebDatadb.Database().Client().StartSession()
	// defer session.EndSession(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	// sctx := mongo.NewSessionContext(ctx, session)
	// if err = session.StartTransaction(); err != nil {
	// 	return 0, util.Errorf("add WebData with Name %s failed", data.Name).WithCause(err)
	// }
	// defer func() {
	// 	if err == nil {
	// 		session.CommitTransaction(context.Background())
	// 	} else {
	// 		session.AbortTransaction(context.Background())
	// 	}
	// }()

	WebDataNum++
	data.ID = WebDataNum
	res, err := WebDatadb.InsertOne(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return 0, util.Errorf("add WebData %s failed to exec.", data.Name).WithCause(err).WithCode(codes.AlreadyExists)
		}
		return 0, util.Errorf("add WebData %s failed to exec.", data.Name).WithCause(err)
	}
	for _, tag := range data.Tags {
		err := IncTagRef(ctx, tag, 1)
		if err != nil {
			return 0, util.Errorf("add WebData with Name %s failed", data.Name).WithCause(err)
		}
	}

	return int(res.InsertedID.(int32)), nil
}

// 删除 WebData 表数据
func DeleteWebData(ID int) (err error) {
	// session, err := WebDatadb.Database().Client().StartSession()
	// defer session.EndSession(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	// sctx := mongo.NewSessionContext(ctx, session)
	// if err = session.StartTransaction(); err != nil {
	// 	return util.Errorf("delete WebData with ID %d failed", ID).WithCause(err)
	// }
	// defer func() {
	// 	if err == nil {
	// 		session.CommitTransaction(context.Background())
	// 	} else {
	// 		session.AbortTransaction(context.Background())
	// 	}
	// }()

	filter := bson.M{"_id": ID}
	var deletedWebData WebData
	err = WebDatadb.FindOneAndDelete(ctx, filter).Decode(&deletedWebData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return util.Errorf("delete WebData with ID %d failed", ID).WithCause(err)
	}

	for _, tag := range deletedWebData.Tags {
		err := IncTagRef(ctx, tag, -1)
		if err != nil {
			if util.HaveErrorCode(err, codes.NotFound) {
				continue
			}
			return util.Errorf("delete WebData with ID %d failed", ID).WithCause(err)
		}
	}

	return nil
}

// 更新 WebData 表数据
func UpdateWebData(data WebData) (err error) {
	// session, err := WebDatadb.Database().Client().StartSession()
	// defer session.EndSession(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	// sctx := mongo.NewSessionContext(ctx, session)
	// if err = session.StartTransaction(); err != nil {
	// 	return util.Errorf("delete WebData with name %s failed", data.Name).WithCause(err)
	// }
	// defer func() {
	// 	if err == nil {
	// 		session.CommitTransaction(context.Background())
	// 	} else {
	// 		session.AbortTransaction(context.Background())
	// 	}
	// }()

	filter := bson.M{"_id": data.ID}
	update := bson.M{"$set": data}
	var originData WebData
	err = WebDatadb.FindOneAndUpdate(ctx, filter, update).Decode(&originData)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return util.Errorf("update WebData %s failed to exec.", data.Name).WithCause(err).WithCode(codes.AlreadyExists)
		}
		return util.Errorf("update WebData %s failed", data.Name).WithCause(err)
	}

	originTags := originData.Tags
	newTags := data.Tags
	ignoreTags := make(map[string]interface{})
	for _, ot := range originTags {
		for _, nt := range newTags {
			if ot == nt {
				ignoreTags[ot] = struct{}{}
				continue
			}
		}
	}
	for _, ot := range originTags {
		if _, ok := ignoreTags[ot]; ok {
			continue
		}
		err := IncTagRef(ctx, ot, -1)
		if err != nil {
			if util.HaveErrorCode(err, codes.NotFound) {
				continue
			}
			return util.Errorf("delete WebData with name %s failed", data.Name).WithCause(err)
		}
	}
	for _, nt := range newTags {
		if _, ok := ignoreTags[nt]; ok {
			continue
		}
		err := IncTagRef(ctx, nt, 1)
		if err != nil {
			return util.Errorf("delete WebData with name %s failed", data.Name).WithCause(err)
		}
	}

	return nil
}

func GetWebDataByName(name string) (WebData, error) {
	filter := bson.M{"name": name}
	result := WebData{}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	err := WebDatadb.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, util.Errorf("get %s WebData failed", name).WithCause(err)
	}
	return result, nil
}

func GetWebDataByTags(tags []string) ([]WebData, error) {
	filter := bson.M{"tags": bson.M{"$all": tags}}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	cursor, err := WebDatadb.Find(ctx, filter)
	if err != nil {
		return nil, util.Errorf("get %s WebData failed", strings.Join(tags, ",")).WithCause(err)
	}
	defer cursor.Close(context.Background())

	var datas []WebData
	if err := cursor.All(context.Background(), &datas); err != nil {
		return nil, util.Errorf("get all attr failed").WithCause(err)
	}
	return datas, nil
}
