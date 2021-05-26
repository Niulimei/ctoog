ALTER TABLE task ADD COLUMN svnUrl varchar(256) DEFAULT '';
ALTER TABLE task ADD COLUMN modelType varchar(256) DEFAULT '';

CREATE TABLE svn_name_pair (
    id integer PRIMARY KEY autoincrement,
    task_id integer,
    svn_username varchar (256),
    git_username varchar (256),
    git_email varchar (256)
);