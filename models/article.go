/* models/article.go */
package models

import "time"

type Article struct {
	ID        int
	Title     string
	Content   string
	AuthorID  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
