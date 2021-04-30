/* eslint-disable */
export namespace User {
  export type Base = {
    username: string;
    password?: string;
    role_id?: Role;
  };

  export enum Role {
    ADMIN = 1,
    NORMAL = 2,
  }
}

export namespace Task {
  export type Base = {
    pvob: string;
    component: string;
    ccUser: string;
    ccPassword: string;
    gitURL: string;
    gitUser: string;
    gitPassword: string;
  };

  export type Item = {
    id: number;
    lastCompleteDateTime: string;
    status: string;
  } & Base;

  export type Detail = {
    taskModel: Base & { matchInfo: MatchInfo[] };
    logList: Log[];
  };
  export type MatchInfo = { stream: string; gitBranch: string };

  export type Log = {
    logID: string;
    status: string;
    startTime: string;
    endTime: string;
    duration: string;
  };
}
