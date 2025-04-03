CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  role VARCHAR(10) NOT NULL DEFAULT "user", 
  password VARCHAR(255) NOT NULL
);
DROP TABLE IF EXISTS categories;
CREATE TABLE IF NOT EXISTS categories (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(255) UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS posts ( 
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  image TEXT ,
  user_id INTEGER NOT NULL,
  creat_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS comments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  content TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  date_creation TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  post_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (post_id) REFERENCES posts(id)
);
CREATE TABLE IF NOT EXISTS categorie_report (
  id INTEGER PRIMARY KEY AUTOINCREMENT, 
  name VARCHAR(255) UNIQUE NOT NULL
);


CREATE TABLE IF NOT EXISTS report(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
   post_id INTEGER NOT NULL,
   id_categorie_report INTEGER NOT NULL ,
   FOREIGN KEY(user_id) REFERENCES users(id),
   FOREIGN KEY(post_id) REFERENCES posts(id),
   FOREIGN KEY(id_categorie_report) REFERENCES categorie_report(id)
);

CREATE TABLE IF NOT EXISTS reactPost(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  react_type  VARCHAR(255) NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE TABLE IF NOT EXISTS reactComment(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  comment_id INTEGER NOT NULL,
  react_type VARCHAR(255) NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (comment_id) REFERENCES comments(id)
);


CREATE TABLE IF NOT EXISTS sessionss (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  token VARCHAR(255) NOT NULL,
  expired_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id), 
  UNIQUE(user_id)
);

CREATE TABLE IF NOT EXISTS post_categorie (
  post_id INTEGER NOT NULL,
  categorie_id INTEGER NOT NULL,
  FOREIGN KEY (post_id) REFERENCES posts(id),
  FOREIGN KEY (categorie_id) REFERENCES categories(id),
  PRIMARY KEY (post_id, categorie_id)
);

CREATE TABLE IF NOT EXISTS messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  sender_id INTEGER NOT NULL,
  receiver_id INTEGER NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (sender_id) REFERENCES users(id),
  FOREIGN KEY (receiver_id) REFERENCES users(id)
);

INSERT INTO categories (name) 
VALUES ('Sport'), ('Music'), ('Movies'), ('Science'), ('Politics'), ('Culture'), ('Technology')
ON CONFLICT (name) DO NOTHING;

INSERT INTO categorie_report (name) 
VALUES ('Irrelevant'), ('Obscene'), ('Illegal'), ('Insulting')
ON CONFLICT (name) DO NOTHING;

INSERT INTO users (username, email, role, password) VALUES ('melfarss', 'medlfarssi10@gmail.com', 'admin', 'Myfarssi123') ON CONFLICT DO NOTHING;