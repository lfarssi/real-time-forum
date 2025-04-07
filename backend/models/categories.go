package models


import (
	"real_time_forum/backend/database"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetCategories() ([]Category, error) {
	query := "SELECT id, name FROM categories"
    rows, err := database.Database.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var categories []Category
    for rows.Next() {
        var category Category
        err = rows.Scan(&category.ID, &category.Name)
        if err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }
    return categories, nil
}

func InsertIntoCategoryPost(postId, categorieId int) error {
	query := "INSERT INTO post_categorie (post_id, categorie_id) VALUES (?,?)"
    _, err := database.Database.Exec(query, postId, categorieId)
    if err!= nil {
        return err
    }
    return nil
}

