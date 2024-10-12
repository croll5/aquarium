package navigateur

import(
	"io/ioutil"
)

type Navigateur struct{}

func (n Navigateur) Extraction(cheminProjet string) error {
	files, err := ioutil.ReadDir(cheminProjet)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	for index, element := range someFiles{
		fmt.Println(someFiles.name)
	}
	RecuperationOrcs(files, cheminProjet)
	return nil
}

func (n Navigateur) Description() string {
	return "Historique de navigation"
}

func (n Navigateur) PrerequisOK(cheminORC string) bool {
	return true
}
