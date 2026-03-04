DROP TABLE IF EXISTS expenses;
CREATE TABLE expenses (
	id				INT AUTO_INCREMENT 	NOT NULL,
	purchased		DATE				NOT NULL,
	description		VARCHAR(256)		NOT NULL,
	amount			INT					NOT NULL,
	PRIMARY KEY		(`id`)
);

