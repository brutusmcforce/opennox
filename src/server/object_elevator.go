package server

// ElevatorUpdateData is the update data of elevators (ElevatorUpdate).
type ElevatorUpdateData struct {
	Field_0       uint32  // 0, 0
	Partner       *Object // 1, 4: the elevator shaft
	PartnerExtent uint32  // 2, 8: extent of the shaft, until they are linked
	Field_3       byte    // 3, 12: state
	_             [3]byte // 3, 13
	Field_4       uint32  // 4, 16: height
}

func (obj *Object) UpdateDataElevator() *ElevatorUpdateData {
	return updateDataAs[ElevatorUpdateData](obj)
}

// ElevatorShaftUpdateData is the update data of elevator shafts (ElevatorShaftUpdate).
type ElevatorShaftUpdateData struct {
	Field_0 uint32  // 0, 0
	Partner *Object // 1, 4: the elevator
	Field_2 uint32  // 2, 8
	State   byte    // 3, 12: the last state of the elevator
	_       [3]byte // 3, 13
}

func (obj *Object) UpdateDataElevatorShaft() *ElevatorShaftUpdateData {
	return updateDataAs[ElevatorShaftUpdateData](obj)
}

// TransporterUpdateData is the update data of transporters (PentagramUpdate, InvisiblePentagramUpdate).
type TransporterUpdateData struct {
	State       byte    // 0, 0
	_           [3]byte // 0, 1
	Active      uint32  // 1, 4
	Speed       byte    // 2, 8
	Counter     byte    // 2, 9
	_           [2]byte // 2, 10
	Destination *Object // 3, 12
	DestExtent  uint32  // 4, 16: extent of the destination, until they are linked
	Anim        byte    // 5, 20
	_           [3]byte // 5, 21
}

func (obj *Object) UpdateDataTransporter() *TransporterUpdateData {
	return updateDataAs[TransporterUpdateData](obj)
}
