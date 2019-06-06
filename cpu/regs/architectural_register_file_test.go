package regs

import (
	"fmt"
	"testing"
)

func TestFloatingPointRegisters(t *testing.T) {
	var regs = NewArchitecturalRegisterFile(true)

	regs.Fpr.SetUint32(0, 100)
	regs.Fpr.SetUint64(1, 20000043)

	fmt.Printf("%d\n", regs.Fpr.Uint32(0))
	fmt.Printf("%d\n", regs.Fpr.Uint64(1))

	regs.Fpr.SetFloat32(4, 100.1)
	regs.Fpr.SetFloat64(5, 20000.043)

	fmt.Printf("%f\n", regs.Fpr.Float32(4))
	fmt.Printf("%f\n", regs.Fpr.Float64(5))
}
