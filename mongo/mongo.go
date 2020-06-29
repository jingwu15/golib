package mongo

import (
    "fmt"
    "time"
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Mongo struct {
    Uri     string
    Db      string
    Client  *mongo.Client
    Conn    *mongo.Database
    IsConn  bool
    Colls   map[string]*mongo.Collection
}

var mongos = map[string]*Mongo{}

//新建mongo连接
func New(uri, db string) (mg *Mongo, err error) {
    m := Mongo{Uri: uri, Db: db, Colls: map[string]*mongo.Collection{}, IsConn: false}
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    //opt := options.Client()
    //opt.SetConnectTimeout(10 * time.Second)
    //opt.SetMaxPoolSize(50)
    //opt.SetMinPoolSize(20)
    //opt.ApplyURI(m.Uri)
    m.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.Uri))
    if err != nil { return &m, err }
    m.Conn = m.Client.Database(db)
    mongos[db] = &m
    return mongos[db], nil
}

//mongo连接是否存在
func Exist(db string) bool {
    _, ok := mongos[db]
    return ok
}

//获取已存在的mogno连接
func GetMongo(db string) (mg *Mongo, err error) {
    if mg, ok := mongos[db]; ok {
        return mg, nil
    } else {
        return mg, fmt.Errorf("数据库连接不存在")
    }
}

//获取文档集连接，不存在则创建
func (m *Mongo)coll(doc string) *mongo.Collection {
    if _, ok := m.Colls[doc]; !ok {
        m.Colls[doc] = m.Conn.Collection(doc)
    }
    return m.Colls[doc]
}

//自增ID
func (m *Mongo)GetAutoId(doc string) (id int, err error) {
    step := 1
    docAutoid := "autoids"
    filter := bson.D{{"doc", doc}}
    update := bson.M{"$inc": bson.M{"id": step}}
    var resultUpdate bson.M
    opts := options.FindOneAndUpdate()
    err = m.coll(docAutoid).FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&resultUpdate)
    if err != nil && err != mongo.ErrNoDocuments { return id, err }
    if err == mongo.ErrNoDocuments {
        ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
        _, err := m.coll(docAutoid).InsertOne(ctx, bson.M{"doc":doc, "id": step})
        if err != nil { return 0, err }
        return step, nil
    } else {
        return int(resultUpdate["id"].(int32)) + step, nil
    }
}

func (m *Mongo)Insert(doc string, data interface{}) (objId primitive.ObjectID, id int, err error) {
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    res, err := m.coll(doc).InsertOne(ctx, data)
    if err != nil { return objId, 0, err }
    return res.InsertedID.(primitive.ObjectID), 0, nil
}

func (m *Mongo)Insert_AutoId(doc string, data bson.M) (objId primitive.ObjectID, autoId int, err error) {
    if _, ok := data["id"]; !ok { data["id"], _ = m.GetAutoId(doc) }
    autoId = 0
    if v, ok := data["id"].(float64); ok { autoId = int(v) }
    if v, ok := data["id"].(int);     ok { autoId = v      }
    if v, ok := data["id"].(int32);   ok { autoId = int(v) }
    data["create_at"] = time.Now()
    data["update_at"] = time.Now()
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    res, err := m.coll(doc).InsertOne(ctx, data)
    if err != nil { return objId, 0, err }
    return res.InsertedID.(primitive.ObjectID), autoId, nil
}

////不存在则插入，存在则更新
//func (m *Mongo)Save(doc string, filter interface{}, update interface{}) (objId primitive.ObjectID, id int, err error) {
//    var resultUpdate bson.M
//    opts := options.FindOneAndUpdate()
//    err = m.coll(doc).FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&resultUpdate)
//    if err != nil && err != mongo.ErrNoDocuments { return objId, id, err }
//    if err == mongo.ErrNoDocuments {
//        res, autoid, err := m.Insert_AutoId(doc, update)
//        fmt.Println("res: ", res)
//        if err != nil { return objId, 0, err }
//        return objId, autoid, nil
//    } else {
//        return objId, int(resultUpdate["id"].(int32)) + 1, nil
//    }
//}

