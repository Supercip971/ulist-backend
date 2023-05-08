package db

import (
	"github.com/pocketbase/pocketbase/models"
)

type ListVisibility struct {
	models.BaseModel

	UserOwnership string `db:"userOwnership" json:"userOwnership"`
	Owner         bool   `db:"owner" json:"owner"`
	Readonly      bool   `db:"readonly" json:"readonly"`
	List          string `db:"list" json:"list"`
}

var _ models.Model = (*ListVisibility)(nil)

func (m *ListVisibility) TableName() string {
	return "lists_visibility" // the name of your collection
}

type ShopList struct {
	models.BaseModel

	Name              string `db:"name" json:"name"`
	ReadonlyByDefault bool   `db:"readonly_by_default" json:"readonly_by_default"`
}

var _ models.Model = (*ShopList)(nil)

func (m *ShopList) TableName() string {
	return "lists" // the name of your collection
}

type DbShopListEntries struct {
	models.BaseModel

	Name string `db:"name" json:"name"`

	List           string `db:"list" json:"list"`
	AddedBy        string `db:"addedBy" json:"addedBy"`
	Checked        bool   `db:"checked" json:"checked"`
	CustomRelation string `db:"customRelation" json:"customRelation"`
}

var _ models.Model = (*DbShopListEntries)(nil)

func (m *DbShopListEntries) TableName() string {
	return "list_entries" // the name of your collection
}
