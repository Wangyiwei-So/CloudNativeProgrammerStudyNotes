// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64

package bpf

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

// loadToy_proxy__bpf returns the embedded CollectionSpec for toy_proxy__bpf.
func loadToy_proxy__bpf() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_Toy_proxy__bpfBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load toy_proxy__bpf: %w", err)
	}

	return spec, err
}

// loadToy_proxy__bpfObjects loads toy_proxy__bpf and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*toy_proxy__bpfObjects
//	*toy_proxy__bpfPrograms
//	*toy_proxy__bpfMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadToy_proxy__bpfObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadToy_proxy__bpf()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// toy_proxy__bpfSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type toy_proxy__bpfSpecs struct {
	toy_proxy__bpfProgramSpecs
	toy_proxy__bpfMapSpecs
}

// toy_proxy__bpfSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type toy_proxy__bpfProgramSpecs struct {
	TcEgress  *ebpf.ProgramSpec `ebpf:"tc_egress"`
	TcIngress *ebpf.ProgramSpec `ebpf:"tc_ingress"`
}

// toy_proxy__bpfMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type toy_proxy__bpfMapSpecs struct {
}

// toy_proxy__bpfObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadToy_proxy__bpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type toy_proxy__bpfObjects struct {
	toy_proxy__bpfPrograms
	toy_proxy__bpfMaps
}

func (o *toy_proxy__bpfObjects) Close() error {
	return _Toy_proxy__bpfClose(
		&o.toy_proxy__bpfPrograms,
		&o.toy_proxy__bpfMaps,
	)
}

// toy_proxy__bpfMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadToy_proxy__bpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type toy_proxy__bpfMaps struct {
}

func (m *toy_proxy__bpfMaps) Close() error {
	return _Toy_proxy__bpfClose()
}

// toy_proxy__bpfPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadToy_proxy__bpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type toy_proxy__bpfPrograms struct {
	TcEgress  *ebpf.Program `ebpf:"tc_egress"`
	TcIngress *ebpf.Program `ebpf:"tc_ingress"`
}

func (p *toy_proxy__bpfPrograms) Close() error {
	return _Toy_proxy__bpfClose(
		p.TcEgress,
		p.TcIngress,
	)
}

func _Toy_proxy__bpfClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed toy_proxy__bpf_bpfel_x86.o
var _Toy_proxy__bpfBytes []byte
