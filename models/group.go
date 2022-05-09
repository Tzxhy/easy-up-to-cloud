package models

import "log"

type ResourceGroupItem struct {
	Gid        string
	Name       string
	CreateDate string
}

func GetResourceGroup(uid string) *[]ResourceGroupItem {
	rows, err := DB.Query("select gid, name, create_date from user_group where user_ids like ?", "%"+uid+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []ResourceGroupItem
	for rows.Next() {
		item := new(ResourceGroupItem)
		rows.Scan(
			&item.Gid,
			&item.Name,
			&item.CreateDate,
		)
		items = append(items, *item)
	}
	return &items
}
