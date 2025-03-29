-- migration.sql

-- Tabla para registrar equipos
CREATE TABLE IF NOT EXISTS teams (
   id SERIAL PRIMARY KEY,
   name VARCHAR(100) NOT NULL
);

-- Crear la tabla de partidos
CREATE TABLE IF NOT EXISTS matches (
   id SERIAL PRIMARY KEY,
   home_team VARCHAR(100) NOT NULL,
   away_team VARCHAR(100) NOT NULL,
   score_a INT NOT NULL,
   score_b INT NOT NULL,
   match_date DATE,
   extra_time INT DEFAULT 0
);

-- Tabla para registrar jugadores
CREATE TABLE IF NOT EXISTS players (
   id SERIAL PRIMARY KEY,
   name VARCHAR(100) NOT NULL,
   team_id INT NOT NULL,  -- Relacionado con la tabla teams
   FOREIGN KEY (team_id) REFERENCES teams(id)
);

-- Tabla para registrar goles
CREATE TABLE IF NOT EXISTS goals (
   id SERIAL PRIMARY KEY,
   match_id INT REFERENCES matches(id),
   team_id INT REFERENCES teams(id),  -- Relacionado con el equipo local o visitante
   goals INT
);

-- Tabla para registrar tarjetas (amarillas y rojas)
CREATE TABLE IF NOT EXISTS cards (
   id SERIAL PRIMARY KEY,
   player_id INT NOT NULL,
   match_id INT NOT NULL,
   card_type VARCHAR(10) NOT NULL,  -- Tipo de tarjeta: 'yellow' o 'red'
   card_time INT NOT NULL,  -- El minuto exacto en el que se dio la tarjeta
   FOREIGN KEY (player_id) REFERENCES players(id),
   FOREIGN KEY (match_id) REFERENCES matches(id)
);

-- Tabla para registrar tiempo extra
CREATE TABLE IF NOT EXISTS extra_time (
   id SERIAL PRIMARY KEY,
   match_id INT REFERENCES matches(id),
   time INT
);
