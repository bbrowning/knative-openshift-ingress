// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package google

import (
	"log"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
)

// WithTransport is a functional option for overriding the default transport
// on a remote image
func WithTransport(t http.RoundTripper) ListerOption {
	return func(l *lister) error {
		l.transport = t
		return nil
	}
}

// WithAuth is a functional option for overriding the default authenticator
// on a remote image
func WithAuth(auth authn.Authenticator) ListerOption {
	return func(l *lister) error {
		l.auth = auth
		return nil
	}
}

// WithAuthFromKeychain is a functional option for overriding the default
// authenticator on a remote image using an authn.Keychain
func WithAuthFromKeychain(keys authn.Keychain) ListerOption {
	return func(l *lister) error {
		auth, err := keys.Resolve(l.repo.Registry)
		if err != nil {
			return err
		}
		if auth == authn.Anonymous {
			log.Println("No matching credentials were found, falling back on anonymous")
		}
		l.auth = auth
		return nil
	}
}
