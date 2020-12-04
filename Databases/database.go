package main

import (
	"database/sql"
	"strconv"
)

type database struct {
	db *sql.DB
}

func (d *database) getAll () []*Korzina {
	row, _ := d.db.Query("SELECT * FROM bd.public.tovary")
	tovary := make([]*Korzina, 0)
	for row.Next() {
		tovar := new(Korzina)
		err := row.Scan(&tovar.id, &tovar.name,&tovar.otdel,&tovar.marka,&tovar.srok, &tovar.price)
		if err != nil {
			panic(err)
		}
		tovary = append(tovary, tovar)
	}
	return tovary
}

func (d *database) remove (sold []Korzina) {
	for _, i := range sold {
		_, err := d.db.Exec("DELETE FROM bd.public.tovary WHERE id = $1", i.id)
		if err != nil {
			panic(err)
		}
	}
}
// функция создания нового адреса в БД
func (d *database) newAddress (u []string) {
	_, err := d.db.Exec("INSERT INTO bd.public.tovary (NAME, OTDEL, MARKA, SROK, PRICE) VALUES($1, $2, $3, $4, $5)",
		u[0], u[2], u[3], u[4], u[1])
	if err != nil {
		panic(err)
	}
}

func (d *database) newAddressUser (u []string) {
	_, err := d.db.Exec("INSERT INTO bd.public.client (number, address, fio) VALUES(nextval($1, $2, $3))",
		u[1], u[2], u[0])
	if err != nil {
		panic(err)
	}
}


// функция извлечения из БД нужного адреса
func (d *database) getNeeded (id string) *[6]string {
	idd, _ := strconv.Atoi(id)
	row := d.db.QueryRow("SELECT * FROM bd.public.tovary WHERE id = ($1)", idd)
	tovar := &[6]string{}
	err := row.Scan(&tovar[0], &tovar[1],&tovar[2],&tovar[3],&tovar[4], &tovar[5])
	if err != nil{
		tovar[0] = "ID не существует. Попробуйте снова"
	}
	return tovar
}


func (d *database) getNeededUser (id string) *[6]string {
	idd, _ := strconv.Atoi(id)
	row := d.db.QueryRow("SELECT * FROM bd.public.client WHERE id = ($1)", idd)
	client := &[6]string{}
	err := row.Scan(&client[0], &client[1],&client[2],&client[3])
	if err != nil{
		client[0] = "ID не существует. Попробуйте снова"
	}
	return client
}
/*
// функция получения максимального (последнего ID) из БД
func (d *database) getMax () int {
	row := d.db.QueryRow("SELECT max(id) FROM url")
	var maxId int
	err := row.Scan(&maxId)
	if err != nil {
		panic(err)
	}
	return maxId
}*/

