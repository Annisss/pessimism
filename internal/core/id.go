package core

import (
	"fmt"

	"github.com/google/uuid"
)

// UUID ... third-party wrapper struct
type UUID struct {
	uuid.UUID
}

// newUUID ... Constructor
func newUUID() UUID {
	return UUID{
		uuid.New(),
	}
}

// nilUUID ... Returns a zero'd out 16 byte array
func nilUUID() UUID {
	return UUID{[16]byte{0}}
}

// ShortString ... Short string representation for easier
// debugging and ensuring conformance with pessimism specific abstractions
// https://pkg.go.dev/github.com/google/UUID#UUID.String
func (id UUID) ShortString() string {
	uid := id.UUID
	// Only render first 8 bytes instead of entire sequence
	return fmt.Sprintf("%d%d%d%d%d%d%d%d%d",
		uid[0],
		uid[1],
		uid[2],
		uid[2],
		uid[3],
		uid[4],
		uid[5],
		uid[6],
		uid[7])
}

// ComponentPID ... Component Primary ID
type ComponentPID [4]byte

// Represents a non-deterministic ID that's assigned to
// every uniquely constructed ETL component
type ComponentUUID struct {
	PID  ComponentPID
	UUID UUID
}

// Used for local lookups to look for active collisions
type PipelinePID [9]byte

// Represents a non-deterministic ID that's assigned to
// every uniquely constructed ETL pipeline
type PipelineUUID struct {
	PID  PipelinePID
	UUID UUID
}

// PipelineType ... Returns pipeline type decoding from encoded pid byte
func (uuid PipelineUUID) PipelineType() PipelineType {
	return PipelineType(uuid.PID[0])
}

// InvSessionPID ... Invariant session Primary ID
type InvSessionPID [3]byte

// Represents a non-deterministic ID that's assigned to
// every uniquely constructed invariant session
type InvSessionUUID struct {
	PID  InvSessionPID
	UUID UUID
}

// Network ... Returns network decoding from encoded pid byte
func (pid InvSessionPID) Network() Network {
	return Network(pid[0])
}

// InvType ... Returns invariant type decoding from encoded pid byte
func (pid InvSessionPID) InvType() InvariantType {
	return InvariantType(pid[2])
}

// NOTE - This is useful for error handling with functions that
// also return a ComponentID
// NilCompID ... Returns a zero'd out or empty component UUID
func NilComponentUUID() ComponentUUID {
	return ComponentUUID{
		PID:  ComponentPID{0},
		UUID: nilUUID(),
	}
}

// NilCompID ... Returns a zero'd out or empty pipeline UUID
func NilPipelineUUID() PipelineUUID {
	return PipelineUUID{
		PID:  PipelinePID{0},
		UUID: nilUUID(),
	}
}

// NilInvariantUUID ... Returns a zero'd out or empty invariant UUID
func NilInvariantUUID() InvSessionUUID {
	return InvSessionUUID{
		PID:  InvSessionPID{0},
		UUID: nilUUID(),
	}
}

// MakeComponentUUID ... Constructs a component PID sequence & random UUID
func MakeComponentUUID(pt PipelineType, ct ComponentType, rt RegisterType, n Network) ComponentUUID {
	cID := ComponentPID{
		byte(n),
		byte(pt),
		byte(ct),
		byte(rt),
	}

	return ComponentUUID{
		PID:  cID,
		UUID: newUUID(),
	}
}

// MakePipelineUUID ... Constructs a pipeline PID sequence & random UUID
func MakePipelineUUID(pt PipelineType, firstCID, lastCID ComponentUUID) PipelineUUID {
	cID1, cID2 := firstCID.PID, lastCID.PID

	pID := PipelinePID{
		byte(pt),
		cID1[0],
		cID1[1],
		cID1[2],
		cID1[3],
		cID2[0],
		cID2[1],
		cID2[2],
		cID2[3],
	}

	return PipelineUUID{
		PID:  pID,
		UUID: newUUID(),
	}
}

// MakeInvSessionUUID ... Constructs an invariant PID sequence & random UUID
func MakeInvSessionUUID(n Network, pt PipelineType, invType InvariantType) InvSessionUUID {
	pID := InvSessionPID{
		byte(n),
		byte(pt),
		byte(invType),
	}

	return InvSessionUUID{
		PID:  pID,
		UUID: newUUID(),
	}
}

// String ... Returns string representation of a component PID
func (pid ComponentPID) String() string {
	return fmt.Sprintf("%s:%s:%s:%s",
		Network(pid[0]).String(),
		PipelineType(pid[1]).String(),
		ComponentType(pid[2]).String(),
		RegisterType(pid[3]).String(),
	)
}

// String ... Returns string representation of a component UUID
func (uuid ComponentUUID) String() string {
	return fmt.Sprintf("%s::%s",
		uuid.PID.String(),
		uuid.UUID.ShortString(),
	)
}

// Type ... Returns component type byte value from component UUID
func (uuid ComponentUUID) Type() ComponentType {
	return ComponentType(uuid.PID[2])
}

// String ... Returns string representation of a pipeline PID
func (pid PipelinePID) String() string {
	pt := PipelineType(pid[0]).String()
	cID1 := ComponentPID(*(*[4]byte)(pid[1:5])).String()
	cID2 := ComponentPID(*(*[4]byte)(pid[5:9])).String()

	return fmt.Sprintf("%s::%s::%s", pt, cID1, cID2)
}

// String ... Returns string representation of a pipeline UUID
func (uuid PipelineUUID) String() string {
	return fmt.Sprintf("%s:::%s",
		uuid.PID.String(), uuid.UUID.ShortString(),
	)
}

// String ... Returns string representation of an invariant sesion PID
func (pid InvSessionPID) String() string {
	return fmt.Sprintf("%s:%s:%s",
		Network(pid[0]).String(),
		PipelineType(pid[1]).String(),
		InvariantType(pid[2]).String(),
	)
}

// String ... Returns string reprsentation of an invariant session UUID
func (uuid InvSessionUUID) String() string {
	return fmt.Sprintf("%s::%s",
		uuid.PID.String(), uuid.UUID.ShortString())
}
