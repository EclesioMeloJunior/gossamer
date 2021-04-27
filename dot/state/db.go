// Copyright 2019 ChainSafe Systems (ON) Corp.
// This file is part of gossamer.
//
// The gossamer library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The gossamer library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the gossamer library. If not, see <http://www.gnu.org/licenses/>.

package state

import (
	"encoding/json"
	"fmt"

	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/genesis"
	"github.com/ChainSafe/gossamer/lib/trie"

	database "github.com/ChainSafe/chaindb"
)

// SetupDatabase will return an instance of database based on basepath
func SetupDatabase(basepath string) (database.Database, error) {
	// initialise database using data directory
	db, err := database.NewBadgerDB(&database.Config{
		DataDir: basepath,
	})

	if err != nil {
		logger.Error(
			"failed to setup database",
			"basepath", basepath,
			"error", err,
		)
		return nil, err
	}

	return db, nil
}

// StoreNodeGlobalName stores the current node name to avoid create new ones after each initialization
func StoreNodeGlobalName(db database.Database, nodeName string) error {
	return db.Put(common.NodeNameKey, []byte(nodeName))
}

// LoadNodeGlobalName loads the latest stored node global name
func LoadNodeGlobalName(db database.Database) (string, error) {
	nodeName, err := db.Get(common.NodeNameKey)
	if err != nil {
		return "", err
	}

	return string(nodeName), nil
}

// StoreBestBlockHash stores the hash at the BestBlockHashKey
func StoreBestBlockHash(db database.Database, hash common.Hash) error {
	return db.Put(common.BestBlockHashKey, hash[:])
}

// LoadBestBlockHash loads the hash stored at BestBlockHashKey
func LoadBestBlockHash(db database.Database) (common.Hash, error) {
	hash, err := db.Get(common.BestBlockHashKey)
	if err != nil {
		return common.Hash{}, err
	}

	return common.NewHash(hash), nil
}

// StoreGenesisData stores the given genesis data at the known GenesisDataKey.
func StoreGenesisData(db database.Database, gen *genesis.Data) error {
	enc, err := json.Marshal(gen)
	if err != nil {
		return fmt.Errorf("cannot scale encode genesis data: %s", err)
	}

	return db.Put(common.GenesisDataKey, enc)
}

// LoadGenesisData retrieves the genesis data stored at the known GenesisDataKey.
func LoadGenesisData(db database.Database) (*genesis.Data, error) {
	enc, err := db.Get(common.GenesisDataKey)
	if err != nil {
		return nil, err
	}

	data := &genesis.Data{}
	err = json.Unmarshal(enc, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// StoreLatestStorageHash stores the current root hash in the database at LatestStorageHashKey
func StoreLatestStorageHash(db database.Database, root common.Hash) error {
	return db.Put(common.LatestStorageHashKey, root[:])
}

// LoadLatestStorageHash retrieves the hash stored at LatestStorageHashKey from the DB
func LoadLatestStorageHash(db database.Database) (common.Hash, error) {
	hashbytes, err := db.Get(common.LatestStorageHashKey)
	if err != nil {
		return common.Hash{}, err
	}

	return common.NewHash(hashbytes), nil
}

// StoreTrie encodes the entire trie and writes it to the DB
// The key to the DB entry is the root hash of the trie
func StoreTrie(db database.Database, t *trie.Trie) error {
	return t.Store(db)
}

// LoadTrie loads an encoded trie from the DB where the key is `root`
func LoadTrie(db database.Database, t *trie.Trie, root common.Hash) error {
	return t.Load(db, root)
}
