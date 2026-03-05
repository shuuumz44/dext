DROP TABLE IF EXISTS expenses;
CREATE TABLE expenses (
	id				INT AUTO_INCREMENT 	NOT NULL,
	purchased		DATE,
	description		VARCHAR(256),
	amount			INT					NOT NULL	DEFAULT 0,
	PRIMARY KEY		(`id`)
);

