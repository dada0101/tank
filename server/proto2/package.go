package proto2

import "TankDemo/network"


// Read(ProtocolBytes)
// Write(ProtocolBytes)
// Exec(agent)

type (
	Package interface {
		Init([]byte)
		Decode() bool
		Encode() bool
		Exec(a *network.Agent)
	}

	PackcageFactory struct {
		mapping map[string]*Package
	}
)

func NewPackageFactory() *PackcageFactory {
	return & PackcageFactory{
		make(map[string]*Package),
	}
}

func(pf *PackcageFactory)NewPackage(name string) *Package {
	pkg, ok := pf.mapping[name]
	if !ok {
		return nil
	}
	return pkg
}

type NilPackage struct {

}


