package aoi

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMger := NewAOIManager(100,300,4,200,450,5)
	fmt.Println(aoiMger)
}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	aoiMgr := NewAOIManager(0,250, 5, 0, 250, 5)
	for k, _ := range aoiMgr.grids {
		grids := aoiMgr.GetSurroundGridsByGid(k)
		fmt.Println("gid:",k,"grids len = ",len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Printf("grid ID: %d, surrounding grid IDs are %v\n", k, gIDs)
	}
}