//不存在则插入，存在则更新      //强制所有数据使用bson.M
func (m *Mongo)Update(doc string, filter bson.M, data bson.M) (result bson.M, id int, err error) {
    var resultUpdate bson.M
    //处理更新时间问题
    if _, ok := data["$set"]; ok {
        if _, ok := data["$set"].(bson.M); ok {
            data["$set"].(bson.M)["update_at"] = time.Now()
        }
        if _, ok := data["$set"].(map[string]interface{}); ok {
            data["$set"].(map[string]interface{})["update_at"] = time.Now()
        }
    } else {
        data["$set"] = bson.M{"update_at": time.Now()}
    }
    opts := options.FindOneAndUpdate()
    err = m.coll(doc).FindOneAndUpdate(context.TODO(), filter, data, opts).Decode(&resultUpdate)
    if err != nil { return resultUpdate, 0, err }
    if _, ok := resultUpdate["id"].(int32); ok {
        return resultUpdate, int(resultUpdate["id"].(int32)), nil
    } else {
        return resultUpdate, int(resultUpdate["id"].(float64)), nil
    }
}

//不存在则插入，存在则更新      //强制所有数据使用bson.M
func (m *Mongo)Updates(doc string, filter bson.M, data bson.M) (updateCount int, err error) {
    //处理更新时间问题
    if _, ok := data["$set"]; ok {
        if _, ok := data["$set"].(bson.M); ok {
            data["$set"].(bson.M)["update_at"] = time.Now()
        }
        if _, ok := data["$set"].(map[string]interface{}); ok {
            data["$set"].(map[string]interface{})["update_at"] = time.Now()
        }
    } else {
        data["$set"] = bson.M{"update_at": time.Now()}
    }
    result, err := m.coll(doc).UpdateMany(context.TODO(), filter, data)
    if err != nil { return 0, err }
    return int(result.MatchedCount), nil
}

//批量插入数据
func (m *Mongo)Inserts(doc string, datas []interface{}) (ids []primitive.ObjectID, err error) {
    ids = []primitive.ObjectID{}
    opts := options.InsertMany().SetOrdered(false)
    res, err := m.coll(doc).InsertMany(context.TODO(), datas, opts)
    if err != nil { return ids, err }
    for _, row := range res.InsertedIDs {
        ids = append(ids, row.(primitive.ObjectID))
    }
    return ids, nil
}

func (m *Mongo)Get(doc string, filter interface{}) (result bson.M, err error) {
    opts := options.FindOne()
    err = m.coll(doc).FindOne(context.TODO(), filter, opts).Decode(&result)
    return result, err
    //if err != nil {
    //    if err == mongo.ErrNoDocuments { return result, nil }
    //    return result, err
    //}
    //return result, nil
}

func (m *Mongo)Gets(doc string, filter interface{}, opts ...map[string]interface{}) (results []bson.M, err error) {
    opt := options.Find()
    opt.SetMaxTime(30 * time.Second)
    if len(opts) > 0 {
        for key, v := range opts[0] {
            if key == "limit" || key == "pagesize" { opt.SetLimit(int64(v.(int))) }
            if key == "skip" || key == "offset" { opt.SetSkip(int64(v.(int))) }
            if key == "sort" { opt.SetSort(v) }
        }
    }
    cursor, err := m.coll(doc).Find(context.TODO(), filter, opt)
    if err != nil { return results, err }
    if err = cursor.All(context.TODO(), &results); err != nil { return results, err }
    return results, nil
}

//删除一条数据
func (m *Mongo)Delete(doc string, filter interface{}) (count int, err error) {
    //opts := options.Delete()
    //res, err := m.coll(doc).DeleteOne(context.TODO(), filter, opts)
    //if err != nil { return 0, err }
    //return int(res.DeletedCount), nil
    opts := options.Delete()
    res, err := m.coll(doc).DeleteMany(context.TODO(), filter, opts)
    if err != nil { return 0, err }
    return int(res.DeletedCount), nil
}

////删除多条数据
//func (m *Mongo)Deletes(doc string, filter interface{}) (count int, err error) {
//    opts := options.Delete()
//    res, err := m.coll(doc).DeleteMany(context.TODO(), filter, opts)
//    if err != nil { return 0, err }
//    return int(res.DeletedCount), nil
//}

func (m *Mongo)Count(doc string, filter interface{}) (count int, err error) {
    opts := options.Count().SetMaxTime(30 * time.Second)
    count64, err := m.coll(doc).CountDocuments(context.TODO(), filter, opts)
    if err != nil { return 0, err }
    return int(count64), err
}

func (m *Mongo)Group(doc string, filter interface{}, opts ...map[string]interface{}) (results []bson.M, err error) {
    opt := options.Aggregate()
    //cursor, err := m.coll(doc).Aggregate(context.TODO(), mongo.Pipeline{filter}, opt)
    cursor, err := m.coll(doc).Aggregate(context.TODO(), filter, opt)
    if err != nil { return results, err }
    err = cursor.All(context.TODO(), &results)
    if err != nil { return results, err }
    return results, nil
}
