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

CREATE TABLE IF NOT EXISTS dungeon(
    dungeon_id              SERIAL,
    created_at              timestamp NOT NULL DEFAULT NOW(),
    created_by              BIGINT NOT NULL,
	selected_character_id   INT DEFAULT null,
    has_started             BOOLEAN DEFAULT false,
	has_ended               BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS dungeon_tile(
    tile_id                 SERIAL,
    dungeon_id              INT NOT NULL,
    x                       INT NOT NULL,
    y                       INT NOT NULL,
	is_discovered           BOOLEAN DEFAULT false,
    is_exit                 BOOLEAN DEFAULT false,
	is_impassable           BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS link_character_tile(
    tile_id                 SERIAL,
    character_instance_id   INT NOT NULL
);

CREATE TABLE IF NOT EXISTS character_instance(
    character_instance_id       SERIAL,
    character_model_id          INT NOT NULL,
    precision                   INT NOT NULL,
    strength                    INT NOT NULL,
    endurance                   INT NOT NULL,
    agility                     INT NOT NULL,
    hitpoints                   INT NOT NULL,
    chosen_action_id            INT
);

CREATE TABLE IF NOT EXISTS character_model(
    character_model_id      SERIAL,
    name                    VARCHAR(255) NOT NULL,
    player_id               BIGINT NOT NULL,
    precision               INT NOT NULL,
    strength                INT NOT NULL,
    endurance               INT NOT NULL,
    agility                 INT NOT NULL,
    hitpoints               INT NOT NULL,
    is_occupied             BOOLEAN DEFAULT false,
    is_alive                BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS link_entity_tile(
    tile_id                 SERIAL,
    entity_instance_id   INT NOT NULL
);

CREATE TABLE IF NOT EXISTS entity_instance(
    entity_instance_id      SERIAL,
    entity_model_id         INT NOT NULL,
    precision               INT NOT NULL,
    strength                INT NOT NULL,
    endurance               INT NOT NULL,
    agility                 INT NOT NULL,
    hitpoints               INT NOT NULL,
    chosen_action_id        INT
);

CREATE TABLE IF NOT EXISTS entity_model(
    entity_model_id         SERIAL,
    name                    VARCHAR(255) NOT NULL,
    precision               INT NOT NULL,
    strength                INT NOT NULL,
    endurance               INT NOT NULL,
    agility                 INT NOT NULL,
    hitpoints               INT NOT NULL,
    is_alive                BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS link_event_tile(
    tile_id                 SERIAL,
    entity_instance_id      INT NOT NULL
);

CREATE TABLE IF NOT EXISTS event_model(
    event_id       SERIAL,
    event_name     VARCHAR(255) NOT NULL
);
