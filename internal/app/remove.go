package app

import "log"

func (a *app) remove(id int) {
	if id < 1 || id > len(a.records) {
		log.Fatal("ID is out of range")
	}

	if id == len(a.records) {
		a.records = a.records[:len(a.records)-1]
	} else {
		a.records = append(a.records[:id-1], a.records[id:]...)
	}

	a.saveRecords(a.records)
	a.renderer.Render(a.records)
}
