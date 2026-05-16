package helper

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	paths "github.com/dogukanmeral/scx-adapt/internal"
)

func LoadBPFScx(filepath string, structName string) error {
	collSpec, err := ebpf.LoadCollectionSpec(filepath)
	if err != nil {
		return err
	}

	coll, err := ebpf.NewCollection(collSpec)
	if err != nil {
		return err
	}

	schedOpsMap := coll.Maps[structName]
	if schedOpsMap == nil {
		return fmt.Errorf("Map not found: %s", structName)
	}

	l, err := link.AttachStructOps(link.StructOpsOptions{
		Map: schedOpsMap,
	})
	if err != nil {
		return err
	}

	err = l.Pin(paths.SCHEDBPFPINPATH)
	if err != nil {
		return err
	}

	return nil
}
