package common

import (
	"crypto/tls"
	"errors"
	"net/rpc"
	"strings"
)

// Returns false if the two major versions mismatch
func CompatibleVersions(v1, v2 string) bool {
	major1 := strings.SplitN(v1, ".", 2)
	major2 := strings.SplitN(v2, ".", 2)
	return major1[0] == major2[0]
}

type RPCServerOpts interface {
	RPCHostAndPort() string
}

type BasicRPCServerOpts string

func (o BasicRPCServerOpts) RPCHostAndPort() string {
	return string(o)
}

type RPCClient struct {
	BaseName     string
	RPCVersion   string
	Opts         RPCServerOpts
	UseTLS       bool
	VersionError error
	VersionOk    bool
}

func NewRPCClient(hostAndPort, baseName, rpcVersion string, useTLS bool) *RPCClient {
	return NewRPCClientWithConfig(BasicRPCServerOpts(hostAndPort), baseName, rpcVersion, useTLS)
}

func NewRPCClientWithConfig(config RPCServerOpts, baseName, rpcVersion string, useTLS bool) *RPCClient {
	return &RPCClient{baseName, rpcVersion, config, useTLS, nil, false}
}

func (r *RPCClient) newClient() (*rpc.Client, error) {
	if r.UseTLS {
		return r.newTLSClient()
	}
	return rpc.DialHTTP("tcp", r.Opts.RPCHostAndPort())
}

func (r *RPCClient) tlsConfig() (*tls.Config, error) {
	var err error
	config := &tls.Config{}
	config.InsecureSkipVerify = true
	return config, err
}

func (r *RPCClient) newTLSClient() (*rpc.Client, error) {
	config, err := r.tlsConfig()
	if err != nil {
		panic(err)
	}
	conn, err := tls.Dial("tcp", r.Opts.RPCHostAndPort(), config)
	if err != nil {
		panic(err)
	}
	c := rpc.NewClient(conn)
	return c, err
}

func (r *RPCClient) checkVersion() error {
	if r.VersionOk {
		return nil
	}
	arg := VersionArg{}
	var reply VersionReply
	err := r.doRequest("Version", arg, &reply)
	if err != nil {
		r.VersionError = err
		r.VersionOk = false
		return err
	}
	if !CompatibleVersions(reply.RPCVersion, r.RPCVersion) {
		err := errors.New("Version Mismatch. Server: " + reply.RPCVersion + ", Client: " + r.RPCVersion)
		r.VersionError = err
		r.VersionOk = false
		return err
	}
	r.VersionOk = true
	return nil
}

func (r *RPCClient) doRequest(name string, arg interface{}, reply interface{}) error {
	client, err := r.newClient()
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Call(r.BaseName+"."+name, arg, reply)
}

func (r *RPCClient) Call(name string, arg interface{}, reply interface{}) error {
	if err := r.checkVersion(); err != nil {
		return err
	}
	return r.doRequest(name, arg, reply)
}
