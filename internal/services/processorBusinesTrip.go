package services

type porocessorBusinesTrip struct {
}

func (proc *porocessorBusinesTrip) Check(data any) error {
	/*
		reportData, ok := data.(*reports.ReportData)
		if !ok {
			return errors.New("INCORRECT BUSINES TRIP DATA STRUCTURE")
		}

		details, ok := reportData.Details.(*reports.BussinesTripDetails)
		if !ok {
			return errors.New("INCORRECT BUSINES TRIP DATA STRUCTURE")
		}
		if (!details.Supervisor.Blank() && !details.ActingSupervisor.Blank()) || (details.Supervisor.Blank() && details.ActingSupervisor.Blank()) {
			return errors.New("SUPERVISOR ERROR")
		}
	*/

	return nil

}
