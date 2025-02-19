package bconf

type fieldSetGroup struct {
	name      string
	fieldSets FieldSets
}

type fieldSetGroups []*fieldSetGroup
