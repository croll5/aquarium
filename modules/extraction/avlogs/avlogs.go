package avlogs

type AvLogs struct{}

func (a AvLogs) Parse(projectLink string) {

}

func (a AvLogs) Description() string {
	return "Parsage des journaux d'antivirus dans le fichier avlogs"
}

func (a AvLogs) PrerequisOK(projectLink string) bool { return false }
