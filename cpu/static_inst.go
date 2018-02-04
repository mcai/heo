package cpu

type StaticInst struct {
	Mnemonic *Mnemonic
	MachInst MachInst

	InputDependencies  []uint32
	OutputDependencies []uint32

	NumPhysicalRegistersToAllocate map[RegisterDependencyType]uint32
}

func NewStaticInst(mnemonic *Mnemonic, machInst MachInst) *StaticInst {
	var staticInst = &StaticInst{
		Mnemonic: mnemonic,
		MachInst: machInst,

		InputDependencies:  mnemonic.GetInputDependencies(machInst),
		OutputDependencies: mnemonic.GetOutputDependencies(machInst),

		NumPhysicalRegistersToAllocate: make(map[RegisterDependencyType]uint32),
	}

	for _, outputDependency := range staticInst.OutputDependencies {
		if outputDependency != 0 {
			var outputDependencyType, _ = RegisterDependencyFromInt(outputDependency)

			if _, ok := staticInst.NumPhysicalRegistersToAllocate[outputDependencyType]; !ok {
				staticInst.NumPhysicalRegistersToAllocate[outputDependencyType] = 0
			}

			staticInst.NumPhysicalRegistersToAllocate[outputDependencyType]++
		}
	}

	return staticInst
}

func (staticInst *StaticInst) Execute(context *Context) {
	staticInst.Mnemonic.Execute(context, staticInst.MachInst)
}

func (staticInst *StaticInst) Disassemble(pc uint32) string {
	return Disassemble(pc, string(staticInst.Mnemonic.Name), staticInst.MachInst)
}
