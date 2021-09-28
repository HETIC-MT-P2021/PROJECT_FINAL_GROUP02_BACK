CREATE TABLE IF NOT EXISTS characters (
    id SERIAL,
    charName VARCHAR(255) NOT NULL,
    player BIGINT NOT NULL,
    weaponSkill INT NOT NULL,
    strength INT NOT NULL,
    endurance INT NOT NULL,
    agility INT NOT NULL,
    hitpoints INT NOT NULL,
    isCharAlive BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS battleChars(
    id SERIAL,
    charName VARCHAR(255) NOT NULL,
    player BIGINT NOT NULL,
    weaponSkill INT NOT NULL,
    balisticSkill INT NOT NULL,
    strength INT NOT NULL,
    endurance INT NOT NULL,
    agility INT NOT NULL,
    hitpoints INT NOT NULL,
    isFighting BOOLEAN,
    isDodging BOOLEAN,
    isFleeing BOOLEAN
);

CREATE TABLE IF NOT EXISTS duelPreparation (
    id SERIAL,
	selectingPlayer    VARCHAR(255) NOT NULL,
    isReady         INT NOT NULL,
    isOver          BOOLEAN NOT NULL,
    turn            Int
);

CREATE TABLE IF NOT EXISTS duelPlayers(
    id SERIAL,
    preparationId INT NOT NULL,
    challenger VARCHAR(255) NOT NULL,
	challenged VARCHAR(255) NOT NULL,
    challengerChar VARCHAR(255),
	challengedChar VARCHAR(255)
);