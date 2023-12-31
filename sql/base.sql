CREATE TABLE user (
    id integer PRIMARY KEY autoincrement,
    username varchar (128),
    password varchar (32),
    role_id integer
);

CREATE TABLE log (
    id integer PRIMARY KEY autoincrement,
    time integer,
    level varchar (256),
    user varchar (256),
    action varchar (256),
    position varchar (256),
    message varchar (256),
    errcode integer
);

CREATE TABLE task (
    id integer PRIMARY KEY autoincrement,
    pvob varchar (256),
    component varchar (256),
    dir varchar (256),
    keep varchar (256),
    cc_user varchar (256),
    cc_password varchar (256),
    git_url varchar (256),
    git_user varchar (256),
    git_password varchar (256),
    git_email varchar (256),
    status varchar (16),
    last_completed_date_time varchar (64),
    creator varchar(128),
    worker_id integer,
    include_empty boolean
);

CREATE TABLE match_info (
    id integer PRIMARY KEY autoincrement,
    task_id integer,
    stream varchar (256),
    git_branch varchar (256)
);

CREATE TABLE task_log (
    log_id integer PRIMARY KEY autoincrement,
    task_id integer,
    status varchar (16),
    start_time varchar (64),
    end_time varchar (64),
    duration integer
);

CREATE TABLE task_command_out (
    log_id integer PRIMARY KEY,
    content text
);

CREATE TABLE worker (
    id integer PRIMARY KEY autoincrement,
    worker_url varchar (256),
    status varchar (16),
    task_count integer,
    register_time varchar (64)
);

CREATE TABLE schedule (
    id integer PRIMARY KEY autoincrement,
    status varchar (16),
    schedule varchar (16),
    task_id integer,
    creator varchar (128)
);

INSERT INTO user (username,password,role_id) VALUES('admin', 'b17eccdc6c06bd8e15928d583503adf9', 1);
