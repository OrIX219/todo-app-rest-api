package todo

type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type ErrNoSuchList struct{}

func (e *ErrNoSuchList) Error() string {
	return "No list with such id"
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type ErrNoSuchItem struct{}

func (e *ErrNoSuchItem) Error() string {
	return "No item with such id"
}

type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}

type UpdateListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type ErrInvalidUpdateListInput struct{}

func (e *ErrInvalidUpdateListInput) Error() string {
	return "Invalid list update input"
}

func (i UpdateListInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return &ErrInvalidUpdateListInput{}
	}
	return nil
}

type UpdateItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

type ErrInvalidUpdateItemInput struct{}

func (e *ErrInvalidUpdateItemInput) Error() string {
	return "Invalid item update input"
}

func (i UpdateItemInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return &ErrInvalidUpdateItemInput{}
	}
	return nil
}
