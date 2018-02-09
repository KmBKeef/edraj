package main

import (
	"fmt"
	"net/http"
	"path"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AddonsMan to manage the addons installed on the system
type AddonsMan struct {
	mongoDb   *mgo.Database
	fileStore Storage
}

func (man *AddonsMan) init(config *Config) (err error) {
	man.mongoDb = mongoSession.DB(addon)
	man.fileStore.RootPath = path.Join(config.dataPath, addon)
	man.fileStore.TrashPath = path.Join(config.dataPath, trash, addon)
	return
}

func (man *AddonsMan) query(request *Request) (response *QueryResponse) {
	response = &QueryResponse{}

	// TODO
	// Use request.EntryQuery
	err := man.mongoDb.C(addon).Find(bson.M{}).All(&response.Entries)
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to run query (%v). err: %s", request, err.Error())
	}

	response.Returned = int64(len(response.Entries))
	response.Status = succeeded
	response.Code = http.StatusFound

	return
}

func (man *AddonsMan) get(request *Request) (response *QueryResponse) {
	response = &QueryResponse{}

	if request.EntryID == "" {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request.entryid must have valid value  (%s)", request.EntryID)
		return
	}
	var entry Entry
	err := man.mongoDb.C(addon).FindId(bson.ObjectIdHex(request.EntryID)).One(&entry)
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to retrive entry (%s). err: %s", request.EntryID, err.Error())
	}

	response.Status = succeeded
	response.Code = http.StatusFound
	response.Entries = append(response.Entries, entry)
	response.Returned = 1
	response.Total = 1

	return
}

func (man *AddonsMan) create(request *Request) (response Response) {
	if request.ObjectType == "" {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request.entry  must have valid value  (%v)", request.Entry)
		return
	}
	entry := Entry{} // request.Entry
	//entry.ID = bson.NewObjectId()
	err := man.mongoDb.C(addon).Insert(&entry)
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to create entry (%v). err: %s", request.Entry, err.Error())
	}

	response.Status = succeeded
	response.Code = http.StatusCreated
	// TODO return the id of the created object

	return
}

func (man *AddonsMan) update(request *Request) (response Response) {
	if request.ObjectType == "" {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request.Entry must have valid value  (%v)", request.Entry)
		return
	}
	err := man.mongoDb.C(addon).Update(request.Entry.ID, &request.Entry)
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to create entry (%v). err: %s", request.Entry, err.Error())
	}

	response.Status = succeeded
	response.Code = http.StatusCreated
	// TODO return the id of the created object

	return
}

func (man *AddonsMan) delete(request *Request) (response Response) {

	if request.EntryID == "" {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request.entryid must have valid value  (%s)", request.EntryID)
		return
	}
	var entry Entry
	err := man.mongoDb.C(addon).Remove(&entry)
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to retrive entry (%s). err: %s", request.EntryID, err.Error())
	}

	response.Status = succeeded
	response.Code = http.StatusGone

	return
}
