CREATE TABLE IF not EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email char(50) UNIQUE NOT NULL,
    username char(13) UNIQUE NOT NULL,
    firstName char(30) NOT NULL,
    lastName char(30) NOT NULL,
    age INTEGER,
    gender char(20) NOT NULL,
    password char(40),
    createdAt DATE NOT NULL,
    session TEXT ,
    expiredAt DATE,
    authType INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title char(50),
    content TEXT, 
    dateCreation DATE,
    image BLOB,
    userID INTEGER,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE 
);

CREATE TABLE IF NOT EXISTS category(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name varchar(20) UNIQUE
);

CREATE TABLE IF NOT EXISTS postCategory(
    postID INTEGER,
    categoryID INTEGER,
    PRIMARY KEY (postID, categoryID),
    FOREIGN KEY (postID) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE, 
    FOREIGN KEY (categoryID) REFERENCES category(id) ON DELETE CASCADE ON UPDATE CASCADE 
);

CREATE TABLE IF NOT EXISTS postLike(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    status varchar(10),
    userID INTEGER,   
    postID INTEGER,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (postID) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS comment(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT,
    dateCreation DATE,
    userID INTEGER,    
    postID INTEGER,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (postID) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS commentLike(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    status varchar(10),
    userID INTEGER,    
    commentID INTEGER,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (commentID) REFERENCES comment(id) ON DELETE CASCADE ON UPDATE CASCADE
);
DROP TRIGGER IF EXISTS check_post_exists_trigger;

CREATE TRIGGER check_post_exists_trigger
BEFORE INSERT ON comment
FOR EACH ROW
BEGIN
    SELECT CASE
        WHEN NOT EXISTS (SELECT 1 FROM posts WHERE id = NEW.postID) 
        THEN RAISE(ABORT, 'Post does not exist')
    END;
END;




DROP TRIGGER IF EXISTS check_post_like_target;

CREATE TRIGGER check_post_like_target
BEFORE INSERT ON postLike
FOR EACH ROW
BEGIN
    SELECT CASE
        WHEN NOT EXISTS (SELECT 1 FROM posts WHERE id = NEW.postID)
        THEN RAISE(ABORT, 'Post does not exist')
    END;
END;



DROP TRIGGER IF EXISTS check_comment_like_target;

CREATE TRIGGER check_comment_like_target
BEFORE INSERT ON commentLike
FOR EACH ROW
BEGIN
    SELECT CASE
        WHEN NOT EXISTS (SELECT 1 FROM comment WHERE id = NEW.commentID)
        THEN RAISE(ABORT, 'Comment does not exist')
    END;
END;


INSERT OR IGNORE INTO category (name) VALUES ('Coding');
INSERT OR IGNORE INTO category (name) VALUES ('Innovation');
INSERT OR IGNORE INTO category (name) VALUES ('Bitcoin');
INSERT OR IGNORE INTO category (name) VALUES ('kids');
INSERT OR IGNORE INTO category (name) VALUES ('movie');
INSERT OR IGNORE INTO category (name) VALUES ('sport');
INSERT OR IGNORE INTO category (name) VALUES ('food');