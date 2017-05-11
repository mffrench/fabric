/*
Copyright IBM Corp. 2016, 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package statecouchdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/version"
	"github.com/hyperledger/fabric/core/ledger/ledgerconfig"
	"github.com/hyperledger/fabric/core/ledger/util/couchdb"
	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("statecouchdb")

var compositeKeySep = []byte{0x00}
var lastKeyIndicator = byte(0x01)

var binaryWrapper = "valueBytes"

//TODO querySkip is implemented for future use by query paging
//currently defaulted to 0 and is not used
var querySkip = 0

// VersionedDBProvider implements interface VersionedDBProvider
type VersionedDBProvider struct {
	couchInstance *couchdb.CouchInstance
	databases     map[string]*VersionedDB
	mux           sync.Mutex
	openCounts    uint64
}

// NewVersionedDBProvider instantiates VersionedDBProvider
func NewVersionedDBProvider() (*VersionedDBProvider, error) {
	logger.Debugf("constructing CouchDB VersionedDBProvider")
	couchDBDef := ledgerconfig.GetCouchDBDefinition()
	couchInstance, err := couchdb.CreateCouchInstance(couchDBDef.URL, couchDBDef.Username, couchDBDef.Password)
	if err != nil {
		return nil, err
	}

	return &VersionedDBProvider{couchInstance, make(map[string]*VersionedDB), sync.Mutex{}, 0}, nil
}

// GetDBHandle gets the handle to a named database
func (provider *VersionedDBProvider) GetDBHandle(dbName string) (statedb.VersionedDB, error) {
	provider.mux.Lock()
	defer provider.mux.Unlock()

	vdb := provider.databases[dbName]
	if vdb == nil {
		var err error
		vdb, err = newVersionedDB(provider.couchInstance, dbName)
		if err != nil {
			return nil, err
		}
		provider.databases[dbName] = vdb
	}
	return vdb, nil
}

// Close closes the underlying db instance
func (provider *VersionedDBProvider) Close() {
	// No close needed on Couch
}

// VersionedDB implements VersionedDB interface
type VersionedDB struct {
	db     *couchdb.CouchDatabase
	dbName string
	dbType string
}

// newVersionedDB constructs an instance of VersionedDB
func newVersionedDB(couchInstance *couchdb.CouchInstance, dbName string) (*VersionedDB, error) {
	// CreateCouchDatabase creates a CouchDB database object, as well as the underlying database if it does not exist
	db, err := couchdb.CreateCouchDatabase(*couchInstance, dbName)
	if err != nil {
		return nil, err
	}
	return &VersionedDB{db, dbName, "CouchDB"}, nil
}

// Get VersionedDB type
func (vdb *VersionedDB) GetVDBType() string {
	return vdb.dbType
}

// Open implements method in VersionedDB interface
func (vdb *VersionedDB) Open() error {
	// no need to open db since a shared couch instance is used
	return nil
}

// Close implements method in VersionedDB interface
func (vdb *VersionedDB) Close() {
	// no need to close db since a shared couch instance is used
}

func removeDataWrapper(wrappedValue []byte, attachments []couchdb.Attachment) ([]byte, version.Height) {

	logger.Debugf("wrappedValue: " + string(wrappedValue))

	//initialize the return value
	returnValue := []byte{} // TODO: empty byte or nil

	//initialize a default return version
	returnVersion := version.NewHeight(0, 0)

	//create a generic map for the json
	jsonResult := make(map[string]interface{})

	//unmarshal the selected json into the generic map
	json.Unmarshal(wrappedValue, &jsonResult)

	// handle binary or json data
	if jsonResult[dataWrapper] == nil && attachments != nil { // binary attachment
		// get binary data from attachment
		for _, attachment := range attachments {
			if attachment.Name == binaryWrapper {
				returnValue = attachment.AttachmentBytes
			}
		}
	} else {
		//place the result json in the data key
		returnMap := jsonResult[dataWrapper]

		//marshal the mapped data.   this wrappers the result in a key named "data"
		returnValue, _ = json.Marshal(returnMap)

	}

	//create an array containing the blockNum and txNum
	logger.Debugf("jsonResult[version]: " + fmt.Sprintf("%s", jsonResult["version"]))
	versionArray := strings.Split(fmt.Sprintf("%s", jsonResult["version"]), ":")

	//convert the blockNum from String to unsigned int
	blockNum, _ := strconv.ParseUint(versionArray[0], 10, 64)

	//convert the txNum from String to unsigned int
	txNum, _ := strconv.ParseUint(versionArray[1], 10, 64)

	//create the version based on the blockNum and txNum
	returnVersion = version.NewHeight(blockNum, txNum)

	return returnValue, *returnVersion

}

// GetState implements method in VersionedDB interface
func (vdb *VersionedDB) GetState(namespace string, key string) (*statedb.VersionedValue, error) {
	logger.Debugf("GetState(). ns=%s, key=%s", namespace, key)

	compositeKey := constructCompositeKey(namespace, key)

	logger.Infof("Get State through ReadDoc: %s", string(compositeKey))
	couchDoc, _, err := vdb.db.ReadDoc(string(compositeKey))
	if err != nil {
		return nil, err
	}
	if couchDoc == nil {
		return nil, nil
	}

	//remove the data wrapper and return the value and version
	returnValue, returnVersion := removeDataWrapper(couchDoc.JSONValue, couchDoc.Attachments)

	return &statedb.VersionedValue{Value: returnValue, Version: &returnVersion}, nil
}

// GetStateMultipleKeys implements method in VersionedDB interface
func (vdb *VersionedDB) GetStateMultipleKeys(namespace string, keys []string) ([]*statedb.VersionedValue, error) {

	// first : define document ids list keys
	allGet := couchdb.DocsAllKeys{}
	for _, key := range keys {
		allGet.Keys = append(allGet.Keys, string(constructCompositeKey(namespace, key)))
	}

	// second : get documents revisions in one shoot (CouchDB OP)
	compositeKeysDocMap, err := vdb.db.ReadDocsKeys(allGet)
	logger.Infof("compositeKeysDocMap: %+v\n", compositeKeysDocMap)
	if err != nil {
		logger.Errorf("Error during ReadDocsKeys(): %s\n", err.Error())
		return nil, err
	}

	// third : build result
	vals := make([]*statedb.VersionedValue, len(keys))
	idx := 0
	for _, couchDoc := range compositeKeysDocMap {
		if couchDoc == nil {
			vals[idx] = nil
		} else if  len(couchDoc.JSONValue) == 0 {
			vals[idx] = nil
		} else {
			returnValue, returnVersion := removeDataWrapper(couchDoc.JSONValue, couchDoc.Attachments)
			vals[idx] = &statedb.VersionedValue{Value: returnValue, Version: &returnVersion}
		}
		idx++
	}
	return vals, nil
}

// GetKStateByMultipleKeys implements method in VersionedDB interface
func (vdb *VersionedDB) GetKStateByMultipleKeys(namespace string, keys []string) (map[string]*statedb.VersionedValue, error) {

	// first : define document ids list keys
	allGet := couchdb.DocsAllKeys{}
	for _, key := range keys {
		allGet.Keys = append(allGet.Keys, string(constructCompositeKey(namespace, key)))
	}

	// second : get documents revisions in one shoot (CouchDB OP)
	compositeKeysDocMap, err := vdb.db.ReadDocsKeys(allGet)
	logger.Infof("compositeKeysDocMap: %+v\n", compositeKeysDocMap)
	if err != nil {
		logger.Errorf("Error during ReadDocsKeys(): %s\n", err.Error())
		return nil, err
	}

	// third : build result
	vals := map[string]*statedb.VersionedValue{}
	for cKey, couchDoc := range compositeKeysDocMap {
		_, key := splitCompositeKey([]byte(cKey))
		if couchDoc == nil {
			vals[key] = nil
		} else if  len(couchDoc.JSONValue) == 0 {
			vals[key] = nil
		} else {
			returnValue, returnVersion := removeDataWrapper(couchDoc.JSONValue, couchDoc.Attachments)
			vals[key] = &statedb.VersionedValue{Value: returnValue, Version: &returnVersion}
		}
	}
	return vals, nil
}

// GetStateRangeScanIterator implements method in VersionedDB interface
// startKey is inclusive
// endKey is exclusive
func (vdb *VersionedDB) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (statedb.ResultsIterator, error) {

	//Get the querylimit from core.yaml
	queryLimit := ledgerconfig.GetQueryLimit()

	compositeStartKey := constructCompositeKey(namespace, startKey)
	compositeEndKey := constructCompositeKey(namespace, endKey)
	if endKey == "" {
		compositeEndKey[len(compositeEndKey)-1] = lastKeyIndicator
	}
	queryResult, err := vdb.db.ReadDocRange(string(compositeStartKey), string(compositeEndKey), queryLimit, querySkip)
	if err != nil {
		logger.Debugf("Error calling ReadDocRange(): %s\n", err.Error())
		return nil, err
	}
	logger.Debugf("Exiting GetStateRangeScanIterator")
	return newKVScanner(namespace, *queryResult), nil

}

// ExecuteQuery implements method in VersionedDB interface
func (vdb *VersionedDB) ExecuteQuery(namespace, query string) (statedb.ResultsIterator, error) {

	//Get the querylimit from core.yaml
	queryLimit := ledgerconfig.GetQueryLimit()

	queryString, err := ApplyQueryWrapper(namespace, query, queryLimit, 0)
	if err != nil {
		logger.Debugf("Error calling ApplyQueryWrapper(): %s\n", err.Error())
		return nil, err
	}

	queryResult, err := vdb.db.QueryDocuments(queryString)
	if err != nil {
		logger.Debugf("Error calling QueryDocuments(): %s\n", err.Error())
		return nil, err
	}
	logger.Debugf("Exiting ExecuteQuery")
	return newQueryScanner(*queryResult), nil
}

// ApplyUpdates implements method in VersionedDB interface
func (vdb *VersionedDB) ApplyUpdates(batch *statedb.UpdateBatch, height *version.Height) error {
	if namespaces := batch.GetUpdatedNamespaces(); len(namespaces) <= 1 {
		if len(namespaces) == 1 {
			if updates := batch.GetUpdates(namespaces[0]); len(updates) > 1 {
				return applyUpdatesBulk(vdb, batch, height)
			} else {
				return applyUpdatesUnit(vdb, batch, height)
			}
		} else {
			return applyUpdatesUnit(vdb, batch, height)
		}
	} else {
		if len(namespaces) == 1 {
			if updates := batch.GetUpdates(namespaces[0]); len(updates) > 1 {
				return applyUpdatesBulk(vdb, batch, height)
			} else {
				return applyUpdatesUnit(vdb, batch, height)
			}
		} else {
			return applyUpdatesBulk(vdb, batch, height)
		}
	}
	return applyUpdatesUnit(vdb, batch, height)
}

// 1.0.0-alpha implementation for ApplyUpdates
func applyUpdatesUnit(vdb *VersionedDB, batch *statedb.UpdateBatch, height *version.Height) error {
	logger.Infof("Entering applyUpdatesUnit")
	updateDocCount := 0
	namespaces := batch.GetUpdatedNamespaces()
	for _, ns := range namespaces {
		updates := batch.GetUpdates(ns)
		for k, vv := range updates {
			compositeKey := constructCompositeKey(ns, k)
			logger.Debugf("Channel [%s]: Applying key=[%#v]", vdb.dbName, compositeKey)

			//convert nils to deletes
			if vv.Value == nil {

				vdb.db.DeleteDoc(string(compositeKey), "")

			} else {
				couchDoc := &couchdb.CouchDoc{}

				//Check to see if the value is a valid JSON
				//If this is not a valid JSON, then store as an attachment
				if couchdb.IsJSON(string(vv.Value)) {
					// Handle it as json
					couchDoc.JSONValue = addVersionAndChainCodeID(vv.Value, ns, vv.Version)
				} else { // if the data is not JSON, save as binary attachment in Couch
					//Create an attachment structure and load the bytes
					attachment := &couchdb.Attachment{}
					attachment.AttachmentBytes = vv.Value
					attachment.ContentType = "application/octet-stream"
					attachment.Name = binaryWrapper
					couchDoc.Attachments = append(couchDoc.Attachments, *attachment)
					couchDoc.JSONValue = addVersionAndChainCodeID(nil, ns, vv.Version)
				}

				// SaveDoc using couchdb client and use attachment to persist the binary data
				rev, err := vdb.db.SaveDoc(string(compositeKey), "", couchDoc)
				if err != nil {
					logger.Errorf("Error during Commit(): %s\n", err.Error())
					return err
				}
				if rev != "" {
					logger.Debugf("Saved document revision number: %s\n", rev)
				}
			}
			updateDocCount++
		}
	}

	// Record a savepoint at a given height
	err := vdb.recordSavepoint(height)
	if err != nil {
		logger.Errorf("Error during recordSavepoint: %s\n", err.Error())
		return err
	}
	logger.Infof("Exiting applyUpdatesUnit %d", updateDocCount)
	return nil
}

// improvement try with CouchDB bulk operation for ApplyUpdates
func applyUpdatesBulk(vdb *VersionedDB, batch *statedb.UpdateBatch, height *version.Height) error {
	logger.Infof("Entering applyUpdatesBulk")
	bulkDocs := couchdb.DocsBulk{}
	namespaces := batch.GetUpdatedNamespaces()
	updateDocCount := 0

	// first : define document ids list from batch
	allInsertUpdate := couchdb.DocsAllKeys{}
	for _, ns := range namespaces {
		updates := batch.GetUpdates(ns)
		for k := range updates {
			allInsertUpdate.Keys = append(allInsertUpdate.Keys, string(constructCompositeKey(ns, k)))
		}
	}

	// second : get documents revisions in one shoot (CouchDB OP)
	compositeKeysRevMap, err := vdb.db.ReadRevsKeys(allInsertUpdate)
	logger.Debugf("compositeKeysRevMap: %+v\n", compositeKeysRevMap)
	if err != nil {
		logger.Errorf("Error during ReadDocsRev(): %s\n", err.Error())
		return err
	}

	// third : build bulkDocs
	for _, ns := range namespaces {
		updates := batch.GetUpdates(ns)
		for k, vv := range updates {
			compositeKey := constructCompositeKey(ns, k)
			if vv.Value == nil {
				docUnit := &couchdb.DocBulkUnit{}
				docUnit.Id = string(compositeKey)
				docUnit.Delete = true
				docUnit.Rev = compositeKeysRevMap[docUnit.Id]
				bulkDocs.Docs = append(bulkDocs.Docs, docUnit)
			} else {
				docUnit := &couchdb.DocBulkUnit{}
				docUnit.Id = string(compositeKey)
				if val, ok := compositeKeysRevMap[docUnit.Id]; ok {
					docUnit.Rev = val
				}
				docUnit.Version = fmt.Sprintf("%v:%v", vv.Version.BlockNum, vv.Version.TxNum)
				docUnit.Chaincodeid = ns
				docUnit.Data = (*json.RawMessage)(&vv.Value)
				logger.Debugf("new docUnit in bulk: %+v\n", docUnit)
				bulkDocs.Docs = append(bulkDocs.Docs, docUnit)
			}
			updateDocCount++
		}
	}

	// fourth : apply operations on bulkDocs (CouchDB OP)
	_, err = vdb.db.BulkDocs(bulkDocs)
	if err != nil {
		logger.Errorf("Error during Commit(): %s\n", err.Error())
		return err
	}

	// finally : record a savepoint at a given height (CouchDB OP)
	err = vdb.recordSavepoint(height)
	if err != nil {
		logger.Errorf("Error during recordSavepoint: %s\n", err.Error())
		return err
	}
	logger.Infof("Exiting applyUpdatesBulk %d", updateDocCount)

	return nil
}

//addVersionAndChainCodeID adds keys for version and chaincodeID to the JSON value
func addVersionAndChainCodeID(value []byte, chaincodeID string, version *version.Height) []byte {

	//create a version mapping
	jsonMap := map[string]interface{}{"version": fmt.Sprintf("%v:%v", version.BlockNum, version.TxNum)}

	//add the chaincodeID
	jsonMap["chaincodeid"] = chaincodeID

	//Add the wrapped data if the value is not null
	if value != nil {

		//create a new genericMap
		rawJSON := (*json.RawMessage)(&value)

		//add the rawJSON to the map
		jsonMap[dataWrapper] = rawJSON

	}

	//marshal the data to a byte array
	returnJSON, _ := json.Marshal(jsonMap)

	return returnJSON

}

// Savepoint docid (key) for couchdb
const savepointDocID = "statedb_savepoint"

// Savepoint data for couchdb
type couchSavepointData struct {
	BlockNum  uint64 `json:"BlockNum"`
	TxNum     uint64 `json:"TxNum"`
	UpdateSeq string `json:"UpdateSeq"`
}

// recordSavepoint Record a savepoint in statedb.
// Couch parallelizes writes in cluster or sharded setup and ordering is not guaranteed.
// Hence we need to fence the savepoint with sync. So ensure_full_commit is called before AND after writing savepoint document
// TODO: Optimization - merge 2nd ensure_full_commit with savepoint by using X-Couch-Full-Commit header
func (vdb *VersionedDB) recordSavepoint(height *version.Height) error {
	var err error
	var savepointDoc couchSavepointData
	// ensure full commit to flush all changes until now to disk
	dbResponse, err := vdb.db.EnsureFullCommit()
	if err != nil || dbResponse.Ok != true {
		logger.Errorf("Failed to perform full commit\n")
		return errors.New("Failed to perform full commit")
	}

	// construct savepoint document
	// UpdateSeq would be useful if we want to get all db changes since a logical savepoint
	dbInfo, _, err := vdb.db.GetDatabaseInfo()
	if err != nil {
		logger.Errorf("Failed to get DB info %s\n", err.Error())
		return err
	}
	savepointDoc.BlockNum = height.BlockNum
	savepointDoc.TxNum = height.TxNum
	savepointDoc.UpdateSeq = dbInfo.UpdateSeq

	savepointDocJSON, err := json.Marshal(savepointDoc)
	if err != nil {
		logger.Errorf("Failed to create savepoint data %s\n", err.Error())
		return err
	}

	// SaveDoc using couchdb client and use JSON format
	_, err = vdb.db.SaveDoc(savepointDocID, "", &couchdb.CouchDoc{JSONValue: savepointDocJSON, Attachments: nil})
	if err != nil {
		logger.Errorf("Failed to save the savepoint to DB %s\n", err.Error())
		return err
	}

	// ensure full commit to flush savepoint to disk
	dbResponse, err = vdb.db.EnsureFullCommit()
	if err != nil || dbResponse.Ok != true {
		logger.Errorf("Failed to perform full commit\n")
		return errors.New("Failed to perform full commit")
	}
	return nil
}

// GetLatestSavePoint implements method in VersionedDB interface
func (vdb *VersionedDB) GetLatestSavePoint() (*version.Height, error) {

	var err error
	logger.Infof("Get latest save point through ReadDoc: %s", savepointDocID)
	couchDoc, _, err := vdb.db.ReadDoc(savepointDocID)
	if err != nil {
		logger.Errorf("Failed to read savepoint data %s\n", err.Error())
		return nil, err
	}

	// ReadDoc() not found (404) will result in nil response, in these cases return height nil
	if couchDoc == nil || couchDoc.JSONValue == nil {
		return nil, nil
	}

	savepointDoc := &couchSavepointData{}
	err = json.Unmarshal(couchDoc.JSONValue, &savepointDoc)
	if err != nil {
		logger.Errorf("Failed to unmarshal savepoint data %s\n", err.Error())
		return nil, err
	}

	return &version.Height{BlockNum: savepointDoc.BlockNum, TxNum: savepointDoc.TxNum}, nil
}

func constructCompositeKey(ns string, key string) []byte {
	compositeKey := []byte(ns)
	compositeKey = append(compositeKey, compositeKeySep...)
	compositeKey = append(compositeKey, []byte(key)...)
	return compositeKey
}

func splitCompositeKey(compositeKey []byte) (string, string) {
	split := bytes.SplitN(compositeKey, compositeKeySep, 2)
	return string(split[0]), string(split[1])
}

type kvScanner struct {
	cursor    int
	namespace string
	results   []couchdb.QueryResult
}

func newKVScanner(namespace string, queryResults []couchdb.QueryResult) *kvScanner {
	return &kvScanner{-1, namespace, queryResults}
}

func (scanner *kvScanner) Next() (statedb.QueryResult, error) {

	scanner.cursor++

	if scanner.cursor >= len(scanner.results) {
		return nil, nil
	}

	selectedKV := scanner.results[scanner.cursor]

	_, key := splitCompositeKey([]byte(selectedKV.ID))

	//remove the data wrapper and return the value and version
	returnValue, returnVersion := removeDataWrapper(selectedKV.Value, selectedKV.Attachments)

	return &statedb.VersionedKV{
		CompositeKey:   statedb.CompositeKey{Namespace: scanner.namespace, Key: key},
		VersionedValue: statedb.VersionedValue{Value: returnValue, Version: &returnVersion}}, nil
}

func (scanner *kvScanner) Close() {
	scanner = nil
}

type queryScanner struct {
	cursor  int
	results []couchdb.QueryResult
}

func newQueryScanner(queryResults []couchdb.QueryResult) *queryScanner {
	return &queryScanner{-1, queryResults}
}

func (scanner *queryScanner) Next() (statedb.QueryResult, error) {

	scanner.cursor++

	if scanner.cursor >= len(scanner.results) {
		return nil, nil
	}

	selectedResultRecord := scanner.results[scanner.cursor]

	namespace, key := splitCompositeKey([]byte(selectedResultRecord.ID))

	//remove the data wrapper and return the value and version
	returnValue, returnVersion := removeDataWrapper(selectedResultRecord.Value, selectedResultRecord.Attachments)

	return &statedb.VersionedQueryRecord{
		Namespace: namespace,
		Key:       key,
		Version:   &returnVersion,
		Record:    returnValue}, nil
}

func (scanner *queryScanner) Close() {
	scanner = nil
}
