package cpu

import "fmt"

type PhysicalRegisterState string

const (
	PhysicalRegisterState_AVAILABLE              = PhysicalRegisterState("AVAILABLE")
	PhysicalRegisterState_RENAME_BUFFER_INVALID  = PhysicalRegisterState("RENAME_BUFFER_INVALID")
	PhysicalRegisterState_RENAME_BUFFER_VALID    = PhysicalRegisterState("RENAME_BUFFER_VALID")
	PhysicalRegisterState_ARCHITECTURAL_REGISTER = PhysicalRegisterState("ARCHITECTURAL_REGISTER")
)

type PhysicalRegister struct {
	PhysicalRegisterFile                         *PhysicalRegisterFile
	Num                                          uint32
	State                                        PhysicalRegisterState
	Allocator                                    *ReorderBufferEntry
	Dependency                                   int32
	EffectiveAddressComputationOperandDependents []*ReorderBufferEntry
	StoreAddressDependents                       []*LoadStoreQueueEntry
	Dependents                                   []GeneralReorderBufferEntry
}

func NewPhysicalRegister(physicalRegisterFile *PhysicalRegisterFile, num uint32) *PhysicalRegister {
	var physicalReg = &PhysicalRegister{
		PhysicalRegisterFile: physicalRegisterFile,
		Num:                  num,
		State:                PhysicalRegisterState_AVAILABLE,
		Dependency:           -1,
	}

	return physicalReg
}

func (physicalReg *PhysicalRegister) Reserve(dependency uint32) {
	if physicalReg.State != PhysicalRegisterState_AVAILABLE {
		panic("Impossible")
	}

	physicalReg.Dependency = int32(dependency)

	physicalReg.State = PhysicalRegisterState_ARCHITECTURAL_REGISTER

	physicalReg.PhysicalRegisterFile.NumFreePhysicalRegisters--
}

func (physicalReg *PhysicalRegister) Allocate(allocator *ReorderBufferEntry, dependency uint32) {
	if physicalReg.State != PhysicalRegisterState_AVAILABLE {
		panic("Impossible")
	}

	physicalReg.Allocator = allocator

	physicalReg.Dependency = int32(dependency)

	physicalReg.State = PhysicalRegisterState_RENAME_BUFFER_INVALID

	physicalReg.PhysicalRegisterFile.NumFreePhysicalRegisters--
}

func (physicalReg *PhysicalRegister) Writeback() {
	if physicalReg.State != PhysicalRegisterState_RENAME_BUFFER_INVALID {
		panic("Impossible")
	}

	physicalReg.State = PhysicalRegisterState_RENAME_BUFFER_VALID

	for _, effectiveAddressComputationOperandDependent := range physicalReg.EffectiveAddressComputationOperandDependents {
		effectiveAddressComputationOperandDependent.EffectiveAddressComputationOperandReady = true
	}

	for _, storeAddressDependent := range physicalReg.StoreAddressDependents {
		storeAddressDependent.StoreAddressReady = true
	}

	for _, dependent := range physicalReg.Dependents {
		dependent.RemoveNotReadyOperand(uint32(physicalReg.Dependency))
	}

	physicalReg.EffectiveAddressComputationOperandDependents = []*ReorderBufferEntry{}
	physicalReg.StoreAddressDependents = []*LoadStoreQueueEntry{}
	physicalReg.Dependents = []GeneralReorderBufferEntry{}
}

func (physicalReg *PhysicalRegister) Commit() {
	if physicalReg.State != PhysicalRegisterState_RENAME_BUFFER_VALID {
		panic("Impossible")
	}

	physicalReg.Allocator = nil

	physicalReg.State = PhysicalRegisterState_ARCHITECTURAL_REGISTER
}

func (physicalReg *PhysicalRegister) Recover() {
	if physicalReg.State != PhysicalRegisterState_RENAME_BUFFER_INVALID &&
		physicalReg.State != PhysicalRegisterState_RENAME_BUFFER_VALID {
		panic("Impossible")
	}

	physicalReg.Allocator = nil

	physicalReg.Dependency = -1

	physicalReg.State = PhysicalRegisterState_AVAILABLE

	physicalReg.PhysicalRegisterFile.NumFreePhysicalRegisters++
}

