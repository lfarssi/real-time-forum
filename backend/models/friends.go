package models

import "real_time_forum/backend/database"

func Friends(userID int)([]*UserAuth, error)  {
    query := `SELECT id, firstName, lastName, gender
        FROM users
        WHERE id != ?
        ORDER BY firstName
    `
    rows, err := database.DB.Query(query, userID) 
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*UserAuth
    for rows.Next() {
        var user UserAuth
        err := rows.Scan(&user.ID,&user.FirstName, &user.LastName, &user.Gender)
        if err != nil {
            return nil, err
        }
        users = append(users, &user)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return users, nil    
}