// Copyright 2017 The go-networkchain Authors
// This file is part of the go-networkchain library.
//
// The go-networkchain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-networkchain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-networkchain library. If not, see <http://www.gnu.org/licenses/>.

package testutil

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/networkchain/go-networkchain/swarm/api"
	httpapi "github.com/networkchain/go-networkchain/swarm/api/http"
	"github.com/networkchain/go-networkchain/swarm/storage"
)

func NewTestSwarmServer(t *testing.T) *TestSwarmServer {
	dir, err := ioutil.TempDir("", "swarm-storage-test")
	if err != nil {
		t.Fatal(err)
	}
	storeparams := &storage.StoreParams{
		ChunkDbPath:   dir,
		DbCapacity:    5000000,
		CacheCapacity: 5000,
		Radius:        0,
	}
	localStore, err := storage.NewLocalStore(storage.MakeHashFunc("SHA3"), storeparams)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	chunker := storage.NewTreeChunker(storage.NewChunkerParams())
	dpa := &storage.DPA{
		Chunker:    chunker,
		ChunkStore: localStore,
	}
	dpa.Start()
	a := api.NewApi(dpa, nil)
	srv := httptest.NewServer(httpapi.NewServer(a))
	return &TestSwarmServer{
		Server: srv,
		Dpa:    dpa,
		dir:    dir,
	}
}

type TestSwarmServer struct {
	*httptest.Server

	Dpa *storage.DPA
	dir string
}

func (t *TestSwarmServer) Close() {
	t.Server.Close()
	t.Dpa.Stop()
	os.RemoveAll(t.dir)
}
