package models

import "real_time_forum/backend/database"

func GetCategories() ([]Category, error) {
	query := "SELECT id, name FROM category"
    rows, err := database.DB.Query(query)
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