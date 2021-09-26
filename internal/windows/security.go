package windows

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	advAPI32DLL              = windows.MustLoadDLL("Advapi32.dll")
	procGetNamedSecurityInfo = advAPI32DLL.MustFindProc("GetNamedSecurityInfoW")
	procSetNamedSecurityInfo = advAPI32DLL.MustFindProc("SetNamedSecurityInfo")
)

func GetNamedSecurityInfo(objectName string, objectType int32, secInfo uint32, owner, group **windows.SID, dacl, sacl, secDesc *windows.Handle) error {
	ret, _, err := procGetNamedSecurityInfo.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(objectName))),
		uintptr(objectType),
		uintptr(secInfo),
		uintptr(unsafe.Pointer(owner)),
		uintptr(unsafe.Pointer(group)),
		uintptr(unsafe.Pointer(dacl)),
		uintptr(unsafe.Pointer(sacl)),
		uintptr(unsafe.Pointer(secDesc)),
	)
	if ret != 0 {
		return err
	}
	return nil
}

func SetNamedSecurityInfo(objectName string, objectType int32, secInfo uint32, owner, group *windows.SID, dacl, sacl windows.Handle) error {
	ret, _, err := procSetNamedSecurityInfo.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(objectName))),
		uintptr(objectType),
		uintptr(secInfo),
		uintptr(unsafe.Pointer(owner)),
		uintptr(unsafe.Pointer(group)),
		uintptr(dacl),
		uintptr(sacl),
	)
	if ret != 0 {
		return err
	}
	return nil
}
