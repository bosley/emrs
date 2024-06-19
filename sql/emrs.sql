
CREATE TABLE user_directory (
	id MEDIUMINT NOT NULL AUTO_INCREMENT,
	username VARCHAR(25),
    authentication VARCHAR(75),
    identity VARCHAR(255),
    PRIMARY KEY (id));
    
CREATE TABLE vouchers (
	id MEDIUMINT NOT NULL AUTO_INCREMENT,
    user_id VARCHAR(25),
    voucher VARCHAR(255),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id)
		REFERENCES user_directory(id));