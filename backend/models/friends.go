package models

import "real_time_forum/backend/database"

func Friends(userID int)([]*UserAuth, error)  {
    query := `SELECT u.id, u.firstName, u.lastName, u.gender, m.sentAt
        FROM users u
        INNER JOIN messages m ON m.receiverID=u.id
        WHERE u.id != ?
        ORDER BY u.firstName AND m.sentAt
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