package aquaframe

import (
	"database/sql"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type Aquaframe struct {
	Table dataframe.DataFrame
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

func (adf Aquaframe) Iloc(r int, c int) string {
	return adf.Table.Elem(r, c).String()
}
