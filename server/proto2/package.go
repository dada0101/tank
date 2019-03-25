package proto2

import "TankDemo/network"

type (
	Package interface {
		Read(stream *BufferStream) bool
		Write(stream *BufferStream) bool
		Exec(a *network.Agent) *BufferStream
	}

	PackageFactory struct {
		mapping map[string]*Package
	}
)

func NewPackageFactory() *PackageFactory {
	return & PackageFactory{
		make(map[string]*Package),
	}
}

func(pf *PackageFactory)NewPackage(name string) *Package {
	pkg, ok := pf.mapping[name]
	if !ok {
		return nil
	}
	return pkg
}

type NilPackage struct {

}

