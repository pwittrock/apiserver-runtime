/*
Copyright 2017 The Kubernetes Authors.

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

package apiserver

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	pkgserver "k8s.io/apiserver/pkg/server"
)

type StorageProvider func(s *runtime.Scheme, g genericregistry.RESTOptionsGetter) (rest.Storage, error)

var (
	GroupName           = "example.com"
	APIs                = map[schema.GroupVersionResource]StorageProvider{}
	GenericAPIServerFns []func(*pkgserver.GenericAPIServer) *pkgserver.GenericAPIServer
)

// buildStorageMap gets all of the registered APIs
func BuildStorageMap(s *runtime.Scheme, g genericregistry.RESTOptionsGetter) (map[string]map[string]rest.Storage, error) {
	apis := map[string]map[string]rest.Storage{}
	var err error
	for k, v := range APIs {
		if _, found := apis[k.Version]; !found {
			apis[k.Version] = map[string]rest.Storage{}
		}
		apis[k.Version][k.Resource], err = v(s, g)
		if err != nil {
			return nil, err
		}
	}
	return apis, nil
}

func ApplyGenericAPIServerFns(in *pkgserver.GenericAPIServer) *pkgserver.GenericAPIServer {
	for i := range GenericAPIServerFns {
		in = GenericAPIServerFns[i](in)
	}
	return in
}
