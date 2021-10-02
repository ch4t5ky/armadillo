package windows

import (
	"golang.org/x/sys/windows"
	"os"
)

func Chmod(path string, fileMode os.FileMode) {
	mode := uint32(fileMode)
	creatorOwnerSID, _ := windows.StringToSid("S-1-3-0")
	creatorGroupSID, _ := windows.StringToSid("S-1-3-1")
	everyoneSID, _ := windows.StringToSid("S-1-1-0")
	_ = Apply(
		path,
		true,
		false,
		GrantSid(((mode&0700)<<23)|((mode&0200)<<9), creatorOwnerSID),
		GrantSid(((mode&0070)<<26)|((mode&0020)<<12), creatorGroupSID),
		GrantSid(((mode&0007)<<29)|((mode&0002)<<15), everyoneSID),
	)

}
