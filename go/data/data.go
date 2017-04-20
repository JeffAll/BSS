package data

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	clearQuery  = "DELETE FROM data"
	appendQuery = "INSERT INTO data('order','value') values(?,?)"
	queryQuery  = "SELECT `value` FROM data ORDER BY `order`"
)

type Items struct {
	Items []string `json:"items"`
}

type Data struct {
	db    *sql.DB
	query *sql.Stmt
}

func BuildData(
	driver string,
	address string,
) (
	*Data,
	error,
) {
	toReturn := Data{}

	var err error
	toReturn.db, err = sql.Open(
		driver,
		address,
	)
	if err != nil {
		log.Printf(
			"Error Opening Connection to Database\n\t:%s",
			err,
		)
		return nil, err
	}

	toReturn.query, err = toReturn.db.Prepare(
		queryQuery,
	)
	if err != nil {
		log.Printf(
			"Error Preparing Query Query\n\t:%s",
			err,
		)
		return nil, err
	}
	return &toReturn, nil
}

func (d *Data) AppendArrayAndClear(
	items Items,
) error {
	log.Printf(
		"AppendArrayAndClear",
	)
	tx, err := d.db.Begin()
	if err != nil {
		log.Printf(
			"Error Insantiating Transaction\n\t:%s",
			err,
		)
		return err
	}
	err = clear(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = appendItems(
		tx,
		items,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Printf(
			"Error Committing the transaction\n\t:%s",
			err,
		)
		tx.Rollback()
	}
	return err
}

// Clear executes the clear sql statement to delete all the rows
// from the 'data' table
func clear(
	tx *sql.Tx,
) error {
	_, err := tx.Exec(clearQuery)
	if err != nil {
		log.Printf(
			"Error Executing Clear Statement\n\t:%s",
			err,
		)
	}
	return err
}

// Append appends a single row to the 'data' table
func appendItem(
	stmt *sql.Stmt,
	order int,
	item string,
) error {
	_, err := stmt.Exec(
		order,
		item,
	)
	if err != nil {
		log.Printf(
			"Error Executing Append Statement\n\t:%s",
			err,
		)
	}
	return err
}

// AppendItems executes an append for each value in
// toAppend
func appendItems(
	tx *sql.Tx,
	items Items,
) error {
	stmt, err := tx.Prepare(
		appendQuery,
	)
	if err != nil {
		log.Printf(
			"Error Preparing Append Statement\n\t:%s",
			err,
		)
		return err
	}
	for i := 0; i < len(items.Items); i++ {
		err := appendItem(
			stmt,
			i,
			items.Items[i],
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Query returns the 'value' columns in data
func (d *Data) Query() (
	Items,
	error,
) {
	rows, err := d.query.Query()
	if err != nil {
		log.Printf(
			"Error Executing Query Statement\n\t:%s",
			err,
		)
		return Items{}, err
	}
	var toReturn Items
	for rows.Next() {
		var toAdd string
		err = rows.Scan(
			&toAdd,
		)
		if err != nil {
			log.Printf(
				"Error Scanning Returned Rows\n\t:%s",
				err,
			)
		}
		toReturn.Items = append(toReturn.Items, toAdd)
	}
	return toReturn, nil
}
