CREATE TABLE users (
	id UUID PRIMARY KEY,
	login VARCHAR(100) NOT NULL UNIQUE,
	password TEXT NOT NULL, 
	email VARCHAR(255) UNIQUE,
	name VARCHAR(100) NOT NULL DEFAULT 'user',
	created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE desks (
	id UUID PRIMARY KEY,
	name VARCHAR(100) NOT NULL DEFAULT 'new desk',
	password TEXT NOT NULL,
	owner_id UUID NOT NULL,
	created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc'),

	CONSTRAINT fk_owner FOREIGN KEY (owner_id) 
		REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE desk_members (
	user_id UUID NOT NULL,
	desk_id UUID NOT NULL,
			
	PRIMARY KEY (user_id, desk_id),

	CONSTRAINT fk_member_user FOREIGN KEY (user_id) 
		REFERENCES users(id) ON DELETE CASCADE,
			
	CONSTRAINT fk_member_desk FOREIGN KEY (desk_id) 
		REFERENCES desks(id) ON DELETE CASCADE
);

CREATE TABLE tasks (
	id UUID PRIMARY KEY,
	author_id UUID NOT NULL,
	desk_id UUID NOT NULL,
	name VARCHAR(255) NOT NULL,
	description TEXT DEFAULT '',
	done BOOLEAN NOT NULL DEFAULT FALSE,
	deadline TIMESTAMPTZ,
	created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc'),

	CONSTRAINT fk_task_desk FOREIGN KEY (desk_id) 
		REFERENCES desks(id) ON DELETE CASCADE,

	CONSTRAINT fk_task_author FOREIGN KEY (author_id) 
		REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_desk_members_user_desk ON desk_members(user_id, desk_id);

CREATE INDEX IF NOT EXISTS idx_tasks_desk_created ON tasks(desk_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_tasks_author_id ON tasks(author_id);
