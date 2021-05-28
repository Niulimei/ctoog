ALTER TABLE task ADD COLUMN svn_url varchar(256) DEFAULT '';
ALTER TABLE task ADD COLUMN model_type varchar(256) DEFAULT 'clearcase';
ALTER TABLE task ADD COLUMN gitignore text DEFAULT '';

CREATE TABLE svn_name_pair (
    id integer PRIMARY KEY autoincrement,
    task_id integer,
    svn_username varchar (256),
    git_username varchar (256),
    git_email varchar (256)
);
