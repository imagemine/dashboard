// Copyright 2017 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
  "github.com/kubernetes/dashboard/src/app/backend/client/envvar"
  "time"

	"k8s.io/client-go/tools/clientcmd/api"
)

// Resource information that are used as encryption key storage. Can be accessible by multiple dashboard replicas.
var EncryptionKeyHolderName = envvar.EnvVariable("KEY_HOLDER_NAME", "kubernetes-dashboard-key-holder")
// Resource information that are used as certificate storage for custom certificates used by the user.
var CertificateHolderSecretName = envvar.EnvVariable("DASHBOARD_CERTS_NAME", "kubernetes-dashboard-certs")

const (


	// Expiration time (in seconds) of tokens generated by dashboard. Default: 15 min.
	DefaultTokenTTL = 900
)

// AuthenticationModes represents auth modes supported by dashboard.
type AuthenticationModes map[AuthenticationMode]bool

// ProtectedResource represents basic information about resource that should be filtered out from Dashboard UI.
type ProtectedResource struct {
	// ResourceName is a name of the protected resource.
	ResourceName string
	// ResourceNamespace is a namespace of the protected resource. Should be empty if resource is non-namespaced.
	ResourceNamespace string
}

// IsEnabled returns true if given auth mode is supported, false otherwise.
func (self AuthenticationModes) IsEnabled(mode AuthenticationMode) bool {
	_, exists := self[mode]
	return exists
}

// Array returns array of auth modes supported by dashboard.
func (self AuthenticationModes) Array() []AuthenticationMode {
	modes := []AuthenticationMode{}
	for mode := range self {
		modes = append(modes, mode)
	}

	return modes
}

// Add adds given auth mode to AuthenticationModes map
func (self AuthenticationModes) Add(mode AuthenticationMode) {
	self[mode] = true
}

// AuthenticationMode represents auth mode supported by dashboard, i.e. basic.
type AuthenticationMode string

// String returns string representation of auth mode.
func (self AuthenticationMode) String() string {
	return string(self)
}

// Authentication modes supported by dashboard should be defined below.
const (
	Token AuthenticationMode = "token"
	Basic AuthenticationMode = "basic"
)

// AuthManager is used for user authentication management.
type AuthManager interface {
	// Login authenticates user based on provided LoginSpec and returns AuthResponse. AuthResponse contains
	// generated token and list of non-critical errors such as 'Failed authentication'.
	Login(*LoginSpec) (*AuthResponse, error)
	// Refresh takes valid token that hasn't expired yet and returns a new one with expiration time set to TokenTTL. In
	// case provided token has expired, token expiration error is returned.
	Refresh(string) (string, error)
	// AuthenticationModes returns array of auth modes supported by dashboard.
	AuthenticationModes() []AuthenticationMode
	// AuthenticationSkippable tells if the Skip button should be enabled or not
	AuthenticationSkippable() bool
}

// TokenManager is responsible for generating and decrypting tokens used for authorization. Authorization is handled
// by K8S apiserver. Token contains AuthInfo structure used to create K8S api client.
type TokenManager interface {
	// Generate secure token based on AuthInfo structure and save it tokens' payload.
	Generate(api.AuthInfo) (string, error)
	// Decrypt generated token and return AuthInfo structure that will be used for K8S api client creation.
	Decrypt(string) (*api.AuthInfo, error)
	// Refresh returns refreshed token based on provided token. In case provided token has expired, token expiration
	// error is returned.
	Refresh(string) (string, error)
	// SetTokenTTL sets expiration time (in seconds) of generated tokens.
	SetTokenTTL(time.Duration)
}

// Authenticator represents authentication methods supported by Dashboard. Currently supported types are:
//   - Token based - Any bearer token accepted by apiserver
//   - Basic - Username and password based authentication. Requires that apiserver has basic auth enabled also
//   - Kubeconfig based - Authenticates user based on kubeconfig file. Only token/basic modes are supported within
//     the kubeconfig file.
type Authenticator interface {
	// GetAuthInfo returns filled AuthInfo structure that can be used for K8S api client creation.
	GetAuthInfo() (api.AuthInfo, error)
}

// LoginSpec is extracted from request coming from Dashboard frontend during login request. It contains all the
// information required to authenticate user.
type LoginSpec struct {
	// Username is the username for basic authentication to the kubernetes cluster.
	Username string `json:"username,omitempty"`
	// Password is the password for basic authentication to the kubernetes cluster.
	Password string `json:"password,omitempty"`
	// Token is the bearer token for authentication to the kubernetes cluster.
	Token string `json:"token,omitempty"`
	// KubeConfig is the content of users' kubeconfig file. It will be parsed and auth data will be extracted.
	// Kubeconfig can not contain any paths. All data has to be provided within the file.
	KubeConfig string `json:"kubeconfig,omitempty"`
}

// AuthResponse is returned from our backend as a response for login/refresh requests. It contains generated JWEToken
// and a list of non-critical errors such as 'Failed authentication'.
type AuthResponse struct {
	// Name is a user/subject name if available
	Name string `json:"name,omitempty"`
	// JWEToken is a token generated during login request that contains AuthInfo data in the payload.
	JWEToken string `json:"jweToken"`
	// Errors are a list of non-critical errors that happened during login request.
	Errors []error `json:"errors"`
}

// TokenRefreshSpec contains token that is required by token refresh operation.
type TokenRefreshSpec struct {
	// JWEToken is a token generated during login request that contains AuthInfo data in the payload.
	JWEToken string `json:"jweToken"`
}

// LoginModesResponse contains list of auth modes supported by dashboard.
type LoginModesResponse struct {
	Modes []AuthenticationMode `json:"modes"`
}

// LoginSkippableResponse contains a flag that tells the UI not to display the Skip button.
// Note that this only hides the button, it doesn't disable unauthenticated access.
type LoginSkippableResponse struct {
	Skippable bool `json:"skippable"`
}
