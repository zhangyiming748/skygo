### ObjectId
primitive.ObjectID

primitive.ObjectIDFromHex(idStr)


### 使用mongo2提供的widget
```go
collection := mongo2.GetDefaultMongoDatabase().Collection(common.MC_PROJECT)
id, _ := primitive.ObjectIDFromHex("60769c4be13823b1df47eca4")

widget := mongo2.Widget{}
param := qmap.QM{
    "e__id" : id,
}
widget.SetTransformerFunc(func(qm qmap.QM) qmap.QM{
    qm["aaaa"] = "aaaaaaaaaaa"
    return qm
})
widget.SetQueryArray(param)
a , b := widget.One(collection)
fmt.Println(a,b)
```

### mongo包支持事务
```go
sess, err := orm_mongo.GetMongoClient().StartSession()
if err != nil {
	panic(err)
}
defer sess.EndSession(context.TODO())
sessCtx := mongo.NewSessionContext(context.TODO(), sess)

// Start a transaction and sessCtx as the Context parameter to InsertOne and
// FindOne so both operations will be run in the transaction.
if err = sess.StartTransaction(); err != nil {
	panic(err)
}

coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_PROJECT)
res, err := coll.InsertOne(sessCtx, bson.M{"_id" : primitive.NewObjectID(), "a" : "a"})
if err != nil {
	_ = sess.AbortTransaction(context.Background())
	panic(err)
}

var result bson.M
err = coll.FindOne(sessCtx, bson.M{"_id": res.InsertedID}).Decode(&result)
if err != nil {
	panic(err)
}

res, err = coll.InsertOne(sessCtx, bson.M{"_id" : primitive.NewObjectID(), "bbbb" : "bbbb"})
if err != nil {
	_ = sess.AbortTransaction(context.Background())
	return
}

if err = sess.CommitTransaction(context.Background()); err != nil {
	panic(err)
}

return
```




		       
