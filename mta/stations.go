package mta

type Stations map[StopID]*Station

type Station struct {
	ID          StopID
	Name        string
	Coordinates *Coordinates
	Stops       map[StopID]struct{}
}

func (s *Station) StopIDs() []StopID {
	vv := make([]StopID, 0, len(s.Stops))
	for k := range s.Stops {
		vv = append(vv, k)
	}
	return vv
}