func (physicalReg *PhysicalRegister) Reclaim() {
	if physicalReg.State != PhysicalRegisterState_ARCHITECTURAL_REGISTER {
		panic("Impossible")
	}

	physicalReg.Allocator = nil

	physicalReg.Dependency = -1

	physicalReg.State = PhysicalRegisterState_AVAILABLE

	physicalReg.PhysicalRegisterFile.NumFreePhysicalRegisters++
}

func (physicalReg *PhysicalRegister) Ready() bool {
	return physicalReg.State == PhysicalRegisterState_RENAME_BUFFER_VALID ||
		physicalReg.State == PhysicalRegisterState_ARCHITECTURAL_REGISTER
}

type PhysicalRegisterFile struct {
	RegisterDependencyType   RegisterDependencyType
	PhysicalRegisters        []*PhysicalRegister
	NumFreePhysicalRegisters uint32
}

func NewPhysicalRegisterFile(registerDependencyType RegisterDependencyType, size uint32) *PhysicalRegisterFile {
	var physicalRegs = &PhysicalRegisterFile{
		RegisterDependencyType:   registerDependencyType,
		NumFreePhysicalRegisters: size,
	}

	for i := uint32(0); i < size; i++ {
		physicalRegs.PhysicalRegisters = append(
			physicalRegs.PhysicalRegisters,
			NewPhysicalRegister(physicalRegs, i),
		)
	}

	return physicalRegs
}

func (physicalRegs *PhysicalRegisterFile) Allocate(allocator *ReorderBufferEntry, dependency uint32) *PhysicalRegister {
	for _, physicalReg := range physicalRegs.PhysicalRegisters {
		if physicalReg.State == PhysicalRegisterState_AVAILABLE {
			physicalReg.Allocate(allocator, dependency)
			return physicalReg
		}
	}

	panic("Impossible")
}

func (physicalRegs *PhysicalRegisterFile) Full() bool {
	return physicalRegs.NumFreePhysicalRegisters == 0
}

func (physicalRegs *PhysicalRegisterFile) Dump() {
	for i, physicalReg := range physicalRegs.PhysicalRegisters {
		fmt.Printf("physicalRegister[%d]={type=%s, num=%d, dependency=%d, state=%s, ready=%t}\n",
			i, physicalReg.PhysicalRegisterFile.RegisterDependencyType, physicalReg.Num, physicalReg.Dependency, physicalReg.State, physicalReg.Ready())

		if physicalReg.Allocator != nil {
			var reorderBufferEntry = physicalReg.Allocator

			var loadStoreQueueEntryId = int32(-1)

			if reorderBufferEntry.LoadStoreQueueEntry != nil {
				loadStoreQueueEntryId = reorderBufferEntry.LoadStoreQueueEntry.Id()
			}

			fmt.Printf(
				"physicalRegister[%d].allocator=ReorderBufferEntry{id=%d, dispatched=%t, issued=%t, completed=%t, squashed=%t, loadStoreQueueEntry.id=%d, notReadyOperands=%+v, allOperandReady=%t}\n",
				i,
				reorderBufferEntry.Id(),
				reorderBufferEntry.Dispatched(),
				reorderBufferEntry.Issued(),
				reorderBufferEntry.Completed(),
				reorderBufferEntry.Squashed(),
				loadStoreQueueEntryId,
				reorderBufferEntry.NotReadyOperands(),
				reorderBufferEntry.AllOperandReady(),
			)
		}

		for j, dependent := range physicalReg.Dependents {
			fmt.Printf("physicalRegister[%d].dependent[%d]={id=%d}\n", i, j, dependent.Id())
		}
		for j, storeAddressDependent := range physicalReg.StoreAddressDependents {
			fmt.Printf("physicalRegister[%d].storeAddressDependent[%d]={id=%d}\n", i, j, storeAddressDependent.Id())
		}
		for j, EffectiveAddressComputationOperandDependent := range physicalReg.EffectiveAddressComputationOperandDependents {
			fmt.Printf("physicalRegister[%d].EffectiveAddressComputationOperandDependent[%d]={id=%d}\n", i, j, EffectiveAddressComputationOperandDependent.Id())
		}
	}
}
