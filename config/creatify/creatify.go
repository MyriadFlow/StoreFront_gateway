package creatify

import (
	"errors"

	"github.com/TheLazarusNetwork/marketplace-engine/config/smartcontract"
	"github.com/TheLazarusNetwork/marketplace-engine/util/pkg/logwrapper"
)

type tRole int

var (
	ErrRoleNotExist = errors.New("role does not exist")
)

const (
	CREATOR_ROLE  tRole = iota
	ADMIN_ROLE    tRole = iota
	OPERATOR_ROLE tRole = iota
)

type tRoles map[tRole][32]byte

var roles tRoles = tRoles{}
var initiated = false

func GetRole(role tRole) ([32]byte, error) {
	if !initiated {
		InitRolesId()
	}
	v, ok := roles[role]
	if !ok {
		return [32]byte{}, ErrRoleNotExist
	}
	return v, nil
}
func InitRolesId() {
	client, err := smartcontract.GetClient()
	if err != nil {
		logwrapper.Fatalf("failed to client, error: %v", err.Error())
	}
	instance, err := GetInstance(client)
	if err != nil {
		logwrapper.Fatalf("failed to get instance for %v , error: %v", "CREATIFY", err.Error())
	}
	creatorRoleId, err := instance.CREATIFYCREATORROLE(nil)
	if err != nil {
		logwrapper.Fatalf("Failed to get %v, error: %v", "CREATIFYCREATORROLE", err.Error())
	}
	roles[CREATOR_ROLE] = creatorRoleId
	adminRoleId, err := instance.CREATIFYADMINROLE(nil)
	if err != nil {
		logwrapper.Fatalf("Failed to get %v, error: %v", "CREATIFYADMINROLE", err.Error())
	}
	roles[ADMIN_ROLE] = adminRoleId

	operatorRoleId, err := instance.CREATIFYOPERATORROLE(nil)
	if err != nil {
		logwrapper.Fatalf("Failed to get %v, error: %v", "CREATIFYOPERATORROLE", err.Error())
	}
	roles[OPERATOR_ROLE] = operatorRoleId
}
