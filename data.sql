

INSERT INTO TGBot_db (user_name, url) VALUES ('Kosnstantin_Savin', 'https://beatfilmfestival.ru/');

DROP TABLE tgbot_db

CREATE TABLE  public.tgbot_db (id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    create_time DATE DEFAULT CURRENT_DATE,
    user_name VARCHAR(255),
    url VARCHAR(255),
	description VARCHAR(255),
	category VARCHAR(255),
	price INT,
	time_duration INT)