package cpu

type OoODriver interface {
}

type OoO struct {
	Driver                      OoODriver
	CurrentReorderBufferEntryId int32
	CurrentDecodeBufferEntryId  int32
}

func NewOoO(driver OoODriver) *OoO {
	var ooo = &OoO{
		Driver: driver,
	}

	return ooo
}

func (ooo *OoO) ResetStats() {
}
