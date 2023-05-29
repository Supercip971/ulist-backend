package db

type ShoppingListEntry struct {
	Id             string `db:"id" json:"id"`
	Created        string `db:"created" json:"created"`
	Name           string `db:"name" json:"name"`
	List           string `db:"list" json:"list"`
	AddedBy        string `db:"addedBy" json:"addedBy"`
	Checked        bool   `db:"checked" json:"checked"`
	CustomRelation string `db:"customRelation" json:"customRelation"`
}

type ShoppingList struct {
	Lists []ShoppingListEntry `json:"entries"`
	Count int                 `json:"count"`
}

type UserListEntry struct {
	Id      string `db:"id" json:"id"`
	Created string `db:"created" json:"created"`
	ListId  string `db:"listId" json:"listId"`
	UserId  string `db:"userId" json:"userId"`
	Owner   bool   `db:"owner" json:"owner"`
}

type CreateShoppingList struct {
	Name string `json:"name"`
}

type UserLists struct {
	Lists []UserListEntry `json:"lists"`
	Count int             `json:"count"`
}

type PutShoppingList struct {
	Name    string `db:"name" json:"name"`
	Checked bool   `db:"checked" json:"checked"`
}

type UpdateShoppingList struct {
	EntryId string `db:"id" json:"id"`     // entry id, using 0 for new entry
	Name    string `db:"name" json:"name"` //
	Checked bool   `db:"checked" json:"checked"`
}

type PostShoppingList struct {
	ListId string `db:"id" json:"id"`

	Entries []PutShoppingList    `db:"pushes" json:"pushes"`
	Update  []UpdateShoppingList `db:"updates" json:"updates"`
}
