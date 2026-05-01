DROP TABLE IF EXISTS expenses;
CREATE TABLE expenses (
	id				INT AUTO_INCREMENT	NOT NULL,
	category		INT,
	name			VARCHAR(255),
	amount			DECIMAL(20,2)		NOT NULL	DEFAULT 0,
	purchased		DATE 				NOT NULL	DEFAULT (CURRENT_DATE()),
	created			DATETIME 			NOT NULL	DEFAULT (CURRENT_DATE()),
	PRIMARY KEY		(id)
	FOREIGN KEY		(category)						REFERENCES categories(id)
);

DROP TABLE IF EXISTS budget;
CREATE TABLE budget (
	id 				INT					DEFAULT 0,
	threshold		DECIMAL(20,2)
);
INSERT INTO budget (id, threshold) VALUES(0, 0.0);

DROP TABLE IF EXISTS categories;
CREATE TABLE categories (
	id				INT AUTO_INCREMENT	NOT NULL,
	name			VARCHAR(255)

	PRIMARY KEY (id)
);
