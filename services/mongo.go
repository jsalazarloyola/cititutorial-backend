package services

import (
    "context"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type TasksService struct {
    DatabaseName string
    CollectionName string
    Client *mongo.Client
    Collection *mongo.Collection
}

// Crea una estructura  para almacenar información del cliente y la tarea
func NewTasksService(client *mongo.Client, dbname string, name string) *TasksService {
    return &TasksService{
        DatabaseName: dbname,
        CollectionName: name,
        Client: client,
        Collection: GetCollection(client, dbname, name),
    }
}

// Estima la cantidad de documentos en la colección
func (ts TasksService) EstimateTotalDocs() int64 {
    count, err := ts.Collection.EstimatedDocumentCount(context.TODO())
    if err != nil {
        log.Print("Error al contar documentos: ", err.Error())
    }

    return count
}


// Inicia y conecta con MongoDB
func InitMongo() *mongo.Client {
    // 10 segundos de timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    var client *mongo.Client
    client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal("MongoDB connection error:", err)
        return nil
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to MongoDB")

    return client
}

// Oobtiene la colección, aunque el nombre está algo hardcodeaedo
func GetCollection(client *mongo.Client, dbname string, collectionName string) *mongo.Collection {
    collection := client.Database(dbname).Collection(collectionName)
    return collection
}
