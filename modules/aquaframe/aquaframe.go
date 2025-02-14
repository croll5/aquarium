/*
Copyright ou © ou Copr. Charles Mailley, (21 janvier 2025)

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à l'analyse des collectes
traçologiques effectuées avec le logiciel DFIR-ORC.

Ce logiciel est régi par la licence CeCILL soumise au droit français et
respectant les principes de diffusion des logiciels libres. Vous pouvez
utiliser, modifier et/ou redistribuer ce programme sous les conditions
de la licence CeCILL telle que diffusée par le CEA, le CNRS et l'INRIA
sur le site "http://www.cecill.info".

En contrepartie de l'accessibilité au code source et des droits de copie,
de modification et de redistribution accordés par cette licence, il n'est
offert aux utilisateurs qu'une garantie limitée.  Pour les mêmes raisons,
seule une responsabilité restreinte pèse sur l'auteur du programme,  le
titulaire des droits patrimoniaux et les concédants successifs.

A cet égard  l'attention de l'utilisateur est attirée sur les risques
associés au chargement,  à l'utilisation,  à la modification et/ou au
développement et à la reproduction du logiciel par l'utilisateur étant
donné sa spécificité de logiciel libre, qui peut le rendre complexe à
manipuler et qui le réserve donc à des développeurs et des professionnels
avertis possédant  des  connaissances  informatiques approfondies.  Les
utilisateurs sont donc invités à charger  et  tester  l'adéquation  du
logiciel à leurs besoins dans des conditions permettant d'assurer la
sécurité de leurs systèmes et ou de leurs données et, plus généralement,
à l'utiliser et l'exploiter dans les mêmes conditions de sécurité.

Le fait que vous puissiez accéder à cet en-tête signifie que vous avez
pris connaissance de la licence CeCILL, et que vous en avez accepté les
termes.
*/

package aquaframe

import (
	"database/sql"
	"fmt"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type Aquaframe struct {
	Table dataframe.DataFrame
	Error error
}

// Implémentation de l'interface fmt.Stringer
func (adf Aquaframe) String() string {
	return adf.Table.String()
}

/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ----------------------------------------   INITIALISATION   ---------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */

func Df(dataframe dataframe.DataFrame) *Aquaframe {
	adf := Aquaframe{}
	adf.Table = dataframe
	return &adf

}

func RowsToAquaframe(rows *sql.Rows) *Aquaframe {
	// Take column names
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("adb.WARNING - SelectFrom failed getting columns" + err.Error())
		return nil
	}
	// Store data
	var records [][]string
	records = append(records, columns)

	// Read data
	for rows.Next() {
		columnPointers := make([]interface{}, len(columns))
		columnValues := make([]sql.NullString, len(columns))
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Println("adb.WARNING - SelectFrom failed scanning: " + err.Error())
			return nil
		}
		row := make([]string, len(columns))
		for i, col := range columnValues {
			if col.Valid {
				row[i] = col.String
			} else {
				row[i] = ""
			}
		}
		records = append(records, row)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("adb.WARNING - SelectFrom failed during rows iteration: " + err.Error())
		return nil
	}
	// Convert and return the dataframe
	adf := Aquaframe{}
	adf.Table = dataframe.LoadRecords(records)
	return &adf
}

/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ----------------------------------------      METHODS       ---------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */

func (adf Aquaframe) Head(nFirstRows int) *Aquaframe {
	// Return it
	df := Aquaframe{}
	df.Table = adf.Table.Subset([]int{0, nFirstRows})
	return &df
}

func (adf Aquaframe) AddColumn(colname string, colvalues interface{}) {
	var newCol series.Series
	// Create a Series with the good type
	switch v := colvalues.(type) {
	case []int:
		newCol = series.New(v, series.Int, colname)
	case []float64:
		newCol = series.New(v, series.Float, colname)
	case []string:
		newCol = series.New(v, series.String, colname)
	case []bool:
		newCol = series.New(v, series.Bool, colname)
	// Ajoutez plus de cas pour d'autres types si nécessaire
	default:
		fmt.Printf("Unsupported type: %T\n", v)
		return
	}
	// Add the new column to the dataframe
	adf.Table = adf.Table.CBind(dataframe.New(newCol))
}

func (adf Aquaframe) Strloc(r int, c int) string {
	return adf.Table.Elem(r, c).String()
}

func (adf Aquaframe) Intloc(r int, c int) (int, error) {
	return adf.Table.Elem(r, c).Int()
}

func (adf Aquaframe) ToMap() []map[string]interface{} {
	var result []map[string]interface{}

	for i := 0; i < adf.Table.Nrow(); i++ {
		row := make(map[string]interface{})
		for _, colName := range adf.Table.Names() {
			val := adf.Table.Col(colName).Elem(i)
			switch val.Type() {
			case series.Int:
				row[colName], _ = val.Int()
			case series.Float:
				row[colName] = val.Float()
			case series.String:
				row[colName] = val.String()
			case series.Bool:
				row[colName], _ = val.Bool()
			default:
				row[colName] = val
			}
		}
		result = append(result, row)
	}
	return result
}
