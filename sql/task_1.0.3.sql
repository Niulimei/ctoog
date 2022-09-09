ALTER TABLE task ADD COLUMN gitlab_group varchar(256) DEFAULT '';
ALTER TABLE task ADD COLUMN gitlab_project varchar(256) DEFAULT '';
ALTER TABLE task ADD COLUMN gitlab_token varchar(256) DEFAULT '';
ALTER TABLE task ADD COLUMN gitee_group varchar(256) DEFAULT '';
ALTER TABLE task ADD COLUMN gitee_token varchar(256) DEFAULT '';
ALTER TABLE task ADD COLUMN gitee_project varchar(256) DEFAULT '';
ALTER TABLE task ADD COLUMN source_url varchar(256) DEFAULT '';
ALTER TABLE task ADD COLUMN target_url varchar(256) DEFAULT '';
