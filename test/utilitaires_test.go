package aquarium

import (
	"aquarium/modules/extraction/utilitaires"
	"encoding/binary"
	"fmt"
	"testing"
	"time"
)

func TestFiletimeVersGoHivers(t *testing.T) {
	var depart int64 = 133847143860000000
	var departEncode []byte
	departEncode = binary.LittleEndian.AppendUint64(departEncode, uint64(depart))
	var attendu time.Time = time.Date(2025, 2, 22, 17, 13, 6, 0, time.FixedZone("CET", 3600))
	var recu time.Time = utilitaires.FileTimeVersGo(departEncode)
	if attendu.Compare(recu) != 0 {
		t.Fatal(fmt.Printf("FileTimeVersGo(%d) a renvoyé %s alors que %s était attendu.", depart, recu, attendu))
	}
}

func TestFiletimeVersGoEte(t *testing.T) {
	var depart int64 = 132433170450000000
	var departEncode []byte
	departEncode = binary.LittleEndian.AppendUint64(departEncode, uint64(depart))
	var attendu time.Time = time.Date(2020, 8, 31, 5, 10, 45, 0, time.FixedZone("CET", 7200))
	var recu time.Time = utilitaires.FileTimeVersGo(departEncode)
	if attendu.Compare(recu) != 0 {
		t.Fatal(fmt.Printf("FileTimeVersGo(%d) a renvoyé %s alors que %s était attendu.", depart, recu, attendu))
	}
}
