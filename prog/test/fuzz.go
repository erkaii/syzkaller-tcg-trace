// Copyright 2019 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package test

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/google/syzkaller/prog"
	_ "github.com/google/syzkaller/sys"
	"github.com/google/syzkaller/sys/targets"
)

func FuzzDeserialize(data []byte) int {
	p0, err0 := fuzzTarget.Deserialize(data, prog.NonStrict)
	p1, err1 := fuzzTarget.Deserialize(data, prog.Strict)
	if p0 == nil {
		if p1 != nil {
			panic("NonStrict is stricter than Strict")
		}
		if err0 == nil || err1 == nil {
			panic("no error")
		}
		return 0
	}
	if err0 != nil {
		panic("got program and error")
	}
	data0 := p0.Serialize()
	if p1 != nil {
		if err1 != nil {
			panic("got program and error")
		}
		if !bytes.Equal(data0, p1.Serialize()) {
			panic("got different data")
		}
	}
	p2, err2 := fuzzTarget.Deserialize(data0, prog.NonStrict)
	if err2 != nil {
		panic(fmt.Sprintf("failed to parse serialized: %v\n%s", err2, data0))
	}
	if !bytes.Equal(data0, p2.Serialize()) {
		panic("got different data")
	}
	p3 := p0.Clone()
	if !bytes.Equal(data0, p3.Serialize()) {
		panic("got different data")
	}
	if prodData, err := p0.SerializeForExec(); err == nil {
		if _, err := fuzzTarget.DeserializeExec(prodData, nil); err != nil {
			panic(err)
		}
	}
	p3.Mutate(rand.NewSource(0), 3, fuzzChoiceTable, nil, nil)
	return 0
}

func FuzzParseLog(data []byte) int {
	if len(fuzzTarget.ParseLog(data, prog.NonStrict)) != 0 {
		return 1
	}
	return 0
}

var fuzzTarget, fuzzChoiceTable = func() (*prog.Target, *prog.ChoiceTable) {
	target, err := prog.GetTarget(targets.TestOS, targets.TestArch64)
	if err != nil {
		panic(err)
	}
	return target, target.DefaultChoiceTable()
}()
