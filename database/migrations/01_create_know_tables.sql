-- +goose Up
CREATE TABLE IF NOT EXISTS UrlReferences (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			originalUrl TEXT NOT NULL,
			short_code TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS UrlAdditonalParameters (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            value TEXT NOT NULL,
            type TEXT NOT NULL CHECK (type IN ("Query", "Route")),
            urlReferenceId INT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (urlReferenceId) REFERENCES UrlReferences(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Clicks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            vistedUrl TEXT NOT NULL,
            userAgent TEXT NULL,
            publicIp TEXT NULL,
            urlReferenceId INTEGER NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,   
            FOREIGN KEY (urlReferenceId) REFERENCES UrlReferences(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ClickAddtionalParamaters (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            clickId INTEGER NOT NULL,
            urlAdditonalParameterId INTEGER NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,  
            FOREIGN KEY (clickId) REFERENCES Clicks(id) ON DELETE CASCADE
            FOREIGN KEY (urlAdditonalParameterId) REFERENCES UrlAdditonalParameters(id)

);

CREATE TABLE IF NOT EXISTS ClickInvalidParamaters (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            clickId INTEGER NOT NULL,
            type TEXT NOT NULL CHECK (type IN ("Query", "Route")),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (clickId) REFERENCES Clicks(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS UrlAdditionalParameters;
DROP TABLE IF EXISTS UrlReferences;